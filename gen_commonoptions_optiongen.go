// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package promptx

var _ = CommonOptionsOptionDeclareWithDefault()

type CommonOptions struct {
	TipText         string
	TipTextColor    Color
	TipBGColor      Color
	PrefixText      string
	PrefixTextColor Color
	PrefixBGColor   Color
	// check input valid
	ValidFunc      func(status int, in *Document) error
	ValidTextColor Color
	ValidBGColor   Color
	// exec input command
	ExecFunc   func(ctx Context, command string)
	FinishKey  Key
	CancelKey  Key
	Completion []CompleteOption
	// if command slice size &gt; 0. it will ignore ExecFunc and ValidFunc options
	Commands []*Cmd
	// alway check input command
	AlwaysCheckCommand bool
}

func (cc *CommonOptions) SetOption(opt CommonOption) {
	_ = opt(cc)
}

func (cc *CommonOptions) ApplyOption(opts ...CommonOption) {
	for _, opt := range opts {
		_ = opt(cc)
	}
}

func (cc *CommonOptions) GetSetOption(opt CommonOption) CommonOption {
	return opt(cc)
}

type CommonOption func(cc *CommonOptions) CommonOption

func WithCommonOptionTipText(v string) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.TipText
		cc.TipText = v
		return WithCommonOptionTipText(previous)
	}
}

func WithCommonOptionTipTextColor(v Color) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.TipTextColor
		cc.TipTextColor = v
		return WithCommonOptionTipTextColor(previous)
	}
}

func WithCommonOptionTipBGColor(v Color) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.TipBGColor
		cc.TipBGColor = v
		return WithCommonOptionTipBGColor(previous)
	}
}

func WithCommonOptionPrefixText(v string) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.PrefixText
		cc.PrefixText = v
		return WithCommonOptionPrefixText(previous)
	}
}

func WithCommonOptionPrefixTextColor(v Color) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.PrefixTextColor
		cc.PrefixTextColor = v
		return WithCommonOptionPrefixTextColor(previous)
	}
}

func WithCommonOptionPrefixBGColor(v Color) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.PrefixBGColor
		cc.PrefixBGColor = v
		return WithCommonOptionPrefixBGColor(previous)
	}
}

func WithCommonOptionValidFunc(v func(status int, in *Document) error) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.ValidFunc
		cc.ValidFunc = v
		return WithCommonOptionValidFunc(previous)
	}
}

func WithCommonOptionValidTextColor(v Color) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.ValidTextColor
		cc.ValidTextColor = v
		return WithCommonOptionValidTextColor(previous)
	}
}

func WithCommonOptionValidBGColor(v Color) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.ValidBGColor
		cc.ValidBGColor = v
		return WithCommonOptionValidBGColor(previous)
	}
}

func WithCommonOptionExecFunc(v func(ctx Context, command string)) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.ExecFunc
		cc.ExecFunc = v
		return WithCommonOptionExecFunc(previous)
	}
}

func WithCommonOptionFinishKey(v Key) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.FinishKey
		cc.FinishKey = v
		return WithCommonOptionFinishKey(previous)
	}
}

func WithCommonOptionCancelKey(v Key) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.CancelKey
		cc.CancelKey = v
		return WithCommonOptionCancelKey(previous)
	}
}

func WithCommonOptionCompletion(v ...CompleteOption) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.Completion
		cc.Completion = v
		return WithCommonOptionCompletion(previous...)
	}
}

func WithCommonOptionCommands(v ...*Cmd) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.Commands
		cc.Commands = v
		return WithCommonOptionCommands(previous...)
	}
}

func WithCommonOptionAlwaysCheckCommand(v bool) CommonOption {
	return func(cc *CommonOptions) CommonOption {
		previous := cc.AlwaysCheckCommand
		cc.AlwaysCheckCommand = v
		return WithCommonOptionAlwaysCheckCommand(previous)
	}
}

func NewCommonOptions(opts ...CommonOption) *CommonOptions {
	cc := newDefaultCommonOptions()
	for _, opt := range opts {
		_ = opt(cc)
	}
	if watchDogCommonOptions != nil {
		watchDogCommonOptions(cc)
	}
	return cc
}

func InstallCommonOptionsWatchDog(dog func(cc *CommonOptions)) {
	watchDogCommonOptions = dog
}

var watchDogCommonOptions func(cc *CommonOptions)

func newDefaultCommonOptions() *CommonOptions {

	cc := &CommonOptions{}

	for _, opt := range [...]CommonOption{
		WithCommonOptionTipText(""),
		WithCommonOptionTipTextColor(Yellow),
		WithCommonOptionTipBGColor(DefaultColor),
		WithCommonOptionPrefixText(">>> "),
		WithCommonOptionPrefixTextColor(Green),
		WithCommonOptionPrefixBGColor(DefaultColor),
		WithCommonOptionValidFunc(nil),
		WithCommonOptionValidTextColor(Red),
		WithCommonOptionValidBGColor(DefaultColor),
		WithCommonOptionExecFunc(nil),
		WithCommonOptionFinishKey(Enter),
		WithCommonOptionCancelKey(ControlC),
		WithCommonOptionCompletion(nil...),
		WithCommonOptionCommands(nil...),
		WithCommonOptionAlwaysCheckCommand(false),
	} {
		_ = opt(cc)
	}

	return cc
}
