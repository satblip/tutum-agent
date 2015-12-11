package agent

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func SendReq(method, url, token string, data []byte) (res []byte, err error) {
	headers := []string{"Authorization TutumAgentToken " + token,
		"Content-Type application/json",
		"User-Agent tutum-agent/" + VERSION}
	for i := 1; ; i *= 2 {
		if i > MAX_WAIT_TIME {
			i = 1
		}
		res, err := send(method, url, data, headers)
		if err == nil {
			return res, err
		}
		if err.Error() == "401" || (err.Error() == "404" && method == "PATCH") {
			return nil, err
		}
		log.Printf("Http error: %s. retry in %d seconds", err, i)
		time.Sleep(time.Duration(i) * time.Second)
	}
	return
}

func send(method, url string, data []byte, headers []string) ([]byte, error) {
	//log.Print("url: ",url)
	//log.Print("data: ", string(data))
	//log.Print("header: ", string(headers))
	var dataReader io.Reader
	if data == nil {
		dataReader = nil
	} else {
		dataReader = bytes.NewReader(data)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, dataReader)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for _, header := range headers {
			terms := strings.SplitN(header, " ", 2)
			if len(terms) == 2 {
				req.Header.Add(terms[0], terms[1])
			}
		}
	}

	resp, err := client.Do(req)
	//log.Print("resp: ", resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201, 202:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	default:
		err_msg := fmt.Sprintf("%d", resp.StatusCode)
		return nil, errors.New(err_msg)
	}
}
