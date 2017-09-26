// config.go
// Part of the py65816 package
// Scot W. Stevenson scot.stevenson@gmail.com
// First version: 26. Sep 2017
// Second version: 26. Sep 2017

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"strings"
)

func IsComment(s string) bool {
	cs := strings.TrimSpace(s)
	return strings.HasPrefix(cs, ";")
}

func IsEmpty(s string) bool {
	cs := strings.TrimSpace(s)
	return cs == ""
}

// TODO code test
func IsMemBlockDef(s string) bool {
	return strings.ToLower(s) == "ram" || strings.ToLower(s) == "rom"
}
