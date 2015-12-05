// toolbag - dynamic tools for my website
// Copyright (C) 2015 Austin Adams
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"

	"code.austinjadams.com/toolbag/tools"
)

func usage(msg string) {
	log.Println(msg + "\n")
	flag.Usage()
	os.Exit(1)
}

func args(usefcgi, usehttp *bool, unix, tcp *string) {
	flag.BoolVar(usefcgi, "fcgi", false, "use fastcgi")
	flag.BoolVar(usehttp, "http", false, "use http")
	flag.StringVar(unix, "unix", "", "path to a unix socket")
	flag.StringVar(tcp, "tcp", "", "path to a unix socket")
	flag.Parse()

	// xnor
	if *usefcgi == *usehttp {
		usage("specify -fcgi or -http, but not both")
	} else if (*unix != "") == (*tcp != "") {
		usage("specify -unix or -tcp, but not both")
	}
}

func listen(unix, tcp string) (sock net.Listener, path string, err error) {
	if unix != "" {
		path = unix
		sock, err = net.ListenUnix("unix", &net.UnixAddr{unix, "unix"})
	} else {
		path = tcp
		sock, err = net.Listen("tcp", tcp)
	}
	return
}

func main() {
	// the systemd journal will include the program name and a
	// timestamp, so don't print one
	log.SetFlags(0)

	// allow tools to add flags before we call flag.Parse()
	tb := tools.All()

	// args
	var usefcgi, usehttp bool
	var unix, tcp string
	args(&usefcgi, &usehttp, &unix, &tcp)

	if err := tb.Init(); err != nil {
		usage(err.Error())
	}

	sock, path, err := listen(unix, tcp)
	if err != nil {
		log.Fatalln("can't open socket", path, err)
	}
	log.Println("listening on socket", path)

	var serve func(net.Listener, http.Handler) error
	if usehttp {
		log.Println("using http...")
		serve = http.Serve
	} else if usefcgi {
		log.Println("using fastcgi...")
		serve = fcgi.Serve
	}

	if err := serve(sock, tb); err != nil {
		log.Fatalln("can't serve", err)
	}
}
