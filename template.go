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
	DockerfileTemp = `FROM ubuntu
RUN apt-get -y update && apt-get install -y wget unzip curl

WORKDIR /mc

RUN wget {{.}} -O bedrock-server.zip 

RUN unzip bedrock-server.zip && rm bedrock-server.zip

ENV LD_LIBRARY_PATH=/mc

CMD ["./bedrock_server"]`

	DockercomposeTemp = `version: "3.7"
services:
  registry:
    image: {{.Image}}
    container_name: bedrock_server
    volumes:
      - {{.InstallDir}}/worlds:/mc/worlds
      - {{.InstallDir}}/server.properties:/mc/server.properties
      - {{.InstallDir}}/permissions.json:/mc/permissions.json
    ports:
      - 19132:19132/udp
    restart: always`

	PermissionsJsonTemp = `[
	{
		"permission": "operator",
		"xuid": "{{.XUID}}"
	}
]
`
	//go:embed server.properties
	ServerPropertiesTemp string
)
