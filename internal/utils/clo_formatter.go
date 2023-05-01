package utils

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"time"
)

// ToConsole print information to console output with convenient app-format
// Add DateTime at the beginning of line
func ToConsole(title string, data ...interface{}) {
	fmt.Println(
		Gray(12-1, time.Now().Format("2006-01-02 15:04:05")),
		Blue("["+title+"]").Bold(),
		data,
	)
}
