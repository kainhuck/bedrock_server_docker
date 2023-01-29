.PHONY: docker preview

docker:
	docker build -t kainhuck/bedrock_server:1.19.51.01 -f Dockerfile .  --network=host
	docker push kainhuck/bedrock_server:1.19.51.01

preview:
	docker build -t kainhuck/bedrock_server_preview:1.19.70.20 -f Dockerfile .  --network=host
	docker push kainhuck/bedrock_server_preview:1.19.70.20