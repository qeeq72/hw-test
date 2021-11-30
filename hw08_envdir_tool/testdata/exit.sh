#!/usr/bin/env bash

if [[ $1 -eq 1 && ${PORT} -gt 8000 ]];
then
    echo "Equal"
    exit ${TAB}
else
    echo "Not equal"
    exit 1
fi
