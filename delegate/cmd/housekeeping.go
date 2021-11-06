package cmd

import (
	"github.com/thomasduchatelle/dphoto/delegate/backup"
	"github.com/thomasduchatelle/dphoto/delegate/catalog"
	"github.com/thomasduchatelle/dphoto/delegate/cmd/printer"
	"github.com/thomasduchatelle/dphoto/delegate/cmd/screen"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"strings"
)

var (
	housekeepingArgs = struct {
		limit int
		list  bool
	}{}
)
var housekeepingCmd = &cobra.Command{
	Use:   "housekeeping",
	Short: "Run housekeeping script to perform delayed operations",
	Long:  "Run housekeeping script to perform delayed operations",
	Run: func(cmd *cobra.Command, args []string) {
		transactions, err := catalog.FindMoveTransactions()
		printer.FatalIfError(err, 1)

		if len(transactions) == 0 {
			printer.Success("No more media to move, all done.")
			return
		}

		if housekeepingArgs.list {
			printTransactions(transactions)

		} else {
			if housekeepingArgs.limit > 0 && housekeepingArgs.limit < len(transactions) {
				transactions = transactions[:housekeepingArgs.limit]
			}
			err := startHousekeepingTransactions(transactions)
			printer.FatalIfError(err, 2)

			printer.Success("Housekeeping complete.")
		}
	},
}

func startHousekeepingTransactions(transactions []*catalog.MoveTransaction) error {
	total := 0
	for _, t := range transactions {
		total += t.Count
	}

	operator := &housekeepingOperator{
		grandTotal:        total,
		totalTransactions: len(transactions),
	}
	segments := make([]screen.Segment, 3)

	tableGenerator := screen.NewTable(" ", 2, 0, 80, 25)
	operator.globalProgress, segments[0] = screen.NewProgressLine(tableGenerator, "")
	operator.currentProgress, segments[1] = screen.NewProgressLine(tableGenerator, "")
	segments[2], operator.setRollingLog = screen.NewUpdatableSegment("")

	operator.updateGlobalTransactionBar(0)
	operator.refreshScreen = screen.NewAutoRefreshScreen(screen.RenderingOptions{Width: 180}, segments...)

	for _, transaction := range transactions {
		_, err := catalog.RelocateMovedMedias(operator, transaction.TransactionId)
		if err != nil {
			return err
		}
		operator.updateGlobalTransactionBar(transaction.Count)
	}

	operator.refreshScreen.Refresh()
	operator.refreshScreen.Stop()

	return nil
}

func printTransactions(transactions []*catalog.MoveTransaction) {
	table := simpletable.New()
	table.Header = &simpletable.Header{Cells: []*simpletable.Cell{
		{Text: "Transaction ID"},
		{Text: "Count"},
	}}
	table.Body = &simpletable.Body{Cells: make([][]*simpletable.Cell, len(transactions))}

	for i, t := range transactions {
		table.Body.Cells[i] = []*simpletable.Cell{
			{Text: t.TransactionId},
			{Text: fmt.Sprint(t.Count), Align: simpletable.AlignRight},
		}
	}

	fmt.Println(table.String())
}

func init() {
	rootCmd.AddCommand(housekeepingCmd)

	housekeepingCmd.Flags().IntVarP(&housekeepingArgs.limit, "number", "n", 0, "Limit number of transaction to physically move")
	housekeepingCmd.Flags().BoolVarP(&housekeepingArgs.list, "list", "l", false, "Only display the transactions, do not move anything")
}

type housekeepingOperator struct {
	refreshScreen                 *screen.AutoRefreshScreen
	globalProgress                *screen.ProgressLine
	currentProgress               *screen.ProgressLine
	setRollingLog                 func(string)
	lastMoves                     []string
	grandTotal                    int
	countFromCompletedTransaction int
	countTransaction              int
	totalTransactions             int
	currentCount                  int
	currentTotal                  int
}

func (h *housekeepingOperator) updateGlobalTransactionBar(countFromLastCompletedTransaction int) {
	h.countFromCompletedTransaction += countFromLastCompletedTransaction
	h.globalProgress.SetExplanation(fmt.Sprintf("%d / %d transactions", h.countTransaction, h.totalTransactions))
	h.countTransaction++
	if h.countTransaction > h.totalTransactions {
		h.globalProgress.SwapSpinner(1)
		h.currentProgress.SwapSpinner(1)
	}
}

func (h *housekeepingOperator) Move(source, dest catalog.MediaLocation) (string, error) {
	if len(h.lastMoves) > 5 {
		h.lastMoves = h.lastMoves[1:]
	}
	h.lastMoves = append(h.lastMoves, aurora.Gray(10, fmt.Sprintf("%-70s -> %s", source.FolderName+"/"+source.Filename, dest.FolderName)).String())
	h.setRollingLog(strings.Join(h.lastMoves, "\n"))

	defer func() {
		h.currentCount++
		h.updateProgress(h.currentCount, h.currentTotal)
	}()
	return backup.MovePhysicalStorage(Owner, source.FolderName, source.Filename, dest.FolderName)
}

func (h *housekeepingOperator) UpdateStatus(done, total int) error {
	h.currentCount = done
	h.currentTotal = total

	h.updateProgress(done, total)

	return nil
}

func (h *housekeepingOperator) updateProgress(done int, total int) {
	h.globalProgress.SetBar(uint(done+h.countFromCompletedTransaction), uint(h.grandTotal))
	h.currentProgress.SetBar(uint(done), uint(total))
	h.currentProgress.SetExplanation(fmt.Sprintf("%d / %d medias", done, total))
}

func (h *housekeepingOperator) Continue() bool {
	return true
}
