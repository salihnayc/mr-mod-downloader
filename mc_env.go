package main

import "fmt"

// Minecraft environment info which mods downloaded for
type EnvInfo struct {
	modloader        string
	minecraftVersion string
}

// Creates a new env struct and returns it
func NewEnv(modloader string, minecraftVersion string) EnvInfo {
	return EnvInfo{modloader: modloader, minecraftVersion: minecraftVersion}
}

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
