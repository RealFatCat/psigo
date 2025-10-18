package psigo

import (
	"os"
	"strings"
	"syscall"
	"testing"
)

func TestNewSignals(t *testing.T) {
	signals := NewSignals()

	// Check that all expected signal types are present
	expectedKeys := []string{SigBlk, SigPnd, SigIgn, SigCgt}
	for _, key := range expectedKeys {
		if _, exists := signals[key]; !exists {
			t.Errorf("Expected signal type %s not found in NewSignals()", key)
		}
	}

	// Check that all values are initialized to nil
	for _, key := range expectedKeys {
		if signals[key] != nil {
			t.Errorf("Expected signal type %s to be nil, got %v", key, signals[key])
		}
	}

	// Check that no unexpected keys are present
	if len(signals) != len(expectedKeys) {
		t.Errorf("Expected %d signal types, got %d", len(expectedKeys), len(signals))
	}
}

func TestSignalsString(t *testing.T) {
	signals := NewSignals()
	signals[SigPnd] = []syscall.Signal{syscall.SIGHUP}

	result := signals.String()
	if !strings.Contains(result, "SigPnd:") {
		t.Error("Expected result to contain 'SigPnd:'")
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		value    uint64
		expected []syscall.Signal
	}{
		{
			name:     "zero value",
			value:    0,
			expected: []syscall.Signal{},
		},
		{
			name:     "single bit set",
			value:    1,
			expected: []syscall.Signal{syscall.SIGHUP}, // Signal 1
		},
		{
			name:     "multiple bits set",
			value:    3, // Binary: 11, signals 1 and 2
			expected: []syscall.Signal{syscall.SIGHUP, syscall.SIGINT},
		},
		{
			name:     "hex value",
			value:    0x5, // Binary: 101, signals 1 and 3
			expected: []syscall.Signal{syscall.SIGHUP, syscall.SIGQUIT},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Decode(tt.value)

			// Check length
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d signals, got %d", len(tt.expected), len(result))
			}

			// Check that all expected signals are present
			for _, expectedSig := range tt.expected {
				found := false
				for _, resultSig := range result {
					if resultSig == expectedSig {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected signal %v not found in result", expectedSig)
				}
			}
		})
	}
}

func TestConvertAndDecode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expectedLen int
	}{
		{
			name:        "valid input",
			input:       "SigPnd: 0000000000000003",
			expectError: false,
			expectedLen: 2, // Signals 1 and 2
		},
		{
			name:        "valid input with spaces",
			input:       "SigBlk:   0000000000000001  ",
			expectError: false,
			expectedLen: 1, // Signal 1
		},
		{
			name:        "invalid format - no colon",
			input:       "SigPnd 0000000000000003",
			expectError: true,
		},
		{
			name:        "invalid format - multiple colons",
			input:       "SigPnd: 0000000000000003: extra",
			expectError: true,
		},
		{
			name:        "invalid hex value",
			input:       "SigPnd: invalid_hex",
			expectError: true,
		},
		{
			name:        "empty hex value",
			input:       "SigPnd: ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertAndDecode(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
				if len(result) != tt.expectedLen {
					t.Errorf("Expected %d signals, got %d", tt.expectedLen, len(result))
				}
			}
		})
	}
}

func TestFromPID(t *testing.T) {
	// Test invalid PID
	_, err := FromPID(0)
	if err == nil {
		t.Error("Expected error for PID 0, but got none")
	}

	_, err = FromPID(-1)
	if err == nil {
		t.Error("Expected error for negative PID, but got none")
	}

	// Test with non-existent PID
	_, err = FromPID(999999)
	if err == nil {
		t.Error("Expected error for non-existent PID, but got none")
	}

	// Test with current process PID (should work if /proc/self/status exists)
	pid := os.Getpid()
	signals, err := FromPID(pid)
	if err != nil {
		t.Logf("Could not test with PID %d: %v", pid, err)
		return
	}

	// Verify that we got a valid Signals map
	if signals == nil {
		t.Error("Expected non-nil signals map")
	}

	// Check that all expected keys are present
	expectedKeys := []string{SigBlk, SigPnd, SigIgn, SigCgt}
	for _, key := range expectedKeys {
		if _, exists := signals[key]; !exists {
			t.Errorf("Expected signal type %s not found", key)
		}
	}
}

// Benchmark tests
func BenchmarkDecode(b *testing.B) {
	value := uint64(0xFFFFFFFFFFFFFFFF)
	for b.Loop() {
		Decode(value)
	}
}

func BenchmarkNewSignals(b *testing.B) {
	for b.Loop() {
		NewSignals()
	}
}

func BenchmarkSignalsString(b *testing.B) {
	signals := NewSignals()
	signals[SigPnd] = []syscall.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT}
	signals[SigIgn] = []syscall.Signal{syscall.SIGTERM}

	for b.Loop() {
		_ = signals.String()
	}
}
