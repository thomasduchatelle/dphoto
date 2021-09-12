package ui

import (
	"fmt"
)

type StaticSession struct {
	existingRepository   ExistingRecordRepositoryPort
	suggestionRepository SuggestionRecordRepositoryPort
	renderer             *recordsRenderer
}

func NewSimpleSession(existingRepository ExistingRecordRepositoryPort, suggestionRepository SuggestionRecordRepositoryPort) *StaticSession {
	return &StaticSession{
		existingRepository:   existingRepository,
		suggestionRepository: suggestionRepository,
		renderer:             new(recordsRenderer),
	}
}

func (s *StaticSession) Render() error {
	existing, err := s.existingRepository.FindExistingRecords()
	if err != nil {
		return err
	}

	suggestions := s.suggestionRepository.FindSuggestionRecords()

	content, err := s.renderer.Render(&RecordsState{
		Records:  createFlattenTree(existing, suggestions),
		Rejected: s.suggestionRepository.Rejects(),
		Selected: -1,
	})

	fmt.Println()
	fmt.Println(content)

	return err
}
