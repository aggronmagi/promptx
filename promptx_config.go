// config.go
package promptx

import (
	"github.com/aggronmagi/promptx/input"
	"github.com/aggronmagi/promptx/output"
)

// PromptxConfigs 是Promptx的链式配置入口
// 提供了流畅的API用于配置Promptx的各个方面
type PromptxConfigs struct {
	common       []CommonOption
	input        []InputOption
	selects      []SelectOption
	complete     []CompleteOption
	manager      BlocksManager
	inputParser  input.ConsoleParser
	outputWriter output.ConsoleWriter
	stderrWriter output.ConsoleWriter
	commandSets  []*Cmd
}

// NewConfig 创建并返回一个新的Promptx链式配置器
func NewConfig() *PromptxConfigs {
	return &PromptxConfigs{}
}

// Build 根据当前配置构建并返回Promptx实例
func (c *PromptxConfigs) Build() PromptMain {
	return New(
		WithInputs(c.input...),
		WithSelects(c.selects...),
		WithCommon(c.common...),
		WithManager(c.manager),
		WithInput(c.inputParser),
		WithOutput(c.outputWriter),
		WithStderr(c.stderrWriter),
	)
}

// Theme 返回主题/样式配置器
// 用于配置颜色、背景等视觉样式
func (c *PromptxConfigs) Theme() *ThemeConfig {
	return &ThemeConfig{inner: c}
}

// Keys 返回键盘按键配置器
// 用于配置各种操作的快捷键
func (c *PromptxConfigs) Keys() *KeysConfig {
	return &KeysConfig{inner: c}
}

// Input 返回输入框配置器
// 用于配置输入框的行为和内容
func (c *PromptxConfigs) Input() *InputConfig {
	return &InputConfig{inner: c}
}

// Select 返回选择器配置器
// 用于配置选择列表的行为和内容
func (c *PromptxConfigs) Select() *SelectConfig {
	return &SelectConfig{inner: c}
}

// Complete 返回自动补全配置器
// 用于配置自动补全的行为和样式
func (c *PromptxConfigs) Complete() *CompleteConfig {
	return &CompleteConfig{inner: c}
}

// Common 返回通用配置器
// 用于配置Promptx的通用设置
func (c *PromptxConfigs) Common() *CommonConfig {
	return &CommonConfig{inner: c}
}

// Hardware 返回硬件/IO配置器
// 用于配置输入输出设备和管理器
func (c *PromptxConfigs) Hardware() *HardwareConfig {
	return &HardwareConfig{inner: c}
}

// Commands 返回命令配置器
// 用于配置可执行的命令
func (c *PromptxConfigs) Commands() *CommandsConfig {
	return &CommandsConfig{inner: c}
}

// KeysConfig 键盘按键配置器
// 用于统一管理各种操作的快捷键
type KeysConfig struct {
	inner *PromptxConfigs
}

// Common 返回通用快捷键配置器
// 用于配置通用的确认和取消快捷键
func (k *KeysConfig) Common() *CommonKeysConfig {
	return &CommonKeysConfig{inner: k.inner}
}

// Input 返回输入框快捷键配置器
// 用于配置输入框专用的快捷键
func (k *KeysConfig) Input() *InputKeysConfig {
	return &InputKeysConfig{inner: k.inner}
}

// Select 返回选择器快捷键配置器
// 用于配置选择器专用的快捷键
func (k *KeysConfig) Select() *SelectKeysConfig {
	return &SelectKeysConfig{inner: k.inner}
}

// CommonKeysConfig 通用快捷键配置器
type CommonKeysConfig struct {
	inner *PromptxConfigs
}

// Finish 设置通用确认键
func (c *CommonKeysConfig) Finish(key Key) *CommonKeysConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionFinish(key))
	return c
}

// Cancel 设置通用取消键
func (c *CommonKeysConfig) Cancel(key Key) *CommonKeysConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionCancel(key))
	return c
}

// InputKeysConfig 输入框快捷键配置器
type InputKeysConfig struct {
	inner *PromptxConfigs
}

