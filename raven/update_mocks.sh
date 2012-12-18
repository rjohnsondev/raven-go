#!/bin/sh

# raven.HttpClient
$GOPATH/bin/mockgen -source="raven.go" \
                    -package="raven" \
                    -destination="mock_raven.go"
