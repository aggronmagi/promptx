// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package promptx

import (
	"github.com/aggronmagi/promptx/input"
	"github.com/aggronmagi/promptx/output"
)

var _ = PromptOptionsOptionDeclareWithDefault()

type PromptOptions struct {
	// default global input options
	InputOptions []InputOption
	// default global select options
	SelectOptions []SelectOption
	// default common options. use to create default optins
	CommonOpions []CommonOption
	// default manager. if it is not nil, ignore CommonOpions.
	BlocksManager BlocksManager
	// input
	Input input.ConsoleParser
	// output
	Output output.ConsoleWriter
	Stderr output.ConsoleWriter
}

func (cc *PromptOptions) SetOption(opt PromptOption) {
	_ = opt(cc)
}

func (cc *PromptOptions) ApplyOption(opts ...PromptOption) {
	for _, opt := range opts {
		_ = opt(cc)
	}
}

func (cc *PromptOptions) GetSetOption(opt PromptOption) PromptOption {
	return opt(cc)
}

type PromptOption func(cc *PromptOptions) PromptOption

func WithInputOptions(v ...InputOption) PromptOption {
	return func(cc *PromptOptions) PromptOption {
		previous := cc.InputOptions
		cc.InputOptions = v
		return WithInputOptions(previous...)
	}
}

func WithSelectOptions(v ...SelectOption) PromptOption {
	return func(cc *PromptOptions) PromptOption {
		previous := cc.SelectOptions
		cc.SelectOptions = v
		return WithSelectOptions(previous...)
	}
}

func WithCommonOpions(v ...CommonOption) PromptOption {
	return func(cc *PromptOptions) PromptOption {
		previous := cc.CommonOpions
		cc.CommonOpions = v
		return WithCommonOpions(previous...)
	}
}

func WithBlocksManager(v BlocksManager) PromptOption {
	return func(cc *PromptOptions) PromptOption {
		previous := cc.BlocksManager
		cc.BlocksManager = v
		return WithBlocksManager(previous)
	}
}

func WithInput(v input.ConsoleParser) PromptOption {
	return func(cc *PromptOptions) PromptOption {
		previous := cc.Input
		cc.Input = v
		return WithInput(previous)
	}
}

func WithOutput(v output.ConsoleWriter) PromptOption {
	return func(cc *PromptOptions) PromptOption {
		previous := cc.Output
		cc.Output = v
		return WithOutput(previous)
	}
}

func WithStderr(v output.ConsoleWriter) PromptOption {
	return func(cc *PromptOptions) PromptOption {
		previous := cc.Stderr
		cc.Stderr = v
		return WithStderr(previous)
	}
}

func NewPromptOptions(opts ...PromptOption) *PromptOptions {
	cc := newDefaultPromptOptions()
	for _, opt := range opts {
		_ = opt(cc)
	}
	if watchDogPromptOptions != nil {
		watchDogPromptOptions(cc)
	}
	return cc
}

func InstallPromptOptionsWatchDog(dog func(cc *PromptOptions)) {
	watchDogPromptOptions = dog
}

var watchDogPromptOptions func(cc *PromptOptions)

func newDefaultPromptOptions() *PromptOptions {

	cc := &PromptOptions{}

	for _, opt := range [...]PromptOption{
		WithInputOptions(nil...),
		WithSelectOptions(nil...),
		WithCommonOpions(nil...),
		WithBlocksManager(nil),
		WithInput(input.NewStandardInputParser()),
		WithOutput(output.NewStandardOutputWriter()),
		WithStderr(output.NewStderrWriter()),
	} {
		_ = opt(cc)
	}

	return cc
}
