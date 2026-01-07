package promptx

import completion "github.com/aggronmagi/promptx/completion"

// RawSelect get select value with raw option
func (p *Promptx) RawSelect(tip string, list []string, opts ...SelectOption) (result int) {
	// copy new config
	newCC := *p.selectCC
	newCC.ApplyOption(opts...)
	// reset options internal
	newCC.SetOption(WithSelectOptionOnFinish(
		func(sels []int) {
			if len(sels) < 1 {
				result = -1
				return
			}
			result = sels[0]
		},
	))
	newCC.SetOption(WithSelectOptionMulti(false))
	newCC.SetOption(WithSelectOptionTip(tip))

	// if opts set options. use opts values
	if len(newCC.Options) < 1 {
		if len(list) < 1 {
			return -1
		}
		for _, v := range list {
			newCC.Options = append(newCC.Options, &completion.Suggest{
				Text: v,
			})
		}
	}
	sel := NewSelectManager(&newCC)

	sel.SetExecContext(p)
	sel.SetWriter(p.cc.Output)
	sel.UpdateWinSize(p.cc.Input.GetWinSize())
	p.console.Run(sel)
	return
}

// RawMulSel get multiple value with raw option
func (p *Promptx) RawMulSel(tip string, list []string, opts ...SelectOption) (result []int) {
	// copy new config
	newCC := *p.selectCC
	newCC.ApplyOption(opts...)
	newCC.SetOption(WithSelectOptionOnFinish(
		func(sels []int) {
			result = sels
		},
	))
	newCC.SetOption(WithSelectOptionMulti(true))

	// if opts set options. use opts values
	if len(newCC.Options) < 1 {
		if len(list) < 1 {
			return
		}
		for _, v := range list {
			newCC.Options = append(newCC.Options, &completion.Suggest{
				Text: v,
			})
		}
	}
	sel := NewSelectManager(&newCC)
	sel.SetExecContext(p)
	sel.SetWriter(p.cc.Output)
	sel.UpdateWinSize(p.cc.Input.GetWinSize())
	p.console.Run(sel)
	return
}
