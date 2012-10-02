// A little web server for the desktop display
package main

import (
	"code.google.com/p/gorilla/color"
	. "github.com/choffee/gofirmata"
	"html/template"
	"log"
	"net/http"
)

type LED struct {
	color color.Hex
}

type Page struct {
	Title  string
	Color  string
	Status string
}

type Status struct {
	Arduino bool
}

func (s *Status) String() string {
	if s.Arduino {
		return "Arduino connected"
	}
	return "Arduino not connected"
}

// Set the LED to red at start
// Need to do it here so the first run of the
// template works
var led LED = LED{color: color.Hex("FF0000")}
var status Status = Status{Arduino: false}

func (led *LED) SetRGB(r, g, b byte) {
	led.color = color.RGBToHex(r, g, b)
}
func (led *LED) SetString(c color.Hex) {
	led.color = c
}
func (led *LED) Color() color.Hex {
	return led.color
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "My Little LED page",
		Color:  string(led.Color()),
		Status: status.String()}
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, p)
}

func colorHandler(w http.ResponseWriter, r *http.Request) {
	led.SetString(color.Hex(r.FormValue("color")))
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	led.SetString(color.Hex("FF0000"))
	status.Arduino = false
	board, err := NewBoard("/dev/ttyUSB1", 57600)
	if err != nil {
		log.Println("Could not connect to Arduino, do you have the right port?")
	} else {
		status.Arduino = true
	}
	board.SetPinMode(9, MODE_PWM)
	board.SetPinMode(10, MODE_PWM)
	board.SetPinMode(11, MODE_PWM)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/setColor", colorHandler)
	http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("js"))))
	http.ListenAndServe(":8080", nil)
}
