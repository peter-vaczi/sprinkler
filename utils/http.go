package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

func GetRequest(url string, responseData interface{}) error {
	return httpRequest("GET", url, nil, responseData)
}

func DeleteRequest(url string) error {
	return httpRequest("DELETE", url, nil, nil)
}

func PutRequest(url string, data interface{}) error {
	return httpRequest("PUT", url, data, nil)
}

func PostRequest(url string, data interface{}) error {
	return httpRequest("POST", url, data, nil)
}

func EncodeJson(data interface{}) *bytes.Buffer {
	if nil != data {
		outgoingJSON, err := json.Marshal(data)
		if err == nil {
			return bytes.NewBuffer(outgoingJSON)
		}
	}

	return bytes.NewBuffer([]byte{})
}

func httpRequest(method string, url string, data interface{}, responseData interface{}) error {

	body := EncodeJson(data)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if nil != data {
		req.Header.Set("Content-Type", "application/json")
	}

	var client = &http.Client{
		Timeout: time.Second * 5,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < 300 {
		if responseData != nil {
			return json.NewDecoder(res.Body).Decode(responseData)
		}
		return nil
	} else {
		msg, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(msg))
	}
}
