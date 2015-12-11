#!/bin/bash

TOKEN=${TOKEN:-$1}
SERVER_HOST=${SERVER_HOST:-"https://dashboard.tutum.co"}
DOCKER_PID=${DOCKER_PID:-"/var/run/docker.pid"}

if [ -z "$TOKEN" ] ; then
	echo "* Please provide a tutum token"
	echo "  example: $0 <token>"
	exit 1
fi

if ! docker version > /dev/null 2>&1; then
	echo "=> Cannot connect to the Docker daemon. Is the docker daemon running on this host?"
	exit 1
fi

if [ -z "$DOCKER_VERSION" ]; then
	DOCKER_VERSION=$(docker version 2>/dev/null | grep -iF "Version" | head -n 1 | cut -d ":" -f2 | xargs)
fi

if [ -z "$DOCKER_VERSION" ]; then
	echo "=> Error: environment variable DOCKER_VERSION is empty, or docker client($DOCKER_BIN_FILE) is properly mounted"
	exit 1
fi

if [ ! -f ${DOCKER_PID} ]; then 
	echo "=> Cannot find docker pid ${DOCKER_PID}"
	exit 1
fi

echo "=> Pulling the latest version of tutum/agent image"
#docker pull tutum/agent

docker rm -vf tutum-agent > /dev/null 2>&1 || true

echo "=> Configuring docker"
docker run \
	--net host \
	--pid host \
	--privileged \
	--restart always \
	--env TOKEN="$TOKEN" \
	--env SERVER_HOST="$SERVER_HOST" \
	--env DOCKER_VERSION="$DOCKER_VERSION" \
	--volume /etc:/etc \
	--volume ${DOCKER_PID}:${DOCKER_PID} \
	--name tutum-agent \
	tutum/agent

echo "=> Restarting docker"
RETRIES=60
for (( i=0 ; ; i++ )); do
	if [ ${i} -eq ${RETRIES} ]; then
		echo "Time out. Please see docker logs and container tutum-agent logs for error messages"
		exit 1
	fi
	sleep 5
	docker ps > /dev/null 2>&1 && break
done

echo "=> Installation logs"
timeout 10s docker logs -f tutum-agent