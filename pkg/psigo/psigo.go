package psigo

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// Signals represents signal masks for a process.
type Signals map[string][]syscall.Signal

// NewSignals creates a new Signals map with all signal types initialized.
func NewSignals() Signals {
	return Signals{
		SigBlk: nil,
		SigPnd: nil,
		SigIgn: nil,
		SigCgt: nil,
	}
}

// String returns a formatted string representation of the signal masks.
func (s Signals) String() string {
	var b strings.Builder
	for _, k := range OrderedMasks {
		b.Write([]byte(k))
		b.Write([]byte(": "))
		for _, sig := range s[k] {
			b.Write(fmt.Appendf(nil, "%d) %s ", sig, unix.SignalName(sig)))
		}
		b.Write([]byte("\n"))
	}
	return b.String()
}

// Decode converts a signal bitmask to a list of signals.
func Decode(value uint64) []syscall.Signal {
	res := make([]syscall.Signal, 0, len(KnownSignals))

	for i := range 64 { // Check up to 64 bits
		if value&(1<<i) != 0 { // Check if bit i is set
			if sig, exists := KnownSignals[i+1]; exists { // Signal numbers start at 1
				res = append(res, sig)
			}
		}
	}
	return res
}

// FromPID reads signal masks from /proc/[pid]/status for the given process.
func FromPID(pid int) (Signals, error) {
	if pid <= 0 {
		return nil, fmt.Errorf("invalid PID: %d", pid)
	}
	f, err := os.OpenFile(filepath.Join("/proc", strconv.Itoa(pid), "status"), os.O_RDONLY, 0o444)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := NewSignals()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()

		for sig := range maps.Keys(m) {
			if strings.HasPrefix(text, sig) {
				decoded, err := convertAndDecode(text)
				if err != nil {
					return nil, fmt.Errorf("converting and decoding: %w", err)
				}

				m[sig] = decoded
			}
		}
	}
	return m, nil
}

func convertAndDecode(sig string) ([]syscall.Signal, error) {
	split := strings.Split(sig, ":")
	if len(split) != 2 {
		return nil, fmt.Errorf("failed to get signals value from %q", sig)
	}
	sigValue := (strings.TrimSpace(split[1]))

	toDecode, err := strconv.ParseUint(sigValue, 16, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uint %q: %w", sigValue, err)
	}

	return Decode(toDecode), nil
}
