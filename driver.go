package go_html_ui

import (
	"context"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Message struct {
	Action    string        `json:"action"`
	Arguments []interface{} `json:"arguments"`
}

type Attribute struct {
	Id        string
	Attribute string
	Value     interface{}
}

type Value struct {
	Id    string
	Value interface{}
}

type wsDriver struct {
	sendChan chan Message
	//attributeChan - reasponse for request getAttribute
	attributeChan chan Attribute
	//valueChan - reasponse for request getValue
	valueChan chan Value
	eventChan chan string
}

func (g wsDriver) Run(address string) {

	wsf := http.HandlerFunc(g.wsHandler)
	http.Handle("/gui", wsf)

	fs := http.FileServer(http.Dir("./gui"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(address, nil))
}

// web-socket handler
func (g wsDriver) wsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Print(err)
	}
	//new connection
	defer conn.CloseNow()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go g.sender(ctx, conn)
	var mes Message

	for {
		err = wsjson.Read(ctx, conn, &mes)
		if err != nil {
			log.Print(err)
			ctx.Done()
			return
		}

		//recieve mes
		switch mes.Action {
		case "event":
			if evt, ok := mes.Arguments[0].(string); ok {
				g.eventChan <- evt
			}
		case "getValue":
			val := Value{
				Id:    mes.Arguments[0].(string),
				Value: mes.Arguments[1],
			}
			g.valueChan <- val
		case "getAttribute":
			val := Attribute{
				Id:        mes.Arguments[0].(string),
				Attribute: mes.Arguments[1].(string),
				Value:     mes.Arguments[2],
			}
			g.attributeChan <- val
		case "quit":
			cancel()
			conn.CloseNow()
			return
		default: //unknown command
			time.Sleep(time.Millisecond * 50)
		}
	}
}

// send messages to ui
func (g wsDriver) sender(ctx context.Context, ws *websocket.Conn) {
	var mes Message
	for {
		select {
		case <-ctx.Done():
			return
		case mes = <-g.sendChan:
			_ = wsjson.Write(ctx, ws, mes) //send: %s", mes
		}
	}
}
