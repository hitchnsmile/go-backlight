// +build linux
// +build 386

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func err(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
func main() {

	out, e := exec.Command("ls", "/sys/class/backlight/").Output()
	err(e)

	display := string(out[:len(out)-1])
	maxBrightness := getData(display, "max_brightness")
	actualBrightness := getData(display, "actual_brightness")

	targetBrightness := flag.Int("set", actualBrightness, "The value has to be a number")
	adjBrightness := flag.Int("adjust", 0, "The value has to be a number")

	flag.Parse()

	if *targetBrightness != actualBrightness {
		c := changeValue(*targetBrightness, maxBrightness)
		writeData(display, c)
	} else if *adjBrightness != 0 {
		c := adjustValue(actualBrightness, *adjBrightness, maxBrightness)
		writeData(display, c)
	} else {
		fmt.Println("Use -set int or -adjust int commands as sudo")
	}
}

func adjustValue(current int, change int, max int) int {
	cur := current * 100 / max
	chPoint := max / 100
	if v := (cur + change); v > 100 {
		return max
	} else if d := (v * chPoint); d > max {
		return max
	} else if v < 5 {
		fmt.Println("Don't play with darkness...")
		return 2 * chPoint
	} else {
		return v * chPoint
	}
}

func changeValue(target int, max int) int {
	chPoint := max / 100
	if v := target; v >= 100 {
		return max
	} else if v < 2 {
		fmt.Println("Don't play with darkness...")
		return 2 * chPoint
	} else {
		return v * chPoint
	}
}

func getData(display string, target string) int {
	data := make([]byte, 100)
	f, e := os.Open("/sys/class/backlight/" + display + "/" + target) // For read access.
	err(e)
	defer f.Close()
	count, e := f.Read(data)
	err(e)
	r, e := strconv.Atoi(string(data[:count-1]))
	err(e)
	return r
}

func getDisplay() string {
	data := make([]byte, 100)
	f, e := os.Open("go-backlight-display") // For read access.
	err(e)
	defer f.Close()
	count, e := f.Read(data)
	err(e)
	return string(data[:count-1])
}

func writeData(display string, value int) int {
	bs := []byte(strconv.Itoa(value))
	f, e := os.OpenFile("/sys/class/backlight/"+display+"/brightness", os.O_RDWR|os.O_CREATE, 0755) // For read access.
	err(e)
	defer f.Close()

	r, e := f.Write(bs)
	err(e)
	return r
}
