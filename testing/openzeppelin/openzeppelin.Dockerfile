FROM node:lts-bullseye

RUN npm install -g truffle --loglevel verbose
RUN npm install -g mocha-spec-json-output-reporter
RUN mkdir -p openzeppelin-contracts
WORKDIR /openzeppelin-contracts
COPY ./openzeppelin-contracts/package.json /openzeppelin-contracts
COPY ./openzeppelin-contracts/hardhat /openzeppelin-contracts/hardhat
COPY ./openzeppelin-contracts/hardhat.config.js /openzeppelin-contracts
RUN yarn install
COPY ./openzeppelin-contracts/contracts /openzeppelin-contracts/contracts
COPY ./truffle-config.js /
RUN truffle compile
COPY ./openzeppelin-contracts/scripts /openzeppelin-contracts/scripts
COPY ./openzeppelin-contracts/ /openzeppelin-contracts
COPY ./truffle-config.js /

CMD [ "truffle", "test" ]
