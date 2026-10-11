package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/filesystemvolume"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

var miniaturesArgs = struct {
	owner string
}{}

var miniaturesCmd = &cobra.Command{
	Use:   "miniatures [folder]",
	Short: "Generate miniatures for images found in a local directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if miniaturesArgs.owner == "" {
			printer.ErrorText("--owner is mandatory")
			os.Exit(1)
		}

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
					Owner:   miniaturesArgs.owner,
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

func newSmartVolume(volumePath string) (backup.SourceVolume, error) {
	if strings.HasPrefix(volumePath, "s3://") {
		return nil, errors.New("s3:// volumes are not supported by dphotops miniatures")
	}
	return filesystemvolume.New(volumePath), nil
}

func findMediaIDFromSignature(analysedMedia *backup.AnalysedMedia) (catalog.MediaId, error) {
	signature := catalog.MediaSignature{
		SignatureSha256: analysedMedia.Sha256Hash,
		SignatureSize:   analysedMedia.FoundMedia.Size(),
	}

	return catalog.GenerateMediaId(signature)
}

func init() {
	rootCmd.AddCommand(miniaturesCmd)
	miniaturesCmd.Flags().StringVar(&miniaturesArgs.owner, "owner", "", "owner (email) that will own the generated miniatures - mandatory")
}
