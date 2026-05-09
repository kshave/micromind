build-docker:
	docker build -t ghcr.io/kshave/micromind:latest .

publish-ghcr:
	docker push ghcr.io/kshave/micromind:latest