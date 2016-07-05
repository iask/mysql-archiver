package sso

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	SSO_URL = "http://sso.pt.xxx.com"
)

type SSO struct {
	SSODomain  string
	BrokerName string
	SecretKey  string
	Credential string
}

func (this *SSO) GetUser() (string, error) {
	url := fmt.Sprintf("%s/login/broker/%s/broker_cookies/%s/user", this.SSODomain, this.BrokerName, this.Credential)
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	user, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(user), err
}

func (this *SSO) IsLogin() (bool, error) {
	url := fmt.Sprintf("%s/login/broker/%s/broker_cookies/%s/check", this.SSODomain, this.BrokerName, this.Credential)
	res, err := http.Get(url)
	if err != nil {
		return false, err
	}
	isLogin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}
	return string(isLogin) == "1", nil
}

func (this *SSO) GetLogoutUrl() string {
	return fmt.Sprintf("%s/login/logout?broker_name=%s", this.SSODomain, this.BrokerName)
}

func (this *SSO) GetLoginUrl(callback string) string {
	return fmt.Sprintf("%s/login?broker_cookies=%s&callback=%s", this.SSODomain, this.Credential, callback)
}

func (this *SSO) GenerateCredential() error {
	uri := this.SSODomain + "/login/broker_cookies"

	data := make(url.Values)
	data.Set("broker_name", this.BrokerName)
	data.Add("secret_key", this.SecretKey)
	res, err := http.PostForm(uri, data)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	credential, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	this.Credential = string(credential)
	return nil
}

func (this *SSO) GetApiUsername(w http.ResponseWriter, r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth != "" {
		info := strings.Split(auth, ";")
		if len(info) == 3 {
			userUrl := fmt.Sprintf("/mias/api/user_ip/%s/auth/%s/username", r.RemoteAddr, auth)
			url := fmt.Sprintf("%s%s", this.SSODomain, userUrl)
			ret, err := http.Get(url)
			if err != nil {
				// Empty string means still not login yet
				return "", nil
			}
			if ret.StatusCode == 200 {
				data, err := ioutil.ReadAll(ret.Body)
				if err != nil {
					return "", err
				}

				var OkPkg struct {
					UserName string `json:"user_name"`
					Code     int    `json:"http_code"`
				}
				err = json.Unmarshal(data, &OkPkg)
				if err == nil {
					if OkPkg.Code == 200 {
						return OkPkg.UserName, nil
					}
				}

				var ErrPkg struct {
					Msg  string `json:"msg"`
					Code int    `json:"http_code"`
				}
				err = json.Unmarshal(data, &ErrPkg)
				if err == nil {
					http.Error(w, ErrPkg.Msg, ErrPkg.Code)
					return "", nil
				}
				return "", err
			}
		}
	}
	return "", nil
}

func New(broker, key string) *SSO {
	return &SSO{
		SSODomain:  SSO_URL,
		BrokerName: broker,
		SecretKey:  key,
	}
}
