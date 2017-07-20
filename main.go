package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

const (
	version string = "1.0"
)

var (
	n = flag.Int("n", 5, "Specify an interval n seconds to run command")
	h = flag.Bool("h", false, "Display Help")
	v = flag.Bool("v", false, "Display version")
)

func main() {
	flag.Parse()

	if *v {
		fmt.Println(version)
		return
	}
	if *h {
		flag.PrintDefaults()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "command is required")
		os.Exit(1)
	}
	command := args[0]
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}

	repeatCmd(command, args, time.Duration(*n))
}

func repeatCmd(cmd string, args []string, n time.Duration) {
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)

	// first lunch
	if err := executeCmd(cmd, args); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error running '%s' command: \n %v\n", cmd, err)
		return
	}
	for {
		select {
		case <-time.After(n * time.Second):
			if err := executeCmd(cmd, args); err != nil {
				fmt.Fprintf(os.Stderr, "There was an error running '%s' command: \n %v\n", cmd, err)
				return
			}
		case <-done:
			return
		}
	}
}

func executeCmd(command string, args []string) error {
	if err := clearCmd(); err != nil {
		return err
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func clearCmd() error {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
