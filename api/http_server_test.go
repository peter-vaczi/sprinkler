package api_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/peter-vaczi/sprinklerd/api"
	"github.com/peter-vaczi/sprinklerd/core"
)

var httpAPI api.API

func init() {
	mainEvents := make(chan interface{})
	gpioStub := NewGpioStub()
	httpAPI = api.New("http://localhost:9999", mainEvents)
	core.InitGpio(gpioStub)
	go core.Run(context.TODO(), mainEvents)
}

func req(t *testing.T, method, path, body string, expStatus int, expBody string) {
	var req *http.Request
	if len(body) == 0 {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
	}
	w := httptest.NewRecorder()
	httpAPI.ServeHTTP(w, req)

	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, expStatus)
	// assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, string(respBody), expBody)
}

func TestSomething(t *testing.T) {
	req(t, "GET", "/v1", "", 200, "{}")
	req(t, "GET", "/v1/devices", "", 200, "{}")
	req(t, "GET", "/v1/programs", "", 200, "{}")

	req(t, "GET", "/v1/devices/dev1", "", 404, "Not found\n")
	req(t, "GET", "/v1/programs/pr1", "", 404, "Not found\n")
}
