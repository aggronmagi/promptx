package promptx

import (
	"fmt"
	"strconv"
	"strings"
)

func InputInt(p Interaction, tip string, val int, check ...func(in int) error) (_ int, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputInt(p Interaction, tip string, val int, check ...func(in int) error) int {
	num, eof := InputInt(p, tip, val, check...)
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

func InputIntSlice(p Interaction, tip string, val []int, check ...func(in []int) error) (_ []int, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputIntSlice(p Interaction, tip string, val []int, check ...func(in []int) error) []int {
	num, eof := InputIntSlice(p, tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func InputInt64(p Interaction, tip string, val int64, check ...func(in int64) error) (_ int64, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputInt64(p Interaction, tip string, val int64, check ...func(in int64) error) int64 {
	num, eof := InputInt64(p, tip, val, check...)
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

func InputInt64Slice(p Interaction, tip string, val []int64, check ...func(in []int64) error) (_ []int64, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputInt64Slice(p Interaction, tip string, val []int64, check ...func(in []int64) error) []int64 {
	num, eof := InputInt64Slice(p, tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func InputInt32(p Interaction, tip string, val int32, check ...func(in int32) error) (_ int32, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputInt32(p Interaction, tip string, val int32, check ...func(in int32) error) int32 {
	num, eof := InputInt32(p, tip, val, check...)
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

func InputInt32Slice(p Interaction, tip string, val []int32, check ...func(in []int32) error) (_ []int32, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputInt32Slice(p Interaction, tip string, val []int32, check ...func(in []int32) error) []int32 {
	num, eof := InputInt32Slice(p, tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func InputString(p Interaction, tip string, val string, check ...func(in string) error) (_ string, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputString(p Interaction, tip string, val string, check ...func(in string) error) string {
	num, eof := InputString(p, tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func InputStringSlice(p Interaction, tip string, val []string, check ...func(in []string) error) (_ []string, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputStringSlice(p Interaction, tip string, val []string, check ...func(in []string) error) []string {
	num, eof := InputStringSlice(p, tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func InputFloat32(p Interaction, tip string, val float32, check ...func(in float32) error) (_ float32, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputFloat32(p Interaction, tip string, val float32, check ...func(in float32) error) float32 {
	num, eof := InputFloat32(p, tip, val, check...)
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

func InputFloat32Slice(p Interaction, tip string, val []float32, check ...func(in []float32) error) (_ []float32, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputFloat32Slice(p Interaction, tip string, val []float32, check ...func(in []float32) error) []float32 {
	num, eof := InputFloat32Slice(p, tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func InputFloat64(p Interaction, tip string, val float64, check ...func(in float64) error) (_ float64, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputFloat64(p Interaction, tip string, val float64, check ...func(in float64) error) float64 {
	num, eof := InputFloat64(p, tip, val, check...)
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

func InputFloat64Slice(p Interaction, tip string, val []float64, check ...func(in []float64) error) (_ []float64, eof error) {
	text, eof := p.RawInput(tip, WithInputOptionValid(func(d *Document) error {
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

func MustInputFloat64Slice(p Interaction, tip string, val []float64, check ...func(in []float64) error) []float64 {
	num, eof := InputFloat64Slice(p, tip, val, check...)
	if eof != nil {
		panic("user cancel")
	}
	return num
}

func WithInputOptionDefaultTextAny(v interface{}) InputOption {
	if v == nil {
		return WithInputOptionDefault("")
	}
	in := fmt.Sprint(v)
	in = strings.TrimLeft(strings.TrimRight(in, "]"), "[")
	in = strings.TrimLeft(strings.TrimRight(in, ">"), "<")
	in = strings.Trim(strings.TrimSpace(in), "nil")
	return WithInputOptionDefault(in)
}

// Input get input
func Input(p Interaction, tip string, checker InputChecker, defaultValue ...string) (result string, err error) {
	if len(defaultValue) > 0 {
		return p.RawInput(tip, WithInputOptionValid(checker), WithInputOptionDefault(defaultValue[0]))
	}
	return p.RawInput(tip, WithInputOptionValid(checker))
}

// MustInput get input
func MustInput(p Interaction, tip string, checker InputChecker) string {
	result, err := p.RawInput(tip, WithInputOptionValid(checker))
	if err != nil {
		panic("user cancel")
	}
	return result
}

// Select get select value
func Select(p Interaction, tip string, list []string, defaultSelect ...int) (result int) {
	if len(defaultSelect) > 0 {
		return p.RawSelect(tip, list, WithSelectOptionDefaults(defaultSelect[0]))
	}

	return p.RawSelect(tip, list)
}

// MustSelect get select value
func MustSelect(p Interaction, tip string, list []string, defaultSelect ...int) (result int) {
	result = Select(p, tip, list, defaultSelect...)
	if result < 0 {
		panic("user cancel")
	}
	return result
}

// SelectString get select value
func SelectString(p Interaction, tip string, list []string, defaultSelect ...int) (_ string, cancel bool) {
	index := Select(p, tip, list, defaultSelect...)
	if index < 0 {
		cancel = true
		return
	}
	return list[index], false
}

// MustSelectString get select value
func MustSelectString(p Interaction, tip string, list []string, defaultSelect ...int) string {
	index := Select(p, tip, list, defaultSelect...)
	if index < 0 {
		panic("user cancel")
	}
	return list[index]
}

// MulSel get multiple value with raw option
func MulSel(p Interaction, tip string, list []string, defaultSelects ...int) (result []int) {
	return p.RawMulSel(tip, list, WithSelectOptionDefaults(defaultSelects...))
}

// MustMulSel get multiple value with raw option
func MustMulSel(p Interaction, tip string, list []string, defaultSelects ...int) (result []int) {
	result = p.RawMulSel(tip, list, WithSelectOptionDefaults(defaultSelects...))
	if len(result) < 1 {
		panic("user cancel")
	}
	return result
}

// MulSelString get multiple value with raw option
func MulSelString(p Interaction, tip string, list []string, defaultSelects ...int) (result []string) {
	sels := p.RawMulSel(tip, list, WithSelectOptionDefaults(defaultSelects...))
	result = make([]string, 0, len(sels))
	for _, k := range sels {
		result = append(result, list[k])
	}
	return result
}

// MustMulSelString get multiple value with raw option
func MustMulSelString(p Interaction, tip string, list []string, defaultSelects ...int) (result []string) {
	sels := p.RawMulSel(tip, list, WithSelectOptionDefaults(defaultSelects...))
	if len(sels) < 1 {
		panic("user cancel")
	}
	result = make([]string, 0, len(sels))
	for _, k := range sels {
		result = append(result, list[k])
	}
	return result
}
