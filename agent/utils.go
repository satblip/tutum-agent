package agent

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func IsFileExist(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func UrlJoin(url1 string, url2 string) (url string) {
	if strings.HasSuffix(url1, "/") {
		if strings.HasPrefix(url2, "/") {
			url = url1 + url2[1:]
		} else {
			url = url1 + url2
		}
	} else {
		if strings.HasPrefix(url2, "/") {
			url = url1 + url2
		} else {
			url = url1 + "/" + url2
		}
	}
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return
}

func LoadFile(file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(b)
}

func SaveFile(file, data string) {
	if err := ioutil.WriteFile(file, []byte(data), 0644); err != nil {
		log.Fatal(err)
	}
}