// Finish 设置输入框确认键
func (i *InputKeysConfig) Finish(key Key) *InputKeysConfig {
	i.inner.input = append(i.inner.input, WithInputOptionFinish(key))
	return i
}

// Cancel 设置输入框取消键
func (i *InputKeysConfig) Cancel(key Key) *InputKeysConfig {
	i.inner.input = append(i.inner.input, WithInputOptionCancel(key))
	return i
}

// SelectKeysConfig 选择器快捷键配置器
type SelectKeysConfig struct {
	inner *PromptxConfigs
}

// Finish 设置选择器确认键
func (s *SelectKeysConfig) Finish(key Key) *SelectKeysConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionFinish(key))
	return s
}

// Cancel 设置选择器取消键
func (s *SelectKeysConfig) Cancel(key Key) *SelectKeysConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionCancel(key))
	return s
}

// ThemeConfig 主题/样式配置器
// 用于统一配置应用程序的视觉样式
type ThemeConfig struct {
	inner *PromptxConfigs
}

// Common 返回通用主题配置器
// 用于配置通用提示、前缀等样式
func (t *ThemeConfig) Common() *ThemeCommonConfig {
	return &ThemeCommonConfig{inner: t.inner}
}

// Input 返回输入框主题配置器
// 用于配置输入框相关的样式
func (t *ThemeConfig) Input() *ThemeInputConfig {
	return &ThemeInputConfig{inner: t.inner}
}

// Select 返回选择器主题配置器
// 用于配置选择器相关的样式
func (t *ThemeConfig) Select() *ThemeSelectConfig {
	return &ThemeSelectConfig{inner: t.inner}
}

// Complete 返回自动补全主题配置器
// 用于配置自动补全相关的样式
func (t *ThemeConfig) Complete() *ThemeCompleteConfig {
	return &ThemeCompleteConfig{inner: t.inner}
}

// ThemeCommonConfig 通用主题配置器
type ThemeCommonConfig struct {
	inner *PromptxConfigs
}

// TipColor 设置通用提示文字颜色
func (t *ThemeCommonConfig) TipColor(color Color) *ThemeCommonConfig {
	t.inner.common = append(t.inner.common, WithCommonOptionTipColor(color))
	return t
}

// TipBG 设置通用提示文字背景颜色
func (t *ThemeCommonConfig) TipBG(color Color) *ThemeCommonConfig {
	t.inner.common = append(t.inner.common, WithCommonOptionTipBG(color))
	return t
}

// PrefixColor 设置通用前缀文字颜色
func (t *ThemeCommonConfig) PrefixColor(color Color) *ThemeCommonConfig {
	t.inner.common = append(t.inner.common, WithCommonOptionPrefixColor(color))
	return t
}

// PrefixBG 设置通用前缀文字背景颜色
func (t *ThemeCommonConfig) PrefixBG(color Color) *ThemeCommonConfig {
	t.inner.common = append(t.inner.common, WithCommonOptionPrefixBG(color))
	return t
}

// ValidColor 设置通用验证错误文字颜色
func (t *ThemeCommonConfig) ValidColor(color Color) *ThemeCommonConfig {
	t.inner.common = append(t.inner.common, WithCommonOptionValidColor(color))
	return t
}

// ValidBG 设置通用验证错误文字背景颜色
func (t *ThemeCommonConfig) ValidBG(color Color) *ThemeCommonConfig {
	t.inner.common = append(t.inner.common, WithCommonOptionValidBG(color))
	return t
}

// ThemeInputConfig 输入框主题配置器
type ThemeInputConfig struct {
	inner *PromptxConfigs
}

// TipColor 设置输入框提示文字颜色
func (t *ThemeInputConfig) TipColor(color Color) *ThemeInputConfig {
	t.inner.input = append(t.inner.input, WithInputOptionTipColor(color))
	return t
}

// TipBG 设置输入框提示文字背景颜色
func (t *ThemeInputConfig) TipBG(color Color) *ThemeInputConfig {
	t.inner.input = append(t.inner.input, WithInputOptionTipBG(color))
	return t
}

