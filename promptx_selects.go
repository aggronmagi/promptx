package promptx

// RawSelect get select value with raw option
func (p *Promptx) RawSelect(tip string, list []string, opts ...SelectOption) (result int) {
	// copy new config
	newCC := *p.selectCC
	newCC.ApplyOption(opts...)
	// reset options internal
	newCC.SetOption(WithSelectOptionFinishFunc(
		func(sels []int) {
			if len(sels) < 1 {
				result = -1
				return
			}
			result = sels[0]
		},
	))
	newCC.SetOption(WithSelectOptionMulti(false))
	newCC.SetOption(WithSelectOptionTipText(tip))

	// if opts set options. use opts values
	if len(newCC.Options) < 1 {
		if len(list) < 1 {
			return -1
		}
		for _, v := range list {
			newCC.Options = append(newCC.Options, &Suggest{
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

// Select get select value
func (p *Promptx) Select(tip string, list []string, defaultSelect ...int) (result int) {
	if len(defaultSelect) > 0 {
		return p.RawSelect(tip, list, WithSelectOptionDefaultSelects(defaultSelect[0]))
	}

	return p.RawSelect(tip, list)
}

// Select get select value
func (p *Promptx) MustSelect(tip string, list []string, defaultSelect ...int) (result int) {
	result = p.Select(tip, list, defaultSelect...)
	if result < 0 {
		panic("user cancel")
	}
	return result
}

// Select get select value
func (p *Promptx) SelectString(tip string, list []string, defaultSelect ...int) (_ string, cancel bool) {
	index := p.Select(tip, list, defaultSelect...)
	if index < 0 {
		cancel = true
		return
	}
	return list[index], false
}

// Select get select value
func (p *Promptx) MustSelectString(tip string, list []string, defaultSelect ...int) string {
	index := p.Select(tip, list, defaultSelect...)
	if index < 0 {
		panic("user cancel")
	}
	return list[index]
}

// RawMulSel get multiple value with raw option
func (p *Promptx) RawMulSel(tip string, list []string, opts ...SelectOption) (result []int) {
	// copy new config
	newCC := *p.selectCC
	newCC.ApplyOption(opts...)
	newCC.SetOption(WithSelectOptionFinishFunc(
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
			newCC.Options = append(newCC.Options, &Suggest{
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

// MulSel get multiple value with raw option
func (p *Promptx) MulSel(tip string, list []string, defaultSelects ...int) (result []int) {
	return p.RawMulSel(tip, list, WithSelectOptionDefaultSelects(defaultSelects...))
}

// MulSel get multiple value with raw option
func (p *Promptx) MustMulSel(tip string, list []string, defaultSelects ...int) (result []int) {
	result = p.RawMulSel(tip, list, WithSelectOptionDefaultSelects(defaultSelects...))
	if len(result) < 1 {
		panic("user cancel")
	}
	return result
}

// MulSel get multiple value with raw option
func (p *Promptx) MulSelString(tip string, list []string, defaultSelects ...int) (result []string) {
	sels := p.RawMulSel(tip, list, WithSelectOptionDefaultSelects(defaultSelects...))
	result = make([]string, 0, len(sels))
	for _, k := range sels {
		result = append(result, list[k])
	}
	return result
}

// MustMulSelString get multiple value with raw option
func (p *Promptx) MustMulSelString(tip string, list []string, defaultSelects ...int) (result []string) {
	sels := p.RawMulSel(tip, list, WithSelectOptionDefaultSelects(defaultSelects...))
	if len(sels) < 1 {
		panic("user cancel")
	}
	result = make([]string, 0, len(sels))
	for _, k := range sels {
		result = append(result, list[k])
	}
	return result
}
