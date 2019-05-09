package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

type session struct {
	ws      *websocket.Conn
	rl      *readline.Instance
	errChan chan error
}

func connect(url, origin, auth string, rlConf *readline.Config) error {
	headers := make(http.Header)
	headers.Add("Origin", origin)
	headers.Add("Authorization", auth)

	ws, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		return err
	}

	rl, err := readline.NewEx(rlConf)
	if err != nil {
		return err
	}
	defer rl.Close()

	sess := &session{
		ws:      ws,
		rl:      rl,
		errChan: make(chan error),
	}

	go sess.readConsole()
	go sess.readWebsocket()

	return <-sess.errChan
}

func (s *session) readConsole() {
	for {
		line, err := s.rl.Readline()
		if err != nil {
			s.errChan <- err
			return
		}

		if strings.HasPrefix(line, "@audio:") {
			tokens := strings.Split(line, ":")
			if len(tokens) != 2 {
				s.printWarning("Invalid command iput: " + line)
				continue
			}
			filePath := tokens[1]
			audioData, err := ioutil.ReadFile(filePath)
			if err != nil {
				s.printWarning("Fail to read file '" + filePath + "' : " + err.Error())
				continue
			}

			err = s.ws.WriteMessage(websocket.BinaryMessage, audioData)
			if err != nil {
				s.errChan <- err
				return
			}
		} else {
			err = s.ws.WriteMessage(websocket.TextMessage, []byte(line))
			if err != nil {
				s.errChan <- err
				return
			}
		}
	}
}

func (s *session) printWarning(text string) {
	rxSprintf := color.New(color.FgYellow).SprintfFunc()
	fmt.Fprint(s.rl.Stdout(), rxSprintf("< %s\n", text))
}

func bytesToFormattedHex(bytes []byte) string {
	text := hex.EncodeToString(bytes)
	return regexp.MustCompile("(..)").ReplaceAllString(text, "$1 ")
}

func (s *session) readWebsocket() {
	rxSprintf := color.New(color.FgGreen).SprintfFunc()

	for {
		msgType, buf, err := s.ws.ReadMessage()
		if err != nil {
			s.errChan <- err
			return
		}

		var text string
		switch msgType {
		case websocket.TextMessage:
			text = string(buf)
		case websocket.BinaryMessage:
			text = bytesToFormattedHex(buf)
		default:
			s.errChan <- fmt.Errorf("unknown websocket frame type: %d", msgType)
			return
		}

		fmt.Fprint(s.rl.Stdout(), rxSprintf("< %s\n", text))
	}
}
