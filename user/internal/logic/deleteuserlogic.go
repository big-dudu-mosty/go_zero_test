// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"user-demo/user/internal/svc"
	"user-demo/user/internal/types"

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
	// 1. 先查询用户是否存在
	_, err = l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return nil, err
	}

	// 2. 删除用户
	err = l.svcCtx.UserModel.Delete(l.ctx, req.Id)
	if err != nil {
		logx.Errorf("删除用户失败: %v", err)
		return nil, err
	}

	logx.Infof("用户删除成功: id=%d", req.Id)

	return &types.DeleteUserResp{
		Message: "删除成功",
	}, nil
}
