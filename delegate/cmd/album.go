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
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// albumCmd represents the album command
var albumCmd = &cobra.Command{
	Use:   "album",
	Short: "Manage your albums",
	Long: `Manage your albums.

Create and list albums in the collections.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("album called")
	},
}

func init() {
	rootCmd.AddCommand(albumCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// albumCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// albumCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
