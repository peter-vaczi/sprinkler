package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/peter-vaczi/sprinklerd/api"
	"github.com/peter-vaczi/sprinklerd/core"
	"github.com/stretchr/testify/assert"
)

var httpAPI api.API

func init() {
	mainEvents := make(chan interface{})
	gpioStub := NewGpioStub()
	httpAPI = api.New("http://localhost:9999", mainEvents)
	core.InitGpio(gpioStub)
	go core.Run(context.TODO(), mainEvents)
}

func jsonContains(t *testing.T, a, b string) {
	ma := make(map[string]interface{})
	mb := make(map[string]interface{})
	err := json.NewDecoder(strings.NewReader(a)).Decode(&ma)
	if err != nil {
		assert.Contains(t, a, b)
		return
	}
	err = json.NewDecoder(strings.NewReader(b)).Decode(&mb)
	if err != nil {
		assert.Contains(t, a, b)
		return
	}
	for k, v := range mb {
		value, found := ma[k]
		if !assert.True(t, found, fmt.Sprintf("field %s not found in json body", k)) {
			return
		}
		switch v := v.(type) {
		case map[string]interface{}:
			ja, _ := json.Marshal(value)
			jb, _ := json.Marshal(v)
			jsonContains(t, string(ja), string(jb))
		default:
			assert.Equal(t, value, v)
		}
	}
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

	assert.Equal(t, expStatus, resp.StatusCode)
	// assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
	//assert.Equal(t, string(respBody), expBody)
	jsonContains(t, string(respBody), expBody)
}

func TestApiEmptyDb(t *testing.T) {
	req(t, "GET", "/v1", "", 200, "{}")
	req(t, "GET", "/v1/devices", "", 200, "{}")
	req(t, "GET", "/v1/programs", "", 200, "{}")

	req(t, "GET", "/v1/devices/unknown-device", "", 404, "Not found")
	req(t, "GET", "/v1/programs/unknown-program", "", 404, "Not found")
}

func TestApiAddDelDevice(t *testing.T) {
	req(t, "POST", "/v1/devices", "{\"name\":\"dev1\", \"pin\":1}", 200, "")
	req(t, "POST", "/v1/devices", "{\"name\":\"dev1\", \"pin\":1}", 406, "Already exists")
	req(t, "GET", "/v1/devices/dev1", "", 200, "{\"name\":\"dev1\", \"pin\":1}")

	req(t, "PUT", "/v1/devices/dev1", "{\"pin\":42}", 200, "")
	req(t, "GET", "/v1/devices/dev1", "", 200, "{\"name\":\"dev1\", \"pin\":42}")
	req(t, "GET", "/v1/devices", "", 200, "{\"dev1\":{\"name\":\"dev1\", \"pin\":42}}")
	req(t, "DELETE", "/v1/devices/dev1", "", 200, "")
	req(t, "GET", "/v1/devices/dev1", "", 404, "Not found")
}

func TestApiAddDelProgram(t *testing.T) {
	req(t, "POST", "/v1/programs", "{\"name\":\"pr1\"}", 200, "")
	req(t, "GET", "/v1/programs/pr1", "", 200, "{\"name\":\"pr1\"}")

	req(t, "GET", "/v1/programs", "", 200, "{\"pr1\":{\"name\":\"pr1\"}}")

	req(t, "DELETE", "/v1/programs/pr1", "", 200, "")
	req(t, "GET", "/v1/programs/pr1", "", 404, "Not found")
}

