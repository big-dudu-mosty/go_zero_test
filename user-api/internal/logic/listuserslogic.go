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

	// 调用 RPC 获取用户列表
	rpcResp, err := l.svcCtx.UserRpc.ListUsers(l.ctx, &user.ListUsersReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		logx.Errorf("RPC 查询用户列表失败: %v", err)
		return nil, err
	}

	logx.Infof("RPC 返回结果: users数量=%d, total=%d", len(rpcResp.List), rpcResp.Total)

	// 转换返回格式
	list := make([]types.UserItem, 0, len(rpcResp.List))
	for _, u := range rpcResp.List {
		list = append(list, types.UserItem{
			Id:    u.Id,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	result := &types.ListUsersResp{
		List:     list,
		Total:    rpcResp.Total,
		Page:     rpcResp.Page,
		PageSize: rpcResp.PageSize,
	}

	logx.Infof("准备返回: List长度=%d, Total=%d, Page=%d, PageSize=%d",
		len(result.List), result.Total, result.Page, result.PageSize)

	return result, nil
}
