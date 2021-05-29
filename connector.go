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

	Version string

	UserID int
	Cookie *http.Cookie
	Token  string
}

func New(hostname, username, password string) (*Connector, error) {
	conn := &Connector{
		Hostname: hostname,
		Username: username,
		Password: password,
	}

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

	// TODO get version

	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		username,
		password,
	}

	result, err := conn.Post("login", data)
	if err != nil {
		return nil, err
	}
	fmt.Println("LOGIN: " + string(result))

	token, err := conn.Get("persons/348/logintoken", true)
	if err != nil {
		return nil, err
	}
	fmt.Println("TOKEN: " + string(token))

	return conn, nil
}

func (conn *Connector) Get(endpoint string, needToken bool) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/%s", conn.Hostname, endpoint)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")

	if needToken && conn.Token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Login %s", conn.Token))
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

	if conn.Token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Login %s", conn.Token))
	}

	response, err := conn.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if len(response.Cookies()) > 0 {
		conn.Cookie = response.Cookies()[0]
	}

	return ioutil.ReadAll(response.Body)
}
