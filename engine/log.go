package engine

import (
	"io"
	"os"
	"log"
)

var logfile *os.File

func SetLog(logname string) {
	logfile, _ = os.Create(logname)
	w := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(0)
	log.SetOutput(w)
}

func CloseLog() {
	logfile.Close()
}