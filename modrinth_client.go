package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func (mc ModrinthClient) downloadFile(req *http.Request, fileLocation string, fileName string) error {
	fmt.Printf("Doing Request %s\n", req.URL.String())

	resp, err := mc.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		out, err := os.Create(fileLocation + fileName)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}

	return nil
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

func (mc ModrinthClient) GetFile(file File, location string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := newRequest(ctx, file.Url)
	if err != nil {
		return err
	}

	err = mc.downloadFile(req, location, file.Filename)
	if err != nil {
		return err
	}

	return nil
}
