package cmd

import (
	"duchatelle.io/dphoto/dphoto/catalog"
	"duchatelle.io/dphoto/dphoto/cmd/printer"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

const ukDateLayout = "02 Jan 06 15:04"

var (
	listArgs = struct {
		stats bool
	}{}
)

var albumCmd = &cobra.Command{
	Use:     "album [--stats]",
	Aliases: []string{"albums", "alb"},
	Short:   "Organise your collection into albums",
	Long:    `Organise your collection into albums.`,
	Run: func(cmd *cobra.Command, args []string) {

		table := simpletable.New()
		table.SetStyle(simpletable.StyleDefault)
		table.Header = &simpletable.Header{
			Cells: listHeaders(),
		}

		if listArgs.stats {
			albums, err := catalog.FindAllAlbumsWithStats()
			printer.FatalWithMessageIfError(err, 2, "Couldn't fetch albums from DPhoto Catalog")

			table.Body = &simpletable.Body{
				Cells: make([][]*simpletable.Cell, len(albums)),
			}

			for i, a := range albums {
				table.Body.Cells[i] = append(listAlbumRow(&a.Album), &simpletable.Cell{Align: simpletable.AlignRight, Text: fmt.Sprint(a.TotalCount())})
			}

		} else {
			albums, err := catalog.FindAllAlbums()
			printer.FatalWithMessageIfError(err, 1, "Couldn't fetch albums from DPhoto Catalog")

			table.Body = &simpletable.Body{
				Cells: make([][]*simpletable.Cell, len(albums)),
			}

			for i, a := range albums {
				table.Body.Cells[i] = listAlbumRow(a)
			}
		}

		if len(table.Body.Cells) == 0 {
			fmt.Println("No album are present.")

		} else {
			fmt.Println(table.String())
		}
	},
}

func listAlbumRow(a *catalog.Album) []*simpletable.Cell {
	return []*simpletable.Cell{
		{Text: a.Start.Format(ukDateLayout)},
		{Text: a.End.Format(ukDateLayout)},
		{Text: aurora.Cyan(a.Name).Bold().String()},
		{Text: aurora.Italic("/" + a.FolderName).String()},
	}
}

func listHeaders() []*simpletable.Cell {
	headers := []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Text: "Start"},
		{Align: simpletable.AlignCenter, Text: "End"},
		{Align: simpletable.AlignLeft, Text: "Name"},
		{Align: simpletable.AlignLeft, Text: "Folder"},
	}

	if listArgs.stats {
		return append(headers, &simpletable.Cell{Align: simpletable.AlignLeft, Text: "Count"})
	}

	return headers
}

func init() {
	rootCmd.AddCommand(albumCmd)

	albumCmd.Flags().BoolVarP(&listArgs.stats, "stats", "s", false, "count number of media in each album")
}
