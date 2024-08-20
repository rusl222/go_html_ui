package main

import (
	"fmt"
	"math/rand"
	"os/exec"

	"github.com/rusl222/go_html_ui"
)

func main() {

	exec.Command("explorer", "http://127.0.0.1:8080").Run()

	//1 - make
	ui := go_html_ui.New()

	//2 - add listeners
	ui.AddEventListener("click", "button1", func() {
		name := ui.Value("edit1")
		ui.SetValue("label1", fmt.Sprintf("Привет %s", name))
	})

	ui.AddEventListener("click", "button2", func() {
		r := int64(rand.Intn(16))
		g := int64(rand.Intn(16))
		b := int64(rand.Intn(16))
		ui.SetAttribute("label1", "style", fmt.Sprintf("color: #%0X%0X%0X", r, g, b))
	})

	//3 - run
	ui.Run(":8080")
}
