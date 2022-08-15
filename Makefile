ifndef GOBIN
GOBIN := $(GOPATH)/bin
endif

ifdef JANUS_PORT
JANUS_PORT := $(JANUS_PORT)
else
JANUS_PORT := 23889
endif

JANUS_DIR := "/go/src/github.com/qtumproject/janus"
GO_VERSION := "1.18"
ALPINE_VERSION := "3.16"
DOCKER_ACCOUNT := ripply


# Latest commit hash
GIT_SHA=$(shell git rev-parse HEAD)

# If working copy has changes, append `-local` to hash
GIT_DIFF=$(shell git diff -s --exit-code || echo "-local")
GIT_REV=$(GIT_SHA)$(GIT_DIFF)
GIT_TAG=$(shell git describe --tags 2>/dev/null)

ifeq ($(GIT_TAG),)
GIT_TAG := $(GIT_REV)
else
GIT_TAG := $(GIT_TAG)$(GIT_DIFF)
endif

check-env:
ifndef GOPATH
	$(error GOPATH is undefined)
endif

.PHONY: install
install: 
	go install \
		-ldflags "-X 'github.com/qtumproject/janus/pkg/params.GitSha=`./sha.sh``git diff -s --exit-code || echo \"-local\"`'" \
		github.com/qtumproject/janus

.PHONY: release
release: darwin linux

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./build/janus-darwin-amd64 github.com/qtumproject/janus/cli/janus

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build -o ./build/janus-linux-amd64 github.com/qtumproject/janus/cli/janus

.PHONY: quick-start
quick-start-regtest:
	cd docker && ./spin_up.regtest.sh && cd ..

.PHONY: quick-start-testnet
quick-start-testnet:
	cd docker && ./spin_up.testnet.sh && cd ..

.PHONY: quick-start-mainnet
quick-start-mainnet:
	cd docker && ./spin_up.mainnet.sh && cd ..

# docker build -t qtum/janus:latest -t qtum/janus:dev -t qtum/janus:${GIT_TAG} -t qtum/janus:${GIT_REV} --build-arg BUILDPLATFORM="$(BUILDPLATFORM)" .
.PHONY: docker-dev
docker-dev:
	docker build -t qtum/janus:latest -t qtum/janus:dev -t qtum/janus:${GIT_TAG} -t qtum/janus:${GIT_REV} --build-arg GO_VERSION=1.18 .

docker-buildx-inspect:
	docker run \
		--privileged \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:$(JANUS_DIR) \
		-w $(JANUS_DIR) \
		qtum/docker-buildx-bin \
			buildx imagetools inspect qtum/janus:dev

docker-build-multi-arch-images: docker-build-arch-amd64 docker-build-arch-arm64
docker-push-multi-arch-images: push-docker-build-arch-amd64 push-docker-build-arch-arm64
docker-release-multi-arch: docker-build-multi-arch-images docker-push-multi-arch-images docker-manifest-create

docker-manifest-create:
	docker manifest create \
		$(DOCKER_ACCOUNT)/janus:latest \
		--amend $(DOCKER_ACCOUNT)/janus:latest-amd64 \
		--amend $(DOCKER_ACCOUNT)/janus:latest-arm64

docker-build-arch-%:
	docker build \
		--platform linux/$(shell echo $@ | sed 's/docker-build-arch-//') \
		-t $(DOCKER_ACCOUNT)/janus:latest-$(shell echo $@ | sed 's/docker-build-arch-//') \
		-t $(DOCKER_ACCOUNT)/janus:dev-$(shell echo $@ | sed 's/docker-build-arch-//') \
		-t $(DOCKER_ACCOUNT)/janus:${GIT_TAG}-$(shell echo $@ | sed 's/docker-build-arch-//') \
		-t $(DOCKER_ACCOUNT)/janus:${GIT_REV}-$(shell echo $@ | sed 's/docker-build-arch-//') \
		--build-arg GO_VERSION=${GO_VERSION} \
		--build-arg ALPINE_VERSION=${ALPINE_VERSION} \
		--build-arg BBUILDPLATFORM=linux/$(shell uname -p) \
		--build-arg BBBUILDPLATFORM=linux/$(shell echo $@ | sed 's/docker-build-arch-//') \
		--build-arg TARGETARCH=$(shell echo $@ | sed 's/docker-build-arch-//') \
		--build-arg TARGETPLATFORM=linux/$(shell echo $@ | sed 's/docker-build-arch-//') \
		.

