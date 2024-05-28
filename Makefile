PROJECT=$(shell basename $(shell pwd))
TAG=ghcr.io/johnjones4/${PROJECT}
VERSION=$(shell date +%s)

info:
	echo ${PROJECT} ${VERSION}

container:
	docker build --platform linux/amd64 -t ${TAG} ./server
	docker push ${TAG}:latest
	docker image rm ${TAG}:latest

ui:
	rm -rf ./frontend/dist/spa || true
	cd frontend && corepack yarn
	cd frontend && corepack yarn run build
	tar zcvf ui.tar.gz ./frontend/dist/spa
	git tag ${VERSION}
	git push origin ${VERSION}
	gh release create ${VERSION} ui.tar.gz --generate-notes

ci: container ui 

todoedit:
	mkdir bin || true
	cd ./terminal/todoedit/ && go build -o ../../bin/todoedit .

install:
	mv ./bin/todoedit /usr/local/bin/todoedit

tools: todoedit
