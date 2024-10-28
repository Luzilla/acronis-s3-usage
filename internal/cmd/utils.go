package cmd

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func headerFmt() func(format string, a ...interface{}) string {
	return color.New(color.FgGreen, color.Underline).SprintfFunc()
}

func columnFmt() func(format string, a ...interface{}) string {
	return color.New(color.FgYellow).SprintfFunc()
}

func errorNoticeFmt(msg string) (int, error) {
	return color.New(color.FgRed, color.BgHiWhite).Println(msg)
}

func emailFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "email",
		Required: true,
	}
}

func keyIDFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "key-id",
		Required: true,
	}
}
