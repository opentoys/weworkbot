# weworkbot
企微机器人推送脚本

Contains only standard libraries!

## 使用
```
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

Wewrok bot doc: https://developer.work.weixin.qq.com/document/path/91770
```