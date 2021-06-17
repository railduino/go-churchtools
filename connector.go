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

type User struct {
	Hostname string
	Username string
	Password string
}

type Connector struct {
	User
	PersonID int
	Token    string

	Build   string
	Version string

	Client *http.Client
	Cookie *http.Cookie
}

type MetaPerson struct {
	ID   int
	Name string
}

type MetaInfo struct {
	CreatedPerson  MetaPerson
	CreatedDate    string
	ModifiedPerson MetaPerson
	ModifiedDate   string
}

type Permissions struct {
	CanEdit          bool
	CanUseExpertMode bool
	AllowPosting     bool
}

type InfoResult struct {
	Build   string
	Version string
}

type LoginData struct {
	Status   string
	Message  string
	PersonID int
	Location string
}

type LoginResult struct {
	Data LoginData
}

type LoginTokenResult struct {
	Data string
}

func (user User) Login() (*Connector, error) {
	conn := &Connector{User: user}

	if user.Hostname == "" {
		return nil, fmt.Errorf("missing Hostname")
	}
	if user.Username == "" {
		return nil, fmt.Errorf("missing Username")
	}

	// Setup the HTTPS client
	rootCAPool := x509.NewCertPool()
	rootCA, err := ioutil.ReadFile(ROOT_CA_FILE)
	if err != nil {
		return nil, err
	}
	rootCAPool.AppendCertsFromPEM(rootCA)

	conn.Client = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout: 60 * time.Second,
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

	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		conn.Username,
		conn.Password,
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

	// If LoginToken is not known, try to obtain
	if conn.Token == "" {
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
	}

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
