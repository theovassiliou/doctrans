#!/bin/bash
golangci-lint run -v -j 4 -E deadcode -E depguard -E dogsled \
              -E errcheck -E goconst -E golint -E gosec -E gosimple -E govet -E exportloopref -E whitespace \
              -E goprintffuncname
              
golangci-lint run -v -j 4 -E ineffassign -E gocritic -E nakedret \
              -E rowserrcheck -E staticcheck -E structcheck -E typecheck -E unconvert -E unused -E varcheck