// PrefixColor 设置输入框前缀文字颜色
func (t *ThemeInputConfig) PrefixColor(color Color) *ThemeInputConfig {
	t.inner.input = append(t.inner.input, WithInputOptionPrefixColor(color))
	return t
}

// PrefixBG 设置输入框前缀文字背景颜色
func (t *ThemeInputConfig) PrefixBG(color Color) *ThemeInputConfig {
	t.inner.input = append(t.inner.input, WithInputOptionPrefixBG(color))
	return t
}

// ValidColor 设置输入框验证错误文字颜色
func (t *ThemeInputConfig) ValidColor(color Color) *ThemeInputConfig {
	t.inner.input = append(t.inner.input, WithInputOptionValidColor(color))
	return t
}

// ValidBG 设置输入框验证错误文字背景颜色
func (t *ThemeInputConfig) ValidBG(color Color) *ThemeInputConfig {
	t.inner.input = append(t.inner.input, WithInputOptionValidBG(color))
	return t
}

// ResultColor 设置输入框结果显示文字颜色
func (t *ThemeInputConfig) ResultColor(color Color) *ThemeInputConfig {
	t.inner.input = append(t.inner.input, WithInputOptionResultColor(color))
	return t
}

// ResultBG 设置输入框结果显示文字背景颜色
func (t *ThemeInputConfig) ResultBG(color Color) *ThemeInputConfig {
	t.inner.input = append(t.inner.input, WithInputOptionResultBG(color))
	return t
}

// DefaultColor 设置输入框默认值文字颜色
func (t *ThemeInputConfig) DefaultColor(color Color) *ThemeInputConfig {
	t.inner.input = append(t.inner.input, WithInputOptionDefaultColor(color))
	return t
}

// DefaultBG 设置输入框默认值文字背景颜色
func (t *ThemeInputConfig) DefaultBG(color Color) *ThemeInputConfig {
	t.inner.input = append(t.inner.input, WithInputOptionDefaultBG(color))
	return t
}

// ThemeSelectConfig 选择器主题配置器
type ThemeSelectConfig struct {
	inner *PromptxConfigs
}

// TipColor 设置选择器提示文字颜色
func (t *ThemeSelectConfig) TipColor(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionTipColor(color))
	return t
}

// TipBG 设置选择器提示文字背景颜色
func (t *ThemeSelectConfig) TipBG(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionTipBG(color))
	return t
}

// HelpColor 设置选择器帮助文字颜色
func (t *ThemeSelectConfig) HelpColor(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionHelpColor(color))
	return t
}

// HelpBG 设置选择器帮助文字背景颜色
func (t *ThemeSelectConfig) HelpBG(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionHelpBG(color))
	return t
}

// ValidColor 设置选择器验证错误文字颜色
func (t *ThemeSelectConfig) ValidColor(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionValidColor(color))
	return t
}

// ValidBG 设置选择器验证错误文字背景颜色
func (t *ThemeSelectConfig) ValidBG(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionValidBG(color))
	return t
}

// SuggestColor 设置选择器选项文字颜色
func (t *ThemeSelectConfig) SuggestColor(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionSuggestColor(color))
	return t
}

// SuggestBG 设置选择器选项文字背景颜色
func (t *ThemeSelectConfig) SuggestBG(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionSuggestBG(color))
	return t
}

// SelSuggestColor 设置选择器选中选项文字颜色
func (t *ThemeSelectConfig) SelSuggestColor(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionSelSuggestColor(color))
	return t
}

// SelSuggestBG 设置选择器选中选项文字背景颜色
func (t *ThemeSelectConfig) SelSuggestBG(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionSelSuggestBG(color))
	return t
}

// DescColor 设置选择器描述文字颜色
func (t *ThemeSelectConfig) DescColor(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionDescColor(color))
	return t
}

