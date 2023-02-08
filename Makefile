PROJECT=$(shell basename $(shell pwd))
TAG=ghcr.io/johnjones4/${PROJECT}
VERSION=$(shell date +%s)

info:
	echo ${PROJECT} ${VERSION}

container:
	docker build -t ${TAG} ./server
	docker push ${TAG}:latest
	docker image rm ${TAG}:latest

ui:
	cd frontend && npm install
	cd frontend && npm run build
	tar zcvf ui.tar.gz ./frontend/build
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
