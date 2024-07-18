#!/bin/bash

# Kill the server process for main.go in the current directory
if [ -f "main.go" ]; then
    pid=$(lsof -i :8080 | grep main | awk '{print $2}')
    if [ -n "$pid" ]; then
        kill -9 $pid
        echo "Stopped server current directory"
    else
        echo "Server in current directory is not running"
    fi
fi

# Kill the server process for main.go in the current directory
if [ -f "cal/main.go" ]; then
    pid=$(lsof -i :8081 | grep main | awk '{print $2}')
    if [ -n "$pid" ]; then
        kill -9 $pid
        echo "Stopped server in cal directory"
    else
        echo "Server in cal directory is not running"
    fi
fi

# Kill the server process for main.go in the current directory
if [ -f "todo/main.go" ]; then
    pid=$(lsof -i :8082 | grep main | awk '{print $2}')
    if [ -n "$pid" ]; then
        kill -9 $pid
        echo "Stopped server in todo directory"
    else
        echo "Server in todo directory is not running"
    fi
fi
