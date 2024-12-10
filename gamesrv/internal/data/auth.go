package data

import "github.com/StellrisJAY/cloud-emu/gamesrv/internal/biz"

type MemberAuthRepo struct {
	data *Data
}

func (m *MemberAuthRepo) StoreAuthInfo(token string, info *biz.MemberAuthInfo) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemberAuthRepo) GetAuthInfo(token string) (*biz.MemberAuthInfo, error) {
	//TODO implement me
	panic("implement me")
}

func NewMemberAuthRepo(data *Data) biz.MemberAuthRepo {
	return &MemberAuthRepo{data: data}
}
