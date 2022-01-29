build:
	go build github.com/daniilty/sharenote-friends/cmd/server
build_docker:
	docker build -t sharenote-auth:latest -f docker/Dockerfile .
