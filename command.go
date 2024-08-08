package promptx

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	completion "github.com/aggronmagi/promptx/completion"
	"github.com/aggronmagi/promptx/internal/debug"
	"github.com/spf13/pflag"
)

// Cmd is a shell command handler.
//
// NOTE: Args and Flags are mutually exclusive with SubCommands,
// and SubCommands is used first.
// input sequence is:
//    command1 [subcommand ...] [arg... ] [flag ...]
// like:
//    query user     userid --log-lv=1
//      |    |         |       |
//    cmd1  sub-cmds  arg     flags
// Args must input if it defines.
type Cmd struct {
	// Command name.
	name string
	// Command name aliases.
	aliases []string
	// Function to execute for the command.
	execFunc func(c CommandContext)
	// One liner help message for the command.
	help string
	// More descriptive help message for the command.
	longHelp string

	// dynamicCmdTip use for command name dynamic change.
	// NOTE: dynamicCmdTip are mutually exclusive with Name.
	dynamicCmdTip func(line string) []*Suggest

	// subCommands. sub command.
	subCommands []*Cmd

	// args command args.
	args []CommandParameter

	children     map[string]*Cmd
	dynamicCache []*Suggest     // dynamic cmd cache
	set          *pflag.FlagSet // flags cache
}

// NewCommand Create an interactive command
// 
// name: command name
// help: prompt information
// args: command parameters
func NewCommand(name, help string, args ...CommandParameter) *Cmd {
	return &Cmd{
		name: name,
		help: help,
		args: args,
	}
}

// NewCommandWithFunc create one command
func NewCommandWithFunc(name, help string, f func(ctx CommandContext), args ...CommandParameter) *Cmd {
	return &Cmd{
		name:     name,
		help:     help,
		execFunc: f,
		args:     args,
	}
}

// ExecFunc Set command execution function
//
// see CommondContext for detail.
func (c *Cmd) ExecFunc(f func(c CommandContext)) *Cmd {
	c.execFunc = f
	return c
}

// 
func (c *Cmd) DynamicTip(f func(line string) []*Suggest) *Cmd {
	c.dynamicCmdTip = f
	return c
}

// Aliases Set command alias
func (c *Cmd) Aliases(aliases ...string) *Cmd {
	c.aliases = aliases
	return c
}
func (c *Cmd) LogHelp(long string) *Cmd {
	c.longHelp = long
	return c
}

// SubCommands adds cmd as a subcommand.
func (c *Cmd) SubCommands(cmds ...*Cmd) *Cmd {
	for _, cmd := range cmds {
		c.subCommands = append(c.subCommands, cmd)
	}
	sort.Sort(cmdSorter(c.subCommands))
	return c
}

// DeleteSubCommand deletes cmd from subcommands.
func (c *Cmd) DeleteSubCommand(name string) {
	c.fixCmd()
	if c.children == nil {
		return
	}
	delete(c.children, name)
	for i := len(c.subCommands) - 1; i >= 0; i-- {
		if c.subCommands[i].name == name {
			c.subCommands = append(c.subCommands[:i], c.subCommands[i+1:]...)
			break
		}
	}
}

func (c *Cmd) fixCmd() {
	if len(c.subCommands) == len(c.children) {
		return
	}
	c.children = make(map[string]*Cmd, len(c.subCommands))
	for _, v := range c.subCommands {

		// fix command alias
		for k := len(v.aliases) - 1; k >= 0; k-- {
			// contains each other will cause unexpected behavior
			if strings.Contains(v.aliases[k], v.name) ||
				strings.Contains(v.name, v.aliases[k]) {
				v.aliases = append(v.aliases[:k], v.aliases[k+1:]...)
				continue
			}
		}

		v.fixCmd()
		c.children[v.name] = v
	}
	// has repeated name command
	if len(c.subCommands) != len(c.children) {
		c.subCommands = c.subCommands[:0]
		for _, v := range c.children {
			c.subCommands = append(c.subCommands, v)
		}
	}
	sort.Sort(cmdSorter(c.subCommands))
}

// Children returns the subcommands of c.
func (c *Cmd) Children() []*Cmd {
	c.fixCmd()
	return c.subCommands
}

func (c *Cmd) hasSubcommand() bool {
	if len(c.children) > 1 {
		return true
	}
	if _, ok := c.children["help"]; !ok {
		return true
	}
	return false
}

// HelpText returns the computed help of the command and its subcommands.
func (c Cmd) HelpText() string {
	var b bytes.Buffer
	p := func(s ...interface{}) {
		fmt.Fprintln(&b)
		if len(s) > 0 {
			fmt.Fprintln(&b, s...)
		}
	}
	if c.longHelp != "" {
		p(c.longHelp)
	} else if c.help != "" {
		p(c.help)
	} else if c.name != "" {
		p(c.name, "has no help")
	}
	if c.hasSubcommand() {
		p("Commands:")
		w := tabwriter.NewWriter(&b, 0, 4, 2, ' ', 0)
		for _, child := range c.subCommands {
			fmt.Fprintf(w, "\t%s\t\t\t%s\n", child.name, child.help)
		}
		w.Flush()
		p()
	}
	return b.String()
}

