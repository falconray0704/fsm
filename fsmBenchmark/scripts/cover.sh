#!/usr/bin/env bash

#set -x
set -e
echo "" > cover.out


#    go test -race -coverprofile=profile.out $d
    go test -coverprofile=profile.out .
    if [ -f profile.out ]; then
        go tool cover -html=profile.out -o=cover.html
        cat profile.out >> cover.out
        rm profile.out
    fi
