package unixsignals

import "os"

func ListenUnixCloseSignals(c chan<- os.Signal) {}
