package promptx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aggronmagi/promptx/v2/blocks"
	"github.com/aggronmagi/promptx/v2/buffer"
)

// ArgDef 参数定义
type ArgDef struct {
	// 字段索引
	Index int
	// 字段名称
	Name string
	// 提示名称（从 tag 获取）
	Prompt string
	// 字段类型
	Type reflect.Type
	// Check 方法名（从 tag 获取）
	CheckName string
	// 是否为必填项
	Required bool
	// 是否为选择类型（Select）
	IsSelect bool
	// 选择选项列表（如果是 Select 类型）
	SelectOptions []string
}

// parseArgDefs 解析参数定义
func parseArgDefs(arg interface{}) []*ArgDef {
	if arg == nil {
		return nil
	}

	argType := reflect.TypeOf(arg)
	if argType.Kind() == reflect.Ptr {
		argType = argType.Elem()
	}

	if argType.Kind() != reflect.Struct {
		return nil
	}

	var defs []*ArgDef

	for i := 0; i < argType.NumField(); i++ {
		field := argType.Field(i)

		// 跳过未导出的字段
		if !field.IsExported() {
			continue
		}

		def := &ArgDef{
			Index:    i,
			Name:     field.Name,
			Type:     field.Type,
			Required: true, // 默认必填
		}

		// 解析 tag
		if parseArgTag(field, def) {
			continue
		}

		defs = append(defs, def)
	}

	return defs
}

// parseArgTag 解析字段 tag
func parseArgTag(field reflect.StructField, def *ArgDef) (ignore bool) {
	ignore = false
	if field.Tag.Get("ignore") != "" {
		ignore = true
		return
	}
	// 解析 arg tag
	argTag := field.Tag.Get("arg")
	if argTag != "" {
		parts := strings.Split(argTag, ",")
		if len(parts) > 0 && parts[0] != "" {
			def.Prompt = parts[0]
		}
		// 检查是否有 required=false
		for _, part := range parts[1:] {
			if part == "optional" {
				def.Required = false
			}
		}
	}

	// 解析 prompt tag（优先级高于 arg）
	promptTag := field.Tag.Get("prompt")
	if promptTag != "" {
		def.Prompt = promptTag
	}

	// 解析 check tag
	checkTag := field.Tag.Get("check")
	if checkTag != "" {
		def.CheckName = checkTag
	}

	// 解析 select tag（选择类型）
	selectTag := field.Tag.Get("select")
	if selectTag != "" {
		def.IsSelect = true
		// select tag 格式：select:"option1,option2,option3"
		options := strings.Split(selectTag, ",")
		def.SelectOptions = options
	}

	// 如果没有设置 prompt，使用字段名
	if def.Prompt == "" {
		def.Prompt = field.Name
	}
	return ignore
}

// parseArgs 解析参数值
func parseArgs(ctx blocks.Context, defs []*ArgDef) interface{} {
	if len(defs) == 0 {
		return nil
	}

	// 从上下文中获取已解析的参数
	// 这里需要从命令解析中获取参数值
	// 暂时返回 nil，实际实现需要在 command_exec.go 中完成
	return nil
}

// createArgValue 创建参数值
func createArgValue(defs []*ArgDef, argType reflect.Type) reflect.Value {
	if argType.Kind() == reflect.Ptr {
		argType = argType.Elem()
	}

	value := reflect.New(argType).Elem()

	for _, def := range defs {
		if def.Index < value.NumField() {
			field := value.Field(def.Index)
			if field.CanSet() {
				// 设置默认值
				setDefaultValue(field, def)
			}
		}
	}

	return value
}

// setDefaultValue 设置默认值
func setDefaultValue(field reflect.Value, def *ArgDef) {
	switch field.Kind() {
	case reflect.String:
		field.SetString("")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.SetUint(0)
	case reflect.Float32, reflect.Float64:
		field.SetFloat(0)
	case reflect.Bool:
		field.SetBool(false)
	}
}

