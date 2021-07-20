// +build !windows

package terminal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aggronmagi/promptx/input"
	"github.com/aggronmagi/promptx/internal/debug"
)

func (t *TerminalApp) handleSignals(exitCh chan int, winSizeCh chan *input.WinSize, stop chan struct{}) {
	in := t.in
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)

	for {
		select {
		case <-stop:
			debug.Log("stop handleSignals")
			return
		case s := <-sigCh:
			switch s {
			case syscall.SIGINT: // kill -SIGINT XXXX or Ctrl+c
				debug.Log("Catch SIGINT")
				exitCh <- 0

			case syscall.SIGTERM: // kill -SIGTERM XXXX
				debug.Log("Catch SIGTERM")
				exitCh <- 1

			case syscall.SIGQUIT: // kill -SIGQUIT XXXX
				debug.Log("Catch SIGQUIT")
				exitCh <- 0

			case syscall.SIGWINCH:
				debug.Log("Catch SIGWINCH")
				winSizeCh <- in.GetWinSize()
			}
		}
	}
}
