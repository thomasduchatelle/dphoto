#!/bin/sh -l

echo "::set-output name=go-build::$(go env GOCACHE)"
echo "::set-output name=go-mod::$(go env GOMODCACHE)"