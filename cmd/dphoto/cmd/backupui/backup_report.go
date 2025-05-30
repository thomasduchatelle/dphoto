package backupui

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/logrusorgru/aurora/v3"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
)

func PrintBackupStats(tracker backup.Report, volumePath string) {
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

	var totals [3]backup.MediaCounter
	i := 0
	for folderName, counts := range tracker.CountPerAlbum() {
		newMarker := ""
		if counts.IsNew() {
			newMarker = "*"
		}

		table.Body.Cells[i] = []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: newMarker},
			{Text: folderName},
			countAndSize(counts.OfType(backup.MediaTypeImage)),
			countAndSize(counts.OfType(backup.MediaTypeVideo)),
			countAndSize(counts.Total()),
		}

		totals[0] = totals[0].AddCounter(counts.OfType(backup.MediaTypeImage))
		totals[1] = totals[1].AddCounter(counts.OfType(backup.MediaTypeVideo))
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

func countAndSize(counter backup.MediaCounter) *simpletable.Cell {
	if counter.Count == 0 {
		return &simpletable.Cell{Align: simpletable.AlignCenter, Text: "-"}
	}

	return &simpletable.Cell{
		Text: fmt.Sprintf("%d (%s)", counter.Count, byteCountIEC(counter.Size)),
	}
}

func byteCountIEC(b int) string {
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
