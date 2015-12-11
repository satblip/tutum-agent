package main //import "github.com/tutumcloud/tutum-agent"

import (
	"log"
	"os"
	"path"
	"strings"
	"time"

	. "github.com/tutumcloud/tutum-agent/agent"
)

func main() {
	if os.Getenv("SERVER_HOST") != "" {
		SERVER_HOST = os.Getenv("SERVER_HOST")
	}
	fileNgrokCfg := path.Join(WORKDIR, FILENAME_NGROKCFG)
	fileNgrokLog := path.Join(WORKDIR, FILENAME_NGROKLOG)
	fileNgrokBin := path.Join(WORKDIR, FILENAME_NGROKBIN)
	fileUUID := path.Join(WORKDIR, FILENAME_UUID)
	fileCacert := path.Join(CERTDIR, FILENAME_CACERT)
	fileCert := path.Join(CERTDIR, FILENAME_CERT)
	fileKey := path.Join(CERTDIR, FILENAME_KEY)
	regUrl := UrlJoin(SERVER_HOST, REG_URI)

	log.Print("====================")
	_ = os.MkdirAll(WORKDIR, 0755)
	_ = os.MkdirAll(CERTDIR, 0755)
	token := strings.TrimSpace(os.Getenv("TOKEN"))
	dockerPid := strings.TrimSpace(os.Getenv("DOCKER_PID"))
	dockerVersion := strings.TrimSpace(os.Getenv("DOCKER_VERSION"))
	if token == "" {
		log.Fatal("Error: empty token")
	}
	if dockerPid == "" {
		log.Fatal("Error: empty docker pid")
	}
	if dockerVersion == "" {
		log.Fatal("Error: empty docker version")
	}

	opts := GetDockerOpts(fileCacert, fileCert, fileKey)

	if !IsFileExist(fileUUID) {
		fqdn := RegNewNode(regUrl, token, fileUUID, fileCacert)
		CreateCerts(fileKey, fileCert, fqdn)
		SetDockerOpts(opts)
		log.Print("Docker config file is updated")
		RestartDocker(dockerPid)
		time.Sleep(10 * time.Second)
	} else {
		if SetDockerOpts(opts) {
			RestartDocker(dockerPid)
			time.Sleep(10 * time.Second)
		}
	}

	uuid := LoadFile(fileUUID)
	cert := LoadFile(fileCert)

	UpdateNode(regUrl, token, uuid, cert, dockerVersion)
	isReachable := IsNodeReachable(regUrl, token, uuid)

	if !isReachable {
		nodeInfo := GetNodeInfo(regUrl, token, uuid)
		CreateNgrokcfg(fileNgrokCfg, nodeInfo.NgrokServer)
		CreateTunnel(regUrl, token, uuid, fileNgrokBin, fileNgrokCfg, fileNgrokLog)
	}
	log.Print("Done!")
	for {
		time.Sleep(5 * time.Second)
	}
}
