# weworkbot

企微机器人推送脚本

## 使用

```
Usage of main:

  -d, --debug             print debug message.
  -f, --file string       set upload and send file,
                          support local path,url, base64;example: ./a.zip, http://a.com/a.zip, data:base,sdjkashd=
  -i, --image string      set upload and send image,
                          support local path,url, base64;example: ./a.png, http://a.com/a.jpeg, data:base,sdjkashd=
  -k, --key string        set wework bot key, support env X_WEWORK_BOT_KEY= for default (default "4cc96c9e-4a3f-4188-8df7-e6300c8fcca7")
  -m, --markdown string   send markdown v2 message
  -t, --text string       send text message
      --timeout string    set deadline to cancel, support 1s, 2m3s...

 Wewrok bot doc: https://developer.work.weixin.qq.com/document/path/91770.
```
