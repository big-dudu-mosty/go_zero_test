// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"user-demo/user/internal/svc"
	"user-demo/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserReq) (resp *types.UpdateUserResp, err error) {
	// 1. 先查询用户是否存在
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return nil, err
	}

	// 2. 更新用户信息
	user.Name = req.Name
	user.Email = req.Email

	// 3. 保存到数据库
	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		logx.Errorf("更新用户失败: %v", err)
		return nil, err
	}

	logx.Infof("用户更新成功: id=%d, name=%s, email=%s", req.Id, req.Name, req.Email)

	return &types.UpdateUserResp{
		Message: "更新成功",
	}, nil
}
