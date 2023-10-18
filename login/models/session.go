package models

type Session struct {
	JAAuthCookie string `json:"JAAuthCookie,required"`
	JboxCookies  string `json:"jbox_cookies,required"`
	UserToken    string `json:"user_token,required"`
}
