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

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserReq) (resp *types.GetUserResp, err error) {
	// 调用 RPC 获取用户
	rpcResp, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.GetUserReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	// 组装返回数据
	return &types.GetUserResp{
		Id:    rpcResp.Id,
		Name:  rpcResp.Name,
		Email: rpcResp.Email,
	}, nil
}
