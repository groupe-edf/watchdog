package output

import (
	"fmt"
	"strings"

	"github.com/groupe-edf/watchdog/internal/core/models"
)

type Database struct {
	Channel chan models.AnalysisResult
	Store   models.Store
}

func (writer *Database) WriteTo() {
	for {
		if result, ok := <-writer.Channel; ok {
			fmt.Printf("|_ %v · %v · (%v)\n", result.Commit.Hash, strings.Split(result.Commit.Subject, "\n")[0], result.ElapsedTime)
		} else {
			break
		}
	}
}
