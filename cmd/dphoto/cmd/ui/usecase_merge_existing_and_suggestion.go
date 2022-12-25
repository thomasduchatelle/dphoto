package ui

import (
	"sort"
	"time"
)

// createFlattenTree merges existing and suggestion records into a flatten tree that can be rendered as a table
func createFlattenTree(existing []*ExistingRecord, suggestions []*SuggestionRecord) []*Record {
	nodes := buildRecordTree(existing)

	for _, suggestion := range suggestions {
		totalCount := 0
		distribution := make(map[string]int)
		for k, v := range suggestion.Distribution {
			distribution[k] = v
			totalCount += v
		}

		for _, node := range nodes {
			count := 0
			for day, numberInTheDay := range suggestion.Distribution {
				if _, present := node.activeDays[day]; present {
					count += numberInTheDay
					delete(distribution, day) // do not modify the map which is iterated
				}
			}

			if count > 0 {
				node.children = append(node.children, &Record{
					Indent:               1,
					Suggestion:           true,
					FolderName:           suggestion.FolderName,
					Name:                 suggestion.Name,
					Start:                suggestion.Start,
					End:                  suggestion.End,
					Count:                count,
					TotalCount:           totalCount,
					ParentExistingRecord: node.ExistingRecord,
					SuggestionRecord:     suggestion,
				})
			}
		}

		left := 0
		for _, d := range distribution {
			left += d
		}

		if left > 0 {
			nodes = append(nodes, &recordNode{
				record: &Record{
					Indent:           0,
					Suggestion:       true,
					FolderName:       suggestion.FolderName,
					Name:             suggestion.Name,
					Start:            suggestion.Start,
					End:              suggestion.End,
					Count:            left,
					TotalCount:       totalCount,
					SuggestionRecord: suggestion,
				},
			})
		}
	}

	var records []*Record
	sort.Slice(nodes, func(i, j int) bool {
		if !nodes[i].record.Start.Equal(nodes[j].record.Start) {
			return nodes[i].record.Start.After(nodes[j].record.Start)
		}

		return nodes[i].record.End.After(nodes[j].record.End)
	})

	for _, node := range nodes {
		records = append(records, node.record)

		sort.Slice(node.children, func(i, j int) bool {
			if !node.children[i].Start.Equal(node.children[j].Start) {
				return node.children[i].Start.After(node.children[j].Start)
			}

			return node.children[i].End.After(node.children[j].End)
		})
		for _, child := range node.children {
			records = append(records, child)
		}
	}

	return records
}

func buildRecordTree(existing []*ExistingRecord) []*recordNode {
	nodes := make([]*recordNode, len(existing))
	for idx, alb := range existing {
		activeDays := make(map[string]interface{})
		for _, period := range alb.ActivePeriods {
			startOfTheDay := time.Date(period.Start.Year(), period.Start.Month(), period.Start.Day(), 0, 0, 0, 0, time.UTC)
			for it := startOfTheDay; it.Before(period.End); it = it.Add(24 * time.Hour) {
				activeDays[it.Format("2006-01-02")] = nil
			}
		}
		nodes[idx] = &recordNode{
			ExistingRecord: alb,
			record: &Record{
				Indent:     0,
				Suggestion: false,
				FolderName: alb.FolderName,
				Name:       alb.Name,
				Start:      alb.Start,
				End:        alb.End,
				Count:      alb.Count,
				TotalCount: alb.Count,
			},
			activeDays: activeDays,
		}
	}

	return nodes
}
