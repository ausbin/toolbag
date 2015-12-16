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
	"html/template"
	"net/http"
	"strings"

	"code.austinjadams.com/execd"
	tb "code.austinjadams.com/toolbag"
)

type Figlet struct {
	defaultFont string
	fonts       map[string][]string
	templ       *template.Template
	net, addr   string

	args struct {
		template, unix, tcp string
	}
}

func NewFiglet() *Figlet {
	return &Figlet{}
}

func (f *Figlet) Name() string { return "figlet" }
func (f *Figlet) Desc() string { return "a web frontend to figlet" }
func (f *Figlet) Path() string { return "/" + f.Name() }
func (f *Figlet) AddArgs(toolbag *tb.ToolBag) {
	toolbag.StringVar(&f.args.template, tb.Arg(f, "template"), "", "path to template")
	toolbag.StringVar(&f.args.unix, tb.Arg(f, "unix"), "", "path to unix socket to execd")
	toolbag.StringVar(&f.args.tcp, tb.Arg(f, "tcp"), "", "tcp address to execd")
}

func (f *Figlet) makeClient() (*execd.Client, error) {
	return execd.DialClient(f.net, f.addr)
}

func (f *Figlet) fontCategory(needle string) string {
	for category, fonts := range f.fonts {
		for _, font := range fonts {
			if needle == font {
				return category
			}
		}
	}
	// no match
	return ""
}

func (f *Figlet) Init() error {
	if f.args.template == "" {
		return errors.New("missing template arg")
	}
	if (f.args.unix == "") == (f.args.tcp == "") {
		return errors.New("specify either unix or tcp, but not both")
	}

	templ, err := template.ParseFiles(f.args.template)
	if err != nil {
		return err
	}
	f.templ = templ

	if f.args.unix != "" {
		f.net = "unix"
		f.addr = f.args.unix
	} else {
		f.net = "tcp"
		f.addr = f.args.tcp
	}

	client, err := f.makeClient()
	if err != nil {
		return err
	}
	// find default font
	defaultFont, err := client.ExecString("", "fig", "default")
	if err != nil {
		return err
	}
	f.defaultFont = strings.TrimSpace(defaultFont)

	// find categories of fonts
	output, err := client.ExecString("", "fig", "ls")
	if err != nil {
		return err
	}
	f.fonts = make(map[string][]string)
	for _, v := range splitLines(output) {
		f.fonts[v] = nil
	}

	// find the fonts in each category
	for category, _ := range f.fonts {
		output, err = client.ExecString("", "fig", "ls", category)
		if err != nil {
			return err
		}
		f.fonts[category] = splitLines(output)
	}

	return nil
}

// serve
func (f *Figlet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	font := r.PostFormValue("font")
	text := r.PostFormValue("text")
	result := ""

	if font == "" {
		font = f.defaultFont
	}

	if text != "" {
		category := f.fontCategory(font)
		if category == "" {
			tb.WhineString(f, w, "nice try, lad")
			return
		}

		client, err := f.makeClient()
		if err != nil {
			tb.Whine(f, w, err)
			return
		}
		// otherwise
		defer client.Close()

		result, err = client.ExecString(text, "fig", category, font)
		if err != nil {
			tb.Whine(f, w, err)
			return
		}
	}

	err := f.templ.Execute(w, &struct {
		Font         string
		Fonts        map[string][]string
		Text, Result string
	}{font, f.fonts, text, result})
	if err != nil {
		tb.LogWhine(f, err)
		return
	}
}
