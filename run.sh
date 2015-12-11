#!/bin/bash
export DOCKER_PID=$(cat /var/run/docker.pid)
exec /tutum-agent