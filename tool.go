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

package toolbag

import (
	"errors"
	"log"
	"net/http"
)

type Tool interface {
	Name() string
	Desc() string
	Path() string
	AddArgs(*ToolBag)
	Init() error
	http.Handler
}

type ToolFunc struct {
	name, desc, path string
	http.HandlerFunc
}

func NewToolFunc(name, desc string, f http.HandlerFunc) *ToolFunc {
	return NewToolFuncAt(name, desc, "/"+name, f)
}

func NewToolFuncAt(name, desc, path string, f http.HandlerFunc) *ToolFunc {
	return &ToolFunc{name, desc, path, f}
}

func (ft *ToolFunc) Name() string     { return ft.name }
func (ft *ToolFunc) Desc() string     { return ft.desc }
func (ft *ToolFunc) Path() string     { return ft.path }
func (ft *ToolFunc) AddArgs(*ToolBag) {}
func (ft *ToolFunc) Init() error      { return nil }

// just a way to compose names of flags
func Arg(t Tool, argname string) string {
	return t.Name() + ":" + argname
}

// the same, but for error messages
func toolMsg(t Tool, msg string) string {
	return t.Name() + ": " + msg
}

// stick the name of a tool before an error message
func PrependToolName(t Tool, err error) error {
	return errors.New(toolMsg(t, err.Error()))
}

// error handling during ServeHTTP()
func whine(w http.ResponseWriter, msg string) {
	logWhine(msg)
	http.Error(w, msg, http.StatusInternalServerError)
}

func logWhine(msg string) {
	log.Println(msg)
}

func WhineString(t Tool, w http.ResponseWriter, msg string) {
	whine(w, toolMsg(t, msg))
}

func Whine(t Tool, w http.ResponseWriter, err error) {
	WhineString(t, w, err.Error())
}

// for when WriteHeader() has already been called
// (e.g., template errors)
func LogWhineString(t Tool, msg string) {
	logWhine(toolMsg(t, msg))
}

func LogWhine(t Tool, err error) {
	LogWhineString(t, err.Error())
}
