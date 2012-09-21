// Small app to show working gofirmata working
// Should put something on the screen to start with.
//

package main

import (
	"bytes"
	"fmt"
	"github.com/choffee/gofirmata"
	"math/rand"
	"time"
)

const (
	DISPMOVE byte = 3
)

// Display
type disp struct {
	Board   firmata.Board
	Addr    byte
	width   byte
	height  byte
	Content [][]byte
	cursorR byte
	cursorC byte
}

func (disp *disp) Write(msg []byte) {
	newmsg := make([]byte, len(msg)+1)
	newmsg[0] = 0
	for l, v := range msg {
		newmsg[l+1] = v
	}
	disp.Board.I2CWrite(disp.Addr, firmata.I2C_MODE_WRITE, newmsg)
}

func (disp *disp) SetSize(c, r byte) {
	disp.width = c
	disp.height = r
	// Create a new blank copy of the content
	for l := byte(0); l < r; l++ {
		disp.Content = append(disp.Content, make([]byte, c))
	}
}

func (disp *disp) clear() {
	msg := []byte{12}
	disp.Write(msg)
	for rk, _ := range disp.Content {
		for ck, _ := range disp.Content[rk] {
			disp.Content[rk][ck] = 0
		}
	}
	disp.moveTo(0, 0)
}

func (disp *disp) moveTo(c, r byte) {
	msg := []byte{DISPMOVE, c, r}
	disp.Write(msg)
	disp.cursorR = r
	disp.cursorC = c
}

// Update the display from a new array
func (disp *disp) updateScreen(newscreen [][]byte) {
	// For now just do it via brute force
	for rk, _ := range newscreen {
		disp.putText(string(newscreen[rk]), byte(rk), 0)
	}
}

func (disp *disp) write(s string) {
	msg := []byte(s)
	disp.Write(msg)
	for _, v := range s {
		disp.Content[disp.cursorR][disp.cursorC] = byte(v)
		disp.cursorC++
	}
}

func (disp *disp) putText(s string, x, y byte) {
	disp.moveTo(x, y)
	disp.write(s)
}

// Keep updating the time on the board
func (disp *disp) showTime() {
	for {
		now := time.Now()
		clock := now.Format(time.Kitchen)
		disp.putText("Time "+clock, 1, 5)
		time.Sleep(1000 * time.Millisecond)
	}
}

func update_bubbles(screen *[][]byte) {
	for rk, _ := range *screen {
		// First replace all the bubbles with the next stage
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("*"), []byte(" "), -1)
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("Q"), []byte("*"), -1)
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("0"), []byte("Q"), -1)
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("O"), []byte("0"), -1)
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("o"), []byte("O"), -1)
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("."), []byte("o"), -1)
		// Now move them up if they can
		if rk > 0 { // not the top row
			for ck, _ := range (*screen)[rk] {
				if bytes.Contains([]byte(".oO0Q*"), []byte{(*screen)[rk][ck]}) {
					(*screen)[rk][ck] = []byte(" ")[0]
					(*screen)[rk-1][ck] = (*screen)[rk][ck]
				}
			}
		}
	}
}

// A better idea may be to add a channel to the function above for doing this.
func add_bubbles(screen *[][]byte) {
	for k, _ := range (*screen)[len(*screen)] {
		if rand.Intn(10) > 4 {
			(*screen)[len(*screen)][k] = []byte(".")[0]
		}
	}
}

func main() {
	board := new(firmata.Board)
	board.Device = "/dev/ttyUSB0"
	board.Baud = 57600
	board.Debug = 2
	err := board.Setup()
	go func() {
		for msg := range *board.Reader {
			fmt.Println(msg)
		}
	}()
	println(err)
	board.I2CConfig(0)
	disp := new(disp)
	disp.SetSize(20, 4)
	disp.Board = *board
	disp.Addr = 0xC6 >> 1
	disp.clear()
	go disp.showTime()
	go func() {
		update_bubbles(&(disp.Content))
		time.Sleep(1000 * time.Millisecond)
	}()
	go func() {
		add_bubbles(&(disp.Content))
		time.Sleep(10000 * time.Millisecond)
	}()
	time.Sleep(1000000 * time.Millisecond)
}
