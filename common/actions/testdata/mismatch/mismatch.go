// Package mismatch pairs a model with another model's request, which must
// not compile. TestGenericActionsRefuseAMismatchedPair builds it and expects
// the build to fail; testdata keeps it out of ./... .
package mismatch

import (
	"go-admin/app/demo/models"
	jobdto "go-admin/app/jobs/service/dto"
	"go-admin/common/actions"
)

var _ = actions.Create[models.DemoProduct, jobdto.SysJobControl]
