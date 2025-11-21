package logic

import (
	"context"

	"user-demo/user-rpc/internal/svc"
	"user-demo/user-rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *user.UpdateUserReq) (*user.UpdateUserResp, error) {
	// 1. 先查询用户是否存在
	userEntity, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return nil, err
	}

	// 2. 更新用户信息
	userEntity.Name = in.Name
	userEntity.Email = in.Email

	// 3. 保存到数据库
	err = l.svcCtx.UserModel.Update(l.ctx, userEntity)
	if err != nil {
		logx.Errorf("更新用户失败: %v", err)
		return nil, err
	}

	logx.Infof("用户更新成功: id=%d, name=%s, email=%s", in.Id, in.Name, in.Email)

	return &user.UpdateUserResp{
		Message: "更新成功",
	}, nil
}
