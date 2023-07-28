package provider

import (
	"fmt"
	"log"
)

type Request struct {
	url     string
	headers map[string]string
	body    interface{}
}

type TestClient struct {
	logs      []Request
	responses map[string]interface{}
}

func (client *TestClient) Post(url string, headers map[string]string, body interface{}) (interface{}, error) {
	return client.respond("POST", url, headers, body)
}

func (client *TestClient) Get(url string, headers map[string]string) (interface{}, error) {
	return client.respond("GET", url, headers, nil)
}

func (client *TestClient) Put(url string, headers map[string]string, body interface{}) (interface{}, error) {
	return client.respond("PUT", url, headers, body)
}

func (client *TestClient) Delete(url string, headers map[string]string) error {
	_, err := client.respond("DELETE", url, headers, nil)
	return err
}

func (client *TestClient) respond(method string, url string, headers map[string]string, body interface{}) (interface{}, error) {
	endpoint := fmt.Sprintf("%s %s", method, url)
	log.Printf("Locating response for %s", endpoint)
	response, ok := client.responses[endpoint]
	if !ok {
		log.Printf("Failed to find response for %s", endpoint)
		return nil, fmt.Errorf("Unable to find response for %s", endpoint)
	}
	log.Printf("Successfully located response for %s", endpoint)
	return response, nil
}