push-docker-build-arch-%:
	docker push $(DOCKER_ACCOUNT)/janus:${GIT_TAG}-$(shell echo $@ | sed 's/push-docker-build-arch-//')
	docker push $(DOCKER_ACCOUNT)/janus:latest-$(shell echo $@ | sed 's/push-docker-build-arch-//')

.PHONY: docker-release
# Builds multiarch image for release
docker-release: docker-build-buildx
	docker run \
		--privileged \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:$(JANUS_DIR) \
		-w $(JANUS_DIR) \
		qtum/docker-buildx-bin
	docker run \
		--privileged \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:$(JANUS_DIR) \
		-w $(JANUS_DIR) \
		qtum/docker-buildx-bin \
			buildx \
			build \
			--platform linux/arm64,linux/arm64 \
			-t $(DOCKER_ACCOUNT)/janus:latest \
			-t $(DOCKER_ACCOUNT)/janus:dev \
			-t $(DOCKER_ACCOUNT)/janus:${GIT_TAG} \
			-t $(DOCKER_ACCOUNT)/janus:${GIT_REV} \
			--build-arg GO_VERSION=${GO_VERSION} \
			--build-arg ALPINE_VERSION=${ALPINE_VERSION} \
			.
	
.PHONY: local-dev
local-dev: check-env install
	docker run --rm --name qtum_testchain -d -p 3889:3889 qtum/qtum qtumd -regtest -rpcbind=0.0.0.0:3889 -rpcallowip=0.0.0.0/0 -logevents=1 -rpcuser=qtum -rpcpassword=testpasswd -deprecatedrpc=accounts -printtoconsole | true
	sleep 3
	docker cp ${GOPATH}/src/github.com/qtumproject/janus/docker/fill_user_account.sh qtum_testchain:.
	docker exec qtum_testchain /bin/sh -c ./fill_user_account.sh
	QTUM_RPC=http://qtum:testpasswd@localhost:3889 QTUM_NETWORK=auto $(GOBIN)/janus --port $(JANUS_PORT) --accounts ./docker/standalone/myaccounts.txt --dev

.PHONY: local-dev-https
local-dev-https: check-env install
	docker run --rm --name qtum_testchain -d -p 3889:3889 qtum/qtum qtumd -regtest -rpcbind=0.0.0.0:3889 -rpcallowip=0.0.0.0/0 -logevents=1 -rpcuser=qtum -rpcpassword=testpasswd -deprecatedrpc=accounts -printtoconsole | true
	sleep 3
	docker cp ${GOPATH}/src/github.com/qtumproject/janus/docker/fill_user_account.sh qtum_testchain:.
	docker exec qtum_testchain /bin/sh -c ./fill_user_account.sh > /dev/null&
	QTUM_RPC=http://qtum:testpasswd@localhost:3889 QTUM_NETWORK=auto $(GOBIN)/janus --port $(JANUS_PORT) --accounts ./docker/standalone/myaccounts.txt --dev --https-key https/key.pem --https-cert https/cert.pem

.PHONY: local-dev-logs
local-dev-logs: check-env install
	docker run --rm --name qtum_testchain -d -p 3889:3889 qtum/qtum:dev qtumd -regtest -rpcbind=0.0.0.0:3889 -rpcallowip=0.0.0.0/0 -logevents=1 -rpcuser=qtum -rpcpassword=testpasswd -deprecatedrpc=accounts -printtoconsole | true
	sleep 3
	docker cp ${GOPATH}/src/github.com/qtumproject/janus/docker/fill_user_account.sh qtum_testchain:.
	docker exec qtum_testchain /bin/sh -c ./fill_user_account.sh
	QTUM_RPC=http://qtum:testpasswd@localhost:3889 QTUM_NETWORK=auto $(GOBIN)/janus --port $(JANUS_PORT) --accounts ./docker/standalone/myaccounts.txt --dev > janus_dev_logs.txt

.PHONY: unit-tests
unit-tests: check-env
	go test -v ./... -timeout 50s

# Translates BUILDPLATFORM into GOARCH/GOHOST
docker-build-xcputranslate:
	docker build -t qtum/xcputranslate -f ./docker/xcputranslate.Dockerfile --build-arg GO_VERSION=$(GO_VERSION) .

# Multiplatform builds
docker-build-buildx: docker-build-xcputranslate
	docker build -t qtum/docker-buildx-bin -f ./docker/buildx-bin.Dockerfile .

docker-build-unit-tests:
	docker build -t qtum/tests.janus -f ./docker/unittests.Dockerfile --build-arg GO_VERSION=$(GO_VERSION) .

docker-unit-tests:
	docker run --rm -v `pwd`:/go/src/github.com/qtumproject/janus qtum/tests.janus

