package grace

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	ErrAlreadyClosed        = errors.New("Listener already closed")
	errRestartListener      = errors.New("No listener for restart")
	errListenerCloseTimeout = errors.New("Listener close timeout")
	errNotSupportedNetwork  = errors.New("Network type is not supported")
)

const (
	envRestartKey       = "_GRACE_RESTART"
	envRestartKeyPrefix = envRestartKey + "="
	errClosed           = "use of closed network connection"
)

// GraceListener requires the file descriptor of listener could be got by File() function.
// When service restarts, the listener will be passed to child process by file descriptor.
// So only TCPListener or UNIXListener is supported.
type GraceListener interface {
	net.Listener                   // Inherit original TCP/UNIX listener interface
	File() (f *os.File, err error) // Get file descriptor
	SetDeadline(t time.Time) error // Needed by close listener
}

type gListener struct {
	GraceListener
	closed      bool
	closedMutex sync.RWMutex
	wg          sync.WaitGroup
}

type conn struct {
	net.Conn
	wg   *sync.WaitGroup
	once sync.Once
}

func (c *conn) Close() error {
	defer c.once.Do(c.wg.Done)
	return c.Conn.Close()
}

func NewGraceListener(l GraceListener) GraceListener {
	return &gListener{GraceListener: l}
}

func (l *gListener) Close() error {
	l.closedMutex.Lock()
	l.closed = true
	l.closedMutex.Unlock()

	var err error
	if os.Getppid() == 1 {
		err = l.GraceListener.SetDeadline(time.Now())
	} else {
		err = l.GraceListener.Close()
	}
	l.wg.Wait()
	return err
}

func (l *gListener) Accept() (net.Conn, error) {
	var c net.Conn
	l.wg.Add(1)
	defer func() {
		if c == nil {
			l.wg.Done()
		}
	}()

	l.closedMutex.RLock()
	if l.closed {
		l.closedMutex.RUnlock()
		return nil, ErrAlreadyClosed
	}
	l.closedMutex.RUnlock()

	c, err := l.GraceListener.Accept()
	if err != nil {
		if strings.HasSuffix(err.Error(), errClosed) {
			return nil, ErrAlreadyClosed
		}

		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			l.closedMutex.RLock()
			if l.closed {
				l.closedMutex.RUnlock()
				return nil, ErrAlreadyClosed
			}
			l.closedMutex.RUnlock()
		}
		return nil, err
	}
	return &conn{Conn: c, wg: &l.wg}, nil
}

type GraceService struct {
	ListenerCloseTimeout time.Duration
	srv                  *http.Server
}

func (gs *GraceService) closeListener(l GraceListener) (err error) {
	gs.srv.SetKeepAlivesEnabled(false)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(l GraceListener) {
		err = l.Close()
		wg.Done()
	}(l)

	if gs.ListenerCloseTimeout == 0 {
		// wait forever... T_T
		wg.Wait()
	} else {
		// wait in background to allow for implementing a timeout
		done := make(chan struct{})
		go func() {
			defer close(done)
			wg.Wait()
		}()

		// wait for graceful termination or timeout
		select {
		case <-done:
			// fmt.Println("wg.Wait done")
		case <-time.After(gs.ListenerCloseTimeout * time.Second):
			return errListenerCloseTimeout
		}
	}

	return
}

// CloseParentService Send TERM signal to old service processor.
func (gs *GraceService) CloseParentService() error {
	parentPID := os.Getppid()
	if parentPID == 1 {
		return nil
	}
	return syscall.Kill(parentPID, syscall.SIGQUIT)
}

func (gs *GraceService) GetListener(network, addr string) (gl GraceListener, err error) {
	isRestart := os.Getenv(envRestartKey)
	if isRestart == "1" {
		gl, err = gs.InheritListener()
	}

	if err != nil || isRestart != "1" {
		switch network {
		case "tcp", "tcp4", "tcp6":
			tcpAddr, err := net.ResolveTCPAddr(network, addr)
			if err != nil {
				return nil, err
			}
			l, err := net.ListenTCP(network, tcpAddr)
			if err != nil {
				return nil, err
			}
			gl = NewGraceListener(l)
		default:
			return nil, errNotSupportedNetwork
		}

	}
	return
}

// InheritListener inherits listener from old processor.
// File descriptor number of listener is 3 for only stdin, stdout, stderr and listener
// are opened when restart.
func (gs *GraceService) InheritListener() (gl GraceListener, err error) {
	f := os.NewFile(uintptr(3), "")
	tmp, err := net.FileListener(f)
	f.Close()
	if err != nil {
		return nil, err
	}
	l := tmp.(GraceListener)
	gl = NewGraceListener(l)
	return
}

func (gs *GraceService) Serve(gl GraceListener, handler http.Handler) (err error) {
	gs.srv = &http.Server{Handler: handler}
	return gs.srv.Serve(gl)
}

func (gs *GraceService) Restart(gl GraceListener) (err error) {
	if gl == nil {
		return errRestartListener
	}

	// Extract the file descriptor from the listener.
	f, err := gl.File()
	if err != nil {
		return err
	}
	defer f.Close()                  // Close listener file descriptor when old processor exit
	syscall.CloseOnExec(int(f.Fd())) // Make sure file descriptor for listener in new process is closed

	// Use the original binary location. This works with symlinks such that if
	// the file it points to has been changed we will use the updated symlink.
	argv0, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}

	// In order to keep the working directory the same as when we started.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	var env []string
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, envRestartKeyPrefix) {
			env = append(env, v)
		}
	}
	env = append(env, fmt.Sprintf("%s%d", envRestartKeyPrefix, 1))

	allFiles := append([]*os.File{os.Stdin, os.Stdout, os.Stderr}, f)
	_, err = os.StartProcess(argv0, os.Args, &os.ProcAttr{
		Dir:   wd,
		Env:   env,
		Files: allFiles,
	})
	return err
}

// WaitSignal waits for signals to gracefully terminate or restart the process.
func (gs *GraceService) WaitSignal(gl GraceListener) error {

	ch := make(chan os.Signal, 6)
	signal.Notify(ch,
		syscall.SIGTERM, // TERM : Shutdown
		syscall.SIGINT,  // INT  : Shutdown
		syscall.SIGQUIT, // QUIT : Gracefully shutdown
		syscall.SIGHUP,  // HUP  : Gracefully reload configure and restart
	// syscall.SIGUSR1, // USR1 : Reopen log file
	// syscall.SIGUSR2, // USR2 : Gracefully restart
	)

	for {
		sig := <-ch
		switch sig {
		case syscall.SIGTERM:
			fallthrough
		case syscall.SIGINT:
			// this ensures a subsequent TERM will trigger standard go behaviour of
			// terminating.
			signal.Stop(ch)
			return nil
		case syscall.SIGQUIT:
			signal.Stop(ch)
			return gs.closeListener(gl)
		case syscall.SIGHUP:
			// we only return here if there's an error, otherwise the new process
			// will send us a TERM when it's ready to trigger the actual shutdown.
			if err := gs.Restart(gl); err != nil {
				return err
			}
		}
	}
}
