#! /bin/bash

if ! [ -d ../.bin ]; then
 curl -sSfL -C - https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./.bin v1.54.2
else
  echo "golangci retrieved from cache"
fi