// +build !windows

package unixsignals

import (
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

func ListenUnixCloseSignals(c chan<- os.Signal) {
	signal.Notify(c, unix.SIGTERM)
}
