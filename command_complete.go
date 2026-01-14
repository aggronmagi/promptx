package promptx

import (
	"strings"
	"unicode"

	"github.com/aggronmagi/promptx/v2/buffer"
	"github.com/aggronmagi/promptx/v2/completion"
)

// FindSuggest 查找建议（用于自动补全）
func (c *Command) FindSuggest(doc buffer.Document) []*completion.Suggest {
	c.fixChildren()
	line := []rune(doc.TextBeforeCursor())
	return c.findSuggest(line, doc.Text, nil)
}

// findSuggest 查找建议（内部实现）
func (c *Command) findSuggest(line []rune, origLine string, cmds []*Command) []*completion.Suggest {
	// 确保子命令映射已更新
	c.fixChildren()

	// 去除左侧空格
	line = trimSpaceLeft(line)
	if len(line) == 0 && c.config != nil {
		// 如果输入为空，返回所有子命令
		var suggests []*completion.Suggest
		for _, child := range c.subCommands {
			suggests = append(suggests, child.suggest())
			// 也添加别名
			for _, alias := range child.aliases {
				suggests = append(suggests, &completion.Suggest{
					Text:        alias,
					Description: child.help,
				})
			}
		}
		return suggests
	}

	var offset int
	goNext := false
	var nextCmd *Command
	var suggest []*completion.Suggest

	// 匹配命令
	matchCmd := func(name []rune, cmd *Command, s *completion.Suggest) {
		// 如果输入长度小于命令名长度，进行模糊匹配
		if len(line) < len(name) {
			if !completion.FuzzyMatchRunes(name, line) {
				return
			}
			if s == nil {
				s = cmd.suggest()
			}
			suggest = append(suggest, s)
			offset = len(line)
			nextCmd = cmd
			return
		}

		// 需要查找子命令或匹配当前命令
		if !hasPrefix(line, name) {
			return
		}

		cname := trimFirstSpace(line)
		if !equal(name, cname) {
			return
		}

		if s == nil {
			s = cmd.suggest()
		}

		if len(line) > len(name) {
			nextCmd = cmd
			goNext = true
		}
		suggest = append(suggest, s)
		offset = len(name)
	}

	// 遍历子命令
	for _, child := range c.subCommands {
		// 命令名
		if child.name != "" {
			matchCmd([]rune(child.name), child, nil)
		}
		// 命令别名
		for _, alias := range child.aliases {
			matchCmd([]rune(alias), child, &completion.Suggest{
				Text:        alias,
				Description: child.help,
			})
		}
	}

	cmds = append(cmds, c)

	// 如果没有找到建议，返回
	if len(suggest) == 0 {
		return nil
	}

	// 如果有多个建议，返回所有建议
	if len(suggest) != 1 {
		return suggest
	}

	// 只有一个建议，继续查找子命令
	// 跳过当前命令名，尝试查找子命令
	for i := offset; i < len(line); i++ {
		if line[i] == ' ' {
			continue
		}
		if nextCmd != nil {
			return nextCmd.findSuggest(line[i:], origLine, cmds)
		}
		return nil
	}

	// 匹配当前命令，查找子命令
	if goNext && nextCmd != nil {
		return nextCmd.findSuggest(nil, origLine, cmds)
	}

	return suggest
}

// suggest 创建建议
func (c *Command) suggest() *completion.Suggest {
	return &completion.Suggest{
		Text:        c.name,
		Description: c.help,
	}
}

// 辅助函数

// trimSpaceLeft 去除左侧空格
func trimSpaceLeft(in []rune) []rune {
	for i, r := range in {
		if !unicode.IsSpace(r) {
			return in[i:]
		}
	}
	return nil
}

// hasPrefix 检查前缀
func hasPrefix(r, prefix []rune) bool {
	if len(r) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if r[i] != prefix[i] {
			return false
		}
	}
	return true
}

// equal 比较两个 rune 切片
func equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// trimFirstSpace 去除第一个空格及其之前的内容
func trimFirstSpace(in []rune) []rune {
	for i, r := range in {
		if unicode.IsSpace(r) {
			if i+1 < len(in) {
				return in[i+1:]
			}
			return nil
		}
	}
	return in
}

// createCompleter 创建自动补全器
// root: 根命令（可能是命令组的根命令）
// commandPrefix: 命令前缀（如果有）
func createCompleter(root *Command, commandPrefix string) func(buffer.Document) []*completion.Suggest {
	return func(doc buffer.Document) []*completion.Suggest {
		text := doc.Text

		// 如果有命令前缀，需要先处理前缀
		if commandPrefix != "" {
			if !strings.HasPrefix(text, commandPrefix) {
				// 输入还没有输入前缀，不提供补全
				return nil
			}
			// 计算前缀的字节长度和 rune 长度
			prefixLen := len(commandPrefix)
			prefixRunes := len([]rune(commandPrefix))

			// 移除前缀
			newText := text[prefixLen:]
			// 调整光标位置
			newCursor := doc.CursorPosition() - prefixRunes
			if newCursor < 0 {
				newCursor = 0
			}
			// 创建新的 Document
			doc = *buffer.NewDocumentWithCursor(newText, newCursor)
		}

		// 使用命令的 FindSuggest
		return root.FindSuggest(doc)
	}
}
