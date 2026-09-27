package actions

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/v2/response"
	"github.com/go-admin-team/go-admin-core/v2/sdk/api"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"go-admin/common/dto"
)

// The generic CRUD actions. They serve the same five kinds of route as
// IndexAction and its siblings, and answer them the same way; what changes is
// what the caller has to get right.
//
// The older actions take instances at registration and serve every request
// from them, so each model and DTO has to implement a Generate that returns a
// copy - return the receiver and concurrent requests share one struct - and a
// list route takes a func() interface{} whose element type nothing checks.
// Here each request declares its own values of the types it was given, so
// neither convention exists, and a model paired with the wrong DTO does not
// compile.
//
// Three things differ from the older actions, all on paths they got wrong:
//
//   - A request whose database connection cannot be found is answered with
//     500. The older actions logged it and wrote nothing, which reaches the
//     client as an empty 200.
//   - The key is matched against the model's primary-key column as a value.
//     The older actions passed GetId() to Where, which GORM reads as a SQL
//     condition when it is a string.
//   - There is no GenerateM, so no error of its to drop.
//
// The type constraints are unexported: callers never name them, and an
// exported constraint would be a promise that could not be taken back.

// record is a model the actions can create, read, change and delete.
type record[T any] interface {
	*T
	schema.Tabler
	SetCreateBy(int)
	SetUpdateBy(int)
	GetId() any
}

// searchReq is a list request: it binds itself, reports its page, and hands
// back the struct whose search tags MakeCondition turns into WHERE clauses.
type searchReq[S any] interface {
	*S
	Bind(*gin.Context) error
	GetPageIndex() int
	GetPageSize() int
	GetNeedSearch() any
}

// idReq is a request naming rows by key: one from the URI, or several from a
// DELETE body, as dto.ObjectById binds them.
type idReq[I any] interface {
	*I
	Bind(*gin.Context) error
	GetId() any
}

// controlReq is a create or update request: it binds itself and builds the
// model row it describes.
type controlReq[T, C any] interface {
	*C
	Bind(*gin.Context) error
	ToModel() (*T, error)
}

// Index lists T, filtered by S and by the caller's data permission.
func Index[T, S any, PT record[T], PS searchReq[S]]() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		db, ok := genericOrm(c)
		if !ok {
			return
		}
		msgID := pkg.GenerateMsgIDFromContext(c)

		var req S
		if err := PS(&req).Bind(c); err != nil {
			response.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
			return
		}
		var object T
		list := make([]T, 0)
		var count int64
		p := GetPermissionFromContext(c)
		err := db.WithContext(c.Request.Context()).Model(&object).
			Scopes(
				dto.MakeCondition(PS(&req).GetNeedSearch()),
				dto.Paginate(PS(&req).GetPageSize(), PS(&req).GetPageIndex()),
				Permission(PT(&object).TableName(), p),
			).
			Find(&list).Limit(-1).Offset(-1).
			Count(&count).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("MsgID[%s] Index error: %s", msgID, err)
			response.Error(c, 500, err, "查询失败")
			return
		}
		response.PageOK(c, &list, int(count), PS(&req).GetPageIndex(), PS(&req).GetPageSize(), "查询成功")
		c.Next()
	}
}

// View answers one T, named by I's key, if the caller's data permission
// reaches it.
func View[T, I any, PT record[T], PI idReq[I]]() gin.HandlerFunc {
	return ViewAs[T, I, T, PT, PI]()
}

// ViewAs is View answering with R instead of T: the row is looked up as T and
// scanned into R, for a detail route that shows a different shape than the
// model stores.
func ViewAs[T, I, R any, PT record[T], PI idReq[I]]() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		db, ok := genericOrm(c)
		if !ok {
			return
		}
		msgID := pkg.GenerateMsgIDFromContext(c)

		var req I
		if err := PI(&req).Bind(c); err != nil {
			response.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
			return
		}
		var object T
		var rsp R
		p := GetPermissionFromContext(c)
		err := db.Model(&object).WithContext(c.Request.Context()).Scopes(
			Permission(PT(&object).TableName(), p),
		).Where(keyIs(PI(&req).GetId())).First(&rsp).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, nil, "查看对象不存在或无权查看")
			return
		}
		if err != nil {
			log.Errorf("MsgID[%s] View error: %s", msgID, err)
			response.Error(c, 500, err, "查看失败")
			return
		}
		response.OK(c, &rsp, "查询成功")
		c.Next()
	}
}

