// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"user-demo/user-api/internal/svc"
	"user-demo/user-api/internal/types"
	"user-demo/user-rpc/user"

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
	// 调用 RPC 创建用户
	rpcResp, err := l.svcCtx.UserRpc.CreateUser(l.ctx, &user.CreateUserReq{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}

	// 返回响应
	return &types.RegisterResp{
		Id:   rpcResp.Id,
		Name: req.Name,
	}, nil
}
