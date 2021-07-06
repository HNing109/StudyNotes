#!/bin/bash
set -ue
cd `dirname $0`
protoc --proto_path=./:  --go_out=plugins=grpc:. *.proto