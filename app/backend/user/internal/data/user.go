package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"maniverse/api/backend/user"
	"maniverse/app/backend/user/internal/biz"
	"maniverse/pkg/errors/mysql"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) Save(ctx context.Context, g *biz.User) (*biz.User, error) {
	u := User{
		UUID:     g.UUID,
		Phone:    g.Phone,
		Password: g.Password,
		NickName: g.Nickname,
	}
	if err := r.data.getDB(ctx).Create(&u).Error; err != nil {
		if mysql.JudgeRecordDuplicate(err) {
			e := mysql.ErrMySQLDataDuplicate
			e = e.WithMetadata(map[string]string{
				"phone": g.Phone,
			})
			return nil, e
		}
		r.log.Errorf("save failed,err:%s,data:%+v", err.Error(), *g)
		return nil, mysql.ErrMySQL
	}
	return &biz.User{ID: u.ID}, nil
}

func (r *userRepo) Update(ctx context.Context, g *biz.User) (*biz.User, error) {
	return g, nil
}

func (r *userRepo) FindByID(context.Context, int64) (*biz.User, error) {
	return nil, nil
}

func (r *userRepo) ListByHello(context.Context, string) ([]*biz.User, error) {
	return nil, nil
}

func (r *userRepo) ListAll(context.Context) ([]*biz.User, error) {
	return nil, nil
}
func (r *userRepo) FindByPhone(ctx context.Context, phone string) (*biz.User, error) {
	u := &User{}
	err := r.data.getDB(ctx).Where("phone = ?", phone).First(u).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			e := user.ErrorUserNotFound("user not found")
			e = e.WithMetadata(map[string]string{
				"phone": phone,
			})
			return nil, e
		}
		r.log.Errorf("find by phone failed,err:%s,phone:%s", err.Error(), phone)
		return nil, mysql.ErrMySQL
	}

	return &biz.User{
		ID:       u.ID,
		UUID:     u.UUID,
		Phone:    u.Phone,
		Nickname: u.NickName,
		Password: u.Password,
	}, nil
}
