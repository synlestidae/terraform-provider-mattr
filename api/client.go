package api

import (
	"log"
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"bytes"
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
	return send[interface{}]("POST", url, headers, &body)
}

func (client *HttpClient) Get(url string, headers map[string]string) (interface{}, error) {
	return send[interface{}]("GET", url, headers, nil)
}

func (client *HttpClient) Put(url string, headers map[string]string, body interface{}) (interface{}, error) {
	return send[interface{}]("PUT", url, headers, &body)
}

func (client *HttpClient) Delete(url string, headers map[string]string) error {
	_, err := send[interface{}]("DELETE", url, headers, nil)
	return err
}

func send[T any](method string, url string, headers map[string]string, body *interface{}) (*T, error) {
	client := http.DefaultClient 

	var bodyJson []byte
	var err error

	if body != nil {
		bodyJson, err = json.Marshal(*body)
		if err != nil {
			return nil, err
		}
	}

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
		log.Println("Response body: ", string(response_body))
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
