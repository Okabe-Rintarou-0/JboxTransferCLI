package tbox

import (
	"bytes"
	"encoding/json"
	"io"
	"jtrans/tbox/models"
	"jtrans/utils"
	"net/http"
)

type IClient interface {
	GetPersonalSpaceInfo() (*models.PersonalSpaceInfo, error)
	StartChunkUpload(path string, chunkCount int64) (*models.StartChunkUploadResult, error)
	RenewChunkUpload(confirmKey string, partNumberRange []int64) (*models.StartChunkUploadResult, error)
	ConfirmChunkUpload(confirmKey string) (*models.ConfirmChunkUploadResult, error)
	GetChunkUploadInfo(confirmKey string) (*models.ChunkUploadInfo, error)
	Upload(ctx *models.StartChunkUploadResult, data []byte, partNumber int64, onProgress models.UploadProgressHandler) error
	CreateDirectory(dirPath string) (*models.ErrorMessage, error)
	GetFileInfo(filePath string) (*models.FileInfo, error)
	GetDirectoryInfo(dirPath string,
		pagination *models.PaginationOption,
		order *models.OrderOption,
		filter string) (*models.DirectoryInfo, error)
}

type Client struct {
	userToken   string
	baseUrl     string
	headers     map[string]string
	libraryId   string
	spaceId     string
	accessToken string
}

func NewClient(userToken string) (IClient, error) {
	var (
		info *models.PersonalSpaceInfo
		err  error
	)
	cli := &Client{
		userToken: userToken,
		baseUrl:   "https://pan.sjtu.edu.cn",
		headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		},
	}
	t := 3
	for t > 0 {
		t -= 1
		info, err = cli.GetPersonalSpaceInfo()
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	cli.libraryId = info.LibraryID
	cli.spaceId = info.SpaceID
	cli.accessToken = info.AccessToken
	return cli, err
}

func (c *Client) postRequest(url string, query map[string]string, body io.Reader) (*http.Response, error) {
	return utils.DoRequest(http.MethodPost, url, c.headers, query, body)
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

func (c *Client) getRequest(url string, query map[string]string) (*http.Response, error) {
	return utils.DoRequest(http.MethodGet, url, c.headers, query, nil)
}

func (c *Client) putRequest(url string, headers map[string]string, query map[string]string, body io.Reader) (*http.Response, error) {
	if headers != nil {
		for k, v := range c.headers {
			headers[k] = v
		}
	} else {
		headers = c.headers
	}
	return utils.DoRequest(http.MethodPut, url, headers, query, body)
}
