package spice

import (
	"testing"
)

func TestEmptyTable(t *testing.T) {
	setupTable(t)

	if table.Lookup(1) != false {
		t.Fatal("lookup should not return for this key")
	}
	if table.OTP(1) != "" {
		t.Fatal("lookup should not return for this key")
	}
	if _, err := table.Connect(1); err == nil {
		t.Fatal("lookup should not return for this key")
	}
}

func TestTableAdd(t *testing.T) {
	setupTable(t)

	table.Add(1, "dst1", "otp1")

	if table.OTP(1) != "otp1" {
		t.Fatal("lookup should have returned this key")
	}

	var dst string
	var err error
	if dst, err = table.Connect(1); err != nil {
		t.Fatal("lookup should have returned this key")
	}

	if dst != "dst1" {
		t.Fatal("lookup should have returned this key")
	}
}

func TestTableConnect(t *testing.T) {
	setupTable(t)
	table.Add(1, "dst1", "otp1")
	table.Connect(1)
	table.Connect(1)
	if table.entries[1].usageCount != 3 {
		t.Fatal("should have 3 connections")
	}
	table.Disconnect(1)
	table.Disconnect(1)
	if table.entries[1].usageCount != 1 {
		t.Fatal("should have 1 connection")
	}
	table.Disconnect(1)
	if _, ok := table.entries[1]; ok {
		t.Fatal("should not have a connection")
	}
}

var table *sessionTable

func setupTable(t *testing.T) {
	t.Helper()
	table = newSessionTable()
}
