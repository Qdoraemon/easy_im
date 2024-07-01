#! /bin/bash

go mod tidy
go build -o im && ./im 