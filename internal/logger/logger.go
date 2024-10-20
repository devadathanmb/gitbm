package logger

import (
	"github.com/fatih/color"
)

func PrintError(format string, a ...interface{}) {
	color.New(color.FgRed, color.Bold).Printf(format+"\n", a...)
}

func PrintSuccess(format string, a ...interface{}) {
	color.New(color.FgGreen, color.Bold).Printf(format+"\n", a...)
}

func PrintInfo(format string, a ...interface{}) {
	color.New(color.FgBlue, color.Bold).Printf(format+"\n", a...)
}

func PrintWarning(format string, a ...interface{}) {
	color.New(color.FgYellow, color.Bold).Printf(format+"\n", a...)
}

func Print(format string, a ...interface{}) {
	color.New().Printf(format+"\n", a...)
}
