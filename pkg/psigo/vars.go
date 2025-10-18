package psigo

import (
	"syscall"

	"golang.org/x/sys/unix"
)

const maxSignal = 255

func init() {
	for i := range syscall.Signal(maxSignal) {
		sn := unix.SignalName(i)
		if sn != "" {
			KnownSignals[int(i)] = i
		}
	}
}

// KnownSignals maps signal numbers to their corresponding signal values.
var KnownSignals = map[int]syscall.Signal{}

var (
	// SigBlk represents blocked signals.
	SigBlk = "SigBlk"
	// SigPnd represents pending signals.
	SigPnd = "SigPnd"
	// SigIgn represents ignored signals.
	SigIgn = "SigIgn"
	// SigCgt represents caught signals.
	SigCgt = "SigCgt"
)

// OrderedMasks defines the display order for signal mask types.
var OrderedMasks = []string{SigPnd, SigBlk, SigIgn, SigCgt}
