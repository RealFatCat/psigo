package psigo

import (
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

func TestKnownSignalsInitialization(t *testing.T) {
	// Test that KnownSignals is properly initialized
	if len(KnownSignals) == 0 {
		t.Error("KnownSignals map should not be empty")
	}
	
	// Test that common signals are present
	commonSignals := []syscall.Signal{
		syscall.SIGHUP,  // 1
		syscall.SIGINT,  // 2
		syscall.SIGQUIT, // 3
		syscall.SIGTERM, // 15
		syscall.SIGKILL, // 9
	}
	
	for _, sig := range commonSignals {
		if _, exists := KnownSignals[int(sig)]; !exists {
			t.Errorf("Expected signal %d (%s) to be in KnownSignals", sig, unix.SignalName(sig))
		}
	}
}

func TestSignalConstants(t *testing.T) {
	// Test that signal constants are properly defined
	expectedConstants := map[string]string{
		"SigBlk": SigBlk,
		"SigPnd": SigPnd,
		"SigIgn": SigIgn,
		"SigCgt": SigCgt,
	}
	
	for name, value := range expectedConstants {
		if value == "" {
			t.Errorf("Constant %s should not be empty", name)
		}
		if value != name {
			t.Errorf("Expected %s to equal %s, got %s", name, name, value)
		}
	}
}

func TestOrderedMasks(t *testing.T) {
	// Test that OrderedMasks contains all expected signal types
	expectedMasks := []string{SigPnd, SigBlk, SigIgn, SigCgt}
	
	if len(OrderedMasks) != len(expectedMasks) {
		t.Errorf("Expected OrderedMasks to have %d elements, got %d", len(expectedMasks), len(OrderedMasks))
	}
	
	// Test that all expected masks are present
	for _, expected := range expectedMasks {
		found := false
		for _, actual := range OrderedMasks {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected mask %s not found in OrderedMasks", expected)
		}
	}
	
	// Test that the order is correct
	expectedOrder := []string{SigPnd, SigBlk, SigIgn, SigCgt}
	for i, expected := range expectedOrder {
		if OrderedMasks[i] != expected {
			t.Errorf("Expected OrderedMasks[%d] to be %s, got %s", i, expected, OrderedMasks[i])
		}
	}
}

func TestMaxSignalConstant(t *testing.T) {
	// Test that maxSignal is reasonable
	if maxSignal <= 0 {
		t.Error("maxSignal should be positive")
	}
	
	if maxSignal > 255 {
		t.Error("maxSignal should not exceed 255")
	}
	
	// Test that we can iterate up to maxSignal
	count := 0
	for i := range syscall.Signal(maxSignal) {
		if unix.SignalName(i) != "" {
			count++
		}
	}
	
	if count == 0 {
		t.Error("Should have found at least one valid signal")
	}
	
	// KnownSignals should have the same count as our iteration
	if len(KnownSignals) != count {
		t.Errorf("KnownSignals length (%d) should match iteration count (%d)", len(KnownSignals), count)
	}
}

func TestSignalNameConsistency(t *testing.T) {
	// Test that all signals in KnownSignals have valid names
	for signalNum, signal := range KnownSignals {
		name := unix.SignalName(signal)
		if name == "" {
			t.Errorf("Signal %d should have a name", signalNum)
		}
		
		// Test that the signal number matches
		if int(signal) != signalNum {
			t.Errorf("Signal number mismatch: expected %d, got %d", signalNum, int(signal))
		}
	}
}

// Test that the init function works correctly
func TestInitFunction(t *testing.T) {
	// This test verifies that the init() function in vars.go works correctly
	// by checking that KnownSignals is populated
	
	if len(KnownSignals) == 0 {
		t.Error("KnownSignals should be populated by init() function")
	}
	
	// Test that we have some reasonable number of signals
	// On Linux, there are typically around 30-40 standard signals
	if len(KnownSignals) < 10 {
		t.Errorf("Expected at least 10 signals, got %d", len(KnownSignals))
	}
	
	if len(KnownSignals) > 100 {
		t.Errorf("Expected at most 100 signals, got %d", len(KnownSignals))
	}
}

// Benchmark the signal name lookup
func BenchmarkSignalNameLookup(b *testing.B) {
	signal := syscall.SIGHUP
	for i := 0; i < b.N; i++ {
		_ = unix.SignalName(signal)
	}
}

// Benchmark the KnownSignals map lookup
func BenchmarkKnownSignalsLookup(b *testing.B) {
	signalNum := 1
	for i := 0; i < b.N; i++ {
		_ = KnownSignals[signalNum]
	}
}
