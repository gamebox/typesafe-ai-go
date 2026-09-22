package main

import "fmt"

func moveCursorHome() {
	fmt.Print("\x1B[H")
}
func eraseScreen() {
	fmt.Print("\x1B[2J")
}

func cursorHidden() {
	fmt.Print("\x1B[?25l")
}
func cursorVisible() {
	fmt.Print("\x1B[?25h")
}
func enterAltScreen() {
	fmt.Print("\x1B[?1049h")
}
func exitAltScreen() {
	fmt.Print("\x1B[?1049l")
}
func moveCursorTo(line, col int) {
	fmt.Printf("\x1B[%d;%dH", line, col)
}
func moveCursorBONL() {
	fmt.Print("\x1B[1E")
}
