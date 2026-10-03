package main

import (
	"fmt"
	"net/url"
	"strings"
)

// -------
// Structs
// -------

type Project struct {
	ProjectID   string `json:"project_id"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Dependency struct {
	VersionID      string `json:"version_id"`
	ProjectID      string `json:"project_id"`
	DependencyType string `json:"dependency_type"`
}

type File struct {
	Url      string `json:"url"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	Hashes   struct {
		Sha512 string `json:"sha512"`
	}
}

type Search struct {
	Hits  []Project `json:"hits"`
	Limit int       `json:"limit"`
	Total int       `json:"total_hits"`
}

type Version struct {
	VersionID    string       `json:"id"`
	ProjectID    string       `json:"project_id"`
	Name         string       `json:"name"`
	Files        []File       `json:"files"`
	Dependencies []Dependency `json:"dependencies"`
}

type ProjectsVersions []Version

// -------------
// Minecraft Env
// ------------

// Minecraft environment info which mods downloaded for
type EnvInfo struct {
	modloader        string
	minecraftVersion string
}

// Creates a new env struct and returns it
func NewEnv(modloader string, minecraftVersion string) EnvInfo {
	return EnvInfo{modloader: modloader, minecraftVersion: minecraftVersion}
}

// ------
// Facets
// ------

// Facet for filtering search results
type FacetBuilder struct {
	env          EnvInfo
	project_type string
}

// Creates a new facet struct and returns it
func NewFacet(env EnvInfo, project_type string) *FacetBuilder {
	return &FacetBuilder{env: env, project_type: project_type}
}

// Converts the facet to string
func (fb *FacetBuilder) String() string {
	return fmt.Sprintf("[[\"categories:%s\"],[\"versions:%s\"],[\"project_type:%s\"]]", fb.env.modloader, fb.env.minecraftVersion, fb.project_type)
}

// -----------
// URL Builder
// -----------

// Struct for building URLs.
type UrlBuilder struct {
	baseUrl    string
	paths      []string
	parameters map[string]string
}

// Creates a new builder.
func NewUrlBuilder(baseUrl string, paths []string, parameters map[string]string) *UrlBuilder {
	ub := &UrlBuilder{}

	if strings.TrimSpace(baseUrl) == "" {
		ub.baseUrl = "https://api.modrinth.com/v2"
	} else {
		ub.baseUrl = baseUrl
	}

	if paths != nil {
		ub.paths = paths
	} else {
		ub.paths = make([]string, 0)
	}

	if parameters != nil {
		ub.parameters = parameters
	} else {
		ub.parameters = make(map[string]string)
	}

	return ub
}

// Adds the given path value to URL.
func (ub *UrlBuilder) AddPath(path string) {
	ub.paths = append(ub.paths, path)
}

// Adds the given parameter to URL.
func (ub *UrlBuilder) AddParameter(key string, value string) {
	ub.parameters[key] = value
}

// Converts the URL to string.
func (ub *UrlBuilder) String() (string, error) {
	u, err := url.Parse(ub.baseUrl)
	if err != nil {
		return "", err
	}

	for _, p := range ub.paths {
		u = u.JoinPath(p)
	}

	v := url.Values{}

	for k, p := range ub.parameters {
		v.Add(k, p)
	}

	u.RawQuery = v.Encode()

	return u.String(), nil
}
