package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

type Client interface {
	Post(url string, headers map[string]string, body interface{}) (interface{}, error)
	Get(url string, headers map[string]string) (interface{}, error)
	Put(url string, headers map[string]string, body interface{}) (interface{}, error)
	Delete(url string, headers map[string]string) error
}

type HttpClient struct {
}

func (client *HttpClient) Post(url string, headers map[string]string, body interface{}) (interface{}, error) {
	responseBod, err := send[interface{}]("POST", url, headers, &body)
	if err != nil {
		return nil, err
	}
	if responseBod == nil {
		return nil, fmt.Errorf("Unable to load data for POST %s", url)
	}
	return *responseBod, err
}

func (client *HttpClient) Get(url string, headers map[string]string) (interface{}, error) {
	responseBod, err := send[interface{}]("GET", url, headers, nil)
	if err != nil {
		return nil, err
	}
	if responseBod == nil {
		return nil, fmt.Errorf("Unable to load data for GET %s", url)
	}
	return *responseBod, err
}

func (client *HttpClient) Put(url string, headers map[string]string, body interface{}) (interface{}, error) {
	responseBod, err := send[interface{}]("PUT", url, headers, &body)
	if err != nil {
		return nil, err
	}
	if responseBod == nil {
		return nil, fmt.Errorf("Unable to load data for PUT %s", url)
	}
	return *responseBod, err
}

func (client *HttpClient) Delete(url string, headers map[string]string) error {
	_, err := send[interface{}]("DELETE", url, headers, nil)
	return err
}

func send[T any](method string, url string, headers map[string]string, body *interface{}) (*T, error) {
	client := http.DefaultClient

	var bodyJson []byte
	var err error

	// raw bytes do not get converted to json

	log.Printf("Body type is %T", body)

	if body != nil {
		if reflect.TypeOf(*body) == reflect.TypeOf([]byte{}) {
			log.Printf("Uploading binary payload")
			// Handle the case when body is of type *[]byte
			byteSlice := (*body).([]byte)
			bodyJson = byteSlice
			headers["Content-Type"] = "application/zip"
		} else {
			log.Printf("Creating JSON body")
			bodyJson, err = json.Marshal(body)
			if err != nil {
				return nil, err
			}
			headers["Content-Type"] = "application/json"
		}
	}

	log.Printf("Uploading %d byte(s)", len(bodyJson))
	bodyBuf := bytes.NewBuffer(bodyJson)
	request, err := http.NewRequest(method, url, bodyBuf)
	if err != nil {
		return nil, err
	}
	for name, value := range headers {
		request.Header.Set(name, value)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if 400 <= resp.StatusCode && resp.StatusCode < 599 {
		response_body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error response body: %s", response_body)
		return nil, fmt.Errorf("%s %s gave error status code %d", method, url, resp.StatusCode)
	}
	if resp.StatusCode == 204 {
		return nil, nil
	}

	result, err := processResponse[T](resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func processResponse[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		response_body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Response body: %s", string(response_body))
		return nil, fmt.Errorf("Got status code %d from API", resp.StatusCode)
	}

	// read raw json body
	response_body, err := ioutil.ReadAll(resp.Body)
	log.Println("Response body: ", string(response_body))
	if err != nil {
		return nil, err
	}

	// parse the body
	var response T
	err = json.Unmarshal(response_body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
