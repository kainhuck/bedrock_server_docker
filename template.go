package main

import (
	_ "embed"
)

type DockerCompose struct {
	Image      string
	InstallDir string
}

type PermissionsJson struct {
	XUID string
}

type ServerProperties struct {
	Mode       string // "survival", "creative", or "adventure"
	Difficulty string // "peaceful", "easy", "normal", or "hard"
	WorldName  string
	WorldSeed  string
}

var (
	//go:embed template/Dockerfile
	DockerfileTemp string

	//go:embed template/docker-compose.yml
	DockercomposeTemp string

	//go:embed template/permissions.json
	PermissionsJsonTemp string

	//go:embed template/server.properties
	ServerPropertiesTemp string
)
