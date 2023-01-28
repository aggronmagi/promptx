// Code generated by "gogen option"; DO NOT EDIT.
// Exec: "gogen option -n InputOption -f -o gen_options_input.go"
// Version: 0.0.3

package promptx

var _ = promptxInputOptions()

// InputOptions promptx options
// generate by https://github.com/aggronmagi/gogen/
type InputOptions struct {
	TipText         string
	TipTextColor    Color
	TipBGColor      Color
	PrefixText      string
	PrefixTextColor Color
	PrefixBGColor   Color
	ValidFunc       func(*Document) error
	ValidTextColor  Color
	ValidBGColor    Color
	FinishFunc      func(input string, eof error)
	FinishKey       Key
	CancelKey       Key
	// result display
	ResultText      InputFinishTextFunc
	ResultTextColor Color
	ResultBGColor   Color
}

func WithInputOptionTipText(v string) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.TipText
		cc.TipText = v
		return WithInputOptionTipText(previous)
	}
}
func WithInputOptionTipTextColor(v Color) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.TipTextColor
		cc.TipTextColor = v
		return WithInputOptionTipTextColor(previous)
	}
}
func WithInputOptionTipBGColor(v Color) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.TipBGColor
		cc.TipBGColor = v
		return WithInputOptionTipBGColor(previous)
	}
}
func WithInputOptionPrefixText(v string) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.PrefixText
		cc.PrefixText = v
		return WithInputOptionPrefixText(previous)
	}
}
func WithInputOptionPrefixTextColor(v Color) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.PrefixTextColor
		cc.PrefixTextColor = v
		return WithInputOptionPrefixTextColor(previous)
	}
}
func WithInputOptionPrefixBGColor(v Color) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.PrefixBGColor
		cc.PrefixBGColor = v
		return WithInputOptionPrefixBGColor(previous)
	}
}
func WithInputOptionValidFunc(v func(*Document) error) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.ValidFunc
		cc.ValidFunc = v
		return WithInputOptionValidFunc(previous)
	}
}
func WithInputOptionValidTextColor(v Color) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.ValidTextColor
		cc.ValidTextColor = v
		return WithInputOptionValidTextColor(previous)
	}
}
func WithInputOptionValidBGColor(v Color) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.ValidBGColor
		cc.ValidBGColor = v
		return WithInputOptionValidBGColor(previous)
	}
}
func WithInputOptionFinishFunc(v func(input string, eof error)) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.FinishFunc
		cc.FinishFunc = v
		return WithInputOptionFinishFunc(previous)
	}
}
func WithInputOptionFinishKey(v Key) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.FinishKey
		cc.FinishKey = v
		return WithInputOptionFinishKey(previous)
	}
}
func WithInputOptionCancelKey(v Key) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.CancelKey
		cc.CancelKey = v
		return WithInputOptionCancelKey(previous)
	}
}

// result display
func WithInputOptionResultText(v InputFinishTextFunc) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.ResultText
		cc.ResultText = v
		return WithInputOptionResultText(previous)
	}
}
func WithInputOptionResultTextColor(v Color) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.ResultTextColor
		cc.ResultTextColor = v
		return WithInputOptionResultTextColor(previous)
	}
}
func WithInputOptionResultBGColor(v Color) InputOption {
	return func(cc *InputOptions) InputOption {
		previous := cc.ResultBGColor
		cc.ResultBGColor = v
		return WithInputOptionResultBGColor(previous)
	}
}

// SetOption modify options
func (cc *InputOptions) SetOption(opt InputOption) {
	_ = opt(cc)
}

// ApplyOption modify options
func (cc *InputOptions) ApplyOption(opts ...InputOption) {
	for _, opt := range opts {
		_ = opt(cc)
	}
}

// GetSetOption modify and get last option
func (cc *InputOptions) GetSetOption(opt InputOption) InputOption {
	return opt(cc)
}

// InputOption option define
type InputOption func(cc *InputOptions) InputOption

// NewInputOptions create options instance.
func NewInputOptions(opts ...InputOption) *InputOptions {
	cc := newDefaultInputOptions()
	for _, opt := range opts {
		_ = opt(cc)
	}
	if watchDogInputOptions != nil {
		watchDogInputOptions(cc)
	}
	return cc
}

// InstallInputOptionsWatchDog install watch dog
func InstallInputOptionsWatchDog(dog func(cc *InputOptions)) {
	watchDogInputOptions = dog
}

var watchDogInputOptions func(cc *InputOptions)

// newDefaultInputOptions new option with default value
func newDefaultInputOptions() *InputOptions {
	cc := &InputOptions{
		TipText:         "",
		TipTextColor:    Yellow,
		TipBGColor:      DefaultColor,
		PrefixText:      ">> ",
		PrefixTextColor: Green,
		PrefixBGColor:   DefaultColor,
		ValidFunc:       nil,
		ValidTextColor:  Red,
		ValidBGColor:    DefaultColor,
		FinishFunc:      nil,
		FinishKey:       Enter,
		CancelKey:       ControlC,
		ResultText:      defaultInputFinishText,
		ResultTextColor: Blue,
		ResultBGColor:   DefaultColor,
	}
	return cc
}
