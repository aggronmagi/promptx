package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/aggronmagi/promptx"
	"github.com/aggronmagi/promptx/internal/debug"
	"github.com/spf13/pflag"
)

// Function constructor - constructs new function for listing given directory
func listFiles(path string) func(string) []*promptx.Suggest {
	return func(line string) []*promptx.Suggest {
		debug.Println("call dynamic completion", line)
		names := make([]*promptx.Suggest, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, &promptx.Suggest{
				Text:        f.Name(),
				Description: path,
			})
		}
		return names
	}
}

var cmds = []*promptx.Cmd{
	&promptx.Cmd{
		Name: "mode",
		Help: "show mode or modify mode",
		SubCommands: []*promptx.Cmd{
			&promptx.Cmd{
				Name: "vi",
				Help: "修改vi模式",
				Func: func(c *promptx.CommandContext) {
					log.Println("sorry! not support vim mode now!")
					return
				},
			},
			&promptx.Cmd{
				Name: "emacs",
				Help: "use emacs edit mode",
				Func: func(c *promptx.CommandContext) {
					log.Println("use emacs mode now!")
					return
				},
			},
		},
		Func: func(c *promptx.CommandContext) {
			log.Println("current is emacs mode")
		},
	},
	&promptx.Cmd{
		Name: "login",
		Help: "登陆游戏",
		Func: func(c *promptx.CommandContext) {
			c.Select("选择登陆的服务器", []string{
				"开发服",
				"测试服",
				"体验服",
			})
			c.Println("xxx -  x")
			c.Select("xxxxx服务器", []string{
				"开发服x2",
				"测试服x2",
				"体验服x2",
			})

			fmt.Println("登陆成功")
			return
		},
		SubCommands: []*promptx.Cmd{
			&promptx.Cmd{
				Name: "hello",
				Help: "say hello",
			},
			&promptx.Cmd{
				Name: "xxx",
				Help: "xxx xxx",
			},
		},
	},
	&promptx.Cmd{
		Name: "login2",
		Help: "测试相似命令(test similar command)",
		Func: func(c *promptx.CommandContext) {
			c.Select(
				"选择登陆的服务器",
				[]string{
					"开发服",
					"测试服",
					"体验服",
				},
				promptx.WithSelectOptionShowHelpText(true),
			)
			c.Select(
				"选择登陆的xxxx服务器",
				[]string{
					"开发服x",
					"测试服x",
					"体验服x",
				},
				promptx.WithSelectOptionShowHelpText(true),
			)

			fmt.Println("登陆成功")
			c.WPrint(&promptx.AskWord, promptx.WordRed(" Red "),
				promptx.WordCyan(" Cyan "),
				promptx.WordBlue(" Blue "),
				promptx.WordGreen(" Green "),
				promptx.WordPurple(" Purple "),
				promptx.WordTurquoise(" Turquoise "),
				promptx.WordWhite(" White "),
				promptx.WordYellow(" Yellow "),
				&promptx.NewLineWord,
			)
			return
		},
	},
	&promptx.Cmd{
		Name: "say",
		Help: "say some words",
		SubCommands: []*promptx.Cmd{
			&promptx.Cmd{
				Name:       "files",
				Help:       "select files for operetion",
				DynamicCmd: listFiles("./"),
				SubCommands: []*promptx.Cmd{
					&promptx.Cmd{
						Name: "with",
						Help: "with command",
						SubCommands: []*promptx.Cmd{
							&promptx.Cmd{
								Name: "following",
								Help: "ssss sub command",
							},
							&promptx.Cmd{
								Name: "items",
								Help: "items command tip",
							},
						},
					},
				},
			},
			&promptx.Cmd{
				Name: "hello",
				Help: "say hello",
			},
			&promptx.Cmd{
				Name: "bye",
				Help: "bye bye",
			},
		},
	},
	&promptx.Cmd{
		Name: "setprompt",
		Help: "set prompt string",
		Func: func(c *promptx.CommandContext) {
			prompt := ""
			if len(c.Args) > 0 {
				prompt = c.Args[0]
			} else {
				var eof error
				prompt, eof = c.Input("input prompt")
				if eof != nil {
					return
				}
			}
			c.SetPrompt(prompt)
		},
	},
	&promptx.Cmd{
		Name: "prompt",
		Help: "update prompt words",
		Func: func(c *promptx.CommandContext) {
			prompt := ""
			c.Println(c.Args)
			if len(c.Args) > 0 {
				prompt = c.Args[0]
			} else {
				var eof error
				prompt, eof = c.Input("input prompt")
				if eof != nil {
					return
				}
			}
			c.Println("---", prompt)

			c.SetPromptWords(&promptx.AskWord, &promptx.Word{
				Text:      prompt + " ",
				TextColor: promptx.Red,
			})
			// NOTE: refresh
			c.Println("xxx")
			c.Println("xxx wait 5 second")
			time.Sleep(time.Second * 5)
			c.Println("fiish")
		},
	},
	&promptx.Cmd{
		Name:    "setpassword",
		Aliases: []string{"password"},
		Help:    "modify password",
	},
	&promptx.Cmd{
		Name: "bye",
		Help: "exit bye bye",
	},
	&promptx.Cmd{
		Name: "ls",
		Help: "ls linux command",
		Func: func(c *promptx.CommandContext) {
			c.ExitRawMode()
			cmd := exec.Command("ls", "-alh")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			c.EnterRawMode()
		},
	},
	&promptx.Cmd{
		Name: "bash",
		Help: "enter linux bash",
		Func: func(c *promptx.CommandContext) {
			err := c.ExitRawMode()
			c.Println("exit rawmode:", err)
			cmd := exec.Command("bash")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			c.Println("run result:", err)
			c.EnterRawMode()
		},
	},
	&promptx.Cmd{
		Name: "help",
		Help: "print help info",
		Func: func(c *promptx.CommandContext) {
			input, err := c.Input("input limit: ",
				promptx.WithInputOptionTipText("tip tip tip"),
				promptx.WithInputOptionValidFunc(func(d *promptx.Document) error {
					_, err := strconv.Atoi(d.Text)
					return err
				}),
			)
			if err != nil {
				return
			}
			limit, _ := strconv.Atoi(input)
			go func() {
				v := 0
				for {
					select {
					case <-time.After(time.Second):
						v++
						fmt.Fprintln(c.Stdout(), "xx---async", v)
						if v >= limit {
							return
						}
					}
				}
			}()
			log.Println("start async print")
		},
	},
	&promptx.Cmd{
		Name: "go",
		Help: "run go command",
		SubCommands: []*promptx.Cmd{
			&promptx.Cmd{
				Name: "build",
				Help: "go command",
				NewFlags: func(set *pflag.FlagSet) interface{} {
					set.BoolP("xxxx", "x", false, "test xxxxx")
					flag.String("o", "", "输出文件名")
					o := flag.Lookup("o")
					set.AddGoFlag(o)
					set.StringP("tip", "t", "", "tip tip tip")
					set.IntP("ssss", "v", 0, "fffffff - int")
					return nil
				},
			},
			&promptx.Cmd{
				Name: "install",
				Help: "go command",
			},
			&promptx.Cmd{
				Name: "test",
				Help: "go command",
			},
		},
	},
	&promptx.Cmd{
		Name: "sleep",
		Help: "sleep some second",
		Func: func(c *promptx.CommandContext) {
			result, err := c.Input("input some second:", promptx.WithInputOptionValidFunc(
				func(d *promptx.Document) error {
					_, err := strconv.Atoi(d.Text)
					return err
				},
			))
			if err != nil {
				fmt.Fprint(c.Stdout(), "user cancel!")
				return
			}
			sec, _ := strconv.Atoi(result)
			if sec > 5 {
				sec = 5
			}
			log.Println(c.Args)
			log.Println("wait for ", time.Duration(sec)*time.Second)
			time.Sleep(time.Second * time.Duration(sec))
			log.Println("finish")
			return
		},
	},
}

func main() {
	// new promptx
	p := promptx.NewPromptx(
		promptx.WithCommonOpions(
			// install commands
			promptx.WithCommonOptionCommands(cmds...),
			// promptx.WithCommonOptionTipText("tip tip tip"),
			promptx.WithCommonOptionHistory(".history"),
		),
	)
	// set log writer
	log.SetOutput(p.Stdout())

	p.ExecCommand([]string{"mode"})

	p.Input("input xx:")
	p.Select("select xx:", []string{
		"x1",
		"x2",
		"x3",
	})

	p.Run()
}
