package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

const (
	version string = "1.0"
)

func main() {
	n := flag.Int("n", 5, "Specify an interval n seconds to run command")
	h := flag.Bool("h", false, "Display Help")
	v := flag.Bool("v", false, "Display version")
	flag.Parse()

	command := strings.Join(flag.Args(), " ")
	fmt.Println(command)

	if *v {
		fmt.Println(version)
	} else if *h || len(command) == 0 || n == nil {
		flag.PrintDefaults()
	} else {
		repeatCmd(command, *n)
	}
}

func repeatCmd(cmd string, n int) {
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)

	// first lunch
	if err := executeCmd(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error running '%s' command: \n %v\n", cmd, err)
	}
	for {
		select {
		case <-time.After(time.Duration(n) * time.Second):
			if err := executeCmd(cmd); err != nil {
				fmt.Fprintf(os.Stderr, "There was an error running '%s' command: \n %v\n", cmd, err)
			}
		case <-done:
			return
		}
	}
}

func executeCmd(command string) error {
	var err error

	if err = clearCmd(); err != nil {
		return err
	}

	cmd := exec.Command(command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}

func clearCmd() error {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
