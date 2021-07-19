#!/bin/sh
if [ $# != 2 ]; then
    echo "please input target and value(e.g. ime, 0x1)"
    exit 1
else
    curl -X POST -d '{"target":"'$1'", "value":"'$2'"}' -H "Content-Type: application/json" localhost:8888/debug/register
fi
