package agent

const (
	WORKDIR            = "/"
	CERTDIR            = "/etc/tutum"
	REG_URI            = "api/agent/node/"
	ERROR_CODE         = 1
	FILENAME_UUID      = "uuid"
	FILENAME_CACERT    = "ca.pem"
	FILENAME_CERT      = "cert.pem"
	FILENAME_KEY       = "key.pem"
	FILENAME_DOCKERCFG = "docker.cfg"
	FILENAME_NGROKCFG  = "ngrok.conf"
	FILENAME_NGROKLOG  = "ngrok.log"
	FILENAME_NGROKBIN  = "ngrok"
	DOCKER_UNIX_HOST   = "unix:///var/run/docker.sock"
	DOCKER_TCP_HOST    = "tcp://0.0.0.0:2375"
	DOCKER_TCP_PORT    = "2375"
	VERSION            = "0.1"
	MAX_WAIT_TIME      = 300 //seconds
	DOCKER_CFG_UPSTART = "/etc/default/docker"
)

var (
	SERVER_HOST = "https://dashboard.tutum.co"
)
