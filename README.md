tutum-agent
===========

	./install.sh <token>

or

	docker run -d \
		--net host \
		--pid host \
		--privileged \
		--restart always \
		--env TOKEN="<token>" \
		--env SERVER_HOST="https://dashboard.tutum.co" \
		--env DOCKER_VERSION="<docker version>" \
		--env DOCKER_PID="<dockerid>" \
		--volume /etc:/etc \
		--volume /tmp/docker.pid:/tmp/docker.pid
		--name tutum-agent \
		tutum/agent
