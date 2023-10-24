package jbox

import (
	"bytes"
	"encoding/json"
	"io"
	"jtrans/jbox/models"
	"jtrans/utils"
	"net/http"
	urllib "net/url"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type IClient interface {
	GetUserInfo() (*models.UserInfo, error)
	GetUserId() (int64, error)
	GetFileInfo(targetPath string) (*models.FileInfo, error)
	GetDirectoryInfo(targetPath string, page int) (*models.DirectoryInfo, error)
	GetChunk(path string, chunkNo, chunkSize int64, onProgress models.DownloadProgressHandler) ([]byte, error)
	DownloadChunk(path string, start, size int64, onProgress models.DownloadProgressHandler) ([]byte, error)
	BatchMove(fromPaths []string, toPath string) (*models.BatchMoveResult, error)
}

type Client struct {
	cookies       string
	baseUrl       string
	headers       map[string]string
	S             string
	uid           int64
	XLENOVOSESSID string
}

func NewClient(cookies string) (IClient, error) {
	cookiesMap := utils.FromCookiesString(cookies)
	cli := &Client{
		cookies: cookies,
		baseUrl: "https://jbox.sjtu.edu.cn",
		headers: map[string]string{
			"Cookie":           cookies,
			"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"X-LENOVO-SESS-ID": cookiesMap["X-LENOVO-SESS-ID"],
		},
		S: cookiesMap["S"],
	}
	var err error
	t := 3
	for t > 0 {
		t -= 1
		cli.uid, err = cli.GetUserId()
		if err == nil {
			break
		}
	}
	return cli, err
}

func (c *Client) postRequest(url string, query map[string]string, body io.Reader) (*http.Response, error) {
	return utils.DoRequest(http.MethodPost, c.baseUrl+url, c.headers, query, body)
}

func (c *Client) postJson(url string, query map[string]string, body any) (*http.Response, error) {
	headers := map[string]string{}
	for k, v := range c.headers {
		headers[k] = v
	}
	headers["Content-Type"] = "application/json"
	marshalled, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return utils.DoRequest(http.MethodPost, url, headers, query, bytes.NewReader(marshalled))
}

func (c *Client) postUrlEncoded(url string, query map[string]string, data map[string]string) (*http.Response, error) {
	headers := map[string]string{}
	for k, v := range c.headers {
		headers[k] = v
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	body := urllib.Values{}
	for k, v := range data {
		body.Add(k, v)
	}
	return utils.DoRequest(http.MethodPost, url, headers, query, strings.NewReader(body.Encode()))
}

func (c *Client) getRequest(url string, headers map[string]string, query map[string]string) (*http.Response, error) {
	if headers != nil {
		for k, v := range c.headers {
			headers[k] = v
		}
	} else {
		headers = c.headers
	}
	return utils.DoRequest(http.MethodGet, url, headers, query, nil)
}

func (c *Client) packTimestamp() string {
	return cast.ToString(time.Now().UnixMilli()) + "000"
}
