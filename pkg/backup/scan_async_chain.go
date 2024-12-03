package backup

//import (
//	"context"
//	"github.com/thomasduchatelle/dphoto/pkg/backup/chain"
//)
//
//type chainController struct {
//	ConcurrencyParameters         ConcurrencyParameters
//	Analyser                      Analyser
//	analyserAnalysedMediaObserver func(analysed) AnalysedMediaObserver
//	analyserRejectsObserver       RejectedMediaObserver
//}
//
//func (m *chainController) Launcher(analyser Analyser, chain2 *analyserObserverChain, tracker scanCompleteObserver) scanningLauncher {
//	launcher := &chain.SliceLauncher[FoundMedia]{
//		Producer: volume.FindMedias,
//		Next: &chain.MultithreadedLink[FoundMedia, *AnalysedMedia]{
//			NumberOfRoutines: l.parent.,
//			Operator:         nil,
//			Next:             nil,
//		},
//	}
//}
//
//type multiCatalogReferencerObserver struct {
//	Observers []CatalogReferencerObserver
//}
//
//func (m *multiCatalogReferencerObserver) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
//	for _, observer := range m.Observers {
//		err := observer.OnMediaCatalogued(ctx, requests)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
