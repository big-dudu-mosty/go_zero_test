package logic

import (
	"context"

	"user-demo/user-rpc/internal/svc"
	"user-demo/user-rpc/model"
	"user-demo/user-rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *user.CreateUserReq) (*user.CreateUserResp, error) {
	// 1. 组装数据库实体
	userEntity := &model.User{
		Name:  in.Name,
		Email: in.Email,
	}

	// 2. 插入数据库（PostgreSQL 使用 RETURNING id 直接返回新增 ID）
	id, err := l.svcCtx.UserModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, err
	}

	// 3. 返回响应
	return &user.CreateUserResp{
		Id: id,
	}, nil
}
