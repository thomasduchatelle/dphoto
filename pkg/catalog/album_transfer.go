package catalog

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"
)

// TransferredMedias is a list of all medias that has be transferred to a different album in the state.
type TransferredMedias struct {
	Transfers  map[AlbumId][]MediaId
	FromAlbums []AlbumId // FromAlbums is a list of potential origins of medias
}

func (t TransferredMedias) IsEmpty() bool {
	count := 0
	for _, ids := range t.Transfers {
		count += len(ids)
	}

	return count == 0
}

func NewTransferredMedias() TransferredMedias {
	return TransferredMedias{
		Transfers: make(map[AlbumId][]MediaId),
	}
}

// TimelineMutationObserver will notify each observer that medias has been transferred to a different album.
type TimelineMutationObserver interface {
	OnTransferredMedias(ctx context.Context, transfers TransferredMedias) error
}

type TimelineMutationObserverFunc func(ctx context.Context, transfers TransferredMedias) error

func (f TimelineMutationObserverFunc) OnTransferredMedias(ctx context.Context, transfers TransferredMedias) error {
	return f(ctx, transfers)
}

// MediaTransferRecords is a description of all medias that needs to be moved accordingly to the Timeline change
type MediaTransferRecords map[AlbumId][]MediaSelector

func (r MediaTransferRecords) String() string {
	if len(r) == 0 {
		return "<no media to transfer>"
	}

	var transfer []string
	for albumId, selectors := range r {
		transfer = append(transfer, fmt.Sprintf("%s<=%s", albumId, selectors))
	}
	return strings.Join(transfer, " ; ")
}

type MediaSelector struct {
	//ExclusiveAlbum *AlbumId  // ExclusiveAlbum is the Album in which medias are NOT (optional)
	FromAlbums []AlbumId // FromAlbums is a list of potential origins of medias ; is mandatory on CreateAlbum case because media are not indexed by date, only per album.
	Start      time.Time // Start is the first date of matching medias, included
	End        time.Time // End is the last date of matching media, excluded at the second
}

func (m MediaSelector) String() string {
	var from []string
	for _, album := range m.FromAlbums {
		from = append(from, album.String())
	}
	return fmt.Sprintf("{from:%s} %s -> %s", strings.Join(from, ","), m.Start.Format(time.DateTime), m.End.Format(time.DateTime))
}

type MediaTransfer interface {
	Transfer(ctx context.Context, records MediaTransferRecords) error
}

type MediaTransferFunc func(ctx context.Context, records MediaTransferRecords) error

func (f MediaTransferFunc) Transfer(ctx context.Context, records MediaTransferRecords) error {
	return f(ctx, records)
}

type TransferMediasRepositoryPort interface {
	TransferMediasFromRecords(ctx context.Context, records MediaTransferRecords) (TransferredMedias, error)
}

type TransferMediasFunc func(ctx context.Context, records MediaTransferRecords) (TransferredMedias, error)

func (f TransferMediasFunc) TransferMediasFromRecords(ctx context.Context, records MediaTransferRecords) (TransferredMedias, error) {
	return f(ctx, records)
}

type MediaTransferExecutor struct {
	TransferMediasRepository  TransferMediasRepositoryPort
	TimelineMutationObservers []TimelineMutationObserver
}

func (d *MediaTransferExecutor) Transfer(ctx context.Context, records MediaTransferRecords) error {
	transfers, err := d.TransferMediasRepository.TransferMediasFromRecords(ctx, records)
	if err != nil || transfers.IsEmpty() {
		return err
	}

	for _, selectors := range records {
		for _, selector := range selectors {
			for _, origin := range selector.FromAlbums {
				if _, isADestination := transfers.Transfers[origin]; !isADestination && !slices.Contains(transfers.FromAlbums, origin) {
					transfers.FromAlbums = append(transfers.FromAlbums, origin)
				}
			}
		}
	}

	for _, observer := range d.TimelineMutationObservers {
		err = observer.OnTransferredMedias(ctx, transfers)
		if err != nil {
			return err
		}
	}

	return nil
}
