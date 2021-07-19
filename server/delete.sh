#!/bin/sh
if [ $# != 1 ]; then
    echo "please input addr(e.g. 0x486)"
    exit 1
else
    curl -X DELETE "localhost:8888/debug/break?addr="$1
fi
