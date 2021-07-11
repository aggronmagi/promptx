// +build !windows

package promptx

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aggronmagi/promptx/internal/debug"
	"golang.org/x/term"
)

func HandleSignals(exitCh chan int, winSizeCh chan *WinSize, stop chan struct{}) {
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
				w, h, err := term.GetSize(syscall.Stdout)
				if err != nil {
					debug.Println("get windows size error:", err)
					continue
				}
				size := &WinSize{
					Row: h,
					Col: w,
				}
				winSizeCh <- size
			}
		}
	}
}
