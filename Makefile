ifndef GOBIN
GOBIN := $(GOPATH)/bin
endif

ifdef JANUS_PORT
JANUS_PORT := $(JANUS_PORT)
else
JANUS_PORT := 23889
endif

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
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
release: darwin linux windows

.PHONY: darwin
darwin: build-darwin-amd64 tar-gz-darwin-amd64 build-darwin-arm64 tar-gz-darwin-arm64

.PHONY: linux
linux: build-linux-386 tar-gz-linux-386 build-linux-amd64 tar-gz-linux-amd64 build-linux-arm tar-gz-linux-arm build-linux-arm64 tar-gz-linux-arm64 build-linux-ppc64 tar-gz-linux-ppc64 build-linux-ppc64le tar-gz-linux-ppc64le build-linux-mips tar-gz-linux-mips build-linux-mipsle tar-gz-linux-mipsle build-linux-riscv64 tar-gz-linux-riscv64 build-linux-s390x tar-gz-linux-s390x

.PHONY: windows
windows: build-windows-386 tar-gz-windows-386 build-windows-amd64 tar-gz-windows-amd64 build-windows-arm64 tar-gz-windows-arm64
	echo hey
#	GOOS=linux GOARCH=arm64 go build -o ./build/janus-linux-arm64 github.com/qtumproject/janus/cli/janus

docker-build-go-build:
	docker build -t qtum/go-build.janus -f ./docker/go-build.Dockerfile --build-arg GO_VERSION=$(GO_VERSION) .

tar-gz-%:
	mv $(ROOT_DIR)/build/bin/janus-$(shell echo $@ | sed s/tar-gz-// | sed 's/-/\n/' | awk 'NR==1')-$(shell echo $@ | sed s/tar-gz-// | sed 's/-/\n/' | awk 'NR==2') $(ROOT_DIR)/build/bin/janus
	tar -czf $(ROOT_DIR)/build/janus-$(GIT_TAG)-$(shell echo $@ | sed s/tar-gz-// | sed 's/-/\n/' | awk 'NR==1' | sed s/darwin/osx/)-$(shell echo $@ | sed s/tar-gz-// | sed 's/-/\n/' | awk 'NR==2').tar.gz $(ROOT_DIR)/build/bin/janus
	mv $(ROOT_DIR)/build/bin/janus $(ROOT_DIR)/build/bin/janus-$(shell echo $@ | sed s/tar-gz-// | sed 's/-/\n/' | awk 'NR==1')-$(shell echo $@ | sed s/tar-gz-// | sed 's/-/\n/' | awk 'NR==2')

# build-os-arch
build-%: docker-build-go-build
	docker run \
		--privileged \
		--rm \
		-v `pwd`/build:/build \
		-v `pwd`:$(JANUS_DIR) \
		-w $(JANUS_DIR) \
		-e GOOS=$(shell echo $@ | sed s/build-// | sed 's/-/\n/' | awk 'NR==1') \
		-e GOARCH=$(shell echo $@ | sed s/build-// | sed 's/-/\n/' | awk 'NR==2') \
		qtum/go-build.janus \
			build \
			-buildvcs=false \
			-ldflags \
				"-X 'github.com/qtumproject/janus/pkg/params.GitSha=`./sha.sh`'" \
			-o /build/bin/janus-$(shell echo $@ | sed s/build-// | sed 's/-/\n/' | awk 'NR==1')-$(shell echo $@ | sed s/build-// | sed 's/-/\n/' | awk 'NR==2') $(JANUS_DIR)

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
		--qtum-rpc=http://${test_user}:${test_user_passwd}@0.0.0.0:3889 \
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