docker-tests: docker-build-unit-tests docker-unit-tests openzeppelin-docker-compose

docker-configure-https: docker-configure-https-build
	docker/setup_self_signed_https.sh

docker-configure-https-build:
	docker build -t qtum/openssl.janus -f ./docker/openssl.Dockerfile ./docker

# -------------------------------------------------------------------------------------------------------------------
# NOTE:
# 	The following make rules are only for local test purposes
# 
# 	Both run-janus and run-qtum must be invoked. Invocation order may be independent, 
# 	however it's much simpler to do in the following order: 
# 		(1) make run-qtum 
# 			To stop Qtum node you should invoke: make stop-qtum
# 		(2) make run-janus
# 			To stop Janus service just press Ctrl + C in the running terminal

# Runs current Janus implementation
run-janus:
	@ printf "\nRunning Janus...\n\n"

	go run `pwd`/main.go \
		--qtum-rpc=http://${test_user}:${test_user_passwd}@0.0.0.0:4889 \
		--qtum-network=auto \
		--bind=0.0.0.0 \
		--port=23889 \
		--accounts=`pwd`/docker/standalone/myaccounts.txt \
		--log-file=janusLogs.txt \
		--dev

run-janus-https:
	@ printf "\nRunning Janus...\n\n"

	go run `pwd`/main.go \
		--qtum-rpc=http://${test_user}:${test_user_passwd}@0.0.0.0:3889 \
		--qtum-network=auto \
		--bind=0.0.0.0 \
		--port=23889 \
		--accounts=`pwd`/docker/standalone/myaccounts.txt \
		--log-file=janusLogs.txt \
		--dev \
		--https-key https/key.pem \
		--https-cert https/cert.pem

test_user = qtum
test_user_passwd = testpasswd

# Runs docker container of qtum locally and starts qtumd inside of it
run-qtum:
	@ printf "\nRunning qtum...\n\n"
	@ printf "\n(1) Starting container...\n\n"
	docker run ${qtum_container_flags} qtum/qtum qtumd ${qtumd_flags} > /dev/null

	@ printf "\n(2) Importing test accounts...\n\n"
	@ sleep 3
	docker cp ${shell pwd}/docker/fill_user_account.sh ${qtum_container_name}:.

	@ printf "\n(3) Filling test accounts wallets...\n\n"
	docker exec ${qtum_container_name} /bin/sh -c ./fill_user_account.sh > /dev/null
	@ printf "\n... Done\n\n"

seed-qtum:
	@ printf "\n(2) Importing test accounts...\n\n"
	docker cp ${shell pwd}/docker/fill_user_account.sh ${qtum_container_name}:.

	@ printf "\n(3) Filling test accounts wallets...\n\n"
	docker exec ${qtum_container_name} /bin/sh -c ./fill_user_account.sh
	@ printf "\n... Done\n\n"

qtum_container_name = test-chain

# TODO: Research -v
qtum_container_flags = \
	--rm -d \
	--name ${qtum_container_name} \
	-v ${shell pwd}/dapp \
	-p 3889:3889

# TODO: research flags
qtumd_flags = \
	-regtest \
	-rpcbind=0.0.0.0:3889 \
	-rpcallowip=0.0.0.0/0 \
	-logevents \
	-addrindex \
	-reindex \
	-txindex \
	-rpcuser=${test_user} \
	-rpcpassword=${test_user_passwd} \
	-deprecatedrpc=accounts \
	-printtoconsole

# Starts continuously printing Qtum container logs to the invoking terminal
follow-qtum-logs:
	@ printf "\nFollowing qtum logs...\n\n"
		docker logs -f ${qtum_container_name}

open-qtum-bash:
	@ printf "\nOpening qtum bash...\n\n"
		docker exec -it ${qtum_container_name} bash

# Stops docker container of qtum
stop-qtum:
	@ printf "\nStopping qtum...\n\n"
		docker kill `docker container ps | grep ${qtum_container_name} | cut -d ' ' -f1` > /dev/null
	@ printf "\n... Done\n\n"

restart-qtum: stop-qtum run-qtum

submodules:
	git submodules init

# Run openzeppelin tests, Janus/QTUM needs to already be running
openzeppelin:
	cd testing && make openzeppelin

# Run openzeppelin tests in docker
# Janus and QTUM need to already be running
openzeppelin-docker:
	cd testing && make openzeppelin-docker

# Run openzeppelin tests in docker-compose
openzeppelin-docker-compose:
	cd testing && make openzeppelin-docker-compose