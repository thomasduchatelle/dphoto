package ui

import "fmt"

type StaticSession struct {
	repository RecordRepositoryPort
	renderer   *recordsRenderer
}

func NewSimpleSession(repositories ...RecordRepositoryPort) *StaticSession {
	return &StaticSession{
		repository: NewRepositoryAggregator(repositories...),
		renderer:   new(recordsRenderer),
	}
}

func (s *StaticSession) Render() error {
	records, err := s.repository.FindRecords()
	if err != nil {
		return err
	}

	content, err := s.renderer.Render(&recordsState{
		Records:  records,
		Selected: -1,
	})

	fmt.Println()
	fmt.Println(content)

	return err
}
