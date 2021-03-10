package service

import (
	"github.com/go-admin-team/go-admin-core/sdk/pkg/logger"
	"gorm.io/gorm"
)

type Service struct {
	Orm   *gorm.DB
	Msg   string
	MsgID string
	Log   *logger.Logger
}
