#!/bin/sh
set -e
SRC_DIR=$(pwd)
BUILD_DIR=/build

mkdir -p ${BUILD_DIR} 2>/dev/null
cd ${BUILD_DIR}

echo "Cloning CoreDNS repo..."
git clone https://github.com/coredns/coredns.git

cd coredns
git checkout v1.8.4

git switch -c fns
mkdir plugin/fns
cp "${SRC_DIR}"/fns.go plugin/fns
cp "${SRC_DIR}"/rrs.go plugin/fns
cp "${SRC_DIR}"/setup.go plugin/fns
cp "${SRC_DIR}"/types.go plugin/fns

echo "Overwrite plugin config..."

echo "metadata:metadata
cancel:cancel
tls:tls
reload:reload
nsid:nsid
bufsize:bufsize
root:root
bind:bind
debug:debug
trace:trace
ready:ready
health:health
pprof:pprof
prometheus:metrics
errors:errors
log:log
acl:acl
loadbalance:loadbalance
cache:cache
rewrite:rewrite
dnssec:dnssec
hosts:hosts
file:file
fns:fns
forward:forward
whoami:whoami
on:github.com/coredns/caddy/onevent
sign:sign" > plugin.cfg

go mod tidy

echo "Building..."
make SHELL='sh -x' CGO_ENABLED=0 coredns