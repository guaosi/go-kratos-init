package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"maniverse/api/frontend/shop"
	tool_jwt "maniverse/pkg/util/jwt"
)

type User struct {
	ID       uint64
	UUID     string
	Phone    string
	Nickname string
	Password string
}
type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}
type UserRepo interface {
	CreateUser(ctx context.Context, u *User) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	log := log.NewHelper(log.With(logger, "module", "usecase/user"))
	return &UserUseCase{
		repo: repo,
		log:  log,
	}
}
func (uc *UserUseCase) CreateUser(ctx context.Context, user *shop.CreateUserRequest) (*shop.CreateUserReply, error) {
	token, _ := jwt.FromContext(ctx)
	data := token.(*tool_jwt.CustomClaims)
	fmt.Println(data.ID, data.AuthorityId)
	u := User{
		Phone:    user.Phone,
		Nickname: user.Nickname,
		Password: user.Password,
	}
	res, err := uc.repo.CreateUser(ctx, &u)
	if err != nil {
		return nil, err
	}
	return &shop.CreateUserReply{
		Data: &shop.CreateUserReply_Data{
			Id: res.ID,
		},
	}, err
}
