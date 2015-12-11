package agent

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func SetDockerOpts(opts string) bool {
	modified := false
	if setUpstart(opts, DOCKER_CFG_UPSTART) {
		modified = true
		log.Print("Docker config file is updated: ", DOCKER_CFG_UPSTART)
	}
	return modified
}

func RestartDocker(dockerPid string) {
	log.Print("Restarting Docker daemon")
	cmd := exec.Command("kill", dockerPid)

	if err := cmd.Start(); err != nil {
		log.Print("Cannot start docker daemon: ", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Print(err)
	}
}

func GetDockerOpts(fileCacert, fileCert, fileKey string) string {
	opts := fmt.Sprintf("\"-H %s -H %s --tlscert %s --tlskey %s --tlscacert %s --tlsverify \"", DOCKER_UNIX_HOST, DOCKER_TCP_HOST, fileCert, fileKey, fileCacert)
	if strings.TrimSpace(os.Getenv("ExtraOpts")) != "" {
		opts += " " + os.Getenv("ExtraOpts")
	}
	return opts
}

func setUpstart(opts, cfg string) bool {
	modified := false
	if IsFileExist(DOCKER_CFG_UPSTART) {
		DOCKEROPTS := "DOCKER_OPTS"
		input, err := ioutil.ReadFile(cfg)
		if err != nil {
			log.Fatal(err)
		}

		optsStr := DOCKEROPTS + "=" + opts
		lines := strings.Split(string(input), "\n")
		for i, line := range lines {
			if !strings.HasPrefix(strings.TrimSpace(line), "#") && strings.Contains(line, DOCKEROPTS) {
				lines[i] = optsStr
			}
		}
		output := strings.Join(lines, "\n")
		if !strings.Contains(output, optsStr) {
			output = output + "\n" + optsStr
		}
		if string(input) != output {
			modified = false
		}
		err = ioutil.WriteFile(cfg, []byte(output), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	return modified
}
