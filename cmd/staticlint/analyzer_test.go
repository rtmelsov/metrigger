package staticlint

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMyAnalyzer(t *testing.T) {
	testsDir := filepath.Join("..", "..")
	for _, a := range getChecks() {
		analysistest.Run(
			t,
			testsDir,
			a,
			"./...", // относительный путь от testsDir
		)
	}
}
