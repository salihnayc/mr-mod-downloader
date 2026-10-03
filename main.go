package main

import (
	"fmt"
)

func main() {
	client := NewModrinthClient("fabric", "26.3")

	var search Search
	if err := client.SearchMods("sodium", &search); err != nil {
		fmt.Println(err.Error())
		return
	}

	modID := search.Hits[1].ProjectID
	var versions ProjectsVersions
	if err := client.SearchProjectVersions(modID, &versions); err != nil {
		fmt.Println(err.Error())
		return
	}

	file := versions[0].Files[0]

	if err := client.DownloadFile(file, ""); err != nil {
		fmt.Println(err.Error())
		return
	}
}
