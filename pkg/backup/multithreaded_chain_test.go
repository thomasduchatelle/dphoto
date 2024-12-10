package backup

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_multithreadedScanRuntime(t *testing.T) {
	simulatedError := errors.New("TEST Simulated error")

	groupThatMustWaitForAProcessInEachStep := newWaitGroup(3)
	simpleReferencerForFile1to3 := CatalogReferencerFakeByName{
		"file-1": nonExistingSimpleReference("file-1"),
		"file-2": nonExistingSimpleReference("file-2"),
		"file-3": nonExistingSimpleReference("file-3"),
	}

	type fields struct {
		options Options
		config  *scanConfiguration
	}
	type args struct {
		volume SourceVolume
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should run through with an empty volume and default values",
			fields: fields{
				config: &scanConfiguration{
					Analyser:             &AnalyserFake{},
					Cataloguer:           CatalogReferencerFakeByName{},
					ScanCompleteObserver: &ScanCompleteObserverFakeAssert{WantCount: 0, WantSize: 0},
				},
			},
			args:    args{volume: &InMemorySourceVolume{}},
			wantErr: assert.NoError,
		},
		{
			name: "it should notify when the volume has been listed",
			fields: fields{
				config: &scanConfiguration{
					Analyser:             &AnalyserFake{},
					Cataloguer:           simpleReferencerForFile1to3,
					ScanCompleteObserver: &ScanCompleteObserverFakeAssert{WantCount: 2, WantSize: 3},
				},
			},
			args: args{volume: &InMemorySourceVolume{
				NewInMemoryMedia("file-1", time.Now(), []byte("a")),
				NewInMemoryMedia("file-2", time.Now(), []byte("ab")),
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should run with 2 analysers concurrently",
			fields: fields{
				options: ReduceOptions(
					OptionsConcurrentAnalyserRoutines(2),
				),
				config: &scanConfiguration{
					Analyser: &AnalyserGroupWaiter{
						group:    newWaitGroup(2),
						delegate: &AnalyserFake{},
					},
					Cataloguer: simpleReferencerForFile1to3,
				},
			},
			args: args{volume: &InMemorySourceVolume{
				NewInMemoryMedia("file-1", time.Now(), []byte{}),
				NewInMemoryMedia("file-2", time.Now(), []byte{}),
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should run the cataloguer on two routines",
			fields: fields{
				options: ReduceOptions(
					OptionsConcurrentCataloguerRoutines(2),
				),
				config: &scanConfiguration{
					Analyser: &AnalyserFake{},
					Cataloguer: &CataloguerGroupWaiter{
						delegate: simpleReferencerForFile1to3,
						group:    newWaitGroup(2),
					},
				},
			},
			args: args{volume: &InMemorySourceVolume{
				NewInMemoryMedia("file-1", time.Now(), []byte{}),
				NewInMemoryMedia("file-2", time.Now(), []byte{}),
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should have a single thread for post-scan reporting",
			fields: fields{
				options: ReduceOptions(
					OptionsConcurrentAnalyserRoutines(2),
					OptionsConcurrentCataloguerRoutines(2),
					OptionsConcurrentUploaderRoutines(2),
				),
				config: &scanConfiguration{
					Analyser:   new(AnalyserFake),
					Cataloguer: simpleReferencerForFile1to3,
					PostCatalogFiltersIn: []CatalogReferencerObserver{
						&SingleThreadedConstrainedCatalogReferencerObserver{
							lock:          sync.Mutex{},
							expectedCalls: 2,
						},
					},
				},
			},
			args: args{volume: &InMemorySourceVolume{
				NewInMemoryMedia("file-1", time.Now(), []byte{}),
				NewInMemoryMedia("file-2", time.Now(), []byte{}),
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should run the analyser, the cataloguer, and the post-scan on 3 different threads",
			fields: fields{
				config: &scanConfiguration{
					Analyser: &AnalyserGroupWaiter{
						allowedToPass: []string{"file-1", "file-2"},
						group:         groupThatMustWaitForAProcessInEachStep,
						delegate:      new(AnalyserFake),
					},
					Cataloguer: &CataloguerGroupWaiter{
						allowedToPass: []string{"file-1", "file-3"},
						delegate:      simpleReferencerForFile1to3,
						group:         groupThatMustWaitForAProcessInEachStep,
					},
					PostCatalogFiltersIn: []CatalogReferencerObserver{&ScanningCompleteGroupWaiter{
						lock:          sync.RWMutex{},
						group:         groupThatMustWaitForAProcessInEachStep,
						allowedToPass: []string{"file-2", "file-3"},
						expectedCalls: 1,
					}},
				},
			},
			args: args{volume: &InMemorySourceVolume{
				NewInMemoryMedia("file-1", time.Now(), []byte{}),
				NewInMemoryMedia("file-2", time.Now(), []byte{}),
				NewInMemoryMedia("file-3", time.Now(), []byte{}),
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should run with 2 cataloguer concurrently",
			fields: fields{
				options: ReduceOptions(
					OptionsConcurrentCataloguerRoutines(2),
				),
				config: &scanConfiguration{
					Analyser: new(AnalyserFake),
					Cataloguer: &CataloguerGroupWaiter{
						delegate: simpleReferencerForFile1to3,
						group:    newWaitGroup(2),
					},
				},
			},
			args: args{volume: &InMemorySourceVolume{
				NewInMemoryMedia("file-1", time.Now(), []byte{}),
				NewInMemoryMedia("file-2", time.Now(), []byte{}),
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should interrupt the process if an error occur during the analyser",
			fields: fields{
				config: &scanConfiguration{
					Analyser: &AnalyserFake{
						ErroredFilename: map[string]error{
							"file-1": simulatedError,
						},
					},
					Cataloguer:          simpleReferencerForFile1to3,
					PostAnalyserRejects: []RejectedMediaObserver{new(analyserFailsFastObserver)},
					PostCatalogFiltersIn: []CatalogReferencerObserver{
						&ScanAssertEndOfHappyPath{
							WantMaxReadyToUploadCount: 10, // usually 0 or 1 before interruption ; sometime 2-3 ; saw once 6 and 7.
						},
					},
				},
			},
			args: args{volume: &InMemorySourceVolume{
				NewInMemoryMedia("file-1", time.Now(), []byte{}),
				NewInMemoryMedia("file-2", time.Now(), []byte{}),
				NewInMemoryMedia("file-3", time.Now(), []byte{}),
				NewInMemoryMedia("file-3", time.Now(), []byte{}),
				NewInMemoryMedia("file-3", time.Now(), []byte{}),
				NewInMemoryMedia("file-3", time.Now(), []byte{}),
				NewInMemoryMedia("file-3", time.Now(), []byte{}),
				NewInMemoryMedia("file-3", time.Now(), []byte{}),
				NewInMemoryMedia("file-3", time.Now(), []byte{}),
				NewInMemoryMedia("file-3", time.Now(), []byte{}),
				NewInMemoryMedia("file-3", time.Now(), []byte{}),
				NewInMemoryMedia("file-3", time.Now(), []byte{}),
			}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, simulatedError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run, err := multithreadedScanRuntime(context.Background(), tt.fields.options, tt.fields.config)

			if !assert.NoError(t, err) {
				return
			}

			err = <-run.Process(context.Background(), tt.args.volume)
			if !tt.wantErr(t, err) {
				return
			}

			if toBeSatisfied, match := tt.fields.config.Analyser.(ToBeSatisfied); match {
				toBeSatisfied.IsSatisfied(t)
			}
			for _, listener := range tt.fields.config.PostCatalogFiltersIn {
				if toBeSatisfied, match := listener.(ToBeSatisfied); match {
					toBeSatisfied.IsSatisfied(t)
				}
			}
			if toBeSatisfied, match := tt.fields.config.ScanCompleteObserver.(ToBeSatisfied); match {
				toBeSatisfied.IsSatisfied(t)
			}
		})
	}
}

type ToBeSatisfied interface {
	IsSatisfied(t *testing.T) bool
}

func newWaitGroup(size int) *sync.WaitGroup {
	group := &sync.WaitGroup{}
	group.Add(size)
	return group
}

type ScanCompleteObserverFake struct {
	count int
	size  int
}

func (s *ScanCompleteObserverFake) OnScanComplete(ctx context.Context, count, size int) error {
	s.count = count
	s.size = size

	return nil
}

type ScanCompleteObserverFakeAssert struct {
	ScanCompleteObserverFake
	WantCount int
	WantSize  int
}

func (s *ScanCompleteObserverFakeAssert) IsSatisfied(t *testing.T) bool {
	return assert.Equal(t, s.WantCount, s.count, "Invalid count") &&
		assert.Equal(t, s.WantSize, s.size, "Invalid size")

}

type AnalyserGroupWaiter struct {
	allowedToPass []string
	delegate      Analyser
	group         *sync.WaitGroup // group must be EXACTLY the number of file that will be received: LESS -> deadlock ; MORE -> negative waitGroup
}

func (a *AnalyserGroupWaiter) Analyse(ctx context.Context, found FoundMedia) (*AnalysedMedia, error) {
	filename := found.MediaPath().Filename
	if slices.Index(a.allowedToPass, filename) < 0 {
		log.Infof("[%d] AnalyserGroupWaiter > %s placed on hold", goid(), filename)

		a.group.Done()
		a.group.Wait()
	}

	log.Infof("[%d] AnalyserGroupWaiter > %s processed", goid(), filename)
	return a.delegate.Analyse(ctx, found)
}

type CataloguerGroupWaiter struct {
	allowedToPass []string
	delegate      Cataloguer
	group         *sync.WaitGroup // group must be EXACTLY the number of file that will be received: LESS -> deadlock ; MORE -> negative waitGroup
}

func (c *CataloguerGroupWaiter) Reference(ctx context.Context, medias []*AnalysedMedia, observer CatalogReferencerObserver) error {

	hold := false
	var filenames []string
	for _, media := range medias {
		filename := media.FoundMedia.MediaPath().Filename
		filenames = append(filenames, filename)

		if slices.Index(c.allowedToPass, filename) < 0 {
			hold = true
		}
	}

	if hold {
		log.Infof("[%d] CataloguerGroupWaiter > %s files placed on hold", goid(), strings.Join(filenames, ", "))

		c.group.Done()
		c.group.Wait()
	}

	log.Infof("[%d] CataloguerGroupWaiter > %s files processed", goid(), strings.Join(filenames, ", "))
	return c.delegate.Reference(ctx, medias, observer)
}

type SingleThreadedConstrainedCatalogReferencerObserver struct {
	lock          sync.Mutex
	expectedCalls int
}

func (s *SingleThreadedConstrainedCatalogReferencerObserver) IsSatisfied(t *testing.T) bool {
	return assert.Equal(t, 0, s.expectedCalls, "Invalid number of calls to SingleThreadedConstrainedCatalogReferencerObserver")
}

func (s *SingleThreadedConstrainedCatalogReferencerObserver) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	if locked := s.lock.TryLock(); locked {
		defer s.lock.Unlock()

		// all good - wait a bit to see if another one would try to process
		s.expectedCalls--
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	return errors.New("SingleThreadedConstrainedCatalogReferencerObserver.OnMediaCatalogued called while processing a different message")
}

type CatalogReferencerFakeByName map[string]CatalogReference

func (c CatalogReferencerFakeByName) Reference(ctx context.Context, medias []*AnalysedMedia, observer CatalogReferencerObserver) error {
	var result []BackingUpMediaRequest
	for _, media := range medias {
		filename := media.FoundMedia.MediaPath().Filename
		if reference, found := c[filename]; found {
			result = append(result, BackingUpMediaRequest{
				AnalysedMedia:    media,
				CatalogReference: reference,
			})
		} else {
			return errors.Errorf("File with name %s doesn't have a reference.", filename)
		}
	}

	return observer.OnMediaCatalogued(ctx, result)
}

func nonExistingSimpleReference(id string) *SimpleReference {
	return &SimpleReference{
		exists:           false,
		albumCreated:     false,
		albumFolderName:  "simple",
		uniqueIdentifier: id,
		mediaId:          id,
	}
}

type SimpleReference struct {
	exists           bool
	albumCreated     bool
	albumFolderName  string
	uniqueIdentifier string
	mediaId          string
}

func (s *SimpleReference) Exists() bool {
	return s.exists
}

func (s *SimpleReference) AlbumCreated() bool {
	return s.albumCreated
}

func (s *SimpleReference) AlbumFolderName() string {
	return s.albumFolderName
}

func (s *SimpleReference) UniqueIdentifier() string {
	return s.uniqueIdentifier
}

func (s *SimpleReference) MediaId() string {
	return s.mediaId
}

type ScanningCompleteGroupWaiter struct {
	lock          sync.RWMutex
	group         *sync.WaitGroup
	expectedCalls int
	allowedToPass []string
}

func (c *ScanningCompleteGroupWaiter) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	c.lock.RLock()
	mustBlock := false

	var filenames []string
	for _, request := range requests {
		filename := request.AnalysedMedia.FoundMedia.MediaPath().Filename
		if slices.Index(c.allowedToPass, filename) < 0 {
			mustBlock = true
		}

		filenames = append(filenames, filename)
	}

	c.lock.RUnlock()

	if mustBlock {
		log.Infof("[%d] ScanningCompleteGroupWaiter > %s placed on hold", goid(), strings.Join(filenames, ", "))

		c.group.Done()
		c.group.Wait()
		c.lock.Lock()
		c.expectedCalls -= len(requests)
		c.lock.Unlock()
	}

	log.Infof("{%d] ScanningCompleteGroupWaiter > %s processed", goid(), strings.Join(filenames, ", "))

	return nil
}

func (c *ScanningCompleteGroupWaiter) IsSatisfied(t *testing.T) bool {
	return assert.Equal(t, 0, c.expectedCalls, "Invalid number of calls to ScanningCompleteGroupWaiter")
}

type ScanAssertEndOfHappyPath struct {
	WantReadyToUpload         []string
	WantMaxReadyToUploadCount int
	gotReadyToUpload          []string
}

func (s *ScanAssertEndOfHappyPath) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	for _, request := range requests {
		s.gotReadyToUpload = append(s.gotReadyToUpload, request.AnalysedMedia.FoundMedia.MediaPath().Filename)
	}

	return nil
}

func (s *ScanAssertEndOfHappyPath) IsSatisfied(t *testing.T) bool {
	if s.WantMaxReadyToUploadCount > 0 {
		return assert.LessOrEqual(t, len(s.gotReadyToUpload), s.WantMaxReadyToUploadCount)
	}
	return assert.ElementsMatch(t, s.WantReadyToUpload, s.gotReadyToUpload)
}

func goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

type AnalyserFake struct {
	ErroredFilename map[string]error
}

func (a *AnalyserFake) Analyse(ctx context.Context, found FoundMedia) (*AnalysedMedia, error) {
	if a.ErroredFilename != nil {
		if err, errored := a.ErroredFilename[found.MediaPath().Filename]; errored {
			return nil, err
		}
	}

	return &AnalysedMedia{
		FoundMedia: found,
		Type:       MediaTypeImage,
		Sha256Hash: found.MediaPath().Filename,
		Details: &MediaDetails{
			DateTime: time.Date(2022, 6, 18, 10, 42, 0, 0, time.UTC),
		},
	}, nil
}
