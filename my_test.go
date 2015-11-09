package MyTest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type loginResponse struct {
	Id           string
	Access_token string
	Expires_in   float64
	Token_type   string
}

func TestRestLogin(t *testing.T) {
	appid := "9ab34d8b"
	appkey := "7a950d78956ed39f3b0815f0f001b43b"
	apphost := "api-jp.kii.com"
	url := fmt.Sprintf("https://%s/api/oauth2/token", apphost)
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
	var login loginResponse
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
