#!/bin/bash

source ./b
cd src
go install $1 cmd/trex-emu.go
