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

	tb "code.austinjadams.com/toolbag"
)

type Tools struct {
	tb    *tb.ToolBag
	templ *template.Template

	args struct {
		source, template string
	}
}

func NewTools() *Tools {
	return &Tools{}
}

func (t *Tools) Name() string { return "tools" }
func (t *Tools) Desc() string { return "list tools" }

// be a catch-all
func (t *Tools) Path() string { return "/" + t.Name() }

func (t *Tools) AddArgs(toolbag *tb.ToolBag) {
	// XXX kind of a hack, but no regrets here
	t.tb = toolbag
	toolbag.StringVar(&t.args.template, tb.Arg(t, "template"), "", "path to template")
	toolbag.StringVar(&t.args.source, tb.Arg(t, "source"), "", "url of toolbag source")
}

func (t *Tools) Init() error {
	if t.args.template == "" {
		return errors.New("missing template arg")
	}
	if t.args.source == "" {
		return errors.New("missing source arg")
	}

	templ, err := template.ParseFiles(t.args.template)
	if err != nil {
		return err
	}
	t.templ = templ
	return nil
}

func (t *Tools) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := t.templ.Execute(w, struct {
		Source string
		Tools  []tb.Tool
	}{t.args.source, t.tb.Tools})

	if err != nil {
		tb.LogWhine(t, err)
	}
}
