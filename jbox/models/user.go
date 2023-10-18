package models

type UserInfo struct {
	ValidStartTime         string     `json:"valid_start_time"`
	DeliverySupport        bool       `json:"delivery_support"`
	Role                   int64      `json:"role"`
	Color                  string     `json:"color"`
	UserName               string     `json:"user_name"`
	Used                   int64      `json:"used"`
	DownloadThresholdSpeed int64      `json:"download_threshold_speed"`
	UploadThresholdSpeed   int64      `json:"upload_threshold_speed"`
	ShareScope             ShareScope `json:"share_scope"`
	Uid                    int64      `json:"uid"`
	MobileChk              bool       `json:"mobile_chk"`
	CloudQuota             int64      `json:"cloud_quota"`
	DocsLimitEnable        int64      `json:"docs_limit_enable"`
	NeManage               int64      `json:"ne_manage"`
	FromDomainAccount      bool       `json:"from_domain_account"`
	Quota                  int64      `json:"quota"`
	Ctime                  string     `json:"ctime"`
	ValidEndTime           string     `json:"valid_end_time"`
	LocalEditSwitch        bool       `json:"local_edit_switch"`
	RegionDesc             string     `json:"region_desc"`
	Email                  string     `json:"email"`
	UserSlug               string     `json:"user_slug"`
	EmailChk               bool       `json:"email_chk"`
	Mobile                 string     `json:"mobile"`
	RegionID               int64      `json:"region_id"`
	Photo                  []string   `json:"photo"`
	LinkSharingEnable      int64      `json:"link_sharing_enable"`
	NetzoneEnable          int64      `json:"netzone_enable"`
	PasswordChangeable     bool       `json:"password_changeable"`
	CloudUsed              int64      `json:"cloud_used"`
	UseCloudQuota          int64      `json:"use_cloud_quota"`
	CloudAllowScan         int64      `json:"cloud_allow_scan"`
	AccountID              int64      `json:"account_id"`
	PersonalSharingEnable  int64      `json:"personal_sharing_enable"`
	UserID                 int64      `json:"user_id"`
	IsBeyondDocsLimit      bool       `json:"is_beyond_docsLimit"`
	PreviewSupport         bool       `json:"preview_support"`
	ValidEnable            int64      `json:"valid_enable"`
	UseLocalQuota          int64      `json:"use_local_quota"`
	Region                 string     `json:"region"`
	CloudAllowShare        int64      `json:"cloud_allow_share"`
	Status                 int64      `json:"status"`
	UseThreshold           int64      `json:"use_threshold"`
}

type ShareScope struct {
	IsAll bool `json:"is_all"`
}
