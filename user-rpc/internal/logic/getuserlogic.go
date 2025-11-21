package logic

import (
	"context"

	"user-demo/user-rpc/internal/svc"
	"user-demo/user-rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserReq) (*user.GetUserResp, error) {
	// 1. 查询用户
	userEntity, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 2. 组装返回数据
	return &user.GetUserResp{
		Id:    userEntity.ID,
		Name:  userEntity.Name,
		Email: userEntity.Email,
		Age:   userEntity.Age,
	}, nil
}
