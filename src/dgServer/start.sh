#! /bin/bash 

killall dgServer
nohup ./dgServer > log.txt &
exit 0