// DescBG 设置选择器描述文字背景颜色
func (t *ThemeSelectConfig) DescBG(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionDescBG(color))
	return t
}

// SelDescColor 设置选择器选中描述文字颜色
func (t *ThemeSelectConfig) SelDescColor(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionSelDescColor(color))
	return t
}

// SelDescBG 设置选择器选中描述文字背景颜色
func (t *ThemeSelectConfig) SelDescBG(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionSelDescBG(color))
	return t
}

// BarColor 设置选择器滚动条颜色
func (t *ThemeSelectConfig) BarColor(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionBarColor(color))
	return t
}

// BarBG 设置选择器滚动条背景颜色
func (t *ThemeSelectConfig) BarBG(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionBarBG(color))
	return t
}

// ResultColor 设置选择器结果显示文字颜色
func (t *ThemeSelectConfig) ResultColor(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionResultColor(color))
	return t
}

// ResultBG 设置选择器结果显示文字背景颜色
func (t *ThemeSelectConfig) ResultBG(color Color) *ThemeSelectConfig {
	t.inner.selects = append(t.inner.selects, WithSelectOptionResultBG(color))
	return t
}

// ThemeCompleteConfig 自动补全主题配置器
type ThemeCompleteConfig struct {
	inner *PromptxConfigs
}

// SuggestionTextColor 设置自动补全建议文字颜色
func (t *ThemeCompleteConfig) SuggestionTextColor(color Color) *ThemeCompleteConfig {
	t.inner.complete = append(t.inner.complete, WithCompleteOptionSuggestionTextColor(color))
	return t
}

// SuggestionBGColor 设置自动补全建议文字背景颜色
func (t *ThemeCompleteConfig) SuggestionBGColor(color Color) *ThemeCompleteConfig {
	t.inner.complete = append(t.inner.complete, WithCompleteOptionSuggestionBGColor(color))
	return t
}

// SelectedSuggestionTextColor 设置自动补全选中建议文字颜色
func (t *ThemeCompleteConfig) SelectedSuggestionTextColor(color Color) *ThemeCompleteConfig {
	t.inner.complete = append(t.inner.complete, WithCompleteOptionSelectedSuggestionTextColor(color))
	return t
}

// SelectedSuggestionBGColor 设置自动补全选中建议文字背景颜色
func (t *ThemeCompleteConfig) SelectedSuggestionBGColor(color Color) *ThemeCompleteConfig {
	t.inner.complete = append(t.inner.complete, WithCompleteOptionSelectedSuggestionBGColor(color))
	return t
}

// DescriptionTextColor 设置自动补全描述文字颜色
func (t *ThemeCompleteConfig) DescriptionTextColor(color Color) *ThemeCompleteConfig {
	t.inner.complete = append(t.inner.complete, WithCompleteOptionDescriptionTextColor(color))
	return t
}

// DescriptionBGColor 设置自动补全描述文字背景颜色
func (t *ThemeCompleteConfig) DescriptionBGColor(color Color) *ThemeCompleteConfig {
	t.inner.complete = append(t.inner.complete, WithCompleteOptionDescriptionBGColor(color))
	return t
}

// SelectedDescriptionTextColor 设置自动补全选中描述文字颜色
func (t *ThemeCompleteConfig) SelectedDescriptionTextColor(color Color) *ThemeCompleteConfig {
	t.inner.complete = append(t.inner.complete, WithCompleteOptionSelectedDescriptionTextColor(color))
	return t
}

// SelectedDescriptionBGColor 设置自动补全选中描述文字背景颜色
func (t *ThemeCompleteConfig) SelectedDescriptionBGColor(color Color) *ThemeCompleteConfig {
	t.inner.complete = append(t.inner.complete, WithCompleteOptionSelectedDescriptionBGColor(color))
	return t
}

// ScrollbarThumbColor 设置自动补全滚动条滑块颜色
func (t *ThemeCompleteConfig) ScrollbarThumbColor(color Color) *ThemeCompleteConfig {
	t.inner.complete = append(t.inner.complete, WithCompleteOptionScrollbarThumbColor(color))
	return t
}

