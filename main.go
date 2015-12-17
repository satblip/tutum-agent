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
	if os.Getenv("SERVER_ADDR") != "" {
		SERVER_ADDR = os.Getenv("SERVER_ADDR")
	}

	fileNgrokCfg := path.Join("/", FILENAME_NGROKCFG)
	fileNgrokLog := path.Join("/", FILENAME_NGROKLOG)
	fileNgrokBin := path.Join("/", FILENAME_NGROKBIN)
	fileUUID := path.Join(WORKDIR, FILENAME_UUID)
	fileCacert := path.Join(WORKDIR, FILENAME_CACERT)
	fileCert := path.Join(WORKDIR, FILENAME_CERT)
	fileKey := path.Join(WORKDIR, FILENAME_KEY)
	regUrl := UrlJoin(SERVER_ADDR, REG_URI)

	_ = os.MkdirAll(WORKDIR, 0755)

	token := strings.TrimSpace(os.Getenv("TOKEN"))
	if token == "" {
		log.Fatal("Error: empty token")
	}
	dockerVersion := strings.TrimSpace(os.Getenv("DOCKER_VERSION"))
	if dockerVersion == "" {
		log.Fatal("Error: empty docker version")
	}

	if !IsFileExist(fileUUID) {
		cacert, fqdn, uuid := RegNewNode(regUrl, token)
		SaveFile(fileCacert, cacert)
		SaveFile(fileUUID, uuid)
		CreateCerts(fileKey, fileCert, fqdn)
		log.Print("Done!")
		os.Exit(0)
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
