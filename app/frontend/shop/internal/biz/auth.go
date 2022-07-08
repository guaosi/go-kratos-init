package biz

import (
	"context"
	shopapi "maniverse/api/frontend/shop"
	"maniverse/app/frontend/shop/internal/conf"
	"maniverse/pkg/util/jwt"
)

type AuthUseCase struct {
	key      string
	timeout  int64
	userRepo UserRepo
}

func NewAuthUseCase(conf *conf.Auth, userRepo UserRepo) *AuthUseCase {
	return &AuthUseCase{
		key:      conf.ApiKey,
		timeout:  conf.Timeout,
		userRepo: userRepo,
	}
}

func (receiver *AuthUseCase) Login(ctx context.Context, req *shopapi.LoginReq) (*shopapi.LoginReply, error) {
	// get user
	user, err := receiver.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	claims := jwt.CustomClaimsConfiguration{
		UserID:      user.ID,
		NickName:    user.Nickname,
		AuthorityId: user.ID,
		BelongTo:    user.ID,
		SecretKey:   receiver.key,
		Timeout:     receiver.timeout,
	}
	token, err := jwt.GenerateToken(&claims)
	if err != nil {
		return nil, err
	}
	return &shopapi.LoginReply{
		Data: &shopapi.LoginReply_Data{
			Token: token,
		},
	}, nil
}
