package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

// ErrNotFound 数据未找到错误
// 当查询数据库记录不存在时返回此错误
// 这是对 go-zero sqlx.ErrNotFound 的重新导出，便于业务层统一使用
var ErrNotFound = sqlx.ErrNotFound
