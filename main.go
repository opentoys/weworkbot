package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

const verison = "1.0.1"
const help = `
The commands are:
    -h,-help       print help info 
    -v,-verison    print verison
    -k,-key        set wework bot key, support env X_WEWORK_BOT_KEY= for default
    -f,-file       set upload and send file, 
                        support local path,url, base64;example: ./a.zip, http://a.com/a.zip, data:base,sdjkashd=
    -i,-image      set upload and send image, 
                        support local path,url, base64;example: ./a.png, http://a.com/a.jpeg, data:base,sdjkashd=
    -md,-markdown  send markdown message
    -t,-text       send text message
    -deadline      set deadline to cancel, support 1s, 2m3s...
    -debug

Wewrok bot doc: https://developer.work.weixin.qq.com/document/path/91770`

var key = ""
var debug = false

func println(args ...any) {
	for _, v := range args {
		fmt.Print(v)
	}
	fmt.Println("")
}

func checkprint(e error) {
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println("success")
}

func main() {
	key = os.Getenv("X_WEWORK_BOT_KEY")
	var data = ParseArgs(os.Args[1:])

	if _, ok := data["debug"]; ok {
		debug = true
	}

	if debug {
		println(data)
	}

	if key != "" {
		println("default key: ", key)
	}

	var ctx = context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var send struct {
		typ     string
		content string
	}

	for k, v := range data {
		switch k {
		case "k", "key":
			key = v
		case "v", "version":
			fmt.Println(verison)
			return
		case "h", "help":
			fmt.Println(help)
			return
		case "t", "text":
			send.typ = "text"
			send.content = v
		case "md", "markdown":
			send.typ = "markdown"
			send.content = v
		case "f", "file":
			send.typ = "file"
			send.content = v
		case "i", "image":
			send.typ = "image"
			send.content = v
		case "deadline":
			du, e := time.ParseDuration(v)
			if e == nil {
				ctx, cancel = context.WithTimeout(ctx, du)
			}
		}
	}
	defer cancel()

	if key == "" {
		fmt.Println("key must be set")
		fmt.Println(help)
		return
	}

	if send.content == "" {
		fmt.Println("send content must be set")
		fmt.Println(help)
		return
	}

	var bot = &Bot{ctx: ctx, key: key}
	switch send.typ {
	case "text":
		fmt.Println("sending text messgae")
		checkprint(bot.Text(send.content, "", ""))
	case "markdown":
		fmt.Println("sending markdown messgae")
		checkprint(bot.Markdown(send.content))
	case "image":
		buf, _, e := Path2bytes(send.content)
		if e != nil {
			fmt.Println(e)
			return
		}
		fmt.Println("sending image")
		checkprint(bot.Image(buf))
	case "file":
		buf, name, e := Path2bytes(send.content)
		if e != nil {
			fmt.Println(e)
			return
		}
		if name == "" {
			name = "unknow"
		}
		fmt.Println("sending file")
		checkprint(bot.File(name, buf))
	}
}
