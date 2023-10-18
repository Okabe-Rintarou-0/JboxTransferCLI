package tbox

import (
	"fmt"
	"jtrans/tbox/models"
	"jtrans/utils"

	"github.com/spf13/cast"
)

func (c *Client) GetDirectoryInfo(dirPath string,
	pagination *models.PaginationOption,
	order *models.OrderOption,
	filter string,
) (*models.DirectoryInfo, error) {
	url := fmt.Sprintf("/api/v1/directory/%s/%s/%s", c.libraryId, c.spaceId, dirPath)
	query := map[string]string{
		"access_token": c.accessToken,
	}

	if pagination != nil {
		query["page"] = cast.ToString(pagination.Page)
		query["page_size"] = cast.ToString(pagination.PageSize)
	}

	if order != nil {
		query["order_by"] = cast.ToString(order.By)
		query["order_by_type"] = cast.ToString(order.Type)
	}

	if len(filter) > 0 {
		query["filter"] = filter
	}

	resp, err := c.getRequest(c.baseUrl+url, query)
	if err != nil {
		return nil, err
	}

	info := models.DirectoryInfo{}
	err = utils.UnmarshalJson[models.DirectoryInfo](resp, &info)
	return &info, err
}

func (c *Client) StartChunkUpload(path string, chunkCount int64) (*models.StartChunkUploadResult, error) {
	url := c.baseUrl + fmt.Sprintf("/api/v1/file/%s/%s/%s", c.libraryId, c.spaceId, path)
	if chunkCount > 50 {
		chunkCount = 50
	}
	chunks := make([]int64, chunkCount)
	var i int64
	for i = 1; i <= chunkCount; i++ {
		chunks[i-1] = i
	}
	data := map[string]interface{}{}
	data["partNumberRange"] = chunks

	resp, err := c.postJson(url, map[string]string{
		"multipart":                    "null",
		"conflict_resolution_strategy": "rename",
		"access_token":                 c.accessToken,
	}, data)
	if err != nil {
		return nil, err
	}

	info := models.StartChunkUploadResult{}
	err = utils.UnmarshalJson[models.StartChunkUploadResult](resp, &info)
	return &info, err
}

func (c *Client) RenewChunkUpload(confirmKey string, partNumberRange []int64) (*models.StartChunkUploadResult, error) {
	url := c.baseUrl + fmt.Sprintf("/api/v1/file/%s/%s/%s", c.libraryId, c.spaceId, confirmKey)
	data := map[string]interface{}{}
	data["partNumberRange"] = partNumberRange

	resp, err := c.postJson(url, map[string]string{
		"renew":        "null",
		"access_token": c.accessToken,
	}, data)
	if err != nil {
		return nil, err
	}

	info := models.StartChunkUploadResult{}
	err = utils.UnmarshalJson[models.StartChunkUploadResult](resp, &info)
	return &info, err
}

func (c *Client) ConfirmChunkUpload(confirmKey string) (*models.ConfirmChunkUploadResult, error) {
	url := c.baseUrl + fmt.Sprintf("/api/v1/file/%s/%s/%s", c.libraryId, c.spaceId, confirmKey)

	resp, err := c.postRequest(url, map[string]string{
		"confirm":                      "null",
		"conflict_resolution_strategy": "rename",
		"access_token":                 c.accessToken,
	}, nil)
	if err != nil {
		return nil, err
	}

	info := models.ConfirmChunkUploadResult{}
	err = utils.UnmarshalJson[models.ConfirmChunkUploadResult](resp, &info)
	return &info, err
}

func (c *Client) GetChunkUploadInfo(confirmKey string) (*models.ChunkUploadInfo, error) {
	url := c.baseUrl + fmt.Sprintf("/api/v1/file/%s/%s/%s", c.libraryId, c.spaceId, confirmKey)

	resp, err := c.getRequest(url, map[string]string{
		"upload":              "null",
		"no_upload_part_info": "1",
		"access_token":        c.accessToken,
	})
	if err != nil {
		return nil, err
	}

	info := models.ChunkUploadInfo{}
	err = utils.UnmarshalJson[models.ChunkUploadInfo](resp, &info)
	return &info, err
}

func (c *Client) Upload(ctx *models.StartChunkUploadResult, data []byte, partNumber int64,
	onProgress models.UploadProgressHandler) error {
	url := fmt.Sprintf("https://%s%s", ctx.Domain, ctx.Path)
	headerInfo := ctx.Parts[cast.ToString(partNumber)].Headers
	bufferSize := 81920 / 2
	body := utils.NewRequestStream(data, bufferSize, onProgress)
	_, err := c.putRequest(url,
		map[string]string{
			"Accept":               "*/*",
			"Accept-Encoding":      "gzip, deflate, br",
			"Accept-Language":      "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
			"x-amz-date":           headerInfo.XAmzDate,
			"authorization":        headerInfo.Authorization,
			"x-amz-content-sha256": headerInfo.XAmzContentSha256,
		},
		map[string]string{
			"uploadId":   cast.ToString(ctx.UploadID),
			"partNumber": cast.ToString(partNumber),
		}, body)
	return err
}

func (c *Client) CreateDirectory(dirPath string) (*models.ErrorMessage, error) {
	url := c.baseUrl + fmt.Sprintf("/api/v1/directory/%s/%s/%s", c.libraryId, c.spaceId, dirPath)

	resp, err := c.putRequest(url, nil, map[string]string{
		"conflict_resolution_strategy": "ask",
		"access_token":                 c.accessToken,
	}, nil)
	if err != nil {
		return nil, err
	}
	errMessage := models.ErrorMessage{}
	if !utils.IsSuccessStatusCode(resp.StatusCode) {
		err = utils.UnmarshalJson[models.ErrorMessage](resp, &errMessage)
		if err != nil {
			return nil, fmt.Errorf("服务器响应%d", resp.StatusCode)
		}
		return &errMessage, fmt.Errorf(errMessage.Message)
	}

	err = utils.UnmarshalJson[models.ErrorMessage](resp, &errMessage)
	if err != nil {
		return nil, err
	}
	if errMessage.Status != 0 {
		return &errMessage, fmt.Errorf("服务器返回失败：%s", errMessage.Message)
	}

	return &errMessage, nil
}

func (c *Client) GetFileInfo(filePath string) (*models.FileInfo, error) {
	url := fmt.Sprintf("/api/v1/directory/%s/%s/%s", c.libraryId, c.spaceId, filePath)
	query := map[string]string{
		"info":         "",
		"access_token": c.accessToken,
	}

	resp, err := c.getRequest(c.baseUrl+url, query)
	if err != nil {
		return nil, err
	}

	info := models.FileInfo{}
	err = utils.UnmarshalJson[models.FileInfo](resp, &info)
	return &info, err
}
