package cmd

import "github.com/fatih/color"

func headerFmt() func(format string, a ...interface{}) string {
	return color.New(color.FgGreen, color.Underline).SprintfFunc()
}

func columnFmt() func(format string, a ...interface{}) string {
	return color.New(color.FgYellow).SprintfFunc()
}

func errorNoticeFmt(msg string) (int, error) {
	return color.New(color.FgRed, color.BgHiWhite).Println(msg)
}
