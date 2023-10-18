package login

import (
	"fmt"
	"jtrans/login/models"
	"jtrans/utils"
	"net/http"
	"regexp"
	"strings"
)

func loginTbox(jaAuthCookie string) (string, error) {
	cli.Jar.SetCookies(authUrl, []*http.Cookie{
		{Name: "JAAuthCookie", Value: jaAuthCookie},
	})
	resp, err := cli.Get("https://pan.sjtu.edu.cn/user/v1/sign-in/sso-login-redirect/xpw8ou8y")
	if err != nil {
		return "", err
	}

	if !utils.IsSuccessStatusCode(resp.StatusCode) {
		return "", fmt.Errorf("服务器响应%d", resp.StatusCode)
	}

	if strings.Contains(resp.Request.URL.Host, "jaccount") {
		return "", fmt.Errorf("未成功认证")
	}

	reg := regexp.MustCompile("code=(.+?)&state=")
	matches := reg.FindStringSubmatch(resp.Request.URL.String())

	if len(matches) == 0 {
		panic(fmt.Errorf("未找到回调code"))
	}
	code := matches[len(matches)-1]
	nextUrl := "https://pan.sjtu.edu.cn/user/v1/sign-in/verify-account-login/xpw8ou8y?device_id=Chrome+116.0.0.0&type=sso&credential=" + code

	resp, err = cli.Post(nextUrl, "", nil)
	if !utils.IsSuccessStatusCode(resp.StatusCode) {
		return "", fmt.Errorf("服务器响应%d", resp.StatusCode)
	}

	loginRes := models.TboxLoginResult{}
	err = utils.UnmarshalJson[models.TboxLoginResult](resp, &loginRes)
	if err != nil {
		return "", err
	}

	if loginRes.Status != 0 {
		return "", fmt.Errorf("服务器返回失败")
	}

	if len(loginRes.UserToken) != 128 {
		return "", fmt.Errorf("UserToken无效")
	}

	return loginRes.UserToken, nil
}
