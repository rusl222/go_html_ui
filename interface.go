package go_html_ui

import (
	"fmt"
)

type gui struct {
	driver    wsDriver
	listeners map[string][]func()
}

type Gui interface {
	Run(address string)
	AddEventListener(event string, id string, callback func())
	Value(id string) interface{}
	SetValue(id string, value interface{})
	Attribute(id string, attribute string) interface{}
	SetAttribute(id string, attribute string, value interface{})
}

func New() Gui {
	g := gui{}
	g.driver.sendChan = make(chan Message)
	g.driver.valueChan = make(chan Value)
	g.driver.attributeChan = make(chan Attribute)
	g.driver.eventChan = make(chan string)

	g.listeners = make(map[string][]func())
	return g
}

func (g gui) Run(address string) {
	go g.driver.Run(address)

	// events handler
	for event := range g.driver.eventChan {
		if fs, ok := g.listeners[event]; ok {
			for _, f := range fs {
				go f() // execute callbacks
			}
		}
	}
}

func (g gui) AddEventListener(event string, id string, callback func()) {
	idEvent := fmt.Sprintf("%s_%s", id, event)
	g.listeners[idEvent] = append(g.listeners[idEvent], callback)
}

// gets the value by tag id.
func (g gui) Value(id string) interface{} {
	arguments := make([]interface{}, 1)
	arguments[0] = id
	g.driver.sendChan <- Message{Action: "getValue", Arguments: arguments}

	value := <-g.driver.valueChan
	return value.Value
}

func (g gui) SetValue(id string, value interface{}) {
	arguments := make([]interface{}, 2)
	arguments[0] = id
	arguments[1] = value
	g.driver.sendChan <- Message{Action: "setValue", Arguments: arguments}
}

// gets the attribute by tag id.
func (g gui) Attribute(id string, attr string) interface{} {
	arguments := make([]interface{}, 2)
	arguments[0] = id
	arguments[1] = attr
	g.driver.sendChan <- Message{Action: "getAttribute", Arguments: arguments}

	value := <-g.driver.attributeChan
	return value.Value
}

func (g gui) SetAttribute(id string, attr string, value interface{}) {
	arguments := make([]interface{}, 3)
	arguments[0] = id
	arguments[1] = attr
	arguments[2] = value
	g.driver.sendChan <- Message{Action: "setAttribute", Arguments: arguments}
}
