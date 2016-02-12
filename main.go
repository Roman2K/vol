package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	if err := processCommand(); err != nil {
		log.Fatal(err)
	}
}

func processCommand() (err error) {
	switch len(os.Args) {
	case 1:
		return printCurrent()
	case 2:
		switch os.Args[1] {
		case "+":
			return increase()
		case "-":
			return decrease()
		}
	}
	usage()
	return
}

func usage() {
	fmt.Printf("usage: %s [+|-]\n", filepath.Base(os.Args[0]))
	os.Exit(1)
}

func printCurrent() (err error) {
	cur, err := get()
	if err != nil {
		return
	}
	fmt.Printf("%d\n", cur)
	return
}

func increase() error {
	return modify(+5)
}

func decrease() error {
	return modify(-5)
}

func modify(delta int) error {
	log.Debugf("modifying volume: %+d", delta)
	cur, err := get()
	if err != nil {
		return err
	}
	target := cur + delta
	log.Debugf("%d => %d", cur, target)
	return set(target)
}

func get() (cur int, err error) {
	const script = "output volume of (get volume settings)"
	out, err := osascript(script)
	if err != nil {
		return
	}
	_, err = fmt.Sscanf(string(out), "%d", &cur)
	return
}

func set(vol int) (err error) {
	const (
		min = 0
		max = 100
		mod = 5
	)
	switch {
	case vol < min:
		log.Warnf("vol %d < %d, adjusting to %d", vol, min, min)
		vol = min
	case vol > max:
		log.Warnf("vol %d > %d, adjusting to %d", vol, max, max)
		vol = max
	}
	if vol%mod != 0 {
		adjusted := vol / mod * mod
		log.Warnf("vol %d not modulo %d, adjusting to %d", vol, mod, adjusted)
		vol = adjusted
	}
	log.Infof("setting volume to %d", vol)
	script := fmt.Sprintf("set volume output volume %d", vol)
	_, err = osascript(script)
	return
}

func osascript(script string) ([]byte, error) {
	log.Debugf("executing osascript `%s`", script)
	return exec.Command("osascript", "-e", script).Output()
}
