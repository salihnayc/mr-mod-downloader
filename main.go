package main

import (
	"fmt"
	"net/http"
)

var client = &http.Client{}
var mr = NewModrinthClient(client)

func main() {
	env := NewEnv("fabric", "26.2")
	var sr Version

	err := mr.GetVersion(env, "wpTNXtBM", &sr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(sr)

	file := sr.Files[0]

	mr.GetFile(file, "./")
}
