package models

type UserInfo struct {
	Entities Entities `json:"entities"`
	Errno    int64    `json:"errno"`
	Error    string   `json:"error"`
}

type Entities struct {
	AccountNo     string        `json:"accountNo"`
	Avatars       interface{}   `json:"avatars"`
	Name          string        `json:"name"`
	UserType      string        `json:"userType"`
	UserStyleName string        `json:"userStyleName"`
	Email         string        `json:"email"`
	Code          string        `json:"code"`
	ExpireDate    string        `json:"expireDate"`
	Mobile        string        `json:"mobile"`
	Identities    []Identity    `json:"identities"`
	OrganizeName  string        `json:"organizeName"`
	Status        string        `json:"status"`
	StatusEN      string        `json:"statusEN"`
	ResponseName  interface{}   `json:"responseName"`
	OrganizeID    string        `json:"organizeId"`
	AuthAccounts  []interface{} `json:"authAccounts"`
}

type Identity struct {
	Kind            string      `json:"kind"`
	IsDefault       bool        `json:"isDefault"`
	DefaultOptional bool        `json:"defaultOptional"`
	Code            string      `json:"code"`
	UserType        string      `json:"userType"`
	Organize        Organize    `json:"organize"`
	TopOrganize     *Organize   `json:"topOrganize"`
	Status          *string     `json:"status"`
	ExpireDate      *string     `json:"expireDate"`
	CreateDate      int64       `json:"createDate"`
	UpdateDate      int64       `json:"updateDate"`
	Gjm             *string     `json:"gjm"`
	FacultyType     interface{} `json:"facultyType"`
	PhotoURL        *string     `json:"photoUrl"`
	Type            *Organize   `json:"type"`
	UserStyleName   string      `json:"userStyleName"`
}

type Organize struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type TboxLoginResult struct {
	UserID        int64          `json:"userId"`
	UserToken     string         `json:"userToken"`
	ExpiresIn     int64          `json:"expiresIn"`
	Organizations []Organization `json:"organizations"`
	IsNewUser     bool           `json:"isNewUser"`
	Status        int            `json:"status"`
}

type Organization struct {
	ID             int64         `json:"id"`
	Name           string        `json:"name"`
	ExtensionData  ExtensionData `json:"extensionData"`
	LibraryID      string        `json:"libraryId"`
	LibraryCERT    string        `json:"libraryCert"`
	OrgUser        OrgUser       `json:"orgUser"`
	IsTemporary    bool          `json:"isTemporary"`
	IsLastSignedIn bool          `json:"isLastSignedIn"`
	Expired        bool          `json:"expired"`
	IPLimitEnabled bool          `json:"ipLimitEnabled"`
}

type ExtensionData struct {
	EnableDocPreview             bool               `json:"enableDocPreview"`
	EnableDocEdit                bool               `json:"enableDocEdit"`
	EnableMediaProcessing        bool               `json:"enableMediaProcessing"`
	Logo                         string             `json:"logo"`
	SsoWay                       string             `json:"ssoWay"`
	IPLimit                      IPLimit            `json:"ipLimit"`
	SyncWay                      string             `json:"syncWay"`
	UserLimit                    int64              `json:"userLimit"`
	ExpireTime                   string             `json:"expireTime"`
	EnableShare                  bool               `json:"enableShare"`
	LibraryFlag                  int64              `json:"libraryFlag"`
	AllowProduct                 string             `json:"allowProduct"`
	EditionConfig                EditionConfig      `json:"editionConfig"`
	EnableYufuLogin              bool               `json:"enableYufuLogin"`
	ShowOrgNameInUI              bool               `json:"showOrgNameInUI"`
	WatermarkOptions             WatermarkOptions   `json:"watermarkOptions"`
	EnableWeworkLogin            bool               `json:"enableWeworkLogin"`
	DefaultTeamOptions           DefaultTeamOptions `json:"defaultTeamOptions"`
	DefaultUserOptions           DefaultUserOptions `json:"defaultUserOptions"`
	AllowChangeNickname          bool               `json:"allowChangeNickname"`
	EnableOpenLDAPLogin          bool               `json:"enableOpenLdapLogin"`
	CacheDocPreviewTypes         string             `json:"cacheDocPreviewTypes"`
	EnableViewAllOrgUser         bool               `json:"enableViewAllOrgUser"`
	EnableWindowsAdLogin         bool               `json:"enableWindowsAdLogin"`
	OfficialDocPreviewTypes      string             `json:"officialDocPreviewTypes"`
	IsAccountNotDependentOnPhone bool               `json:"isAccountNotDependentOnPhone"`
}

type DefaultTeamOptions struct {
	DefaultRoleID  int64       `json:"defaultRoleId"`
	SpaceQuotaSize interface{} `json:"spaceQuotaSize"`
}

type DefaultUserOptions struct {
	Enabled                bool   `json:"enabled"`
	AllowPersonalSpace     bool   `json:"allowPersonalSpace"`
	PersonalSpaceQuotaSize string `json:"personalSpaceQuotaSize"`
}

type EditionConfig struct {
	EditionFlag               string `json:"editionFlag"`
	EnableOverseasPhoneNumber bool   `json:"enableOverseasPhoneNumber"`
	EnableOnlineEdit          bool   `json:"enableOnlineEdit"`
}

type IPLimit struct {
	LimitAdmin bool `json:"limitAdmin"`
}

type WatermarkOptions struct {
	ShareWatermarkType      int64 `json:"shareWatermarkType"`
	EnableShareWatermark    bool  `json:"enableShareWatermark"`
	PreviewWatermarkType    int64 `json:"previewWatermarkType"`
	DownloadWatermarkType   int64 `json:"downloadWatermarkType"`
	EnablePreviewWatermark  bool  `json:"enablePreviewWatermark"`
	EnableDownloadWatermark bool  `json:"enableDownloadWatermark"`
}

type OrgUser struct {
	Nickname           string `json:"nickname"`
	Role               string `json:"role"`
	Avatar             string `json:"avatar"`
	Deregister         bool   `json:"deregister"`
	Enabled            bool   `json:"enabled"`
	NeedChangePassword bool   `json:"needChangePassword"`
}
