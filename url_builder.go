package main

import (
	"net/url"
	"strings"
)

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
