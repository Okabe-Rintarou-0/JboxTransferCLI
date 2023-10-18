package tbox

import (
	"fmt"
	"jtrans/tbox/models"
	"jtrans/utils"
)

func (c *Client) GetPersonalSpaceInfo() (*models.PersonalSpaceInfo, error) {
	url := "/user/v1/space/1/personal"
	resp, err := c.postRequest(c.baseUrl+url, map[string]string{
		"user_token": c.userToken,
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
		return nil, fmt.Errorf(errMessage.Message)
	}
	info := models.PersonalSpaceInfo{}
	err = utils.UnmarshalJson[models.PersonalSpaceInfo](resp, &info)
	if err != nil {
		return nil, err
	}

	if info.Status != 0 {
		return nil, fmt.Errorf("服务器返回失败：%s", info.Message)
	}
	return &info, nil
}
