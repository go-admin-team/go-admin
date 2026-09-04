package models

import contractmodels "github.com/go-admin-team/go-admin-core/v2/sdk/contract/models"

// ActiveRecord is self-referencing (Generate() ActiveRecord), which is why
// it must stay a type alias rather than a defined type: aliasing preserves
// identity with go-admin-core's sdk/contract/models.ActiveRecord, so a
// model whose Generate() returns that interface still satisfies this one. A
// defined type here would break every implementer's method set - see
// go-admin-core's sdk/contract/models package tests for the counterproof
// (PRD 006 counterproof A).
type ActiveRecord = contractmodels.ActiveRecord
