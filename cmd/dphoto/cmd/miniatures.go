package cmd

import (
	"context"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"os"

	"github.com/spf13/cobra"
)

var miniaturesCmd = &cobra.Command{
	Use:   "miniatures",
	Short: "Generate miniatures for images found in a local directory",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		factory.InitArchive(ctx)

		folder := args[0]
		volume, err := newSmartVolume(folder)
		if err != nil {
			printer.ErrorText(err.Error())
			os.Exit(1)
		}

		medias, err := volume.FindMedias(ctx)
		if err != nil {
			printer.ErrorText(err.Error())
			os.Exit(2)
		}

		analyser := &backup.AnalyserFromMediaDetails{
			DetailsReaders: analysers.ListDetailReaders(),
		}

		var images []*archive.ImageToResize
		for _, media := range medias {
			analysedMedia, err := analyser.Analyse(ctx, media)
			if err != nil {
				printer.ErrorText(err.Error())
				os.Exit(4)
			}

			if analysedMedia.Type == backup.MediaTypeImage {
				mediaId, err := findMediaIDFromSignature(analysedMedia)
				if err != nil {
					printer.ErrorText(err.Error())
					os.Exit(5)
				}

				images = append(images, &archive.ImageToResize{
					Owner:   Owner,
					MediaId: string(mediaId),
					Widths:  []int{archive.MiniatureCachedWidth},
					Open:    media.ReadMedia,
				})
			}
		}
		cache, err := archive.LoadImagesInCache(ctx, images...)
		if err != nil {
			printer.ErrorText(err.Error())
			os.Exit(3)
		}
		printer.Success("Generated miniatures for %d images", cache)
	},
}

func findMediaIDFromSignature(analysedMedia *backup.AnalysedMedia) (catalog.MediaId, error) {
	signature := catalog.MediaSignature{
		SignatureSha256: analysedMedia.Sha256Hash,
		SignatureSize:   analysedMedia.FoundMedia.Size(),
	}

	return catalog.GenerateMediaId(signature)
}

func init() {
	opsCmd.AddCommand(miniaturesCmd)
}
