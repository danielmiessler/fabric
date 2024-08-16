package utils

import (
	"fmt"
	"os"

	"gopkg.in/gookit/color.v1"
)

func Print(info string) {
	fmt.Println(info)
}

func PrintWarning (s string) {
	fmt.Println(color.Yellow.Render("Warning: " + s))
}

func LogError(err error) {
	fmt.Fprintln(os.Stderr, color.Red.Render(err.Error()))
}

func LogWarning(err error) {
	fmt.Fprintln(os.Stderr, color.Yellow.Render(err.Error()))
}

func Log(info string) {
	fmt.Println(color.Green.Render(info))
}