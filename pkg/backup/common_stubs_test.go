package backup_test

import "github.com/thomasduchatelle/dphoto/pkg/backup"

func newIdGeneratorWithExclusion(accept func(name string) bool) func(_ string, medias []*backup.AnalysedMedia) map[*backup.AnalysedMedia]string {
	return func(_ string, medias []*backup.AnalysedMedia) map[*backup.AnalysedMedia]string {
		ids := make(map[*backup.AnalysedMedia]string)
		for _, media := range medias {
			filename := media.FoundMedia.MediaPath().Filename
			if accept(filename) {
				ids[media] = "id_" + filename
			}
		}
		return ids
	}
}

type mockVolume []backup.FoundMedia

func (m *mockVolume) String() string {
	return "Mocked Volume"
}

func (m *mockVolume) FindMedias() ([]backup.FoundMedia, error) {
	return *m, nil
}

func (m *mockVolume) Children(backup.MediaPath) (backup.SourceVolume, error) {
	panic("mockVolume cannot generate procreate")
}