// ScrollbarBGColor 设置自动补全滚动条背景颜色
func (t *ThemeCompleteConfig) ScrollbarBGColor(color Color) *ThemeCompleteConfig {
	t.inner.complete = append(t.inner.complete, WithCompleteOptionScrollbarBGColor(color))
	return t
}

// InputConfig 输入框功能配置器
type InputConfig struct {
	inner *PromptxConfigs
}

// Tip 设置输入框提示文字
func (i *InputConfig) Tip(tip string) *InputConfig {
	i.inner.input = append(i.inner.input, WithInputOptionTip(tip))
	return i
}

// Prefix 设置输入框前缀文字
func (i *InputConfig) Prefix(prefix string) *InputConfig {
	i.inner.input = append(i.inner.input, WithInputOptionPrefix(prefix))
	return i
}

// Default 设置输入框默认值
func (i *InputConfig) Default(value string) *InputConfig {
	i.inner.input = append(i.inner.input, WithInputOptionDefault(value))
	return i
}

// Valid 设置输入框验证函数
func (i *InputConfig) Valid(validator func(*Document) error) *InputConfig {
	i.inner.input = append(i.inner.input, WithInputOptionValid(validator))
	return i
}

// OnFinish 设置输入框完成回调函数
func (i *InputConfig) OnFinish(fn func(input string, eof error)) *InputConfig {
	i.inner.input = append(i.inner.input, WithInputOptionOnFinish(fn))
	return i
}

// ResultText 设置输入框结果显示文本函数
func (i *InputConfig) ResultText(fn InputFinishTextFunc) *InputConfig {
	i.inner.input = append(i.inner.input, WithInputOptionResultText(fn))
	return i
}

// SelectConfig 选择器功能配置器
type SelectConfig struct {
	inner *PromptxConfigs
}

// Options 设置选择器选项列表
func (s *SelectConfig) Options(options ...*Suggest) *SelectConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionOptions(options...))
	return s
}

// Rows 设置选择器显示行数
func (s *SelectConfig) Rows(rows int) *SelectConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionRows(rows))
	return s
}

// Multi 设置选择器是否允许多选
func (s *SelectConfig) Multi(multi bool) *SelectConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionMulti(multi))
	return s
}

// Tip 设置选择器提示文字
func (s *SelectConfig) Tip(tip string) *SelectConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionTip(tip))
	return s
}

// ShowHelp 设置选择器是否显示帮助
func (s *SelectConfig) ShowHelp(show bool) *SelectConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionShowHelp(show))
	return s
}

// Help 设置选择器帮助文本函数
func (s *SelectConfig) Help(fn SelHelpTextFunc) *SelectConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionHelp(fn))
	return s
}

// OnFinish 设置选择器完成回调函数
func (s *SelectConfig) OnFinish(fn func(sels []int)) *SelectConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionOnFinish(fn))
	return s
}

// Valid 设置选择器验证函数
func (s *SelectConfig) Valid(validator func(sels []int) error) *SelectConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionValid(validator))
	return s
}

// Defaults 设置选择器默认选中项
func (s *SelectConfig) Defaults(defaults ...int) *SelectConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionDefaults(defaults...))
	return s
}

// ShowItem 设置选择器是否显示项目
func (s *SelectConfig) ShowItem(show bool) *SelectConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionShowItem(show))
	return s
}

// FinishText 设置选择器完成文本函数
func (s *SelectConfig) FinishText(fn SelFinishTextFunc) *SelectConfig {
	s.inner.selects = append(s.inner.selects, WithSelectOptionFinishText(fn))
	return s
}

// CompleteConfig 自动补全功能配置器
type CompleteConfig struct {
	inner *PromptxConfigs
}

// Completer 设置自动补全器
func (c *CompleteConfig) Completer(completer Completer) *CompleteConfig {
	c.inner.complete = append(c.inner.complete, WithCompleteOptionCompleter(completer))
	return c
}

