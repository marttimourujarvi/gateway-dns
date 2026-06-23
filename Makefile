deploy: build
	kind load docker-image gateway-dns:latest --name dev
build: fmt
	docker build . -t gateway-dns
fmt: vet
	go fmt ./...
vet:
	go vet ./...


