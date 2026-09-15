package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/runtime"
)

// freshRuntime hands the test its own Runtime and puts the old one back, the
// same pattern cmd/api/server_test.go uses: sdk.Runtime is a process-wide
// singleton, and a test that registers into it would otherwise leak state
// into every other test in the binary.
func freshRuntime(t *testing.T) {
	t.Helper()
	previous := sdk.Runtime
	t.Cleanup(func() { sdk.Runtime = previous })
	sdk.Runtime = runtime.NewConfig()
}

// TestInitMiddlewareRegistersUsableJwtHandlerFunc is the reverse proof for
// hoisting the JWT instance's construction into InitMiddleware:
// sdk.Runtime.GetHandlerFunc(JwtTokenCheck) must hand back ok=true and a
// non-nil gin.HandlerFunc, not just something GetMiddleware can return as an
// untyped interface{}.
//
// Before this change, InitMiddleware registered the unbound method
// expression (*jwt.GinJWTMiddleware).MiddlewareFunc under this key - a value
// with no receiver bound to it, which is not a gin.HandlerFunc no matter how
// a caller asserts its type. Reverting the registration below to that
// expression makes GetHandlerFunc report ok=false; it does not fail to
// compile, because (*jwt.GinJWTMiddleware).MiddlewareFunc has a well-formed,
// unrelated method-expression type that SetMiddleware's interface{} param
// happily accepts.
func TestInitMiddlewareRegistersUsableJwtHandlerFunc(t *testing.T) {
	freshRuntime(t)

	previousSecret := config.JwtConfig.Secret
	config.JwtConfig.Secret = "test-secret-key"
	t.Cleanup(func() { config.JwtConfig.Secret = previousSecret })

	gin.SetMode(gin.TestMode)
	InitMiddleware(gin.New())

	h, ok := sdk.Runtime.GetHandlerFunc(JwtTokenCheck)
	if !ok {
		t.Fatal("GetHandlerFunc(JwtTokenCheck) reported ok=false after InitMiddleware ran")
	}
	if h == nil {
		t.Fatal("GetHandlerFunc(JwtTokenCheck) reported ok=true but returned a nil handler")
	}
}

// TestInitMiddlewareBuildsOneSharedJwtInstance locks down the fix for the
// four-instances problem: GetAuthMiddleware must return the very instance
// InitMiddleware built and handed to sdk.Runtime, not a lookalike built
// separately by whichever caller asks first.
func TestInitMiddlewareBuildsOneSharedJwtInstance(t *testing.T) {
	freshRuntime(t)

	previousSecret := config.JwtConfig.Secret
	config.JwtConfig.Secret = "test-secret-key"
	t.Cleanup(func() { config.JwtConfig.Secret = previousSecret })

	gin.SetMode(gin.TestMode)
	InitMiddleware(gin.New())

	shared := GetAuthMiddleware()
	if shared == nil {
		t.Fatal("GetAuthMiddleware returned nil after InitMiddleware ran")
	}
	if shared != authMiddleware {
		t.Error("GetAuthMiddleware did not return the package-level instance InitMiddleware built")
	}
}
