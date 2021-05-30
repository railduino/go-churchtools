package churchtools

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ROOT_CA_FILE = "/etc/ssl/certs/ca-certificates.crt"
)

type Connector struct {
	Hostname string
	Username string
	Password string

	Client *http.Client

	Build   string
	Version string

	PersonID int
	Cookie   *http.Cookie
	Token    string
}

type MetaPerson struct {
	ID   int    `json:"id"`
	Name string `json:"-"`
}

type InfoResult struct {
	Build   string `json:"build"`
	Version string `json:"version"`
}

type LoginData struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	PersonID int    `json:"personId"`
	Location string `json:"location"`
}

type LoginResult struct {
	Data LoginData `json:"data"`
}

type LoginTokenResult struct {
	Data string `json:"data"`
}

func New(hostname, username, password string) (*Connector, error) {
	fmt.Printf("CT NEW 5\n")

	conn := &Connector{
		Hostname: hostname,
		Username: username,
		Password: password,
	}

	// Setup the HTTPS client
	rootCAPool := x509.NewCertPool()
	rootCA, err := ioutil.ReadFile(ROOT_CA_FILE)
	if err != nil {
		return nil, err
	}
	rootCAPool.AppendCertsFromPEM(rootCA)

	conn.Client = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAPool,
			},
		},
	}

	// Check connectivity to the ChurchTools host
	result, err := conn.Get("info", false)
	if err != nil {
		return nil, err
	}
	var info InfoResult
	if err := json.Unmarshal(result, &info); err != nil {
		return nil, err
	}
	conn.Build = info.Build
	conn.Version = info.Version

	// Login in order to obtail the Token
	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		username,
		password,
	}
	result, err = conn.Post("login", data)
	if err != nil {
		return nil, err
	}
	var login LoginResult
	if err := json.Unmarshal(result, &login); err != nil {
		return nil, err
	}
	conn.PersonID = login.Data.PersonID

	// Get the LoginToken for further access
	endpoint := fmt.Sprintf("persons/%d/logintoken", conn.PersonID)
	result, err = conn.Get(endpoint, true)
	if err != nil {
		return nil, err
	}
	var token LoginTokenResult
	if err := json.Unmarshal(result, &token); err != nil {
		return nil, err
	}
	conn.Token = token.Data

	return conn, nil
}

func (conn *Connector) Get(endpoint string, needAuth bool) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/%s", conn.Hostname, endpoint)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")

	if needAuth {
		if conn.Token != "" {
			request.Header.Set("Authorization", fmt.Sprintf("Login %s", conn.Token))
		} else if conn.Cookie != nil {
			request.AddCookie(conn.Cookie)
		}
	}

	response, err := conn.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

func (conn *Connector) Post(endpoint string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/%s", conn.Hostname, endpoint)

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")

	if endpoint != "login" {
		if conn.Token != "" {
			request.Header.Set("Authorization", fmt.Sprintf("Login %s", conn.Token))
		} else if conn.Cookie != nil {
			request.AddCookie(conn.Cookie)
		}
	}

	response, err := conn.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if len(response.Cookies()) > 0 {
		conn.Cookie = response.Cookies()[0]
		fmt.Printf("CT COOKIE '%+v'\n", conn.Cookie)
	}

	return ioutil.ReadAll(response.Body)
}

func (conn *Connector) Delete(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/%s", conn.Hostname, endpoint)

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")

	if conn.Token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Login %s", conn.Token))
	} else if conn.Cookie != nil {
		request.AddCookie(conn.Cookie)
	}

	response, err := conn.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
