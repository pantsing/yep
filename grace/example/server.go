package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/pantsing/yep/grace"
	"github.com/qiniu/log"
)

func main() {
	ppid := os.Getppid()
	const msg = "Serving with pid %d ppid %d"
	log.Printf(msg, os.Getpid(), ppid)

	var gs grace.GraceService
	gs.ListenerCloseTimeout = 10

	gl, err := gs.GetListener("tcp", ":6086")
	if err != nil {
		log.Println(err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!"+strconv.Itoa(os.Getpid()))
	})

	go func() {
		err := gs.Serve(gl, mux)
		log.Println(err)
	}()

	err = gs.CloseParentService()
	if err != nil {
		log.Println(err)
	}

	err = gs.WaitSignal(gl)
	if err != nil {
		log.Println(err)
	}
}
