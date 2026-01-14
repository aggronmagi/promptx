package promptx

import (
	"errors"
	"fmt"
	"strings"
)

// CheckFunc 参数检查函数类型
type CheckFunc func(value string) error

var (
	// checkers 注册的检查器
	checkers = make(map[string]CheckFunc)
)

// RegisterChecker 注册检查器
func RegisterChecker(name string, fn CheckFunc) {
	checkers[name] = fn
}

func init() {
	// 注册内置检查器
	RegisterChecker("NotEmpty", CheckerNotEmpty)
	RegisterChecker("NotEmptyAndSpace", CheckerNotEmptyAndSpace)
	RegisterChecker("Integer", CheckerInteger)
	RegisterChecker("NotZeroInteger", CheckerNotZeroInteger)
	RegisterChecker("NaturalNumber", CheckerNaturalNumber)
}

// CheckerNotEmpty 检查非空
func CheckerNotEmpty(value string) error {
	if value == "" {
		return errors.New("empty input")
	}
	return nil
}

// CheckerNotEmptyAndSpace 检查非空且不包含空格
func CheckerNotEmptyAndSpace(value string) error {
	if value == "" {
		return errors.New("empty input")
	}
	if strings.ContainsAny(value, " \n\t") {
		return errors.New("contain invalid char(space,enter,tab)")
	}
	return nil
}

// CheckerInteger 检查整数
func CheckerInteger(value string) error {
	if value == "" {
		return errors.New("empty input")
	}
	var v int64
	_, err := fmt.Sscanf(value, "%d", &v)
	return err
}

// CheckerNotZeroInteger 检查非零整数
func CheckerNotZeroInteger(value string) error {
	if err := CheckerInteger(value); err != nil {
		return err
	}
	var v int64
	fmt.Sscanf(value, "%d", &v)
	if v == 0 {
		return errors.New("zero value invalid")
	}
	return nil
}

// CheckerNaturalNumber 检查自然数（>=1）
func CheckerNaturalNumber(value string) error {
	if err := CheckerInteger(value); err != nil {
		return err
	}
	var v int64
	fmt.Sscanf(value, "%d", &v)
	if v < 1 {
		return fmt.Errorf("%d is not natural number", v)
	}
	return nil
}
