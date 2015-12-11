package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ActiveState/tail"
)

type TunnelReqest struct {
	Tunnel  string `json:"tunnel"`
	Version string `json:"agent_version"`
}

type ReachableResponse struct {
	Reachable bool `json:"reachable"`
}

type NodeInfo struct {
	AgentVersion string `json:"agent_version"`
	DockerUrl    string `json:"docker_url"`
	ExternalFqdn string `json:"external_fqdn"`
	NgrokUrl     string `json:"ngrok_url"`
	PublicCert   string `json:"public_cert"`
	ResourceUri  string `json:"resource_uri"`
	State        string `json:"state"`
	Tunnel       string `json:"tunnel"`
	UserCaCert   string `json:"user_ca_cert"`
	UUID         string `json:"uuid"`
	NgrokServer  string `json:"ngrok_server_addr"`
}

func IsNodeReachable(regUrl, token, uuid string) bool {
	var res ReachableResponse
	url := UrlJoin(regUrl, uuid+"/ping/")
	data, err := SendReq("POST", url, token, nil)
	if err != nil {
		log.Printf("Get node reachable failed: %s", err)
		return false
	}
	if err := json.Unmarshal(data, &res); err != nil {
		log.Fatal("Failed to unmarshal the response: ", err)
	}
	log.Print("Node reachable: ", res.Reachable)
	return res.Reachable
}

func GetNodeInfo(regUrl, token, uuid string) (res NodeInfo) {
	url := UrlJoin(regUrl, uuid)
	data, err := SendReq("GET", url, token, nil)
	if err != nil {
		log.Fatal("Get node infomation failed: s", err)
	}
	if err := json.Unmarshal(data, &res); err != nil {
		log.Fatal("Failed to unmarshal the response: ", err)
	}
	return
}

func CreateNgrokcfg(fileNgrokCfg, NgrokServer string) {
	cfg := fmt.Sprintf("server_addr: %s\ntrust_host_root_certs: true\ninspect_addr: \"disabled\"", NgrokServer)
	SaveFile(fileNgrokCfg, cfg)
}

func CreateTunnel(regUrl, token, uuid, fileNgrokBin, fileNgorkCfg, fileNgrokLog string) {
	os.RemoveAll(fileNgrokLog)
	url := UrlJoin(regUrl, uuid)
	cmd := exec.Command(fileNgrokBin,
		"-config", fileNgorkCfg,
		"-log", "stdout",
		"-proto", "tcp",
		DOCKER_TCP_PORT)

	ngrokLog, err := os.OpenFile(fileNgrokLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	} else {
		defer ngrokLog.Close()
		cmd.Stdout = ngrokLog
	}

	go monitorTunnels(url, token, fileNgrokLog)
	log.Println("Starting NAT tunnel")

	runNgrok(cmd)

	for {
		log.Println("Restarting NAT tunnel in 10 seconds")
		time.Sleep(10 * time.Second)
		runNgrok(cmd)
	}
}

func runNgrok(cmd *exec.Cmd) {
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}
	cmd.Wait()
}

func monitorTunnels(url, token, fileNgrokLog string) {
	update, _ := tail.TailFile(fileNgrokLog, tail.Config{
		Follow: true,
		ReOpen: true})
	for line := range update.Lines {
		if strings.Contains(line.Text, "[INFO] [client] Tunnel established at") {
			terms := strings.Split(line.Text, " ")
			tunnel := terms[len(terms)-1]
			log.Print("Found new tunnel: ", tunnel)
			if tunnel != "" {
				sendTunnel(url, token, tunnel)
			}
		}
	}
}

func sendTunnel(url, token, tunnel string) {
	log.Println("Sending tunnel address to Tutum")
	req := TunnelReqest{}
	req.Version = VERSION
	req.Tunnel = tunnel
	data, err := json.Marshal(req)
	if err != nil {
		log.Fatal("Cannot marshal the post data: ", err)
	}
	_, err = SendReq("PATCH", url, token, data)
	if err != nil {
		log.Print("Failed to send tunnel address to Tutum: ", err)
	} else {
		log.Print("New tunnel has been set up")
		log.Print("Done!")
	}
}
