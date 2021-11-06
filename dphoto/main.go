/*
Copyright © 2020 Thomas Duchatelle <duchatelle.thomas@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	_ "github.com/thomasduchatelle/dphoto/dphoto/backup/adapters"
	_ "github.com/thomasduchatelle/dphoto/dphoto/catalog/adapters"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd"
)

func main() {
	cmd.Execute()
}
