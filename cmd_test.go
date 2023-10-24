package cmd

import (
	"fmt"
	"os"
	"regexp"
	"testing"
)

func Test_StdOut(t *testing.T) {
	var oc = Cmd{
		Command:          "echo test stdout",
		LogFile:          nil,
		DateFormat:       "",
		OutputBufferSize: 8 * 1024,
	}

	out, err := oc.Run()
	if err != nil {
		fmt.Println(err)
	}

	if string(out) != "test stdout\n" {
		t.Fatal("Output not eq 'test stdout'")
	}

}

func Test_File(t *testing.T) {
	// Write file
	fw, err := os.Create("./test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer fw.Close()

	var oc = Cmd{
		Command:          "echo test file",
		LogFile:          fw,
		DateFormat:       "",
		OutputBufferSize: 8 * 1024,
	}

	out, err := oc.Run()
	if err != nil {
		t.Fatal(err)
	}
	fw.Close()

	// Read file and check
	line, err := os.ReadFile("./test.log")
	if err != nil {
		t.Fatal(err)
	}
	// Checking on present 'test file' content in the log file
	match, _ := regexp.Match(`.*\stest\sfile\n`, line)
	if !match {
		t.Fatalf("Line in the log file doesn't content %s", out)
	}
	// Cleanup
	err = os.Remove("./test.log")
	if err != nil {
		t.Fatal(err)
	}
}
