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

package tools

import (
	"net/http"
	"strings"

	tb "code.austinjadams.com/toolbag"
)

var (
	IP = tb.NewToolFunc("ip", "show ip address", func(w http.ResponseWriter, r *http.Request) {
		plainText(w)
		w.Write([]byte(tidyAddress(r.RemoteAddr) + "\n"))
	})

	Headers = tb.NewToolFunc("headers", "write request headers", func(w http.ResponseWriter, r *http.Request) {
		plainText(w)
		r.Header.Write(w)
	})
)

func All() *tb.ToolBag {
	return tb.NewToolBag(IP, Headers, NewFiglet(), NewNotFound(), NewTools())
}

// nasty hack to remove port numbers and the brackets from ipv6 addresses
// it's "nasty" because the go docs (unfortunately) make no guarantee
// about the format of http.Request.RemoteAddr
func tidyAddress(addr string) (result string) {
	result = addr
	if i := strings.LastIndex(result, ":"); i != -1 {
		result = result[:i]
	}
	result = strings.Trim(result, "[]")
	return
}

func splitLines(blob string) []string {
	return strings.Split(strings.TrimSpace(blob), "\n")
}

func plainText(w http.ResponseWriter) {
	// stop go from trying to guess the content-type of a response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}
