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

func stubAnalyserObserverChainForMultithreadedTests(ctx context.Context, controller scanningController, options scanningOptions) *analyserObserverChain {
	return &analyserObserverChain{
		AnalysedMediaObservers: []AnalysedMediaObserver{
			controller.bufferAnalysedMedia(ctx, &analyserToCatalogReferencer{
				CatalogReferencer:          options.cataloguer,
				CatalogReferencerObservers: controller.AppendPreCataloguerFilter(controller.AppendPostCatalogFiltersIn()...),
			}),
		},
		RejectedMediaObservers: controller.AppendPostAnalyserFilterRejects(),
	}
}

func Test_multiThreadedController_Launcher(t *testing.T) {
	simulatedError := errors.New("TEST Simulated error")

	groupThatMustWaitForAProcessInEachStep := newWaitGroup(3)

	type newControllerArgs struct {
		concurrencyParameters ConcurrencyParameters
		monitoringIntegrator  *scanListeners
		bufferSize            int
	}
	type launcherArgs struct {
		analyser   Analyser
		cataloguer Cataloguer
	}
	type args struct {
		volume SourceVolume
	}
	simpleReferencerForFile1to3 := CatalogReferencerFakeByName{
		"file-1": nonExistingSimpleReference("file-1"),
		"file-2": nonExistingSimpleReference("file-2"),
		"file-3": nonExistingSimpleReference("file-3"),
	}
	tests := []struct {
		name              string
		newControllerArgs newControllerArgs
		launcherArgs      launcherArgs
		args              args
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "it should run through with an empty volume and default values",
			newControllerArgs: newControllerArgs{
				concurrencyParameters: ConcurrencyParameters{},
				monitoringIntegrator:  new(scanListeners),
				bufferSize:            1,
			},
			launcherArgs: launcherArgs{
				analyser:   &AnalyserFake{},
				cataloguer: CatalogReferencerFakeByName{},
			},
			args:    args{volume: &InMemorySourceVolume{}},
			wantErr: assert.NoError,
		},
		{
			name: "it should run with 2 analysers concurrently",
			newControllerArgs: newControllerArgs{
				concurrencyParameters: ConcurrencyParameters{ConcurrentAnalyserRoutines: 2},
				monitoringIntegrator:  new(scanListeners),
				bufferSize:            1,
			},
			launcherArgs: launcherArgs{
				analyser: &AnalyserGroupWaiter{
					group:    newWaitGroup(2),
					delegate: &AnalyserFake{},
				},
				cataloguer: simpleReferencerForFile1to3,
			},
			args: args{volume: &InMemorySourceVolume{
				NewInMemoryMedia("file-1", time.Now(), []byte{}),
				NewInMemoryMedia("file-2", time.Now(), []byte{}),
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should run the cataloguer on two routines",
			newControllerArgs: newControllerArgs{
				concurrencyParameters: ConcurrencyParameters{ConcurrentCataloguerRoutines: 2},
				monitoringIntegrator:  new(scanListeners),
				bufferSize:            1,
			},
			launcherArgs: launcherArgs{
				analyser: &AnalyserFake{},
				cataloguer: &CataloguerGroupWaiter{
					delegate: simpleReferencerForFile1to3,
					group:    newWaitGroup(2),
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
			newControllerArgs: newControllerArgs{
				concurrencyParameters: ConcurrencyParameters{
					ConcurrentAnalyserRoutines:   2,
					ConcurrentCataloguerRoutines: 2,
					ConcurrentUploaderRoutines:   2,
				},
				monitoringIntegrator: &scanListeners{
					PostCatalogFiltersIn: []CatalogReferencerObserver{
						&SingleThreadedConstrainedCatalogReferencerObserver{
							lock:          sync.Mutex{},
							expectedCalls: 2,
						},
					},
				},
				bufferSize: 1,
			},
			launcherArgs: launcherArgs{
				analyser:   new(AnalyserFake),
				cataloguer: simpleReferencerForFile1to3,
			},
			args: args{volume: &InMemorySourceVolume{
				NewInMemoryMedia("file-1", time.Now(), []byte{}),
				NewInMemoryMedia("file-2", time.Now(), []byte{}),
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should run the analyser, the cataloguer, and the post-scan on 3 different threads",
			newControllerArgs: newControllerArgs{
				concurrencyParameters: ConcurrencyParameters{},
				monitoringIntegrator: &scanListeners{
					PostCatalogFiltersIn: []CatalogReferencerObserver{&ScanningCompleteGroupWaiter{
						lock:          sync.RWMutex{},
						group:         groupThatMustWaitForAProcessInEachStep,
						allowedToPass: []string{"file-2", "file-3"},
						expectedCalls: 1,
					}},
				},
				bufferSize: 1,
			},
			launcherArgs: launcherArgs{
				analyser: &AnalyserGroupWaiter{
					allowedToPass: []string{"file-1", "file-2"},
					group:         groupThatMustWaitForAProcessInEachStep,
					delegate:      new(AnalyserFake),
				},
				cataloguer: &CataloguerGroupWaiter{
					allowedToPass: []string{"file-1", "file-3"},
					delegate:      simpleReferencerForFile1to3,
					group:         groupThatMustWaitForAProcessInEachStep,
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
			newControllerArgs: newControllerArgs{
				concurrencyParameters: ConcurrencyParameters{ConcurrentCataloguerRoutines: 2},
				monitoringIntegrator:  &scanListeners{},
				bufferSize:            1,
			},
			launcherArgs: launcherArgs{
				analyser: new(AnalyserFake),
				cataloguer: &CataloguerGroupWaiter{
					delegate: simpleReferencerForFile1to3,
					group:    newWaitGroup(2),
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
			newControllerArgs: newControllerArgs{
				concurrencyParameters: ConcurrencyParameters{},
				monitoringIntegrator: &scanListeners{
					PostCatalogFiltersIn: []CatalogReferencerObserver{
						&ScanAssertEndOfHappyPath{
							WantMaxReadyToUploadCount: 10, // usually 0 or 1 ; sometime 2-3 ; saw once 6 and 7.
						},
					},
				},
				bufferSize: 1,
			},
			launcherArgs: launcherArgs{
				analyser: &AnalyserFake{
					ErroredFilename: map[string]error{
						"file-1": simulatedError,
					},
				},
				cataloguer: simpleReferencerForFile1to3,
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
			m := newMultiThreadedController(tt.newControllerArgs.concurrencyParameters, tt.newControllerArgs.monitoringIntegrator, tt.newControllerArgs.bufferSize, new(flushableCollector))
			launcher := m.Launcher(tt.launcherArgs.analyser, stubAnalyserObserverChainForMultithreadedTests(context.Background(), m, scanningOptions{
				Options:    Options{},
				analyser:   nil,
				cataloguer: tt.launcherArgs.cataloguer,
			}), new(ScanCompleteObserverFake))

			err := <-launcher.process(context.Background(), tt.args.volume)
			if tt.wantErr(t, err) {
				if toBeSatisfied, match := tt.launcherArgs.analyser.(ToBeSatisfied); match {
					toBeSatisfied.IsSatisfied(t)
				}
				for _, listener := range tt.newControllerArgs.monitoringIntegrator.PostAnalyserSuccess {
					if toBeSatisfied, match := listener.(ToBeSatisfied); match {
						toBeSatisfied.IsSatisfied(t)
					}
				}
				for _, listener := range tt.newControllerArgs.monitoringIntegrator.PostCatalogFiltersIn {
					if toBeSatisfied, match := listener.(ToBeSatisfied); match {
						toBeSatisfied.IsSatisfied(t)
					}
				}
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

type AnalyserGroupWaiter struct {
	allowedToPass []string
	delegate      Analyser
	group         *sync.WaitGroup // group must be EXACTLY the number of file that will be received: LESS -> deadlock ; MORE -> negative waitGroup
}

func (a *AnalyserGroupWaiter) Analyse(ctx context.Context, found FoundMedia, analysedMediaObserver AnalysedMediaObserver, rejectsObserver RejectedMediaObserver) error {
	filename := found.MediaPath().Filename
	if slices.Index(a.allowedToPass, filename) < 0 {
		log.Infof("[%d] AnalyserGroupWaiter > %s placed on hold", goid(), filename)

		a.group.Done()
		a.group.Wait()
	}

	log.Infof("[%d] AnalyserGroupWaiter > %s processed", goid(), filename)
	return a.delegate.Analyse(ctx, found, analysedMediaObserver, rejectsObserver)
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

type AnalysedMediaGroupWaiter struct {
	lock          sync.RWMutex
	group         *sync.WaitGroup
	allowedToPass []string
	expectedCalls int
}

func (a *AnalysedMediaGroupWaiter) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	a.lock.RLock()
	filename := media.FoundMedia.MediaPath().Filename
	block := slices.Index(a.allowedToPass, filename) < 0
	a.lock.RUnlock()

	if block {
		log.Infof("[%d] AnalysedMediaGroupWaiter > %s placed on hold", goid(), filename)

		a.group.Done()
		a.group.Wait()

		a.lock.Lock()
		a.expectedCalls--
		a.lock.Unlock()
	}

	log.Infof("[%d] AnalysedMediaGroupWaiter > %s processed", goid(), filename)
	return nil
}

func (a *AnalysedMediaGroupWaiter) IsSatisfied(t *testing.T) bool {
	return assert.Equal(t, 0, a.expectedCalls, "Invalid number of calls to AnalysedMediaGroupWaiter")
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
