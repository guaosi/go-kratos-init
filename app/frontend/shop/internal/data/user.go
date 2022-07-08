package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/singleflight"
	userapi "maniverse/api/backend/user"
	"maniverse/app/frontend/shop/internal/biz"
)

var _ biz.UserRepo = (*userRepo)(nil)

type userRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/user")),
		sg:   &singleflight.Group{},
	}
}
func (rp *userRepo) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	reply, err := rp.data.uc.CreateUser(ctx, &userapi.CreateUserRequest{
		Phone:    u.Phone,
		Nickname: u.Nickname,
		Password: u.Password,
	})
	if err != nil {
		return u, err
	}
	u.ID = reply.Id

	return u, err
}
func (rp *userRepo) FindByPhone(ctx context.Context, phone string) (*biz.User, error) {
	result, err, _ := rp.sg.Do(fmt.Sprintf("find_user_by_phone_%s", phone), func() (interface{}, error) {
		u, err := rp.data.uc.GetUserByPhone(ctx, &userapi.GetUserByPhoneReq{
			Phone: phone,
		})
		if err != nil {
			return nil, err
		}
		return &biz.User{
			ID:    uint64(u.Id),
			Phone: u.Phone,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.User), nil
}
