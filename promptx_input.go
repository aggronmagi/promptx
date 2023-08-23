package promptx

import (
	"fmt"
	"strconv"
	"strings"

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

// Input get input
func (p *Promptx) Input(tip string, checker InputChecker, defaultValue ...string) (result string, err error) {
	if len(defaultValue) > 0 {
		return p.RawInput(tip, WithInputOptionValidFunc(checker), WithInputOptionDefaultText(defaultValue[0]))
	}
	return p.RawInput(tip, WithInputOptionValidFunc(checker))
}

// Input get input
func (p *Promptx) MultInput(tip string, checker InputChecker) string {
	result, err := p.RawInput(tip, WithInputOptionValidFunc(checker))
	if err != nil {
		panic("user cancel")
	}
	return result
}

func (p *Promptx) InputInt(tip string, val int, check ...func(in int) error) (_ int, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		n, err := strconv.ParseInt(d.Text, 10, 64)
		if err != nil {
			return err
		}
		for _, cf := range check {
			err = cf(int(n))
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	v, _ := strconv.ParseInt(text, 10, 64)
	return int(v), nil
}

func (p *Promptx) MustInputInt(tip string, val int, check ...func(in int) error) int {
	num, eof := p.InputInt(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func convertToIntSlice(txt string, def []int) ([]int, error) {
	var list []int
	for _, v := range strings.Split(txt, ",") {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return def, err
		}
		list = append(list, int(n))
	}
	return list, nil
}

func (p *Promptx) InputIntSlice(tip string, val []int, check ...func(in []int) error) (_ []int, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		list, err := convertToIntSlice(d.Text, val)
		if err != nil {
			return err
		}

		for _, cf := range check {
			err := cf(list)
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	return convertToIntSlice(text, val)
}

func (p *Promptx) MustInputIntSlice(tip string, val []int, check ...func(in []int) error) []int {
	num, eof := p.InputIntSlice(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func (p *Promptx) InputInt64(tip string, val int64, check ...func(in int64) error) (_ int64, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		n, err := strconv.ParseInt(d.Text, 10, 64)
		if err != nil {
			return err
		}
		for _, cf := range check {
			err = cf(n)
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	v, _ := strconv.ParseInt(text, 10, 64)
	return v, nil
}

func (p *Promptx) MustInputInt64(tip string, val int64, check ...func(in int64) error) int64 {
	num, eof := p.InputInt64(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func convertToInt64Slice(txt string, def []int64) ([]int64, error) {
	var list []int64
	for _, v := range strings.Split(txt, ",") {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return def, err
		}
		list = append(list, n)
	}
	return list, nil
}
func (p *Promptx) InputInt64Slice(tip string, val []int64, check ...func(in []int64) error) (_ []int64, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		list, err := convertToInt64Slice(d.Text, val)
		if err != nil {
			return err
		}
		for _, cf := range check {
			err = cf(list)
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	return convertToInt64Slice(text, val)
}

func (p *Promptx) MustInputInt64Slice(tip string, val []int64, check ...func(in []int64) error) []int64 {
	num, eof := p.InputInt64Slice(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func (p *Promptx) InputInt32(tip string, val int32, check ...func(in int32) error) (_ int32, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		n, err := strconv.ParseInt(d.Text, 10, 32)
		if err != nil {
			return err
		}
		for _, cf := range check {
			err = cf(int32(n))
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	v, _ := strconv.ParseInt(text, 10, 32)
	return int32(v), nil
}

func (p *Promptx) MustInputInt32(tip string, val int32, check ...func(in int32) error) int32 {
	num, eof := p.InputInt32(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func convertToInt32Slice(txt string, def []int32) ([]int32, error) {
	var list []int32
	for _, v := range strings.Split(txt, ",") {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return def, err
		}
		list = append(list, int32(n))
	}
	return list, nil
}

func (p *Promptx) InputInt32Slice(tip string, val []int32, check ...func(in []int32) error) (_ []int32, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		list, err := convertToInt32Slice(d.Text, val)
		if err != nil {
			return err
		}

		for _, cf := range check {
			err = cf(list)
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	return convertToInt32Slice(text, val)
}

func (p *Promptx) MustInputInt32Slice(tip string, val []int32, check ...func(in []int32) error) []int32 {
	num, eof := p.InputInt32Slice(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func (p *Promptx) InputString(tip string, val string, check ...func(in string) error) (_ string, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		for _, cf := range check {
			err := cf(d.Text)
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	return text, nil
}

func (p *Promptx) MustInputString(tip string, val string, check ...func(in string) error) string {
	num, eof := p.InputString(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func (p *Promptx) InputStringSlice(tip string, val []string, check ...func(in []string) error) (_ []string, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		list := strings.Split(d.Text, ",")
		for _, cf := range check {
			err := cf(list)
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	return strings.Split(text, ","), nil
}

func (p *Promptx) MustInputStringSlice(tip string, val []string, check ...func(in []string) error) []string {
	num, eof := p.InputStringSlice(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func (p *Promptx) InputFloat32(tip string, val float32, check ...func(in float32) error) (_ float32, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		n, err := strconv.ParseFloat(d.Text, 32)
		if err != nil {
			return err
		}
		for _, cf := range check {
			err = cf(float32(n))
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	v, _ := strconv.ParseFloat(text, 32)
	return float32(v), nil
}

func (p *Promptx) MustInputFloat32(tip string, val float32, check ...func(in float32) error) float32 {
	num, eof := p.InputFloat32(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func convertToFloat32Slice(txt string, def []float32) ([]float32, error) {
	var list []float32
	for _, v := range strings.Split(txt, ",") {
		n, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return def, err
		}
		list = append(list, float32(n))
	}
	return list, nil
}

func (p *Promptx) InputFloat32Slice(tip string, val []float32, check ...func(in []float32) error) (_ []float32, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		list, err := convertToFloat32Slice(d.Text, val)
		if err != nil {
			return err
		}
		for _, cf := range check {
			err = cf(list)
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	return convertToFloat32Slice(text, val)
}

func (p *Promptx) MustInputFloat32Slice(tip string, val []float32, check ...func(in []float32) error) []float32 {
	num, eof := p.InputFloat32Slice(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func (p *Promptx) InputFloat64(tip string, val float64, check ...func(in float64) error) (_ float64, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		n, err := strconv.ParseFloat(d.Text, 64)
		if err != nil {
			return err
		}
		for _, cf := range check {
			err = cf(float64(n))
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	v, _ := strconv.ParseFloat(text, 64)
	return float64(v), nil
}

func (p *Promptx) MustInputFloat64(tip string, val float64, check ...func(in float64) error) float64 {
	num, eof := p.InputFloat64(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func convertToFloat64Slice(txt string, def []float64) ([]float64, error) {
	var list []float64
	for _, v := range strings.Split(txt, ",") {
		n, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return def, err
		}
		list = append(list, n)
	}
	return list, nil
}

func (p *Promptx) InputFloat64Slice(tip string, val []float64, check ...func(in []float64) error) (_ []float64, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValidFunc(func(d *Document) error {
		list, err := convertToFloat64Slice(d.Text, val)
		if err != nil {
			return err
		}
		for _, cf := range check {
			err = cf(list)
			if err != nil {
				return err
			}
		}
		return nil
	}), WithInputOptionDefaultTextAny(val))
	if eof != nil {
		return val, eof
	}
	return convertToFloat64Slice(text, val)
}

func (p *Promptx) MustInputFloat64Slice(tip string, val []float64, check ...func(in []float64) error) []float64 {
	num, eof := p.InputFloat64Slice(tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func WithInputOptionDefaultTextAny(v interface{}) InputOption {
	if v == nil {
		return WithInputOptionDefaultText("")
	}
	in := fmt.Sprint(v)
	in = strings.TrimLeft(strings.TrimRight(in, "]"), "[")
	in = strings.TrimLeft(strings.TrimRight(in, ">"), "<")
	in = strings.Trim(strings.TrimSpace(in), "nil")
	return WithInputOptionDefaultText(in)
}
