package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/realfatcat/psigo/pkg/psigo"
)

var pid = flag.Int("p", 0, "Decode signals for given PID")
var decode = flag.String("d", "", "Decode this value")

func main() {
	flag.Parse()

	if len(*decode) != 0 {
		toParse := handleHexPrefix(*decode)
		toDecode, err := strconv.ParseUint(toParse, 16, 64)
		if err != nil {
			exit(err)
		}

		for _, sig := range psigo.Decode(toDecode) {
			fmt.Printf("%d) %s ", sig, unix.SignalName(sig))
		}
		fmt.Printf("\n")
	}

	if *pid != 0 {
		signals, err := psigo.FromPID(*pid)
		if err != nil {
			exit(err)
		}

		fmt.Printf("%s\n", signals)
	}
}

func handleHexPrefix(value string) string {
	if strings.HasPrefix(value, "0x") {
		return value[2:]
	}
	return value
}

func exit(err error) {
	fmt.Printf("got error: %s\n", err)
	os.Exit(1)
}
