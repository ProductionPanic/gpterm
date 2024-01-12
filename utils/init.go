package utils

import "github.com/atotto/clipboard"

func CopyToClipboard(text string) {
	err := clipboard.WriteAll(text)
	if err != nil {
		return
	}
}

func ExecuteGeneratedCommand(command string) {

}
