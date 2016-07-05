package web

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	//"github.com/hdksky/autopilot/sso"
	sso "archiver/sso"
)

const (
	BROCKER_NAME = "app"
	SECRET_KEY   = "xxxxxxxxxxxxxxxxxxxxxxxx"
)

// User auth info
// Return by xbox
type User struct {
	Id, Token string
}

type AuthInfo struct {
	Ok   bool
	Msg  string
	Data User
}

type UserInfo struct {
	Id        int
	Name      string
	Email     string
	Cellphone string
}

type Res struct {
	Error  int
	Msg    string
	Detail []UserInfo
}

func isLogin(r *http.Request) bool {
	cookie, err := r.Cookie("user")
	if err != nil {
		ErrorLog(fmt.Sprintf("Read user from cookie fail: %v", err))
		return false
	}

	return cookie.Value != ""
}

func validUser(r *http.Request, user string) bool {
	userInCookie, err := r.Cookie("user")
	if err == nil && userInCookie.Value != "" {
		return cryptToken(userInCookie.Value) == user
	}
	return false
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "credential", Value: "", Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "token", Value: "", Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "user", Value: "", Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)

	url := sso.New(BROCKER_NAME, SECRET_KEY).GetLogoutUrl()
	http.Redirect(w, r, url, 302)
}

type userInfo struct {
	Chinese_name    string
	Last_time       string
	User_name       string
	Id              string
	Session_cookies string
}

func Login(w http.ResponseWriter, r *http.Request) (userInfo, error) {
	var ui userInfo

	ssoCli := sso.New(BROCKER_NAME, SECRET_KEY)
	user, err := ssoCli.GetApiUsername(w, r)
	if err == nil && user != "" {
		cookie := http.Cookie{Name: "user", Value: user, Path: "/", MaxAge: 86400}
		http.SetCookie(w, &cookie)
		cookie = http.Cookie{Name: "token", Value: cryptToken(user), Path: "/", MaxAge: 86400}
		http.SetCookie(w, &cookie)
		return ui, nil
	}

	// Generate credential, if does not exists in cookie
	credential, err := r.Cookie("credential")
	if err == nil && credential.Value != "" {
		ssoCli.Credential = credential.Value
	} else {
		err = ssoCli.GenerateCredential()
		if err != nil {
			return ui, err
		}
		cookie := http.Cookie{Name: "credential", Value: ssoCli.Credential, Path: "/", MaxAge: 86400}
		http.SetCookie(w, &cookie)
	}

	// Whether or not login in SSO
	isLoginInSSO, err := ssoCli.IsLogin()
	if err != nil {
		return ui, err
	}
	if !isLoginInSSO {
		ErrorLog("Not login, redirect to sso")
		redirectToSSO(w, r, ssoCli)
		return ui, nil
	}

	// Already login, get user info
	userJson, err := ssoCli.GetUser()
	if err != nil {
		return ui, err
	}
	err = json.Unmarshal([]byte(userJson), &ui)
	if err != nil {
		return ui, err
	}

	cookie := http.Cookie{Name: "token", Value: cryptToken(ui.User_name), Path: "/", MaxAge: 86400}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "user", Value: ui.User_name, Path: "/", MaxAge: 86400}
	http.SetCookie(w, &cookie)
	ErrorLog(fmt.Sprintf("Login success, redirect to %s", r.RequestURI))
	http.Redirect(w, r, r.RequestURI, 302)

	return ui, nil
}

func cryptToken(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func redirectToSSO(w http.ResponseWriter, r *http.Request, ssoClient *sso.SSO) {
	loginUrl := ssoClient.GetLoginUrl("http://" + r.Host + r.URL.String())
	http.Redirect(w, r, loginUrl, 302)
}
