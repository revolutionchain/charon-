#!/bin/sh
docker-compose -f ${GOPATH}/src/github.com/revolutionchain/charon/docker/quick_start/docker-compose.testnet.yml up -d
# sleep 3 #executing too fast causes some errors
# docker cp ${GOPATH}/src/github.com/revolutionchain/charon/docker/fill_user_account.sh qtum_testchain:.
# docker exec qtum_testnet /bin/sh -c ./fill_user_account.sh