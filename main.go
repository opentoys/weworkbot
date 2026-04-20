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
	"github.com/spf13/pflag"
)

var version = "dev"

var config = &struct {
	Debug    bool
	Version  bool
	Type     string
	Content  string
	Key      string
	File     string
	Image    string
	Markdown string
	Text     string
	Timeout  string
}{}

func init() {
	pflag.BoolVarP(&config.Debug, "debug", "d", false, "print debug message.")
	pflag.StringVarP(&config.Key, "key", "k", os.Getenv("X_WEWORK_BOT_KEY"), "set wework bot key, support env X_WEWORK_BOT_KEY= for default")
	pflag.StringVarP(&config.File, "file", "f", "", "set upload and send file, \r\nsupport local path,url, base64;example: ./a.zip, http://a.com/a.zip, data:base,sdjkashd=")
	pflag.StringVarP(&config.Image, "image", "i", "", "set upload and send image, \r\nsupport local path,url, base64;example: ./a.png, http://a.com/a.jpeg, data:base,sdjkashd=")
	pflag.StringVar(&config.Timeout, "timeout", "", "set deadline to cancel, support 1s, 2m3s...")
	pflag.StringVarP(&config.Markdown, "markdown", "m", "", "send markdown v2 message")
	pflag.StringVarP(&config.Text, "text", "t", "", "send text message")
}

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}
	switch os.Args[1] {
	case "version", "v":
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println(version)
		return
	case "help", "h", "-help", "--help", "-h", "--h":
		printHelp()
		return
	}
	pflag.Parse()

	if config.Key == "" {
		fmt.Println("key must be set")
		printHelp()
		return
	}

	bot := webot.New(config.Key)

	if config.Timeout != "" {
		ts, e := time.ParseDuration(config.Timeout)
		if e != nil {
			panic(e)
		}
		ctx, cancel := context.WithTimeout(context.Background(), ts)
		defer cancel()
		bot = bot.WithContext(ctx)
	}

	if config.Text != "" {
		_, e := bot.SendText(config.Text)
		if e != nil {
			panic(e)
		}
		return
	}

	if config.Markdown != "" {
		_, e := bot.SendMarkdownV2(config.Markdown)
		if e != nil {
			panic(e)
		}
		return
	}

	if config.File != "" {
		buf, name, e := Path2bytes(config.File)
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
		if e != nil {
			panic(e)
		}
		return
	}

	if config.Image != "" {
		buf, _, e := Path2bytes(config.Image)
		if e != nil {
			panic(e)
		}
		_, e = bot.SendImage(base64.RawStdEncoding.EncodeToString(buf), hex.EncodeToString(md5.New().Sum(buf)))
		if e != nil {
			panic(e)
		}
		return
	}
}

func printHelp() {
	fmt.Printf("Usage of %s:\n\n", os.Args[0])
	pflag.PrintDefaults()
	fmt.Println("\r\n Wewrok bot doc: https://developer.work.weixin.qq.com/document/path/91770.")
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
