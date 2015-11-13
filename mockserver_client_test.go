package mock_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type LoginResponse struct {
	Id           string `json:"id"`
	Access_token string `json:"access_token"`
	Expires_in   float64 `json:"expires_in"`
	Token_type   string `json:"token_type"`
}

type QueryParam struct {
	Name string `json:"name"`
	Values []string `json:"values"`
}

type Header struct {
	Name string `json:"name"`
	Values []string `json:"values"`
}

type RequestBody struct {
	Type string `json:"type"`
	Value string `json:"value"`
}

type HttpRequest struct {
	Method string `json:"method"`
	Path string `json:"path"`
	QueryStringParameters []QueryParam `json:"queryStringParameters"`
	Headers []Header `json:"headers"`
	Body RequestBody `json:"body"`
}

type HttpResponse struct {
	StatusCode int `json:"statusCode"`
	Headers []Header `json:"headers"`
	Body string `json:"body"`
}

type MockServerRequest struct {
	HttpRequest HttpRequest `json:"httpRequest"`
	HttpResponse HttpResponse `json:"httpResponse"`
}

// need to setup mockserver before running.
// http://www.mock-server.com/
func setUpMockServer(t *testing.T) {
	mockResponseBytes, _ := json.Marshal(LoginResponse{
		Access_token: "dummyToken",
		Id: "dummyID",
		Expires_in: 24*3600,
		Token_type: "Bearer",
	})
	request := MockServerRequest {
		HttpRequest: HttpRequest {
			Method: "POST",
			Path: "/api/oauth2/token",
			QueryStringParameters: []QueryParam {},
			Headers: []Header {
				Header {
					Name: "content-type",
					Values: []string{"application/json"},
				},
				Header {
					Name: "x-kii-appid",
					Values: []string{"9ab34d8b"},
				},
				Header {
					Name: "x-kii-appkey",
					Values: []string{"7a950d78956ed39f3b0815f0f001b43b"},
				},
			},
			Body: RequestBody {
				Type: "JSON",
				Value: "{\"username\":\"pass1234\", \"password\":\"1234\"}",
			},
		},
		HttpResponse: HttpResponse {
			StatusCode: 200,
			Body: string(mockResponseBytes),
			Headers: []Header {
				Header {
					Name: "content-type",
					Values: []string{"application/json"},
				},
			},
		},
	}
	b, err := json.Marshal(request)
	if err != nil {
		t.Error("failed to marshal json")
		return
	}
	url := "http://localhost:12345/expectation"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	_, err2 := client.Do(req)

	if err2 != nil {
		t.Errorf("got error on api request %s", err)
	}
}

func TestRestLogin(t *testing.T) {
	setUpMockServer(t)
	appid := "9ab34d8b"
	appkey := "7a950d78956ed39f3b0815f0f001b43b"
	apphost := "localhost:12345"
	url := fmt.Sprintf("http://%s/api/oauth2/token", apphost)
	var jsonStr = []byte(`{"username":"pass1234", "password":"1234"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("x-kii-appid", appid)
	req.Header.Set("x-kii-appkey", appkey)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		t.Errorf("got error on api request %s", err)
	}

	defer resp.Body.Close()

	// logging
	t.Logf("response Status:%d\n", resp.Status)
	t.Logf("response Headers:%s\n", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	t.Logf("response Body:%s\n", string(body))

	if resp.Status != "200 OK" {
		t.Errorf("Got unexpected resonse %d", resp.Status)
	}

	// Parse JSON response
	var login LoginResponse
	err2 := json.Unmarshal(body, &login)
	if err2 != nil {
		t.Errorf("got error on response parse: %s", err2)
	}

	if len(login.Id) <= 0 {
		t.Errorf("invalid login Id: %s", login.Id)
	}
	if len(login.Access_token) <= 0 {
		t.Errorf("invalid access token: %s", login.Access_token)
	}
	if login.Expires_in <= 0 {
		t.Errorf("invalid token expiration : %f", login.Expires_in)
	}
	if login.Token_type != "Bearer" {
		t.Errorf("unexpected token type: %f", login.Token_type)
	}
}

// vim: set noet ts=4 sts=4 sw=4 fenc=utf-8 ff=unix :
