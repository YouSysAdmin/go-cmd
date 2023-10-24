package cmd

import (
	"fmt"
	"github.com/armon/circbuf"
	"github.com/mitchellh/go-linereader"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Cmd struct {
	Command          string   // ping 127.0.0.1
	LogFile          *os.File // File descriptor
	DateFormat       string   // Date formatting
	OutputBufferSize int64    // (8 * 1024) - buffer size in bytes
}

func (c *Cmd) Run() (cmdOutput []byte, err error) {
	// The default size for the buffer of a command output
	if c.OutputBufferSize <= 0 {
		c.OutputBufferSize = 8 * 1024 // 8KB
	}
	// Buffer for read/write command output
	buff, _ := circbuf.NewBuffer(c.OutputBufferSize)
	// Pipe of output
	pr, pw := io.Pipe()
	// Read log channel
	pOutputCh := make(chan struct{})

	// Start goroutine for copying command output to STDOUT or/and file
	go c.output(pr, pOutputCh)

	args := strings.Fields(c.Command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = io.MultiWriter(buff, pw)
	cmd.Stdout = io.MultiWriter(buff, pw)

	// Run the command to completion
	err = cmd.Run()

	pw.Close()
	<-pOutputCh

	if err != nil {
		return buff.Bytes(), fmt.Errorf("error running command '%s': %v. Output: %s", c.Command, err, buff.Bytes())
	}

	return buff.Bytes(), nil
}

func (c *Cmd) output(r io.Reader, doneCh chan<- struct{}) {
	defer close(doneCh)

	if c.DateFormat == "" {
		c.DateFormat = time.RFC3339
	}

	lr := linereader.New(r)
	for line := range lr.Ch {
		logLine := fmt.Sprintf("[%s] %s\n", time.Now().Format(c.DateFormat), line)
		// Output to STDOUT
		fmt.Print(logLine)
		// Output to file
		if c.LogFile != nil {
			_, err := c.LogFile.WriteString(logLine)
			if err != nil {
				fmt.Printf("error write to file: %s, Output: %s\n", err, logLine)
			}
		}
	}
}
