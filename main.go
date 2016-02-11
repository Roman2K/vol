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
	if len(os.Args) != 2 {
		usage()
	}
	var err error
	switch os.Args[1] {
	case "+":
		err = increase()
	case "-":
		err = decrease()
	default:
		usage()
	}
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func usage() {
	fmt.Printf("usage: %s +|-\n", filepath.Base(os.Args[0]))
	os.Exit(1)
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
