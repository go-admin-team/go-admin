package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	mycasbin "github.com/go-admin-team/go-admin-core/v2/casbin"
	"github.com/go-admin-team/go-admin-core/v2/jwtauth"
	"github.com/go-admin-team/go-admin-core/v2/response"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/api"
)

// AuthCheckRole 权限检查中间件
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		data, _ := c.Get(jwtauth.JwtPayloadKey)
		v := data.(jwtauth.MapClaims)
		e := sdk.Runtime.GetCasbinByTenant(c.Request.Host)
		var res, casbinExclude bool
		var err error
		//检查权限
		if v["rolekey"] == "admin" {
			res = true
			c.Next()
			return
		}
		casbinExclude, err = excludedFromCasbin(c.Request.Method, c.Request.URL.Path)
		if err != nil {
			log.Errorf("AuthCheckRole: %s", err)
		}
		if casbinExclude {
			log.Infof("Casbin exclusion, no validation method:%s path:%s", c.Request.Method, c.Request.URL.Path)
			c.Next()
			return
		}
		res, err = e.Enforce(v["rolekey"], c.Request.URL.Path, c.Request.Method)
		if err != nil {
			log.Errorf("AuthCheckRole error:%s method:%s path:%s", err, c.Request.Method, c.Request.URL.Path)
			response.Error(c, 500, err, "")
			return
		}

		if res {
			log.Infof("isTrue: %v role: %s method: %s path: %s", res, v["rolekey"], c.Request.Method, c.Request.URL.Path)
			c.Next()
		} else {
			log.Warnf("isTrue: %v role: %s method: %s path: %s message: %s", res, v["rolekey"], c.Request.Method, c.Request.URL.Path, "当前request无权限，请管理员确认！")
			c.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  "对不起，您没有该接口访问权限，请联系管理员",
			})
			c.Abort()
			return
		}

	}
}

// excludedFromCasbin reports whether the route skips the permission check.
//
// It runs for every non-admin request, so the order matters: the method rules
// out most entries with a string compare, where the path test costs a pattern
// match. mycasbin.KeyMatch2 answers what casbin's util.KeyMatch2 answers
// without recompiling the pattern every time, which is what made this loop
// expensive - about 2,500 allocations per request against a 32-entry list.
//
// A pattern that will not compile is a bug in CasbinExclude rather than in the
// request, so the entry is skipped and the scan continues; the error comes
// back for the caller to log.
func excludedFromCasbin(method, path string) (bool, error) {
	var bad error
	for _, i := range CasbinExclude {
		if method != i.Method {
			continue
		}
		ok, err := mycasbin.KeyMatch2(path, i.Url)
		if err != nil {
			bad = fmt.Errorf("CasbinExclude entry %q is not a valid pattern: %w", i.Url, err)
			continue
		}
		if ok {
			return true, bad
		}
	}
	return false, bad
}
