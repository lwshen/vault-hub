package model

import (
	"testing"
	"time"
)

func TestEmailTokenLifecycle(t *testing.T) {
	// DB should already be initialized in test setup; here we only validate helpers
	token, saved, err := CreateEmailToken(1, TokenPurposeMagicLink, time.Millisecond*250)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == 0 {
		t.Fatalf("expected saved token")
	}

	got, err := VerifyAndConsumeEmailToken(token, TokenPurposeMagicLink)
	if err != nil {
		t.Fatalf("verify consume: %v", err)
	}
	if got.ID != saved.ID {
		t.Fatalf("id mismatch")
	}

	// second consume should fail
	if _, err := VerifyAndConsumeEmailToken(token, TokenPurposeMagicLink); err == nil {
		t.Fatalf("expected error on second consume")
	}

	// expired token
	token2, _, err := CreateEmailToken(1, TokenPurposeMagicLink, time.Millisecond*10)
	if err != nil {
		t.Fatalf("create2: %v", err)
	}
	time.Sleep(time.Millisecond * 15)
	if _, err := VerifyAndConsumeEmailToken(token2, TokenPurposeMagicLink); err == nil {
		t.Fatalf("expected expired error")
	}
}
