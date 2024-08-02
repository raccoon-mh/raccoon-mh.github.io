#!/bin/bash

# raccoon-mh-playground 프로그램 경로
BIN_PATH="./raccoon-mh-playground"
# 로그 파일 경로
LOG_FILE="./raccoon-mh-playground.log"
# PID 파일 경로
PID_FILE="./raccoon-mh-playground.pid"

start() {
    echo "Starting raccoon-mh-playground ..."
    nohup $BIN_PATH > $LOG_FILE 2>&1 &
    echo $! > $PID_FILE
    echo "raccoon-mh-playground started with PID $(cat $PID_FILE)"
}

stop() {
    if [ -f $PID_FILE ]; then
        PID=$(cat $PID_FILE)
        echo "Stopping raccoon-mh-playground with PID $PID..."
        kill $PID
        rm $PID_FILE
        echo "raccoon-mh-playground stopped."
    else
        echo "PID file not found. Is the raccoon-mh-playground running?"
    fi
}

status() {
    if [ -f $PID_FILE ]; then
        PID=$(cat $PID_FILE)
        if ps -p $PID > /dev/null; then
            echo "raccoon-mh-playground is running with PID $PID."
        else
            echo "raccoon-mh-playground is not running, but PID file exists."
        fi
    else
        echo "raccoon-mh-playground is not running."
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 {start|stop|status}"
        exit 1
        ;;
esac
