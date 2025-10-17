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

func TestEmailTokenRateLimited(t *testing.T) {
	userID := uint(99)
	window := 200 * time.Millisecond

	limited, _, err := EmailTokenRateLimited(userID, TokenPurposeResetPassword, window)
	if err != nil {
		t.Fatalf("rate limit check before: %v", err)
	}
	if limited {
		t.Fatalf("expected no rate limit before token creation")
	}

	if _, _, err := CreateEmailToken(userID, TokenPurposeResetPassword, time.Minute); err != nil {
		t.Fatalf("create token: %v", err)
	}

	limited, retryAfter, err := EmailTokenRateLimited(userID, TokenPurposeResetPassword, window)
	if err != nil {
		t.Fatalf("rate limit check after: %v", err)
	}
	if !limited {
		t.Fatalf("expected rate limit immediately after token creation")
	}
	if retryAfter <= 0 || retryAfter > window {
		t.Fatalf("unexpected retryAfter: %v", retryAfter)
	}

	time.Sleep(window + time.Millisecond*25)

	limited, _, err = EmailTokenRateLimited(userID, TokenPurposeResetPassword, window)
	if err != nil {
		t.Fatalf("rate limit check post-wait: %v", err)
	}
	if limited {
		t.Fatalf("expected rate limit to expire after wait")
	}
}
