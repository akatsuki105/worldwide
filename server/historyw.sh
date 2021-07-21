#!/bin/sh
if [ $# != 1 ]; then
    echo "please input history count(e.g. 0x20)"
    exit 1
else
    curl -X POST -d '{"history":"'$1'"}' -H "Content-Type: application/json" localhost:8888/debug/history
fi
