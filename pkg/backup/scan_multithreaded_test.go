package backup

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

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
			m := newMultiThreadedController(tt.newControllerArgs.concurrencyParameters, tt.newControllerArgs.monitoringIntegrator)
			launcher := m.Launcher(tt.launcherArgs.analyser, stubAnalyserObserverChainForMultithreadedTests(context.Background(), m, scanningOptions{
				Options:    Options{},
				analyser:   nil,
				cataloguer: tt.launcherArgs.cataloguer,
			}), new(ScanCompleteObserverFake))

			err := <-launcher.Process(context.Background(), tt.args.volume)
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
