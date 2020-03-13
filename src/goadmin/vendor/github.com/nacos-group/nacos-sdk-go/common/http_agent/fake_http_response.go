package http_agent

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
)

/**
*
* @description : mock http response
*
* @author : codezhang
*
* @create : 2019-01-11 12:10
**/
type fakeHttpResponseBody struct {
	body io.ReadSeeker
}

func (body *fakeHttpResponseBody) Read(p []byte) (n int, err error) {
	n, err = body.body.Read(p)
	if err == io.EOF {
		body.body.Seek(0, 0)
	}
	return n, err
}

func (body *fakeHttpResponseBody) Close() error {
	return nil
}

func FakeHttpResponse(status int, body string) (resp *http.Response) {
	return &http.Response{
		Status:     strconv.Itoa(status),
		StatusCode: status,
		Body:       &fakeHttpResponseBody{bytes.NewReader([]byte(body))},
		Header:     http.Header{},
	}
}
