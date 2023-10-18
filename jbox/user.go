package jbox

import (
	"jtrans/jbox/models"
	"jtrans/utils"
)

func (c *Client) GetUserInfo() (*models.UserInfo, error) {
	url := c.baseUrl + "/v2/user/info/get"
	resp, err := c.getRequest(url, nil, map[string]string{
		"S": c.S,
	})
	if err != nil {
		return nil, err
	}

	info := models.UserInfo{}
	err = utils.UnmarshalJson[models.UserInfo](resp, &info)
	return &info, err
}

func (c *Client) GetUserId() (int64, error) {
	info, err := c.GetUserInfo()
	if err != nil {
		return 0, err
	}
	return info.UserID, nil
}
