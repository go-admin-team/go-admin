package database

import (
	"strings"
	"testing"
)

func TestOpenerForRejectsADriverThisBuildDoesNotCarry(t *testing.T) {
	open, err := openerFor("sqlite")
	if err == nil {
		t.Fatalf("misspelled driver was accepted, opener is %v", open != nil)
	}
	// The message has to name the alternatives: the reader is looking at a
	// config file and needs to know what to put there.
	for _, want := range supportedDrivers() {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not mention the %s driver: %s", want, err)
		}
	}
}

func TestOpenerForResolvesTheBuiltInDrivers(t *testing.T) {
	for _, driver := range []string{"mysql", "postgres", "sqlserver"} {
		open, err := openerFor(driver)
		if err != nil {
			t.Fatalf("%s must be available in every build: %s", driver, err)
		}
		if open == nil {
			t.Fatalf("%s resolved to a nil opener", driver)
		}
	}
}
