package backupui

import (
	"github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"
	"github.com/thomasduchatelle/dphoto/delegate/cmd/printer"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/logrusorgru/aurora/v3"
)

func PrintBackupStats(tracker backupmodel.BackupReport, volumePath string) {
	if len(tracker.CountPerAlbum()) == 0 {
		printer.Success("\n\nBackup of %s complete: %s.", aurora.Cyan(volumePath), aurora.Bold(aurora.Yellow("no new medias")))
		return
	}

	printer.Success("\n\nBackup of %s complete\n", aurora.Cyan(volumePath))

	table := simpletable.New()
	table.Header = &simpletable.Header{Cells: []*simpletable.Cell{
		{Text: "New", Align: simpletable.AlignCenter},
		{Text: "Album", Align: simpletable.AlignCenter},
		{Text: "Photo", Align: simpletable.AlignCenter},
		{Text: "Video", Align: simpletable.AlignCenter},
		{Text: "Total", Align: simpletable.AlignCenter},
	}}
	table.Body = &simpletable.Body{Cells: make([][]*simpletable.Cell, len(tracker.CountPerAlbum())+1)}

	newAlbums := make(map[string]interface{})
	for _, album := range tracker.NewAlbums() {
		newAlbums[album] = nil
	}

	var totals [3]backupmodel.MediaCounter
	i := 0
	for folderName, counts := range tracker.CountPerAlbum() {
		newMarker := ""
		if _, ok := newAlbums[folderName]; ok {
			newMarker = "*"
		}

		table.Body.Cells[i] = []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: newMarker},
			{Text: folderName},
			countAndSize(counts.OfType(backupmodel.MediaTypeImage)),
			countAndSize(counts.OfType(backupmodel.MediaTypeVideo)),
			countAndSize(counts.Total()),
		}

		totals[0] = totals[0].AddCounter(counts.OfType(backupmodel.MediaTypeImage))
		totals[1] = totals[1].AddCounter(counts.OfType(backupmodel.MediaTypeVideo))
		totals[2] = totals[2].AddCounter(counts.Total())
		i++
	}

	table.Body.Cells[len(table.Body.Cells)-1] = []*simpletable.Cell{
		{Text: "Totals", Align: simpletable.AlignRight, Span: 2},
		countAndSize(totals[0]),
		countAndSize(totals[1]),
		countAndSize(totals[2]),
	}

	fmt.Println(table.String())
}

func countAndSize(counter backupmodel.MediaCounter) *simpletable.Cell {
	if counter.Count == 0 {
		return &simpletable.Cell{Align: simpletable.AlignCenter, Text: "-"}
	}

	return &simpletable.Cell{
		Text: fmt.Sprintf("%d (%s)", counter.Count, byteCountIEC(counter.Size)),
	}
}

func byteCountIEC(b uint) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
