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
	"errors"
	"flag"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"

	tb "code.austinjadams.com/toolbag"
)

type notFoundImg struct {
	Alt, Base64 string
}

type NotFound struct {
	images []*notFoundImg
	templ  *template.Template

	args struct {
		dir, template string
	}
}

func NewNotFound() *NotFound {
	return &NotFound{}
}

func (nf *NotFound) Name() string { return "notfound" }
func (nf *NotFound) Desc() string { return "not found page showing a random math problem" }

// be a catch-all
func (nf *NotFound) Path() string { return "/" }

func (nf *NotFound) AddArgs() {
	flag.StringVar(&nf.args.dir, tb.Arg(nf, "dir"), "", "directory with the math problems")
	flag.StringVar(&nf.args.template, tb.Arg(nf, "template"), "", "path to template")
}

func (nf *NotFound) Init() error {
	if nf.args.dir == "" {
		return errors.New("missing dir arg")
	}
	if nf.args.template == "" {
		return errors.New("missing template arg")
	}

	templ, err := template.ParseFiles(nf.args.template)
	if err != nil {
		return err
	}
	nf.templ = templ

	f, err := os.Open(nf.args.dir)
	if err != nil {
		return err
	}
	names, err := f.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, name := range names {
		imgpath := path.Join(nf.args.dir, name)
		contents, err := ioutil.ReadFile(imgpath)
		if err != nil {
			return err
		}
		pieces := splitLines(string(contents))

		if len(pieces) != 2 {
			return errors.New("malformed file " + imgpath + ". got " + strconv.Itoa(len(pieces)) +
				" lines, expected 2 (first with alttext, second with base64)")
		}

		nf.images = append(nf.images, &notFoundImg{pieces[0], pieces[1]})
	}
	return nil
}

func (nf *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	// deterministic, but who cares?
	i := rand.Intn(len(nf.images))

	err := nf.templ.Execute(w, nf.images[i])
	// template probably already started writing the request, so just log the error
	if err != nil {
		tb.LogWhine(nf, err)
	}
}
