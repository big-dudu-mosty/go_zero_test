package logic

import (
	"context"

	"user-demo/user-rpc/internal/svc"
	"user-demo/user-rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserLogic) DeleteUser(in *user.DeleteUserReq) (*user.DeleteUserResp, error) {
	// 1. 先查询用户是否存在
	_, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return nil, err
	}

	// 2. 删除用户
	err = l.svcCtx.UserModel.Delete(l.ctx, in.Id)
	if err != nil {
		logx.Errorf("删除用户失败: %v", err)
		return nil, err
	}

	logx.Infof("用户删除成功: id=%d", in.Id)

	return &user.DeleteUserResp{
		Message: "删除成功",
	}, nil
}
