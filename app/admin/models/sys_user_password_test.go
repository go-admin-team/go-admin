package models

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

const knownPassword = "correct-horse-battery-staple"

// A user loaded from the database carries the stored hash in Password, and the
// hooks run on whatever is in the struct. Hashing it a second time produces a
// hash of a hash: the password the user knows stops matching, they cannot log
// in, and nothing reports an error.
//
// Only an Omit("password") on one call site stood between this and every write
// to the model. This is the test that removes the need for it.
func TestEncryptLeavesAnAlreadyHashedPasswordAlone(t *testing.T) {
	fresh := SysUser{Password: knownPassword}
	if err := fresh.Encrypt(); err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	stored := fresh.Password
	if err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(knownPassword)); err != nil {
		t.Fatalf("setup failed: the password was not hashed: %v", err)
	}

	// What a query puts in the struct, and what an update then hands the hook.
	loaded := SysUser{Password: stored}
	if err := loaded.Encrypt(); err != nil {
		t.Fatalf("Encrypt on a loaded user: %v", err)
	}
	if loaded.Password != stored {
		t.Error("Encrypt re-hashed a stored hash; the user can no longer log in")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(loaded.Password), []byte(knownPassword)); err != nil {
		t.Errorf("the user can no longer log in with their password: %v", err)
	}
}

// The other half: a password that is not a hash still gets hashed, on create
// and on update alike.
func TestEncryptHashesAPlaintextPassword(t *testing.T) {
	for _, c := range []struct {
		name string
		hook func(*SysUser) error
	}{
		{"BeforeCreate", func(u *SysUser) error { return u.BeforeCreate(nil) }},
		{"BeforeUpdate", func(u *SysUser) error { return u.BeforeUpdate(nil) }},
	} {
		t.Run(c.name, func(t *testing.T) {
			u := SysUser{Password: knownPassword}
			if err := c.hook(&u); err != nil {
				t.Fatal(err)
			}
			if u.Password == knownPassword {
				t.Fatal("the password was stored as it was typed")
			}
			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(knownPassword)); err != nil {
				t.Errorf("the stored value does not verify the password: %v", err)
			}
		})
	}
}

// An empty Password means "not being set", and must not become a hash of "".
func TestEncryptIgnoresAnEmptyPassword(t *testing.T) {
	u := SysUser{}
	if err := u.Encrypt(); err != nil {
		t.Fatal(err)
	}
	if u.Password != "" {
		t.Errorf("an unset password became %q", u.Password)
	}
}

// Encrypt runs on every update of this model, including the ones that change
// something else entirely. What it costs when there is nothing to do is the
// difference between a profile update and a bcrypt round; the correctness test
// above is what catches a regression, this reports the size of it.
func BenchmarkEncrypt(b *testing.B) {
	fresh := SysUser{Password: knownPassword}
	if err := fresh.Encrypt(); err != nil {
		b.Fatal(err)
	}

	b.Run("already hashed", func(b *testing.B) {
		u := SysUser{Password: fresh.Password}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if err := u.Encrypt(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("plaintext", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			u := SysUser{Password: knownPassword}
			if err := u.Encrypt(); err != nil {
				b.Fatal(err)
			}
		}
	})
}
