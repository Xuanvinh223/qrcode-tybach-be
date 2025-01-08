build:
	go build -o bin/tyxuan-web-printlabel-api cmd/main.go

run:
	go run cmd/main.go

format:
	go fmt tyxuan-web-printlabel-api/...

build-docker:
	docker build . -t tyxuan-web-printlabel-api

run-docker:
	docker run -itd --name tyxuan-web-printlabel-api --restart always -p 8080:8080 tyxuan-web-printlabel-api

exec-docker:
	docker exec -it tyxuan-web-printlabel-api sh