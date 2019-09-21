// Program pinginject is a really simple HTTP server which allows shell
// injection as is commonly found on home routers.
package main

/*
 * pinginject.go
 * Program to allow shell injection in a ping field
 * By J. Stuart McMurray
 * Created 20190920
 * Last Modified 20190921
 */

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

/* handler handles pinging for HTTP queries */
type handler struct {
	param  string /* Vulnerable parameter */
	prefix string /* Start of string passed to shell */
}

// ServeHTTP pings whatever r asks for
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	tag := fmt.Sprintf("[%v] %v %v", r.RemoteAddr, r.Method, r.URL)

	/* Make sure we have a ping field */
	if err := r.ParseForm(); nil != err {
		log.Printf("%v Unable to parse form", tag)
		http.Error(w, "Bad form", http.StatusBadRequest)
		return
	}
	t := r.Form.Get(h.param)
	if "" == t {
		log.Printf("%v Missing %v", tag, h.param)
		http.Error(w, "Missing "+h.param, http.StatusBadRequest)
		return
	}

	/* Ping it.  This is a terrible idea. */
	var (
		o   []byte
		err error
	)
	switch runtime.GOOS {
	case "windows":
		o, err = exec.Command(
			"powershell.exe",
			"-noni", "-nop",
			"-command", h.prefix+t,
		).CombinedOutput()
	default:
		o, err = exec.Command(
			"/bin/sh",
			"-c", h.prefix+t,
		).CombinedOutput()
	}
	if nil != err {
		if 0 != len(o) {
			o = append(o, '\n')
		}
		o = append(o, err.Error()...)
	}
	log.Printf("%v Q: %q (%d)", tag, t, len(o))
	w.Write(o)
}

func main() {
	var (
		pingPrefix = flag.String(
			"command",
			"/sbin/ping -c 4",
			"Ping command `prefix`",
		)
		lAddr = flag.String(
			"listen",
			"0.0.0.0:8080",
			"Listen `address`",
		)
		param = flag.String(
			"parameter",
			"ip",
			"Vulnerable HTTP request `parameter`",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %v [options]

Runs a little webserver which will respond to queries with a parameter "ping"
by putting the contents of the parameter after the command prefix and running
it as a command.  The prefix will be split on whitespace.  This is meant to
be similar to routers which allow shell injection in their ping diagonstic
utility.

Example URL:  http://address/ping.php?ping=127.0.0.1$(callmeback.sh)

The /ping.php in the above example is not significant.  All paths are served.

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Register HTTP handler */
	if !strings.HasSuffix(*pingPrefix, " ") {
		*pingPrefix += " "
	}
	http.Handle("/", handler{
		param:  *param,
		prefix: *pingPrefix,
	})

	/* Listen for connections */
	l, err := net.Listen("tcp", *lAddr)
	if nil != err {
		log.Fatalf("Unable to listen on %v: %v", *lAddr, err)
	}
	log.Printf("Listening on %v for HTTP requests", l.Addr())

	/* Serve HTTP */
	log.Fatalf("HTTP server error: %v", http.Serve(l, nil))
}
