package lib

/*
Console utilities, should know nothing about this application
*/


import (
	"fmt"
    "os"
)

/****/
type winsize struct {
    Row    uint16
    Col    uint16
    Xpixel uint16
    Ypixel uint16
}

const (
    RuneSterling = '£'
    RuneDArrow   = '↓'
    RuneLArrow   = '←'
    RuneRArrow   = '→'
    RuneUArrow   = '↑'
    RuneBullet   = '·'
    RuneBoard    = '░'
    RuneCkBoard  = '▒'
    RuneDegree   = '°'
    RuneDiamond  = '◆'
    RuneGEqual   = '≥'
    RunePi       = 'π'
    RuneHLine    = '─'
    RuneLantern  = '§'
    RunePlus     = '┼'
    RuneLEqual   = '≤'
    RuneLLCorner = '└'
    RuneLRCorner = '┘'
    RuneNEqual   = '≠'
    RunePlMinus  = '±'
    RuneS1       = '⎺'
    RuneS3       = '⎻'
    RuneS7       = '⎼'
    RuneS9       = '⎽'
    RuneBlock    = '█'
    RuneTTee     = '┬'
    RuneRTee     = '┤'
    RuneLTee     = '├'
    RuneBTee     = '┴'
    RuneULCorner = '┌'
    RuneURCorner = '┐'
    RuneVLine    = '│' //'│'
    RuneUVLine   = '╷'
    RuneDVLine   = '╵'
)

const (
    ESC_SAVE_SCREEN = "?47h"
    ESC_RESTORE_SCREEN = "?47l"

    ESC_SAVE_CURSOR = "s"
    ESC_RESTORE_CURSOR = "u"

    ESC_BOLD_ON = "1m"
    ESC_BOLD_OFF = "0m"

    ESC_CURSOR_ON = "?25h"
    ESC_CURSOR_OFF = "?25l"

    ESC_CLEAR_SCREEN = "2J"
    ESC_CLEAR_LINE = "2K"
)

/** print to a specific terminal screen location */
func PrintStrAt(msg string, y, x int) {
    fmt.Printf("\033[%d;%dH%s", y, x, msg)
}

/** print to a terminal code */
func PrintCtrAt(esc string, y, x int) {
    fmt.Printf("\033[%d;%dH\033[%s", y, x, esc)
}

/**
print out an esc control
@param esc control code to print out
*/
func PrintCtrOnErr(esc string) {
    fmt.Fprintf(os.Stderr, "\033[%s", esc)
}

/**
print out an esc control
@param esc control code to print out
*/
func PrintCtrOnOut(esc string) {
    fmt.Fprintf(os.Stdout, "\033[%s", esc)
}

/** Save the screen setup at the start of the app */
func ScrSave() {
    PrintCtrOnErr(ESC_SAVE_SCREEN)
    PrintCtrOnErr(ESC_SAVE_CURSOR)
    //PrintCtrOnErr(ESC_CURSOR_OFF)
    PrintCtrOnErr(ESC_CLEAR_SCREEN)
}

/** print a string in a color */
func ColorText(text string, color int) string {
    encoded := fmt.Sprintf("\033[0;%dm%s\033[0m", color, text)
    return encoded
}

func Red(text string) string {
    return ColorText(text, 31)
}

func Green(text string) string {
    return ColorText(text, 32)
}

func Blue(text string) string {
    return ColorText(text, 34)
}

/** Restore the screen setup from SrcSave() */
func ScrRestore() {
    //PrintCtrOnErr(ESC_CURSOR_ON)
    PrintCtrOnErr(ESC_RESTORE_CURSOR)
    PrintCtrOnErr(ESC_RESTORE_SCREEN)
}
