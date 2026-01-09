package database

import "testing"

func TestConectarMongoDB_InvalidURI(t *testing.T) {
	_, err := ConectarMongoDB("invalid://uri", "db")
	if err == nil {
		t.Fatalf("expected error on invalid mongo uri")
	}
}
