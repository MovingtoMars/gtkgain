package main

import (
	"github.com/MovingtoMars/gotk3/gtk"
	"github.com/MovingtoMars/gtkgain/src/library"
	//"github.com/davecheney/profile"

	//"fmt"
	"log"
	"runtime"
)

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()

	runtime.GOMAXPROCS(runtime.NumCPU())

	gtk.Init(nil)

	createWindow(library.New())

	gtk.Main()
}

func crashIf(mess string, err error) {
	if err != nil {
		log.Fatal(mess, ": ", err)
	}
}