func (c *Cmd) isCmd(name string) bool {
	if c.dynamicCache != nil {
		for _, v := range c.dynamicCache {
			if v.Text == name {
				return true
			}
		}
	} else {
		if c.name == name {
			return true
		}
	}
	for _, v := range c.aliases {
		if v == name {
			return true
		}
	}
	return false
}

// findChildCmd returns the subcommand with matching name or alias.
func (c *Cmd) findChildCmd(name string) *Cmd {

	// find perfect matches first
	for _, cmd := range c.subCommands {
		if cmd.isCmd(name) {
			return cmd
		}
	}

	return nil
}

// ParseInput parse input,and check valid
func (c *Cmd) ParseInput(line string) (cmds []*Cmd, args []string, err error) {
	fields := strings.Fields(line)
	father := c
	discard := -1
	for k, arg := range fields {
		if cmd := father.findChildCmd(arg); cmd != nil {
			cmds = append(cmds, cmd)
			father = cmd
			discard = k + 1
			continue
		}
		discard = k
		break
	}
	debug.Println("discard", discard)
	if discard >= 0 {
		fields = fields[discard:]
	}
	args = fields
	return
}

func (c *Cmd) FixCommandLine(line string, args []string) string {
	fields := strings.Fields(line)
	father := c
	fixs := make([]string, 0, len(fields))
	for _, arg := range fields {
		if cmd := father.findChildCmd(arg); cmd != nil {
			fixs = append(fixs, arg)
			father = cmd
			continue
		}
		break
	}
	fixs = append(fixs, args...)
	return strings.Join(fixs, " ")
}

func (c *Cmd) buildCache(line string) {
	if c.dynamicCmdTip == nil || c.dynamicCache != nil {
		return
	}
	c.dynamicCache = c.dynamicCmdTip(line)
	for _, s := range c.dynamicCache {
		s.Text = c.name + " " + s.Text
	}
}

func (c *Cmd) suggest() *Suggest {
	return &Suggest{
		Text:        c.name,
		Description: c.help,
	}
}

// findSugest find suggest
func (c *Cmd) findSugest(line []rune, pos int, origLine string, cmds []*Cmd) (suggest []*Suggest) {
	// trim left space
	line = TrimSpaceLeft(line)
	//
	var offset int
	goNext := false
	var nextCmd *Cmd
	// match cmd completion
	matchCmd := func(name []rune, cmd *Cmd, s *Suggest) {
		// complete current suggest
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
		// need find sub command or match current command
		if !HasPrefix(line, name) {
			return
		}
		cname := TrimFirstSpace(line)
		debug.Println("check ", string(cname), string(name))
		if !Equal(name, cname) {
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
		return
	}
	for _, child := range c.subCommands {
		// command name
		if child.dynamicCmdTip != nil {
			child.buildCache(origLine)
			for _, v := range child.dynamicCache {
				matchCmd([]rune(v.Text), child, v)
			}
		} else if child.name != "" {
			matchCmd([]rune(child.name), child, nil)
		}
		// command alias
		// if nextCmd != child {
		for _, alias := range child.aliases {
			matchCmd([]rune(alias), child, &Suggest{
				Text:        alias,
				Description: child.help,
			})
		}
		//}
	}
	cmds = append(cmds, c)
	if len(suggest) == 0 {
		if len(cmds) < 1 {
			return
		}
		if len(cmds) == 1 && cmds[0].set == nil {
			return
		}

		// lastest input chars
		left := make([]rune, 0, len(line))
		for _, v := range line {
			if v == ' ' {
				left = left[:0]
				continue
			}
			left = append(left, v)
		}

		if len(left) >= 1 && left[0] != '-' {
			return
		}

		// no more command. find args and falgs
		set := pflag.NewFlagSet("root", pflag.ContinueOnError)
		for k := len(cmds) - 1; k >= 0; k-- {
			if cmds[k].set == nil {
				continue
			}
			set.AddFlagSet(cmds[k].set)
		}
		set.VisitAll(func(f *pflag.Flag) {
			suggest = append(suggest, &Suggest{
				Text:        "--" + f.Name,
				Description: f.Usage,
			})
			if len(left) >= 2 && left[0] == '-' && left[1] == '-' {
			} else {
				suggest = append(suggest, &Suggest{
					Text:        "-" + f.Shorthand,
					Description: f.Usage,
				})
			}
		})
		suggest = completion.FilterFuzzy(suggest, string(left), true)
		return
	} else if len(suggest) != 1 {
		// find sub command completion
		return
	}
	// cut current command name.try find sub commands
	for i := offset; i < len(line); i++ {
		if line[i] == ' ' {
			continue
		}
		return nextCmd.findSugest(line[i:], len(line[i:]), origLine, cmds)
	}
	// match current command. find sub commands
	if goNext {
		return nextCmd.findSugest(nil, 0, origLine, cmds)
	}

	return
}

func (c *Cmd) FindSuggest(doc *Document) []*Suggest {
	c.fixCmd()
	return c.findSugest([]rune(doc.TextBeforeCursor()), doc.CursorPositionCol(), doc.Text, nil)
}

type cmdSorter []*Cmd

func (c cmdSorter) Len() int           { return len(c) }
func (c cmdSorter) Less(i, j int) bool { return c[i].name < c[j].name }
func (c cmdSorter) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
