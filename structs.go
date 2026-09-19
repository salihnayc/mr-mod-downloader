package main

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
