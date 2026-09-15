package models

import contractmodels "github.com/go-admin-team/go-admin-core/v2/sdk/contract/models"

// Migration is the sys_migration row model (data). It is unrelated to
// cmd/migrate/migration.Migration, the in-process registration table this
// package's TableName has nothing to do with - see
// go-admin-core's sdk/contract/models.Migration doc comment for why the two
// share a name.
type Migration = contractmodels.Migration
