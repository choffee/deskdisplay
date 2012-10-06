// A little web server for the desktop display
package main

import (
	// "code.google.com/p/gorilla/color"
	. "github.com/choffee/gofirmata"
	"html/template"
	"log"
	"net/http"
)

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
var led *RGBLED = NewRGBLED(9, 10, 11)
var status Status = Status{Arduino: false}
var board *Board

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "My Little LED page",
		Color:  led.HexString(),
		Status: status.String()}
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, p)
}

func colorHandler(w http.ResponseWriter, r *http.Request) {
	err := led.QuickColor(r.FormValue("color"))
	if err != nil {
		log.Println(err)
		log.Printf("Bad color %s\n", r.FormValue("color"))
	}
	led.SendColor(board)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	status.Arduino = false
	var err error
	board, err = NewBoard("/dev/ttyUSB0", 57600)
	if err != nil {
		log.Println("Could not connect to Arduino, do you have the right port?")
	} else {
		status.Arduino = true
	}
	board.Debug = 99
	// Wait for the version message before we start anything
	log.Printf("Msg: %v", <-board.Reader)
	// For now just print the msg we receive
	go func() {
		for {
			log.Printf("Msg: %v", <-board.Reader)
		}
	}()
	log.Printf("Ver: %d.%d", board.Version()["major"], board.Version()["minor"])
	led.SetupPins(board)
	led.Invert = true
	led.QuickColor("red")
	led.SendColor(board)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/setColor", colorHandler)
	http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("js"))))
	http.ListenAndServe(":8080", nil)
}
