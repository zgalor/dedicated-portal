#!/bin/bash
ret=0
bad_files=$(gofmt -s -l cmd/ pkg/)
if [[ -n "${bad_files}" ]]; then
    echo "gofmt needs to be run on the listed files"
    echo "${bad_files}"
    echo "Try running 'gofmt -w -d [path]'"
    ret=1
fi


bad_files=$(goimports -l cmd/ pkg/)
if [[ -n "${bad_files}" ]]; then
    echo "goimports needs to be run on the listed files"
    echo "${bad_files}"
    echo "Try running 'goimports -w -d [path]'"
    ret=1
fi

exit $ret
