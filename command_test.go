package promptx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newCmd(name string, help string) *Cmd {
	return &Cmd{
		name: name,
		help: help,
	}
}

func TestAddCommand(t *testing.T) {
	cmd := newCmd("root", "")
	assert.Equal(t, len(cmd.subCommands), 0, "should be empty")
	cmd.SubCommands(newCmd("child", ""))
	assert.Equal(t, len(cmd.subCommands), 1, "should include one child command")
}

func TestDeleteCommand(t *testing.T) {
	cmd := newCmd("root", "")
	cmd.SubCommands(newCmd("child", ""))
	assert.Equal(t, len(cmd.subCommands), 1, "should include one child command")
	cmd.DeleteSubCommand("child")
	assert.Equal(t, len(cmd.subCommands), 0, "should be empty")
}

func TestParseCmdInput(t *testing.T) {
	cmd := newCmd("root", "")
	cmd.SubCommands(newCmd("child1", ""))
	cmd.SubCommands(newCmd("child2", ""))
	// res, err := cmd.FindCmd([]string{"child1"})
	// if err != nil {
	// 	t.Fatal("finding should work")
	// }
	// assert.Equal(t, res.Name, "child1")

	// res, err = cmd.FindCmd([]string{"child2"})
	// if err != nil {
	// 	t.Fatal("finding should work")
	// }
	// assert.Equal(t, res.Name, "child2")

	// res, err = cmd.FindCmd([]string{"child3"})
	// if err == nil {
	// 	t.Fatal("should not find this child!")
	// }
	// assert.Nil(t, res)
}

func TestHelpText(t *testing.T) {
	cmd := newCmd("root", "help for root command")
	cmd.SubCommands(newCmd("child1", "help for child1 command"))
	cmd.SubCommands(newCmd("child2", "help for child2 command"))
	res := cmd.HelpText()
	expected := "\nhelp for root command\n\nCommands:\n  child1      help for child1 command\n  child2      help for child2 command\n\n"
	assert.Equal(t, res, expected)
}

func TestChildrenSortedAlphabetically(t *testing.T) {
	cmd := newCmd("root", "help for root command")
	cmd.SubCommands(newCmd("child2", "help for child1 command"))
	cmd.SubCommands(newCmd("child1", "help for child2 command"))
	children := cmd.subCommands
	assert.Equal(t, children[0].name, "child1", "must be first")
	assert.Equal(t, children[1].name, "child2", "must be second")
}
