package ui

import "sort"

type recordRepositoryAggregator struct {
	repositories []RecordRepositoryPort
}

func NewRepositoryAggregator(repositories ...RecordRepositoryPort) RecordRepositoryPort {
	return &recordRepositoryAggregator{repositories: repositories}
}

func (r *recordRepositoryAggregator) FindRecords() ([]*Record, error) {
	var records []*Record
	for _, repo := range r.repositories {
		slice, err := repo.FindRecords()
		if err != nil {
			return nil, err
		}
		records = append(records, slice...)
	}

	sort.Slice(records, func(i, j int) bool {
		if records[i].Start == records[j].Start {
			return records[i].End.Before(records[j].End)
		}

		return records[i].Start.Before(records[j].Start)
	})

	return records, nil
}
