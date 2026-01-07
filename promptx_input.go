package promptx

import (
	"github.com/aggronmagi/promptx/internal/debug"
)

// RawInput get input
func (p *Promptx) RawInput(tip string, opts ...InputOption) (result string, err error) {
	// copy a new config
	newCC := (*p.inputCC)

	// apply input opts
	newCC.ApplyOption(opts...)
	// set internal options
	newCC.SetOption(WithInputOptionFinishFunc(func(input string, eof error) {
		result, err = input, eof
	}))
	if tip != "" {
		newCC.SetOption(WithInputOptionPrefixText(tip))
	}
	//
	input := NewInputManager(&newCC)

	input.SetExecContext(p)
	input.SetWriter(p.cc.Output)
	input.UpdateWinSize(p.cc.Input.GetWinSize())
	debug.Println("enter input")
	p.console.Run(input)
	debug.Println("exit input")
	return
}
