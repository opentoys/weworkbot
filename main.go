package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/opentoys/wecombot/webot"
	"github.com/urfave/cli/v3"
)

var version = "dev"

var config = &struct {
	Debug   bool
	Version bool
	Type    string
	Content string
	Key     string
	Timeout int64
}{}

func main() {
	cmd := &cli.Command{
		Name:                          "webot",
		Usage:                         "企业微信机器人 cli",
		UseShortOptionHandling:        true,
		ArgsUsage:                     "content",
		CustomRootCommandHelpTemplate: cli.RootCommandHelpTemplate + "\r\nWewrok bot doc: https://developer.work.weixin.qq.com/document/path/91770.",
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Hidden:  true,
				Action: func(ctx context.Context, c *cli.Command) (e error) {
					fmt.Println(version)
					os.Exit(0)
					return
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Destination: &config.Debug,
			},
			&cli.BoolFlag{
				Name:     "version",
				Aliases:  []string{"v"},
				OnlyOnce: true,
				Action: func(ctx context.Context, c *cli.Command, _ bool) (e error) {
					fmt.Println(version)
					os.Exit(0)
					return
				},
			},
			&cli.StringFlag{
				Name:        "key",
				Aliases:     []string{"k"},
				DefaultText: os.Getenv("X_WEWORK_BOT_KEY"),
				Usage:       "set wework bot key, support env X_WEWORK_BOT_KEY=",
				Destination: &config.Key,
			},
			&cli.StringFlag{
				Name:        "type",
				Aliases:     []string{"t"},
				Usage:       "send message type: text|markdown,md|file|image",
				DefaultText: "text",
				Destination: &config.Type,
			},
			&cli.Int64Flag{
				Name:        "timeout",
				Usage:       "send context timeout",
				Destination: &config.Timeout,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) (e error) {
			if cmd.NArg() == 0 {
				cli.DefaultShowRootCommandHelp(cmd)
				return
			}
			if config.Timeout > 0 {
				var cancel func()
				ctx, cancel = context.WithTimeout(ctx, time.Duration(config.Timeout)*time.Second)
				defer cancel()
			}
			config.Content = cmd.Args().First()
			e = send(ctx, config.Key, config.Type, config.Content)
			return
		},
	}
	e := cmd.Run(context.Background(), os.Args)
	if e != nil {
		panic(e)
	}
}

func send(_ context.Context, key, typ, content string) (e error) {
	bot := webot.New(key)
	switch typ {
	case "file":
		e = send_file(bot, content)
	case "image":
		e = send_image(bot, content)
	case "md", "markdown":
		_, e = bot.SendMarkdownV2(content)
	default:
		_, e = bot.SendText(content)
	}
	return
}

func send_file(bot *webot.Bot, content string) (e error) {
	buf, name, e := Path2bytes(content)
	if e != nil {
		panic(e)
	}
	resp, e := bot.UploadMedia(bytes.NewBuffer(buf), name, http.DetectContentType(buf))
	if e != nil {
		panic(e)
	}
	if !resp.IsOK() {
		fmt.Println(resp.ErrMsg)
	}
	_, e = bot.SendFile(resp.MediaID)
	return
}

func send_image(bot *webot.Bot, content string) (e error) {
	buf, _, e := Path2bytes(content)
	if e != nil {
		panic(e)
	}
	_, e = bot.SendImage(base64.RawStdEncoding.EncodeToString(buf), hex.EncodeToString(md5.New().Sum(buf)))
	return
}

func Path2bytes(s string) (buf []byte, name string, e error) {
	if s == "" {
		e = errors.New("input file error")
		return
	}
	if strings.HasPrefix(s, "http") {
		var resp *http.Response
		resp, e = http.DefaultClient.Get(s)
		if e != nil {
			return nil, "", e
		}
		buf, e = io.ReadAll(resp.Body)
		name = path.Base(resp.Request.URL.Path)
		return
	}

	if strings.HasPrefix(s, "data:") {
		var lst = strings.Split(s, ",")
		if len(lst) == 0 {
			e = errors.New("base64 error")
			return
		}
		buf, e = base64.RawStdEncoding.DecodeString(lst[1])
		if e != nil {
			buf, e = base64.RawURLEncoding.DecodeString(lst[1])
		}
		return
	}

	buf, e = os.ReadFile(s)
	name = filepath.Base(s)
	return
}
