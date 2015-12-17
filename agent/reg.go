package agent

import (
	"encoding/json"
	"log"
)

type RegResponse struct {
	CaCert    string `json:"user_ca_cert"`
	UUID      string `json:"uuid"`
	FQDN      string `json:"external_fqdn"`
	DockerURL string `json:"docker_url"`
	NgrokURL  string `json:"ngrok_url"`
	PublicIP  string `json:"public_ip"`
}

type RegPost struct {
	Version string `json:"agent_version"`
}

type RegPatch struct {
	Cert          string `json:"public_cert"`
	Version       string `json:"agent_version"`
	DockerVersion string `json:"docker_version,omitempty"`
}

func RegNewNode(regUrl, token string) (string, string, string) {
	log.Print("Register node to ", regUrl)
	res := postReg(regUrl, token)
	log.Print("Fqdn: ", res.FQDN)
	log.Print("Uuid: ", res.UUID)
	return res.CaCert, res.FQDN, res.UUID
}

func UpdateNode(regUrl, token, uuid, cert, dockerVersion string) {
	url := UrlJoin(regUrl, uuid)
	log.Print("Update node information to ", url)
	patchReg(url, token, cert, dockerVersion)
}

func postReg(url, token string) (res RegResponse) {
	req, err := json.Marshal(RegPost{Version: VERSION})
	if err != nil {
		log.Fatal("Cannot marshal the post data: ", err)
	}
	data, err := SendReq("POST", url, token, req)
	if err != nil {
		log.Printf("Registration failed: %s", err)
		log.Fatalf("Token(%s) is invalid", token)
	}
	if err := json.Unmarshal(data, &res); err != nil {
		log.Fatal("Failed to unmarshal the response: ", err)
	}
	return
}

func patchReg(url, token, cert, dockerVersion string) (res RegResponse) {
	req, err := json.Marshal(RegPatch{Version: VERSION, Cert: cert, DockerVersion: dockerVersion})
	if err != nil {
		log.Fatal("Cannot marshal the patch data: ", err)
	}
	data, err := SendReq("PATCH", url, token, req)
	if err != nil {
		log.Printf("Update node information failed: %s", err)
		log.Fatalf("Either token or UUID is invalid, please remove this container and run the script with a new token")
	}
	if err := json.Unmarshal(data, &res); err != nil {
		log.Fatal("Failed to unmarshal the response: ", err)
	}
	return
}
