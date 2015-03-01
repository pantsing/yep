package grace_test

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/pantsing/yep/grace"
)

func ExampleServer() {
	addr, err := net.ResolveTCPAddr("tcp", ":6086")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	gl := grace.NewGraceListener(l)

	err = http.Serve(gl, nil)
	fmt.Println(err)

	time.Sleep(3 * time.Minute)
}
