package plugins

import (
	"bufio"
	"log"
	"os"
)

func RecordError(errorfile string) (wlog *bufio.Writer, f *os.File) {
	f, err := os.OpenFile(errorfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//defer f.Close()
	if err != nil {
		log.Print(err)
		return
	}
	wlog = bufio.NewWriter(f)
	return

}
