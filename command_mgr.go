package promptx

import "fmt"

type commandManager struct {
	commander Commander
	sets      map[string]*CommandSetOptions
}

func newCommandManager(commander Commander) *commandManager {
	return &commandManager{
		commander: commander,
		sets:      make(map[string]*CommandSetOptions),
	}
}

func (m *commandManager) AddCommandSet(name string, cmds []*Cmd, opts ...CommandSetOption) {
	if len(cmds) < 0 {
		panic(fmt.Sprintf("commandset %s do not have any commad", name))
	}
	if _, ok := m.sets[name]; ok {
		panic(fmt.Sprintf("commandset %s register repeated", name))
	}
	set := NewCommandSetOptions(opts...)
	set.ApplyOption(
		WithName(name),
		WithCmds(cmds...),
	)
	m.sets[set.Name] = set
	if len(m.sets) == 1 {
		m.SwitchCommandSet(name)
	}
}

func (m *commandManager) SwitchCommandSet(name string, args ...interface{}) {
	set, ok := m.sets[name]
	if !ok {
		fmt.Printf("commandset %s not exists\n", name)
		return
	}
	m.commander.ResetHistoryFile(set.History)
	m.commander.SetCommandPreCheck(set.PreCheck)
	m.commander.ResetCommands(set.Cmds...)
	m.commander.SetPromptWords(set.Prompt...)
	if set.OnChange != nil {
		set.OnChange(m.commander.(Context), args)
	}
}
