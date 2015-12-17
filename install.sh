#!/bin/bash

set -e

TOKEN=${TOKEN:-$1}
SERVER_HOST=${SERVER_HOST:-"https://dashboard.tutum.co"}
WORKDIR="/etc/tutum"
DOCKER_OPTS="-H unix:///var/run/docker.sock -H tcp://0.0.0.0:2375 --tlscert /etc/tutum/cert.pem --tlskey /etc/tutum/key.pem --tlscacert /etc/tutum/ca.pem --tlsverify ${EXTRA_OPTS}"
UPSTART_DOCKER_CFG="/etc/default/docker"
SYSTEMD_DOCKER_CFG="/etc/sysconfig/docker"
SYSTEMD_DOCKER_SRV1="/usr/lib/systemd/system/docker.service"
SYSTEMD_DOCKER_SRV2="/lib/systemd/system/docker.service"

set_docker_env(){
	if [ -f $1 ]; then
		if cat $1 | grep '^DOCKER_OPTS' >/dev/null; then
			sed -i "s%^DOCKER_OPTS.*%DOCKER_OPTS=\"${DOCKER_OPTS}\"%" $1
		else
			echo "DOCKER_OPTS=\"${DOCKER_OPTS}\"" >> $1
		fi
	fi
}

set_docker_srv(){
	if [ -f ${UPSTART_DOCKER_CFG} ]; then
		ENVFILE=${UPSTART_DOCKER_CFG}
	elif [ -f ${SYSTEMD_DOCKER_CFG} ]; then
		ENVFILE=${SYSTEMD_DOCKER_CFG}
	else
		ENVFILE=${UPSTART_DOCKER_CFG}
		mkdir -p "$(dirname ${UPSTART_DOCKER_CFG})" && touch ${UPSTART_DOCKER_CFG}
		set_docker_env ${UPSTART_DOCKER_CFG}
	fi

	if [ -f $1 ]; then
		if ! cat $1 | grep "^ExecStart=" | grep -iF 'DOCKER_OPTS' > /dev/null; then
			sed -i '/^ExecStart=.*/s/$/ $DOCKER_OPTS/' $1
		fi
		if ! cat $1 | grep "EnvironmentFile=" > /dev/null; then
			sed -i "/\[Service\]/a EnvironmentFile=-$ENVFILE" $1
		fi
	fi
}

if [ -z "$TOKEN" ] ; then
	echo "* Please provide a tutum token"
	exit 1
fi

if [[ $EUID -ne 0 ]]; then
	echo "* This script must be run as root"
	exit 1
fi

if ! docker version > /dev/null 2>&1; then
	echo "=> Cannot connect to the Docker daemon. Is the docker daemon running on this host?"
	exit 1
fi

if [ -z "$DOCKER_VERSION" ]; then
	DOCKER_VERSION=$(docker version 2>/dev/null | grep -iF "Version" | head -n 1 | cut -d ":" -f2 | xargs)
fi

echo "=> Pulling the latest image"
docker pull tutum/agent

echo "=> Initializing"
rm -f ${WORKDIR}/*
mkdir -p ${WORKDIR}
docker rm -vf tutum-agent > /dev/null 2>&1 || true

echo "=> Registering as a new node"
docker run --rm \
	--env TOKEN="$TOKEN" \
	--env SERVER_ADDR="$SERVER_ADDR" \
	--env DOCKER_VERSION="$DOCKER_VERSION" \
	--volume /etc/tutum:/etc/tutum \
	--name tutum-agent \
	tutum/agent

echo "=> Configuring docker startup options"
set_docker_env ${UPSTART_DOCKER_CFG}
set_docker_env ${SYSTEMD_DOCKER_CFG}
set_docker_srv ${SYSTEMD_DOCKER_SRV1}
set_docker_srv ${SYSTEMD_DOCKER_SRV2}

echo "=> Restarting docker"
if systemctl >/dev/null 2>&1 ; then
	systemctl daemon-reload
fi
service docker restart

echo "=> Finishing registration"
RETRIES=60
for (( i=0 ; ; i++ )); do
	if [ ${i} -eq ${RETRIES} ]; then
		echo "Time out. Please see docker logs error messages"
		exit 1
	fi
	sleep 5
	docker ps > /dev/null 2>&1 && break
done

sleep 5
docker run -d \
	--env TOKEN="$TOKEN" \
	--env SERVER_ADDR="$SERVER_ADDR" \
	--env DOCKER_VERSION="$DOCKER_VERSION" \
	--volume /etc/tutum:/etc/tutum \
	--name tutum-agent \
	--restart always \
	--net host \
	tutum/agent

timeout 10s docker logs -f tutum-agent
