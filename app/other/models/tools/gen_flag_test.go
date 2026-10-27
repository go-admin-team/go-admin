package tools

import (
	"encoding/json"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestGenFlagUnmarshal(t *testing.T) {
	cases := map[string]GenFlag{
		`1`: GenFlagYes, `2`: GenFlagNo, `"1"`: GenFlagYes, `"2"`: GenFlagNo,
		`true`: GenFlagUnset, `false`: GenFlagUnset, `null`: GenFlagUnset,
	}
	for in, want := range cases {
		var got struct {
			F GenFlag `json:"f"`
		}
		if err := json.Unmarshal([]byte(`{"f":`+in+`}`), &got); err != nil {
			t.Fatalf("%s: %v", in, err)
		}
		if got.F != want {
			t.Errorf("%s: got %d, want %d", in, got.F, want)
		}
	}
	var bad GenFlag
	if err := json.Unmarshal([]byte(`"yes"`), &bad); err == nil {
		t.Error(`"yes": want an error`)
	}
}

// A v3.4.0 front end sends all three flags as false; the save must keep what is stored.
func TestSysTablesUpdateKeepsFlagsSentAsBooleans(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(new(SysTables), new(SysColumns)); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	stored := SysTables{TBName: "sys_post", IsDataScope: GenFlagYes, IsActions: GenFlagNo, IsAuth: GenFlagYes}
	if err := db.Create(&stored).Error; err != nil {
		t.Fatalf("create: %v", err)
	}

	var sent SysTables
	body := `{"tableId":1,"tableComment":"岗位","isDataScope":false,"isActions":false,"isAuth":false}`
	if err := json.Unmarshal([]byte(body), &sent); err != nil {
		t.Fatalf("unmarshal the v3.4.0 body: %v", err)
	}
	if _, err := sent.Update(db); err != nil {
		t.Fatalf("update: %v", err)
	}

	var got SysTables
	if err := db.First(&got, stored.TableId).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if got.IsDataScope != GenFlagYes || got.IsActions != GenFlagNo || got.IsAuth != GenFlagYes {
		t.Errorf("flags changed to %d/%d/%d", got.IsDataScope, got.IsActions, got.IsAuth)
	}
	if got.TableComment != "岗位" {
		t.Errorf("tableComment not saved: %q", got.TableComment)
	}
}
