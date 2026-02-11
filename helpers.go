package promptx

import (
	"errors"
	"strconv"

	"github.com/aggronmagi/promptx/v2/blocks"
	"github.com/aggronmagi/promptx/v2/buffer"
)

func Input(ctx Context, prompt string, validator func(input string) error) (string, error) {
	return ctx.RawInput(prompt, blocks.WithInputOptionValid(func(doc *buffer.Document) error {
		return validator(doc.Text)
	}))
}

func InputInt[T int | int8 | int16 | int32 | int64](ctx Context, prompt string, validator func(input string) error) (T, error) {
	val, err := ctx.RawInput(prompt, blocks.WithInputOptionValid(func(doc *buffer.Document) error {
		return validator(doc.Text)
	}))
	if err != nil {
		return 0, err
	}
	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return T(num), nil
}

func InputUint[T uint | uint8 | uint16 | uint32 | uint64](ctx Context, prompt string, validator func(input string) error) (T, error) {
	val, err := ctx.RawInput(prompt, blocks.WithInputOptionValid(func(doc *buffer.Document) error {
		return validator(doc.Text)
	}))
	if err != nil {
		return 0, err
	}
	num, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return T(num), nil
}

func InputFloat[T float32 | float64](ctx Context, prompt string, validator func(input string) error) (T, error) {
	val, err := ctx.RawInput(prompt, blocks.WithInputOptionValid(func(doc *buffer.Document) error {
		return validator(doc.Text)
	}))
	if err != nil {
		return 0, err
	}
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	return T(num), nil
}

func InputBool(ctx Context, prompt string, validator func(input string) error) (bool, error) {
	val, err := ctx.RawInput(prompt, blocks.WithInputOptionValid(func(doc *buffer.Document) error {
		return validator(doc.Text)
	}))
	if err != nil {
		return false, err
	}
	num, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}
	return num, nil
}

func SelectString(ctx Context, prompt string, options []string) (string, error) {
	idx := ctx.RawSelect(prompt, options)
	if idx < 0 || idx >= len(options) {
		return "", errors.New("user cancel")
	}
	return options[idx], nil
}
