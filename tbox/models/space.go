package models

type PersonalSpaceInfo struct {
	LibraryID   string `json:"libraryId"`
	SpaceID     string `json:"spaceId"`
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
	Status      int64  `json:"status"`
	Message     string `json:"message"`
}
