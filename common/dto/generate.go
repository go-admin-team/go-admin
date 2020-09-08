package dto

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ObjectById struct {
	Id  int           `uri:"id" validate:"required"`
	Ids []json.Number `json:"ids"`
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
		ids := make([]int64, 0)
		var i int64
		for _, id := range s.Ids {
			i, _ = id.Int64()
			ids = append(ids, i)
		}
		return ids
	}
	return s.Id
}
