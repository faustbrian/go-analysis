package globalgoroutine_test

import (
	"testing"

	"github.com/faustbrian/go-analysis/analyzers/globalgoroutine"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	analysistest.Run(
		t,
		analysistest.TestData(),
		globalgoroutine.Analyzer,
		"globalgoroutine",
		"globalother",
	)
}
