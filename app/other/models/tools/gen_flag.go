package tools

import (
	"bytes"
	"encoding/json"
	"strconv"
)

const (
	// GenFlagUnset is the zero value; Updates skips it, so the stored flag stays.
	GenFlagUnset GenFlag = 0
	GenFlagYes   GenFlag = 1
	GenFlagNo    GenFlag = 2
)

// GenFlag is a generator yes/no setting stored as 1 (yes) or 2 (no).
// Front ends up to v3.4.0 always send false here, so booleans are ignored, not converted.
type GenFlag int

// UnmarshalJSON accepts a number or numeric string; a boolean or null leaves the flag unset.
func (f *GenFlag) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	var b bool
	if bytes.Equal(data, []byte("null")) || json.Unmarshal(data, &b) == nil {
		*f = GenFlagUnset
		return nil
	}
	n, err := strconv.Atoi(string(bytes.Trim(data, `"`)))
	if err != nil {
		return err
	}
	*f = GenFlag(n)
	return nil
}
