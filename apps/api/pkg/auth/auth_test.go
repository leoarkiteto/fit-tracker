package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write([]byte("password-reset"))
	want := mac.Sum(nil)

	got := DeriveKey("secret", "password-reset")
	if !bytes.Equal(got, want) {
		t.Fatalf("DeriveKey = %x, want %x", got, want)
	}
	if len(got) != sha256.Size {
		t.Fatalf("len = %d, want %d", len(got), sha256.Size)
	}
	if bytes.Equal(got, DeriveKey("secret", "other")) {
		t.Fatal("different purposes must derive different keys")
	}
	if bytes.Equal(got, DeriveKey("other", "password-reset")) {
		t.Fatal("different secrets must derive different keys")
	}
}
