# base image
FROM ubuntu:latest

# configure env
ENV DEBIAN_FRONTEND 'noninteractive'

# update apt, install core apt dependencies and delete the apt-cache
# note: this is done in one command in order to keep down the size of intermediate containers
RUN apt update && \
    apt install -y locales git-core sqlite3 libsqlite3-mod-spatialite golang && \
    rm -rf /var/lib/apt/lists/*

# configure locale
RUN locale-gen 'en_US.UTF-8'
ENV LANG 'en_US.UTF-8'
ENV LANGUAGE 'en_US:en'
ENV LC_ALL 'en_US.UTF-8'

# configure git
RUN git config --global 'user.email' 'null@null.com'
RUN git config --global 'user.name' 'Missinglink PBF'

# set GOPATH
ENV GOPATH='/tmp/go'

# change working dir
WORKDIR "$GOPATH/src/github.com/missinglink/pbf"

# copy files
COPY . "$GOPATH/src/github.com/missinglink/pbf"

# fetch dependencies
RUN go get

# build binary
RUN go build
