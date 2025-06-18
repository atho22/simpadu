#!/bin/bash

if pgrep -f 'binary_app'
then
    pgrep -f 'binary_app' | xargs kill &&  echo "baru stop" > ./out_test.txt
else
    echo "udah di-stop" > ./out.txt
fi
