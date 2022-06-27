package backup

import (
	"time"
)

// ScannedFolder represents a (sub)folder in the scanned target
type ScannedFolder struct {
	Name         string
	RelativePath string                  // RelativePath can be used for display purpose
	FolderName   string                  // FolderName is the original folder name (Name with date prefix that have been removed)
	AbsolutePath string                  // AbsolutePath is used to create a new SourceVolume
	Start, End   time.Time               // Start and End are the beginning of the day of the first media, and the beginning of the day following the last media.
	Distribution map[string]MediaCounter // Distribution is the number of media found for each day (format YYYY-MM-DD)
}

// PushBoundaries is updating the ScannedFolder dates, and update the counter.
func (f *ScannedFolder) PushBoundaries(date time.Time, size int) {
	if f.Start.IsZero() || f.Start.After(date) {
		f.Start = atStartOfDay(date)
	}

	if f.End.IsZero() || !f.End.After(date) {
		f.End = atStartOfFollowingDay(date)
	}

	day := distributionKey(date)

	current, _ := f.Distribution[day]
	f.Distribution[day] = current.Add(1, size)
}

func distributionKey(day time.Time) string {
	return day.Format("2006-01-02")
}

func atStartOfDay(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}

func atStartOfFollowingDay(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day+1, 0, 0, 0, 0, value.Location())
}
