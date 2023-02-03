package main

const (
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
		image: {{.image}}
		container_name: bedrock_server
		volumes:
		  - {{.install_root}}/minecraft/worlds:/mc/worlds
		  - {{.install_root}}/minecraft/server.properties:/mc/server.properties
		  - {{.install_root}}/minecraft/permissions.json:/mc/permissions.json
		ports:
		  - 19132:19132/udp
		restart: always
	`

	ServerPropertiesTemp = ``
)