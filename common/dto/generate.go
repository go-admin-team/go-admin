package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ObjectById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
}

func (s *ObjectById) Bind(ctx *gin.Context) error {
	if ctx.Request.Method == http.MethodDelete {
		err := ctx.Bind(s)
		if err != nil {
			return err
		}
		if len(s.Ids) > 0 {
			return nil
		}
	}
	return ctx.BindUri(s)
}

func (s *ObjectById) GetId() interface{} {
	if len(s.Ids) > 0 {
		return s.Ids
	}
	return s.Id
}
