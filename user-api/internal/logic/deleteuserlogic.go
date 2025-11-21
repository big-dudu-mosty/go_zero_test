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

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp *types.DeleteUserResp, err error) {
	// 调用 RPC 删除用户
	rpcResp, err := l.svcCtx.UserRpc.DeleteUser(l.ctx, &user.DeleteUserReq{
		Id: req.Id,
	})
	if err != nil {
		logx.Errorf("RPC 删除用户失败: %v", err)
		return nil, err
	}

	logx.Infof("用户删除成功: id=%d", req.Id)

	return &types.DeleteUserResp{
		Message: rpcResp.Message,
	}, nil
}
