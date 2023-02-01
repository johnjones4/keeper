PROJECT=$(shell basename $(shell pwd))
TAG=ghcr.io/johnjones4/${PROJECT}
VERSION=$(shell date +%s)

.PHONY: cli

info:
	echo ${PROJECT} ${VERSION}

container:
	docker build -t ${TAG} ./server
	docker push ${TAG}:latest
	docker image rm ${TAG}:latest

# cli:
# 	mkdir build || true
# 	rm build/* || true
# 	cd cli && go build -o ../build/note .

# install:
# 	mv -f build/note /usr/local/bin/note

ui:
	cd frontend && npm install
	cd frontend && npm run build
	tar zcvf ui.tar.gz ./frontend/build
	git tag ${VERSION}
	git push origin ${VERSION}
	gh release create ${VERSION} ui.tar.gz --generate-notes

ci: container ui 