// getArgValueFromString 从字符串获取参数值
func getArgValueFromString(str string, fieldType reflect.Type) (reflect.Value, error) {
	switch fieldType.Kind() {
	case reflect.String:
		return reflect.ValueOf(str), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var v int64
		_, err := fmt.Sscanf(str, "%d", &v)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v).Convert(fieldType), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var v uint64
		_, err := fmt.Sscanf(str, "%d", &v)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v).Convert(fieldType), nil
	case reflect.Float32, reflect.Float64:
		var v float64
		_, err := fmt.Sscanf(str, "%f", &v)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v).Convert(fieldType), nil
	case reflect.Bool:
		v := strings.ToLower(str) == "true" || str == "1"
		return reflect.ValueOf(v), nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type: %v", fieldType)
	}
}

// createArgValueFromStrings 从字符串数组创建参数值
func createArgValueFromStrings(defs []*ArgDef, argType reflect.Type, values []string) reflect.Value {
	if argType.Kind() == reflect.Ptr {
		argType = argType.Elem()
	}

	value := reflect.New(argType).Elem()

	for i, def := range defs {
		if i >= len(values) {
			break
		}

		if def.Index < value.NumField() {
			field := value.Field(def.Index)
			if field.CanSet() {
				strValue := values[i]
				if strValue != "" {
					fieldValue, err := getArgValueFromString(strValue, def.Type)
					if err == nil {
						field.Set(fieldValue)
					}
				} else {
					// 设置默认值
					setDefaultValue(field, def)
				}
			}
		}
	}

	return value
}

// checkArg 检查参数
func checkArg(ctx blocks.Context, def *ArgDef, value string, args []string, index int) (string, error) {
	// 如果已经有值，先检查
	if index < len(args) && args[index] != "" {
		value = args[index]
		// 执行检查
		if def.CheckName != "" {
			checker, ok := checkers[def.CheckName]
			if ok {
				if err := checker(value); err != nil {
					return "", fmt.Errorf("invalid value for %s: %v", def.Prompt, err)
				}
			}
		}
		return value, nil
	}

	// 如果没有值且是必填项，需要交互输入
	if def.Required {
		if def.IsSelect && len(def.SelectOptions) > 0 {
			// 使用 Select
			sel := ctx.RawSelect(def.Prompt, def.SelectOptions)
			if sel < 0 {
				return "", errors.New("user cancel")
			}
			return def.SelectOptions[sel], nil
		} else {
			// 使用 Input
			opts := []blocks.InputOption{}
			// 如果有检查器，设置验证函数
			if def.CheckName != "" {
				checker, ok := checkers[def.CheckName]
				if ok {
					opts = append(opts, blocks.WithInputOptionValid(func(doc *buffer.Document) error {
						return checker(doc.Text)
					}))
				}
			}
			result, err := ctx.RawInput(def.Prompt, opts...)
			if err != nil {
				return "", err
			}
			return result, nil
		}
	}

	// 可选参数，返回空值
	return "", nil
}

// checkArgs 检查所有参数
func checkArgs(ctx blocks.Context, defs []*ArgDef, args []string) ([]string, error) {
	checkedArgs := make([]string, len(defs))

	for i, def := range defs {
		value, err := checkArg(ctx, def, "", args, i)
		if err != nil {
			return nil, err
		}
		checkedArgs[i] = value
	}

	return checkedArgs, nil
}

// createArgValueForCommander 为 Commander 类型创建参数值
// originalType 是原始类型（可能是指针或值）
func createArgValueForCommander(defs []*ArgDef, originalType reflect.Type, values []string) reflect.Value {
	// 确定结构体类型
	structType := originalType
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	// 创建结构体值
	value := reflect.New(structType).Elem()

	// 填充字段
	for i, def := range defs {
		if i >= len(values) {
			break
		}

		if def.Index < value.NumField() {
			field := value.Field(def.Index)
			if field.CanSet() {
				strValue := values[i]
				if strValue != "" {
					fieldValue, err := getArgValueFromString(strValue, def.Type)
					if err == nil {
						field.Set(fieldValue)
					}
				} else {
					setDefaultValue(field, def)
				}
			}
		}
	}

	// 根据原始类型返回值或指针
	if originalType.Kind() == reflect.Ptr {
		// 原始是指针，返回指针
		return value.Addr()
	} else {
		// 原始是值，返回值
		return value
	}
}