// Max 设置自动补全最大建议数
func (c *CompleteConfig) Max(max int) *CompleteConfig {
	c.inner.complete = append(c.inner.complete, WithCompleteOptionCompleteMax(max))
	return c
}

// FillSpace 设置自动补全是否填充空格
func (c *CompleteConfig) FillSpace(fill bool) *CompleteConfig {
	c.inner.complete = append(c.inner.complete, WithCompleteOptionCompletionFillSpace(fill))
	return c
}

// WordSeparator 设置自动补全单词分隔符
func (c *CompleteConfig) WordSeparator(separator string) *CompleteConfig {
	c.inner.complete = append(c.inner.complete, WithCompleteOptionWordSeparator(separator))
	return c
}

// CommonConfig 通用功能配置器
type CommonConfig struct {
	inner *PromptxConfigs
}

// Tip 设置通用提示文字
func (c *CommonConfig) Tip(tip string) *CommonConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionTip(tip))
	return c
}

// Prefix 设置通用前缀文字
func (c *CommonConfig) Prefix(prefix string) *CommonConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionPrefix(prefix))
	return c
}

// Valid 设置通用验证函数
func (c *CommonConfig) Valid(validator func(status int, in *Document) error) *CommonConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionValid(validator))
	return c
}

// Exec 设置通用执行函数
func (c *CommonConfig) Exec(executor func(ctx Context, command string)) *CommonConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionExec(executor))
	return c
}

// AlwaysCheck 设置是否始终检查输入
func (c *CommonConfig) AlwaysCheck(check bool) *CommonConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionAlwaysCheck(check))
	return c
}

// History 设置历史记录文件路径
func (c *CommonConfig) History(file string) *CommonConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionHistory(file))
	return c
}

// HistoryMaxSize 设置历史记录最大大小
func (c *CommonConfig) HistoryMaxSize(size int) *CommonConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionHistoryMaxSize(size))
	return c
}

// CommandPrefix 设置命令前缀
func (c *CommonConfig) CommandPrefix(prefix string) *CommonConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionCommandPrefix(prefix))
	return c
}

// PreCheck 设置执行前检查函数
func (c *CommonConfig) PreCheck(check func(ctx Context) error) *CommonConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionPreCheck(check))
	return c
}

// OnNonCommand 设置非命令处理函数
func (c *CommonConfig) OnNonCommand(handler func(ctx Context, command string) error) *CommonConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionOnNonCommand(handler))
	return c
}

// Complete 设置自动补全选项
func (c *CommonConfig) Complete(options ...CompleteOption) *CommonConfig {
	c.inner.common = append(c.inner.common, WithCommonOptionComplete(options...))
	return c
}

// HardwareConfig 硬件/IO配置器
type HardwareConfig struct {
	inner *PromptxConfigs
}

// InputParser 设置输入解析器
func (h *HardwareConfig) InputParser(parser input.ConsoleParser) *HardwareConfig {
	h.inner.inputParser = parser
	return h
}

// OutputWriter 设置输出写入器
func (h *HardwareConfig) OutputWriter(writer output.ConsoleWriter) *HardwareConfig {
	h.inner.outputWriter = writer
	return h
}

// StderrWriter 设置错误输出写入器
func (h *HardwareConfig) StderrWriter(writer output.ConsoleWriter) *HardwareConfig {
	h.inner.stderrWriter = writer
	return h
}

// Manager 设置块管理器
func (h *HardwareConfig) Manager(manager BlocksManager) *HardwareConfig {
	h.inner.manager = manager
	return h
}

// CommandsConfig 命令配置器
type CommandsConfig struct {
	inner *PromptxConfigs
}

// Add 添加命令到命令集
func (c *CommandsConfig) Add(cmds ...*Cmd) *CommandsConfig {
	c.inner.commandSets = append(c.inner.commandSets, cmds...)
	c.inner.common = append(c.inner.common, WithCommonOptionCmds(c.inner.commandSets...))
	return c
}
