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
	"flag"
	"net/http"
)

type ToolBag struct {
	Tools []Tool
	*http.ServeMux
	*flag.FlagSet
}

func NewToolBag(tools ...Tool) *ToolBag {
	tb := &ToolBag{tools, http.NewServeMux(),
		// neither the program name nor choice of error handling for
		// FlagSet.Parse() matters because we're collecting flags, not
		// parsing them. so, pass some random (but valid) stuff
		flag.NewFlagSet("toolbag", flag.PanicOnError)}

	for _, tool := range tools {
		tool.AddArgs(tb)
		tb.ServeMux.Handle(tool.Path(), tool)
	}

	// for now, just pass the args onto the global flag instance
	tb.FlagSet.VisitAll(func(f *flag.Flag) {
		flag.Var(f.Value, f.Name, f.Usage)
	})

	return tb
}

func (tb *ToolBag) Init() error {
	for _, t := range tb.Tools {
		err := t.Init()

		if err != nil {
			return PrependToolName(t, err)
		}
	}

	return nil
}
