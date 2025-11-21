// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"user-demo/user/internal/svc"
	"user-demo/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUsersLogic) ListUsers(req *types.ListUsersReq) (resp *types.ListUsersResp, err error) {
	// 添加日志：查看接收到的参数
	logx.Infof("收到请求参数: page=%d, pageSize=%d", req.Page, req.PageSize)

	// page 参数处理
	page := req.Page
	size := req.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	logx.Infof("处理后参数: page=%d, size=%d, offset=%d", page, size, (page-1)*size)

	offset := (page - 1) * size

	// 调用 model
	users, total, err := l.svcCtx.UserModel.List(l.ctx, offset, size)
	if err != nil {
		logx.Errorf("查询数据库失败: %v", err)
		return nil, err
	}

	logx.Infof("查询结果: users数量=%d, total=%d", len(users), total)

	// 转换返回格式
	list := make([]types.UserItem, 0, len(users))
	for _, u := range users {
		list = append(list, types.UserItem{
			Id:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	result := &types.ListUsersResp{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: size,
	}

	logx.Infof("准备返回: List长度=%d, Total=%d, Page=%d, PageSize=%d",
		len(result.List), result.Total, result.Page, result.PageSize)

	return result, nil
}
