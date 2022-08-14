#!/bin/bash
killall go-admin # kill go-admin service
echo "stop go-admin success"
ps -aux | grep go-admin