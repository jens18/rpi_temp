#!/bin/bash

set -x

#
# Assumption: 'rpi_temp' project is located at: $GOPATH/src/$PACKAGE_ROOT/rpi_temp
#

PACKAGE_ROOT=github.com/jens18

# compile Go program
#
# generate static binaries to enable use of (very small) luxas/alpine Docker base image
#
# see the section on etcd and influxdb in the following Kubernetes-on-arm build script:
# https://github.com/luxas/kubernetes-on-arm/blob/master/images/kubernetesonarm/build/inbuild.sh
#

# NOTE: GOPATH is expected to to be set to this directory.

# download dependencies:
# go get github.com/gorilla/mux

# compile:
CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" rpi_temp.go 
CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" kube_temp.go

# 'go install does produce the following error:
# go install net: open /usr/local/go/pkg/linux_amd64/net.a: permission denied

# install:

# ensure bin directory exists
mkdir -p $GOPATH/bin
# copy executable file
cp rpi_temp $GOPATH/bin

# generate Docker image

# 'uname -m' to identify the required Dockerfile:
# RPI 1: armv6l
# RPI 3: armv7l
# PC: x86_64

# RPI 3 and RPI 1 can use the same Dockerfile: Dockerfile.armv6l
MACHINE=`uname -m`
if [ $MACHINE == "armv7l" ]
then
    MACHINE="armv6l"
fi

docker build -f Dockerfile.$MACHINE -t rpi_temp .





