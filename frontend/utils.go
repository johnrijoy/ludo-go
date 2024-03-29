package frontend

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var Blue = color.New(color.FgBlue).SprintfFunc()
var Red = color.New(color.FgRed).SprintfFunc()
var Green = color.New(color.FgGreen).SprintfFunc()
var GreenH = color.New(color.FgHiGreen).SprintfFunc()
var Magenta = color.New(color.FgHiMagenta).SprintfFunc()

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(Blue(label + " "))
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
