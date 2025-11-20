// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"user-demo/user/internal/svc"
	"user-demo/user/internal/types"
	"user-demo/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 组装数据库实体
	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
	}

	// 2. 插入数据库
	result, err := l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		return nil, err
	}

	// 3. 拿到新增 ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 4. 返回响应
	return &types.RegisterResp{
		Id:   id,
		Name: req.Name,
	}, nil
}
