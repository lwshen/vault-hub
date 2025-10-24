package commands

import (
	"testing"
)

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "string shorter than maxLen",
			input:    "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "string equal to maxLen",
			input:    "hello",
			maxLen:   5,
			expected: "hello",
		},
		{
			name:     "string longer than maxLen",
			input:    "hello world",
			maxLen:   5,
			expected: "hello...",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   10,
			expected: "",
		},
		{
			name:     "maxLen is 0 with non-empty string",
			input:    "hello",
			maxLen:   0,
			expected: "...",
		},
		{
			name:     "maxLen is 0 with empty string",
			input:    "",
			maxLen:   0,
			expected: "",
		},
		{
			name:     "unicode characters (emojis)",
			input:    "Hello 🌍🚀✨",
			maxLen:   8,
			expected: "Hello 🌍🚀...",
		},
		{
			name:     "unicode characters shorter than maxLen",
			input:    "Hello 🌍",
			maxLen:   10,
			expected: "Hello 🌍",
		},
		{
			name:     "unicode characters (Chinese)",
			input:    "你好世界",
			maxLen:   2,
			expected: "你好...",
		},
		{
			name:     "multibyte characters (Japanese)",
			input:    "こんにちは世界",
			maxLen:   5,
			expected: "こんにちは...",
		},
		{
			name:     "mixed ASCII and unicode",
			input:    "Test 测试 тест",
			maxLen:   7,
			expected: "Test 测试...",
		},
		{
			name:     "very long string",
			input:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
			maxLen:   20,
			expected: "abcdefghijklmnopqrst...",
		},
		{
			name:     "single character with maxLen 1",
			input:    "a",
			maxLen:   1,
			expected: "a",
		},
		{
			name:     "single character with maxLen 0",
			input:    "a",
			maxLen:   0,
			expected: "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateString(%q, %d) = %q; want %q",
					tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

// TestTruncateStringLength verifies that truncated strings don't exceed expected length
func TestTruncateStringLength(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		maxLen       int
		maxRuneCount int // Maximum rune count expected (maxLen + 3 for "...")
	}{
		{
			name:         "truncated ASCII string",
			input:        "this is a very long string",
			maxLen:       5,
			maxRuneCount: 8, // 5 chars + "..."
		},
		{
			name:         "truncated unicode string",
			input:        "这是一个很长的字符串",
			maxLen:       3,
			maxRuneCount: 6, // 3 chars + "..."
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			runeCount := len([]rune(result))

			if runeCount > tt.maxRuneCount {
				t.Errorf("truncateString(%q, %d) produced %d runes; max expected %d (result: %q)",
					tt.input, tt.maxLen, runeCount, tt.maxRuneCount, result)
			}
		})
	}
}
