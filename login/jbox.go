package login

import (
	"fmt"
	"io"
	"jtrans/utils"
	"net/http"
	"net/url"
	"strings"
)

var (
	jboxUrl *url.URL
	authUrl *url.URL
)

func init() {
	jboxUrl, _ = url.Parse("https://jbox.sjtu.edu.cn/")
	authUrl, _ = url.Parse("https://jaccount.sjtu.edu.cn/jaccount")
}

func loginJbox(jaAuthCookie string) (string, error) {
	cli.Jar.SetCookies(authUrl, []*http.Cookie{
		{Name: "JAAuthCookie", Value: jaAuthCookie},
	})
	resp, err := cli.Get("https://jbox.sjtu.edu.cn/")
	if err != nil {
		return "", err
	}

	if !utils.IsSuccessStatusCode(resp.StatusCode) {
		return "", fmt.Errorf("服务器响应%d", resp.StatusCode)
	}

	if strings.Contains(resp.Request.URL.Host, "jaccount") {
		return "", fmt.Errorf("未成功认证")
	}

	defer resp.Body.Close()
	var content []byte
	content, err = io.ReadAll(resp.Body)
	body := strings.ToLower(string(content))
	if strings.Contains(body, "vpn") {
		return "", fmt.Errorf("校外访问请使用交大vpn")
	}

	cookies := cli.Jar.Cookies(jboxUrl)
	cookiesList := make([]string, len(cookies))
	for i, c := range cookies {
		cookiesList[i] = c.String()
	}
	cookiesString := strings.Join(cookiesList, "; ")
	return cookiesString, nil
}
