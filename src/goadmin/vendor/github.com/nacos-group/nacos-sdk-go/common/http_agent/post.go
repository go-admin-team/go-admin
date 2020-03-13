package http_agent

import (
	"net/http"
	"strings"
	"time"
)

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-07 15:13
**/

func post(path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error) {
	client := http.Client{}
	client.Timeout = time.Millisecond * time.Duration(timeoutMs)
	var body string
	for key, value := range params {
		if len(value) > 0 {
			body += key + "=" + value + "&"
		}
	}
	if strings.HasSuffix(body, "&") {
		body = body[:len(body)-1]
	}
	request, errNew := http.NewRequest(http.MethodPost, path, strings.NewReader(body))
	if errNew != nil {
		err = errNew
		return
	}
	request.Header = header
	resp, errDo := client.Do(request)
	if errDo != nil {
		err = errDo
	} else {
		response = resp
	}
	return
}
