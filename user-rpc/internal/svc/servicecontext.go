package svc

import (
	"log"
	"user-demo/user-rpc/internal/config"
	"user-demo/user-rpc/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 使用 GORM 创建 PostgreSQL 连接
	db, err := gorm.Open(postgres.Open(c.Postgres.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(db),
	}
}
