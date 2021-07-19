#!/bin/sh
if [ $# != 1 ]; then
    echo "please input addr(e.g. 0x486)"
    exit 1
else
    curl -X POST -d '{"addr":"'$1'"}' -H "Content-Type: application/json" localhost:8888/debug/break
fi
