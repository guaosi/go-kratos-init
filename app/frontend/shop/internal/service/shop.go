package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"maniverse/app/frontend/shop/internal/biz"

	pb "maniverse/api/frontend/shop"
)

type ShopService struct {
	pb.UnimplementedShopServer
	uc  *biz.UserUseCase
	ac  *biz.AuthUseCase
	log *log.Helper
}

func NewShopService(uc *biz.UserUseCase, ac *biz.AuthUseCase, logger log.Logger) *ShopService {
	return &ShopService{
		uc:  uc,
		ac:  ac,
		log: log.NewHelper(log.With(logger, "module", "service/shop")),
	}
}

func (s *ShopService) CreateShop(ctx context.Context, req *pb.CreateShopRequest) (*pb.CreateShopReply, error) {
	return &pb.CreateShopReply{}, nil
}
func (s *ShopService) UpdateShop(ctx context.Context, req *pb.UpdateShopRequest) (*pb.UpdateShopReply, error) {
	return &pb.UpdateShopReply{}, nil
}
func (s *ShopService) DeleteShop(ctx context.Context, req *pb.DeleteShopRequest) (*pb.DeleteShopReply, error) {
	return &pb.DeleteShopReply{}, nil
}
func (s *ShopService) GetShop(ctx context.Context, req *pb.GetShopRequest) (*pb.GetShopReply, error) {
	return &pb.GetShopReply{}, nil
}
func (s *ShopService) ListShop(ctx context.Context, req *pb.ListShopRequest) (*pb.ListShopReply, error) {
	return &pb.ListShopReply{}, nil
}
func (s *ShopService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	return s.uc.CreateUser(ctx, req)
}
func (s *ShopService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginReply, error) {
	return s.ac.Login(ctx, req)
}