func TestApiAddDelDeviceToProgram(t *testing.T) {
	req(t, "POST", "/v1/devices", "{\"name\":\"dev1\", \"pin\":1}", 200, "")
	req(t, "POST", "/v1/devices", "{\"name\":\"dev2\", \"pin\":2}", 200, "")
	req(t, "POST", "/v1/programs", "{\"name\":\"pr1\"}", 200, "")

	req(t, "POST", "/v1/programs/pr1/devices", "{\"device\":\"dev1\", \"duration\":\"5s\"}", 200, "")
	req(t, "POST", "/v1/programs/pr1/devices", "{\"device\":\"dev2\", \"duration\":\"8s\"}", 200, "")

	req(t, "POST", "/v1/programs/pr1/devices", "{\"device\":\"dev-whatever\", \"duration\":\"5s\"}", 404, "Not found")
	req(t, "POST", "/v1/programs/pr-whatever/devices", "{\"device\":\"dev1\", \"duration\":\"5s\"}", 404, "Not found")

	req(t, "GET", "/v1/programs", "", 200, "{\"pr1\":{\"name\":\"pr1\", \"elements\":[{\"device\":\"dev1\",\"duration\":5000000000},{\"device\":\"dev2\",\"duration\":8000000000}]}}")

	req(t, "DELETE", "/v1/programs/2", "", 404, "Not found")

	req(t, "DELETE", "/v1/programs/pr1/devices/0", "", 200, "")
	req(t, "DELETE", "/v1/programs/pr1/devices/0", "", 200, "")
	req(t, "DELETE", "/v1/programs/pr1/devices/0", "", 404, "Element index out of range")

	// cleanup
	req(t, "DELETE", "/v1/programs/pr1", "", 200, "")
	req(t, "DELETE", "/v1/devices/dev1", "", 200, "")
	req(t, "DELETE", "/v1/devices/dev2", "", 200, "")

	req(t, "DELETE", "/v1/programs/pr1/devices/0", "", 404, "Not found")
}

func TestApiDelDeviceInUse(t *testing.T) {
	req(t, "POST", "/v1/devices", "{\"name\":\"dev1\", \"pin\":1}", 200, "")
	req(t, "POST", "/v1/programs", "{\"name\":\"pr1\"}", 200, "")

	req(t, "POST", "/v1/programs/pr1/devices", "{\"device\":\"dev1\", \"duration\":\"5s\"}", 200, "")

	req(t, "DELETE", "/v1/devices/dev1", "", 406, "Device is in use")

	req(t, "DELETE", "/v1/programs/pr1/devices/0", "", 200, "")
	req(t, "DELETE", "/v1/programs/pr1", "", 200, "")
	req(t, "DELETE", "/v1/devices/dev1", "", 200, "")
}

func TestApiStartStopProgram(t *testing.T) {
	req(t, "POST", "/v1/devices", "{\"name\":\"dev1\", \"pin\":1}", 200, "")
	req(t, "POST", "/v1/devices", "{\"name\":\"dev2\", \"pin\":2}", 200, "")
	req(t, "POST", "/v1/programs", "{\"name\":\"pr1\"}", 200, "")

	req(t, "POST", "/v1/programs/pr1/devices", "{\"device\":\"dev1\", \"duration\":\"5s\"}", 200, "")
	req(t, "POST", "/v1/programs/pr1/devices", "{\"device\":\"dev2\", \"duration\":\"8s\"}", 200, "")

	// start pr
	req(t, "POST", "/v1/programs/pr-unknown/start", "", 404, "Not found")
	req(t, "POST", "/v1/programs/pr1/start", "", 200, "")
	time.Sleep(100 * time.Millisecond)
	req(t, "GET", "/v1/devices/dev1", "", 200, "{\"name\":\"dev1\", \"on\":true}")
	req(t, "GET", "/v1/devices/dev2", "", 200, "{\"name\":\"dev2\", \"on\":false}")

	// stop pr
	req(t, "POST", "/v1/programs/pr-unknown/stop", "", 404, "Not found")
	req(t, "POST", "/v1/programs/pr1/stop", "", 200, "")
	time.Sleep(100 * time.Millisecond)
	req(t, "GET", "/v1/devices/dev1", "", 200, "{\"name\":\"dev1\", \"on\":false}")
	req(t, "GET", "/v1/devices/dev2", "", 200, "{\"name\":\"dev2\", \"on\":false}")

	// cleanup
	req(t, "DELETE", "/v1/programs/pr1", "", 200, "")
	req(t, "DELETE", "/v1/devices/dev1", "", 200, "")
	req(t, "DELETE", "/v1/devices/dev2", "", 200, "")
}

// func TestApiBadRequests(t *testing.T) {
// 	req(t, "PUT", "/v1/devices", "invalid", 400, "Invalid json")
// 	req(t, "POST", "/v1/devices", "invalid", 400, "Invalid json")
// 	req(t, "POST", "/v1/programs", "invalid", 400, "Invalid json")
// 	req(t, "POST", "/v1/programs/pr1/devices", "invalid", 400, "Invalid json")
// }