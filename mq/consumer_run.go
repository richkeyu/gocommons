package mq

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	application "github.com/richkeyu/gocommons/dispatch"
	"github.com/richkeyu/gocommons/plog"

	"github.com/urfave/cli/v2"
)

var consumerMap []*ConsumerOption
var commands []*cli.Command
var defaultAction cli.ActionFunc

type Run struct {
}

func (r Run) AddCliCommand(c *cli.Command) {
	if c.Name == `` {
		defaultAction = c.Action
	}
	commands = append(commands, c)
}

func (t Run) AddConsumer(c *ConsumerOption) {
	if c != nil {
		consumerMap = append(consumerMap, c)
	}
}

func Go(Register func(Run)) {

	Register(Run{})
	var consumerCli []*cli.Command
	if len(consumerMap) > 0 {
		for _, v := range consumerMap {
			tmpC := v
			consumerCli = append(consumerCli, &cli.Command{
				Name:  v.Name,
				Usage: v.Usage,
				Action: func(c *cli.Context) error {
					consumerHandle := NewConsumer(tmpC)
					err := consumerHandle.start()
					return err
				},
			})
		}
	}

	commands = append(commands, &cli.Command{
		Name:        "consumer",
		Aliases:     []string{"exec"},
		Usage:       "exec consumer task",
		Subcommands: consumerCli,
		Action: func(context *cli.Context) error {
			l := len(consumerMap)
			if l == 0 {
				return fmt.Errorf("consumer task not found")
			}
			wg := sync.WaitGroup{}
			wg.Add(l)
			for _, v := range consumerMap {
				go func(v *ConsumerOption) {
					for {
						c := NewConsumer(v)
						err := c.start()
						if err == nil {
							break
						}
					}
					wg.Done()
				}(v)
			}
			wg.Wait()
			return nil
		},
	})

	app := cli.NewApp()

	helpAction := app.Action
	app.Action = func(context *cli.Context) error {
		if len(os.Args) > 1 {
			return helpAction(context)
		}

		if defaultAction == nil {
			defaultAction = func(c *cli.Context) error {
				fmt.Println("Please add parameters, eg: consumer/help")
				return nil
			}
		}
		return defaultAction(context)
	}
	app.Commands = commands

	sort.Sort(cli.CommandsByName(app.Commands)) // 通过命令函数来排序，在help中进行展示
	err := application.Run(func() error {
		err := app.RunContext(context.Background(), os.Args)
		return err
	})
	if err != nil {
		plog.Error(nil, err.Error())
	}

	plog.Debug(nil, `Wait 1 second`)
	time.Sleep(time.Second * 1)
	plog.Debug(nil, `Exit!`)
}
