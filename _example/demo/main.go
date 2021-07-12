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
	"github.com/aggronmagi/promptx/internal/std"
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
	},
			fmt.Println("登陆成功")
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
			cmd.Stdin = std.Stdin   // os.Stdin
			cmd.Stdout = std.Stdout // os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			c.EnterRawMode()
		},
	},
	&promptx.Cmd{
		Name: "help",
		Help: "print help info",
		Func: func(c *promptx.CommandContext) {
			go func() {
				v := 0
				for {
					select {
					case <-time.After(time.Second):
						v++
						fmt.Fprintln(c.Stdout(), "xx---async", v)
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
		),
	)
	// set log writer
	log.SetOutput(p.Stdout())

	p.Run()
}
