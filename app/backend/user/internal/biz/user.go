package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "maniverse/api/backend/user"
	"maniverse/pkg/util/process"
)

type User struct {
	ID       uint64
	UUID     string
	Phone    string
	Nickname string
	Password string
}

// UserRepo is a Greater repo.
type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	FindByID(context.Context, int64) (*User, error)
	FindByPhone(context.Context, string) (*User, error)
	ListByHello(context.Context, string) ([]*User, error)
	ListAll(context.Context) ([]*User, error)
}

// UserUseCase is a User UseCase.
type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUseCase new a User UseCase.
func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "usecase/user"))}
}

// CreateUser creates a User, and returns the new User.
func (uc *UserUseCase) CreateUser(ctx context.Context, u *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	user := &User{
		UUID:     process.UUID(),
		Phone:    u.Phone,
		Nickname: u.Nickname,
		Password: u.Password,
	}
	user, err := uc.repo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	res := &pb.CreateUserReply{
		Id: user.ID,
	}
	return res, nil
}
func (uc *UserUseCase) GetUserByPhone(ctx context.Context, u *pb.GetUserByPhoneReq) (*pb.GetUserByPhoneReply, error) {
	user, err := uc.repo.FindByPhone(ctx, u.Phone)
	if err != nil {
		return nil, err
	}
	res := &pb.GetUserByPhoneReply{
		Id:    int64(user.ID),
		Phone: user.Phone,
	}
	return res, nil
}
