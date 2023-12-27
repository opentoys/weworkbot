package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

const host = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
const uploadhost = "https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key="

type Bot struct {
	ctx context.Context
	key string
}

func (s *Bot) post(body []byte) (msg response, e error) {
	if s.key == "" {
		return
	}
	req, e := http.NewRequestWithContext(s.ctx, "POST", host+s.key, bytes.NewBuffer(body))
	if e != nil {
		return
	}
	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return
	}
	buf, e := io.ReadAll(resp.Body)
	if e != nil {
		return
	}

	e = json.Unmarshal(buf, &msg)
	return
}

func (s *Bot) Text(content string, user string, mobile string) (e error) {
	var body = map[string]interface{}{
		"content": content,
	}

	if user != "" {
		body["mentioned_list"] = strings.Split(user, "|")
	}

	if mobile != "" {
		body["mentioned_mobile_list"] = strings.Split(mobile, "|")
	}

	msg, e := s.post(JSONBytes(map[string]interface{}{
		"msgtype": "text",
		"text":    body,
	}))
	if e != nil {
		return
	}
	if msg.Errcode != 0 {
		return errors.New(msg.Errmsg)
	}
	return
}

func (s *Bot) Markdown(content string) (e error) {
	msg, e := s.post(JSONBytes(map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"content": content,
		},
	}))
	if e != nil {
		return
	}
	if msg.Errcode != 0 {
		return errors.New(msg.Errmsg)
	}
	return
}

func (s *Bot) media(typ string, name string, buf []byte) (id string, e error) {
	if s.key == "" {
		return
	}
	var bs bytes.Buffer
	w := multipart.NewWriter(&bs)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="media"; filename="%s"`, name))
	h.Set("Content-Type", http.DetectContentType(buf))
	iow, e := w.CreatePart(h)
	if e != nil {
		return
	}
	iow.Write(buf)
	w.Close()
	req, e := http.NewRequestWithContext(s.ctx, "POST", uploadhost+s.key+"&type="+typ, &bs)
	if e != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return
	}
	buf, e = io.ReadAll(resp.Body)
	if e != nil {
		return
	}
	var msg response
	e = json.Unmarshal(buf, &msg)
	id = msg.MediaId
	return
}

func (s *Bot) File(name string, buf []byte) (e error) {
	mid, e := s.media("file", name, buf)
	if e != nil {
		return
	}
	msg, e := s.post(JSONBytes(map[string]interface{}{
		"msgtype": "file",
		"file": map[string]interface{}{
			"media_id": mid,
		},
	}))
	if e != nil {
		return
	}
	if msg.Errcode != 0 {
		return errors.New(msg.Errmsg)
	}
	return
}

func (s *Bot) Voice(name string, buf []byte) (e error) {
	mid, e := s.media("voice", name, buf)
	if e != nil {
		return
	}
	msg, e := s.post(JSONBytes(map[string]interface{}{
		"msgtype": "voice",
		"voice": map[string]interface{}{
			"media_id": mid,
		},
	}))
	if e != nil {
		return
	}
	if msg.Errcode != 0 {
		return errors.New(msg.Errmsg)
	}
	return
}
func (s *Bot) Image(buf []byte) (e error) {
	md := md5.New()
	var bs64 = base64.StdEncoding.EncodeToString(buf)
	_, e = md.Write(buf)
	if e != nil {
		return
	}
	msg, e := s.post(JSONBytes(map[string]interface{}{
		"msgtype": "image",
		"image": map[string]interface{}{
			"base64": bs64,
			"md5":    hex.EncodeToString(md.Sum(nil)),
		},
	}))
	if e != nil {
		return
	}
	if msg.Errcode != 0 {
		return errors.New(msg.Errmsg)
	}
	return
}
func (s *Bot) News(data string) (e error) {
	return
}
func (s *Bot) Card(data string) (e error) {
	return
}

type response struct {
	Errcode int32  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Type    string `json:"type"`
	MediaId string `json:"media_id"`
}
