#! /bin/bash

source ./b

go test emu/core

for p in ../src/emu/plugins/*; do
    if [[ $(find ./$p -type f -name *"_test.go") ]]; then
        go test $p
        RESULT=$?
        if [[ $RESULT != 0 ]]; then
            echo Test of "$p" failed!
            exit $RESULT
        fi
    fi
done