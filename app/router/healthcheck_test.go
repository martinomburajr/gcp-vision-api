package router

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	x := GCPVisionAPIServer{}

	server := httptest.NewServer(x.Routes())
	log.Print(server.URL)
	defer server.Close()

	request := httptest.NewRequest(http.MethodGet, server.URL+"/", nil)
	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error("not 200 OK")
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("error reading body")
	}

	msg := "GCP-VISION-API is Healthy :) - Authored by Martin Ombura Jr.\nGitHub: @martinomburajr"
	if string(data) != msg {
		t.Errorf("Response Body: %s | wanted: %s", data, msg)
	}
}
