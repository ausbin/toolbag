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
	"net/http"
)

type ToolBag struct {
	Tools []Tool
	*http.ServeMux
}

func NewToolBag(tools ...Tool) *ToolBag {
	mux := http.NewServeMux()

	for _, tool := range tools {
		tool.AddArgs()
		mux.Handle(tool.Path(), tool)
	}

	return &ToolBag{tools, mux}
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
