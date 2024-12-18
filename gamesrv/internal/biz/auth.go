package biz

type MemberAuthInfo struct {
	UserId int64  `json:"userId"`
	AppId  string `json:"appId"`
	Ip     string `json:"ip"`
}

type MemberAuthRepo interface {
	StoreAuthInfo(token string, roomId int64, info *MemberAuthInfo) error
	GetAuthInfo(token string, roomId int64) (*MemberAuthInfo, error)
}
