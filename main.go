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

	//fmt.Println("watch is used to run any designated command at regular intervals.")
	command := strings.Join(flag.Args(), " ")

	if *v {
		fmt.Println(version)
	} else if *h {
		flag.PrintDefaults()
	} else {
		repeatCmd(command, *n)
	}
}

func repeatCmd(cmd string, n int) {
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)

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

func executeCmd(cmd string) error {

	var (
		outCmd []byte
		err    error
	)

	if err := clearCmd(); err != nil {
		return err
	}

	if outCmd, err = exec.Command(cmd).Output(); err != nil {
		return err
	}

	output := string(outCmd)
	fmt.Fprintf(os.Stdout, output)
	return nil
}

func clearCmd() error {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	return c.Run()
}
