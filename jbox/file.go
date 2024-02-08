package jbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"jtrans/constants"
	"jtrans/jbox/models"
	"jtrans/utils"

	"github.com/spf13/cast"
)

func (c *Client) GetDirectoryInfo(targetPath string, page int) (*models.DirectoryInfo, error) {
	url := c.baseUrl + "/v2/metadata_page/databox"
	data := map[string]string{
		"path_type":   "self",
		"target_path": targetPath,
		"page_size":   "50",
		"page_num":    cast.ToString(page),
	}

	resp, err := c.postUrlEncoded(url, map[string]string{"S": c.S}, data)
	if err != nil {
		return nil, err
	}
	if !utils.IsSuccessStatusCode(resp.StatusCode) {
		return nil, fmt.Errorf("服务器响应%d", resp.StatusCode)
	}

	info := models.DirectoryInfo{}
	err = utils.UnmarshalJson[models.DirectoryInfo](resp, &info)

	if err != nil {
		return nil, err
	}
	if info.Type == "error" {
		return nil, fmt.Errorf("服务器返回失败：%s", info.Message)
	}
	return &info, nil
}

func (c *Client) GetFileInfo(targetPath string) (*models.FileInfo, error) {
	url := c.baseUrl + "/v2/metadata_page/databox/" + targetPath
	resp, err := c.getRequest(url, nil, map[string]string{
		"S": c.S,
	})
	if err != nil {
		return nil, err
	}
	if !utils.IsSuccessStatusCode(resp.StatusCode) {
		return nil, fmt.Errorf("服务器响应%d", resp.StatusCode)
	}
	info := models.FileInfo{}
	err = utils.UnmarshalJson[models.FileInfo](resp, &info)
	if err != nil {
		return nil, err
	}
	if info.Type == "error" {
		return nil, fmt.Errorf("服务器返回失败：%s", info.Message)
	}
	return &info, nil
}

func (c *Client) GetChunk(path string, chunkNo, chunkSize int64, onProgress models.DownloadProgressHandler) ([]byte, error) {
	return c.DownloadChunk(path, (chunkNo-1)*constants.ChunkSize, chunkSize, onProgress)
}

func (c *Client) DownloadChunk(path string, start, size int64, onProgress models.DownloadProgressHandler) ([]byte, error) {
	resp, err := c.getRequest("https://jbox.sjtu.edu.cn:10081/v2/files/databox"+path,
		map[string]string{
			"Range":   fmt.Sprintf("bytes=%d-%d", start, start+size-1),
			"Referer": c.baseUrl,
		},
		map[string]string{
			"path_type": "self",
			"S":         c.S,
		})
	if err != nil {
		return nil, err
	}
	var bufferSize int64 = 81920 / 2
	dst := bytes.NewBuffer(nil)
	var bytesRead int64
	var downloaded int64 = 0
	for {
		bytesRead, err = io.CopyN(dst, resp.Body, bufferSize)
		downloaded += bytesRead
		if onProgress != nil {
			onProgress(downloaded, size)
		}
		if bytesRead == 0 || errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return dst.Bytes(), nil
}

func (c *Client) BatchMove(fromPaths []string, toPath string) (*models.BatchMoveResult, error) {
	url := c.baseUrl + "/v2/fileops/batch_move"

	dirInfo, err := c.GetDirectoryInfo(toPath, 0)
	if err != nil {
		return nil, err
	}

	data := models.BatchMoveData{
		To: models.ToData{
			Root:     "databox",
			Path:     toPath,
			PathType: "self",
			Neid:     dirInfo.Neid,
			Nsid:     dirInfo.Nsid,
		},
		From: []models.FromData{},
	}

	for _, fromPath := range fromPaths {
		fileInfo, err := c.GetFileInfo(fromPath)
		if err != nil {
			return nil, err
		}
		data.From = append(data.From, models.FromData{
			Root:     "databox",
			Path:     fromPath,
			PathType: "self",
			Nsid:     fileInfo.Nsid,
			Neid:     fileInfo.Neid,
		})
	}

	marshalled, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := c.postUrlEncoded(url, map[string]string{
		"S":          c.S,
		"uid":        cast.ToString(c.uid),
		"account_id": "1",
	}, map[string]string{
		"json": string(marshalled),
	})

	if err != nil {
		return nil, err
	}

	if !utils.IsSuccessStatusCode(resp.StatusCode) {
		return nil, fmt.Errorf("服务器响应%d", resp.StatusCode)
	}

	result := models.BatchMoveResult{}
	err = utils.UnmarshalJson[models.BatchMoveResult](resp, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
