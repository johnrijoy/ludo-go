package frontend

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var Blue = color.New(color.FgHiBlue).SprintFunc()
var Red = color.New(color.FgRed).SprintFunc()
var Green = color.New(color.FgGreen).SprintFunc()
var GreenH = color.New(color.FgHiGreen).SprintFunc()
var GreenD = color.New(color.FgGreen, color.Faint).SprintFunc()
var Magenta = color.New(color.FgHiMagenta).SprintFunc()
var Gray = color.New(color.FgWhite, color.Faint).SprintFunc()
var Yellow = color.New(color.FgYellow).SprintFunc()
var Cyan = color.New(color.FgCyan).SprintFunc()

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

func silentLog(label ...interface{}) {
	fmt.Println(Gray(label...))
}

func warnLog(label ...interface{}) {
	fmt.Println(Yellow(label...))
}

func errorLog(label ...interface{}) {
	prefix := Red("Error:")
	x := make([]interface{}, 0)
	x = append(x, prefix)
	x = append(x, label...)
	fmt.Println(x...)
}
