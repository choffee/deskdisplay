// Small app to show working gofirmata working
// Should put something on the screen to start with.
//

package main

import (
  "fmt"
  "github.com/choffee/gofirmata"
  "time"
)

const (
  DISPMOVE byte = 3
)

// Display
type disp struct {
  Board firmata.Board
  Addr  byte
}

func (disp *disp)Write(msg []byte) {
  newmsg := make([]byte, len(msg) + 1)
  newmsg[0] = 0
  for l,v := range msg {
    newmsg[l+1] = v
  }
  disp.Board.I2CWrite(disp.Addr,firmata.I2C_MODE_WRITE,newmsg)
}

func (disp *disp) clear() {
  msg := []byte{12}
  disp.Write(msg)
}

func (disp *disp) moveTo(x,y byte) {
  msg := []byte{DISPMOVE,x,y}
  disp.Write(msg)
}

func (disp *disp) write(s string){
  msg := []byte(s)
  disp.Write(msg)
}

func (disp *disp) putText(s string, x,y byte) {
  disp.moveTo(x,y)
  disp.write(s)
}

// Keep updating the time on the board
func (disp *disp) showTime() {
  for {
    now := time.Now()
    clock := now.Format(time.Kitchen)
    disp.putText("Time " + clock,1,5)
    time.Sleep(1000 * time.Millisecond)
  }
}


func main () {
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
  board.I2CConfig(9)
  disp := new(disp)
  disp.Board = *board
  disp.Addr  = 0xC6 >> 1
  disp.clear()
  go disp.showTime()
  time.Sleep(1000000 * time.Millisecond)
}

