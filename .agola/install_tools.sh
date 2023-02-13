#!/bin/bash
# Installl inspect tools

# TODO: use a function or something
[ ! -f "$GOPATH/bin/revive" ] && echo "Revive not found, installing..." && go install github.com/mgechev/revive@latest
[ ! -f "$GOPATH/bin/gosec" ] && echo "Gosec not found, installling..." && go install github.com/securego/gosec/v2/cmd/gosec@latest
[ ! -f "$GOPATH/bin/staticcheck" ] && echo "Staticcheck not found, installling..." && go install honnef.co/go/tools/cmd/staticcheck@latest