// Create inserts the T that C builds, recording the caller as its creator.
func Create[T, C any, PT record[T], PC controlReq[T, C]]() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		db, ok := genericOrm(c)
		if !ok {
			return
		}

		var req C
		if err := PC(&req).Bind(c); err != nil {
			response.Error(c, http.StatusUnprocessableEntity, err, err.Error())
			return
		}
		object, err := PC(&req).ToModel()
		if err != nil {
			response.Error(c, 500, err, "模型生成失败")
			return
		}
		PT(object).SetCreateBy(user.GetUserId(c))
		if err = db.WithContext(c.Request.Context()).Create(object).Error; err != nil {
			log.Errorf("Create error: %s", err)
			response.Error(c, 500, err, "创建失败")
			return
		}
		response.OK(c, PT(object).GetId(), "创建成功")
		c.Next()
	}
}

// Update writes the T that C builds over the row with its key, if the
// caller's data permission reaches that row.
func Update[T, C any, PT record[T], PC controlReq[T, C]]() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		db, ok := genericOrm(c)
		if !ok {
			return
		}
		msgID := pkg.GenerateMsgIDFromContext(c)

		var req C
		if err := PC(&req).Bind(c); err != nil {
			response.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
			return
		}
		object, err := PC(&req).ToModel()
		if err != nil {
			response.Error(c, 500, err, "模型生成失败")
			return
		}
		PT(object).SetUpdateBy(user.GetUserId(c))

		p := GetPermissionFromContext(c)
		result := db.WithContext(c.Request.Context()).Scopes(
			Permission(PT(object).TableName(), p),
		).Where(keyIs(PT(object).GetId())).Updates(object)
		if err = result.Error; err != nil {
			log.Errorf("MsgID[%s] Update error: %s", msgID, err)
			response.Error(c, 500, err, "更新失败")
			return
		}
		if result.RowsAffected == 0 {
			response.Error(c, http.StatusForbidden, nil, "无权更新该数据")
			return
		}
		response.OK(c, PT(object).GetId(), "更新成功")
		c.Next()
	}
}

// Delete removes the rows I names, where the caller's data permission
// reaches them.
func Delete[T, I any, PT record[T], PI idReq[I]]() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		db, ok := genericOrm(c)
		if !ok {
			return
		}
		msgID := pkg.GenerateMsgIDFromContext(c)

		var req I
		if err := PI(&req).Bind(c); err != nil {
			log.Errorf("MsgID[%s] Bind error: %s", msgID, err)
			response.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
			return
		}
		var object T
		PT(&object).SetUpdateBy(user.GetUserId(c))

		p := GetPermissionFromContext(c)
		result := db.WithContext(c.Request.Context()).Scopes(
			Permission(PT(&object).TableName(), p),
		).Where(keyIs(PI(&req).GetId())).Delete(&object)
		if err := result.Error; err != nil {
			log.Errorf("MsgID[%s] Delete error: %s", msgID, err)
			response.Error(c, 500, err, "删除失败")
			return
		}
		if result.RowsAffected == 0 {
			response.Error(c, http.StatusForbidden, nil, "无权删除该数据")
			return
		}
		// The key of the empty model, as DeleteAction answers: the rows are
		// named by the request, not by a model.
		response.OK(c, PT(&object).GetId(), "删除成功")
		c.Next()
	}
}

// genericOrm reads the request's database, answering 500 when there is none.
func genericOrm(c *gin.Context) (*gorm.DB, bool) {
	db, err := pkg.GetOrm(c)
	if err != nil {
		api.GetRequestLogger(c).Error(err)
		response.Error(c, 500, err, "数据库连接获取失败")
		return nil, false
	}
	return db, true
}

// keyIs matches the primary-key column against id as a value, or against
// each element when id is a slice or array.
func keyIs(id any) clause.Expression {
	v := reflect.ValueOf(id)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		values := make([]any, v.Len())
		for i := range values {
			values[i] = v.Index(i).Interface()
		}
		return clause.IN{Column: clause.PrimaryColumn, Values: values}
	}
	return clause.Eq{Column: clause.PrimaryColumn, Value: id}
}
