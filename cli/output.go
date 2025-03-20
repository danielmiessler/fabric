package cli

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
)

func CopyToClipboard(message string) (err error) {
	if err = clipboard.WriteAll(message); err != nil {
		err = fmt.Errorf("could not copy to clipboard: %v", err)
	}
	return
}

func CreateOutputFile(message string, fileName string) (err error) {
	var file *os.File
	if file, err = os.Create(fileName); err != nil {
		err = fmt.Errorf("error creating file: %v", err)
		return
	}
	defer file.Close()
	if _, err = file.WriteString(message); err != nil {
		err = fmt.Errorf("error writing to file: %v", err)
	} else {
		fmt.Printf("\n\n... written to %s\n", fileName)
	}
	return
}
