package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

var (
	verbose = flag.Bool("v", false, "verbose")
)

func main() {
	// We have to remove negative ints from the list of args to parse because
	// they're treated as (invalid) flags. Example: vol -2
	args := make([]string, 0, len(os.Args))
	negRe := regexp.MustCompile(`\A-\d+\z`)
	for _, arg := range os.Args[1:] {
		if negRe.MatchString(arg) {
			continue
		}
		args = append(args, arg)
	}
	flag.CommandLine.Parse(args)

	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if err := processCommand(); err != nil {
		log.Fatal(err)
	}
}

func processCommand() (err error) {
	// We can't use flag.NArg() nor flag.Arg() as in main() we parsed a list of
	// args minus negative ints
	switch len(os.Args) {
	case 1:
		return printCurrent()
	case 2:
		s := os.Args[1]
		var i int
		if _, err := fmt.Sscanf(s, "+%d", &i); err == nil {
			fmt.Printf("+x %d\n", i)
			return modify(i)
		} else if _, err := fmt.Sscanf(s, "-%d", &i); err == nil {
			fmt.Printf("-x %d\n", i)
			return modify(-i)
		} else if _, err := fmt.Sscanf(s, "%d", &i); err == nil {
			fmt.Printf("set  x %d\n", i)
			return set(i)
		}
	}
	usage()
	return
}

func usage() {
	fmt.Printf("usage: %s [-v] [vol|+-delta]\n", filepath.Base(os.Args[0]))
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
	)
	switch {
	case vol < min:
		log.Warnf("vol %d < %d, adjusting to %d", vol, min, min)
		vol = min
	case vol > max:
		log.Warnf("vol %d > %d, adjusting to %d", vol, max, max)
		vol = max
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
