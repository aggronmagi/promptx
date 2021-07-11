package promptx

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/pflag"
)

// CommandContext run context
type CommandContext struct {
	Context
	// command and subcommands
	Cmds []*Cmd
	// input args
	Args []string
	// input line
	Line string
	// Root root command. use for dynamic modify command.
	Root *Cmd
}

// // CmdArg command args interface
// type CmdArg interface {
// }

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
	Name string
	// Command name aliases.
	Aliases []string
	// Function to execute for the command.
	Func func(c *CommandContext)
	// One liner help message for the command.
	Help string
	// More descriptive help message for the command.
	LongHelp string

	// DynamicCmd use for command name dynamic change.
	// NOTE: DynamicCmd are mutually exclusive with Name.
	DynamicCmd func(line string) []*Suggest

	// SubCommands. sub command.
	SubCommands []*Cmd

	// // Args command args.
	// Args []CmdArg

	// Command Flag Paramter
	// if return value not nil,it will be stored in flagsValue. Use GetFlagValue() get it.
	// NOTE: parent command's flags is also parsed and can use.
	NewFlags func(set *pflag.FlagSet) interface{}

	flagsValue   interface{}
	children     map[string]*Cmd
	dynamicCache []*Suggest     // dynamic cmd cache
	set          *pflag.FlagSet // flags cache
}

// AddCmd adds cmd as a subcommand.
func (c *Cmd) AddCmd(cmds ...*Cmd) {
	for _, cmd := range cmds {
		c.SubCommands = append(c.SubCommands, cmd)
	}
	sort.Sort(cmdSorter(c.SubCommands))
}

// DeleteCmd deletes cmd from subcommands.
func (c *Cmd) DeleteCmd(name string) {
	c.fixCmd()
	if c.children == nil {
		return
	}
	delete(c.children, name)
	for i := len(c.SubCommands) - 1; i >= 0; i-- {
		if c.SubCommands[i].Name == name {
			c.SubCommands = append(c.SubCommands[:i], c.SubCommands[i+1:]...)
			break
		}
	}
}

func (c *Cmd) fixCmd() {
	if len(c.SubCommands) == len(c.children) {
		return
	}
	c.children = make(map[string]*Cmd, len(c.SubCommands))
	for _, v := range c.SubCommands {
		if v.NewFlags != nil {
			v.set = pflag.NewFlagSet(v.Name, pflag.ContinueOnError)
			v.flagsValue = v.NewFlags(v.set)
		}

		// fix command alias
		for k := len(v.Aliases) - 1; k >= 0; k-- {
			// contains each other will cause unexpected behavior
			if strings.Contains(v.Aliases[k], v.Name) ||
				strings.Contains(v.Name, v.Aliases[k]) {
				v.Aliases = append(v.Aliases[:k], v.Aliases[k+1:]...)
				continue
			}
		}

		v.fixCmd()
		c.children[v.Name] = v
	}
	// has repeated name command
	if len(c.SubCommands) != len(c.children) {
		c.SubCommands = c.SubCommands[:0]
		for _, v := range c.children {
			c.SubCommands = append(c.SubCommands, v)
		}
	}
	sort.Sort(cmdSorter(c.SubCommands))
}

// Children returns the subcommands of c.
func (c *Cmd) Children() []*Cmd {
	c.fixCmd()
	return c.SubCommands
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
	if c.LongHelp != "" {
		p(c.LongHelp)
	} else if c.Help != "" {
		p(c.Help)
	} else if c.Name != "" {
		p(c.Name, "has no help")
	}
	if c.hasSubcommand() {
		p("Commands:")
		w := tabwriter.NewWriter(&b, 0, 4, 2, ' ', 0)
		for _, child := range c.SubCommands {
			fmt.Fprintf(w, "\t%s\t\t\t%s\n", child.Name, child.Help)
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
		if c.Name == name {
			return true
		}
	}
	for _, v := range c.Aliases {
		if v == name {
			return true
		}
	}
	return false
}

// findChildCmd returns the subcommand with matching name or alias.
func (c *Cmd) findChildCmd(name string) *Cmd {

	// find perfect matches first
	for _, cmd := range c.SubCommands {
		if cmd.isCmd(name) {
			return cmd
		}
	}

	return nil
}

// ParseInput parse input,and check valid
func (c *Cmd) ParseInput(line string, reset bool) (cmds []*Cmd, args []string, err error) {
	fields := strings.Fields(line)
	father := c
	for k, arg := range fields {
		if cmd := father.findChildCmd(arg); cmd != nil {
			cmds = append(cmds, cmd)
			father = cmd
			continue
		}
		fields = fields[k:]
		break
	}
	set := pflag.NewFlagSet("root", pflag.ContinueOnError)
	if c.set != nil {
		set.AddFlagSet(c.set)
	}
	for k := len(cmds) - 1; k >= 0; k-- {
		cmd := cmds[k]
		if cmd.set == nil {
			continue
		}
		set.AddFlagSet(cmd.set)
		// tmp := pflag.NewFlagSet(cmd.Name, pflag.ContinueOnError)
		// cmd.flagsValue = cmd.NewFlags(tmp)
		// set.AddFlagSet(tmp)
	}
	if reset {
		// reset default value
		set.VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
		})
	}
	// parse flag failed
	err = set.Parse(fields)
	// args
	args = set.Args()

	return
}

func (c *Cmd) buildCache(line string) {
	if c.DynamicCmd == nil || c.dynamicCache != nil {
		return
	}
	c.dynamicCache = c.DynamicCmd(line)
}

func (c *Cmd) suggest() *Suggest {
	return &Suggest{
		Text:        c.Name,
		Description: c.Help,
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
	matchCmd := func(name []rune, child *Cmd, s *Suggest) {
		if len(line) >= len(name) {
			if HasPrefix(line, name) {

				if s == nil {
					s = child.suggest()
				}
				if len(line) > len(name) {
					nextCmd = child
					goNext = true
				}
				suggest = append(suggest, s)
				offset = len(name)
			}
		} else {
			if fuzzyMatchRunes(name, line) {
				if s == nil {
					s = child.suggest()
				}
				suggest = append(suggest, s)
				offset = len(line)
				nextCmd = child
			}
		}
	}
	for _, child := range c.SubCommands {
		// command name
		if child.DynamicCmd != nil {
			child.buildCache(origLine)
			for _, v := range child.dynamicCache {
				matchCmd([]rune(v.Text), child, v)
			}
		} else if child.Name != "" {
			matchCmd([]rune(child.Name), child, nil)
		}
		// command alias
		// if nextCmd != child {
		for _, alias := range child.Aliases {
			matchCmd([]rune(alias), child, &Suggest{
				Text:        alias,
				Description: child.Help,
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
		suggest = FilterFuzzy(suggest, string(left), true)
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
func (c cmdSorter) Less(i, j int) bool { return c[i].Name < c[j].Name }
func (c cmdSorter) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
