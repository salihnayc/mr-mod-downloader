package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ModrinthClient struct {
	client *http.Client
	env    EnvInfo
}

func NewModrinthClient(modloader string, minecraftVersion string) ModrinthClient {
	client := &http.Client{}
	env := NewEnv(modloader, minecraftVersion)
	return ModrinthClient{client: client, env: env}
}

// Creates a request to given url.
func newRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Set user agent to project link
	req.Header.Set("User-Agent", "github.com/salihnayc/mr-mod-downloader")

	return req, nil
}

// Does the given request and stores returned json in the given struct
func (mc ModrinthClient) doRequest(req *http.Request, s any) error {
	fmt.Printf("Doing Request %s\n", req.URL.String())

	resp, err := mc.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(s)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("Server Returned %d", resp.StatusCode)
	}
}

// Creates a search request with the given search query. Stores the returned json in the given struct.
func (mc ModrinthClient) SearchMods(query string, s *Search) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	facet := NewFacet(mc.env, "mod")

	ub := NewUrlBuilder("", nil, nil)

	ub.AddPath("search")
	ub.AddParameter("query", query)
	ub.AddParameter("facets", facet.String())

	u, err := ub.String()
	if err != nil {
		return err
	}

	req, err := newRequest(ctx, u)
	if err != nil {
		return err
	}

	err = mc.doRequest(req, s)
	if err != nil {
		return err
	}

	return nil
}

func (mc ModrinthClient) SearchProjectVersions(modID string, s *ProjectsVersions) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ub := NewUrlBuilder("", nil, nil)

	ub.AddPath("project")
	ub.AddPath(modID)
	ub.AddPath("version")

	ub.AddParameter("loaders", fmt.Sprintf("[\"%s\"]", mc.env.modloader))
	ub.AddParameter("game_versions", fmt.Sprintf("[\"%s\"]", mc.env.minecraftVersion))
	ub.AddParameter("include_changelog", "false")

	u, err := ub.String()
	if err != nil {
		return err
	}

	req, err := newRequest(ctx, u)
	if err != nil {
		return err
	}

	err = mc.doRequest(req, s)
	if err != nil {
		return err
	}

	return nil
}

func (mc ModrinthClient) GetVersion(projectID string, s *Version) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ub := NewUrlBuilder("", nil, nil)

	ub.AddPath("version")
	ub.AddPath(projectID)

	u, err := ub.String()
	if err != nil {
		return err
	}

	req, err := newRequest(ctx, u)
	if err != nil {
		return err
	}

	err = mc.doRequest(req, s)
	if err != nil {
		return err
	}

	return nil
}

// A custom Writer
type ProgressTracker struct {
	Total      int
	Downloaded int
}

func (pt *ProgressTracker) Write(p []byte) (int, error) {
	n := len(p)
	pt.Downloaded += n
	pt.showProgress()

	return n, nil
}

func (pt *ProgressTracker) showProgress() {
	downloadMB := float64(pt.Downloaded) / (1024 * 1024)
	totalMB := float64(pt.Total) / (1024 * 1024)
	percent := (float64(pt.Downloaded) / float64(pt.Total)) * 100
	fmt.Printf("\rProgress: %05.2fMB/%05.2fMB (%06.2f%%)", downloadMB, totalMB, percent)
}

// Downloads the file to the path. If the location parameter is empty mods is
func (mc ModrinthClient) DownloadFile(file File, location string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := newRequest(ctx, file.Url)
	if err != nil {
		return err
	}

	fmt.Printf("Doing Request %s\n", req.URL.String())

	resp, err := mc.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Server Returned %v", resp.StatusCode)
	}

	if strings.TrimSpace(location) == "" {
		location = "mods"
	}

	stat, err := os.Stat(location)

	if errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(location, 0750); err != nil {
			return err
		}
	} else if !stat.IsDir() {
		return fmt.Errorf("Given location %v is not a dir", location)
	}

	out, err := os.Create(filepath.Join(location, file.Filename))
	if err != nil {
		return err
	}
	defer out.Close()

	tracker := &ProgressTracker{file.Size, 0}

	fmt.Printf("Dowloading %v to %v directory\n", file.Filename, location)

	reader := io.TeeReader(resp.Body, tracker)

	_, err = io.Copy(out, reader)
	fmt.Println()

	return err
}
