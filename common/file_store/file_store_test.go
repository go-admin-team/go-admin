package file_store

import (
	"errors"
	"os"
	"testing"
)

// The three implementations keep their provider client in an interface{} field
// that Setup assigns, so an unconfigured store holds nil. Asserting nil to the
// provider's type panics, and the upload endpoint reaches that path whenever a
// request names a provider the deployment never configured.
func TestUnconfiguredStoresReportItInsteadOfPanicking(t *testing.T) {
	stores := map[DriverType]FileStoreType{
		AliYunOSS: new(ALiYunOSS),
		HuaweiOBS: new(HuaWeiOBS),
		QiNiuKodo: new(QiNiuKODO),
	}

	for driver, store := range stores {
		t.Run(string(driver), func(t *testing.T) {
			err := store.UpLoad("img/x.png", "/tmp/x.png")
			if err == nil {
				t.Fatal("upload on an unconfigured store returned no error")
			}
			var notCfg *ErrNotConfigured
			if !errors.As(err, &notCfg) {
				t.Fatalf("want ErrNotConfigured, got %v", err)
			}
			if notCfg.Driver != driver {
				t.Errorf("error names %s, want %s", notCfg.Driver, driver)
			}
		})
	}
}

func TestUnconfiguredTokenReportsItToo(t *testing.T) {
	if _, err := new(QiNiuKODO).GetTempToken(); err == nil {
		t.Fatal("token from an unconfigured store returned no error")
	}
}

func TestSetupRejectsAnUnknownDriver(t *testing.T) {
	if _, err := (&OXS{}).Setup("NoSuchCloud"); err == nil {
		t.Fatal("unknown driver was accepted")
	}
}

// Setup reaching the provider needs credentials, so it runs only when they are
// supplied. Previously the test carried a comment telling the reader to paste
// their own, which meant it failed for everyone who did not.
func TestSetupWithRealCredentials(t *testing.T) {
	endpoint := os.Getenv("GOADMIN_OSS_ENDPOINT")
	if endpoint == "" {
		t.Skip("set GOADMIN_OSS_ENDPOINT, GOADMIN_OSS_AK, GOADMIN_OSS_SK, GOADMIN_OSS_BUCKET to run")
	}
	oxs := OXS{
		Endpoint:        endpoint,
		AccessKeyID:     os.Getenv("GOADMIN_OSS_AK"),
		AccessKeySecret: os.Getenv("GOADMIN_OSS_SK"),
		BucketName:      os.Getenv("GOADMIN_OSS_BUCKET"),
	}
	store, err := oxs.Setup(AliYunOSS)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if store == nil {
		t.Fatal("setup returned no store and no error")
	}
}
