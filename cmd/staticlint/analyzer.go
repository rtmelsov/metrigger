package staticlint

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"honnef.co/go/tools/staticcheck"
)

// Multichecker тип для анализа
var Multichecker = &analysis.Analyzer{
	Name: "mutichecker",
	Doc:  "checks my project to errors",
	Run:  run,
}

// Run функция для запуска проверки кода
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			return false
		})
	}
	return nil, nil
}

var checks = map[string]bool{
	// SA1000 - SA1032
	"SA1000": true, "SA1001": true, "SA1002": true, "SA1003": true, "SA1004": true,
	"SA1005": true, "SA1006": true, "SA1007": true, "SA1008": true, "SA1010": true,
	"SA1011": true, "SA1012": true, "SA1013": true, "SA1014": true, "SA1015": true,
	"SA1016": true, "SA1017": true, "SA1018": true, "SA1019": true, "SA1020": true,
	"SA1021": true, "SA1023": true, "SA1024": true, "SA1025": true, "SA1026": true,
	"SA1027": true, "SA1028": true, "SA1029": true, "SA1030": true, "SA1031": true,
	"SA1032": true,

	// SA2000 - SA2003
	"SA2000": true, "SA2001": true, "SA2002": true, "SA2003": true,

	// SA3000 - SA3001
	"SA3000": true, "SA3001": true,

	// SA4000 - SA4032
	"SA4000": true, "SA4001": true, "SA4003": true, "SA4004": true, "SA4005": true,
	"SA4006": true, "SA4008": true, "SA4009": true, "SA4010": true, "SA4011": true,
	"SA4012": true, "SA4013": true, "SA4014": true, "SA4015": true, "SA4016": true,
	"SA4017": true, "SA4018": true, "SA4019": true, "SA4020": true, "SA4021": true,
	"SA4022": true, "SA4023": true, "SA4024": true, "SA4025": true, "SA4026": true,
	"SA4027": true, "SA4028": true, "SA4029": true, "SA4030": true, "SA4031": true,
	"SA4032": true,

	// SA5000 - SA5012
	"SA5000": true, "SA5001": true, "SA5002": true, "SA5003": true, "SA5004": true,
	"SA5005": true, "SA5007": true, "SA5008": true, "SA5009": true, "SA5010": true,
	"SA5011": true, "SA5012": true,

	// SA6000 - SA6006
	"SA6000": true, "SA6001": true, "SA6002": true, "SA6003": true,
	"SA6005": true, "SA6006": true,

	// SA9001 - SA9009
	"SA9001": true, "SA9002": true, "SA9003": true, "SA9004": true, "SA9005": true,
	"SA9006": true, "SA9007": true, "SA9008": true, "SA9009": true,
	// OTHER
	"ST1000": true,
}

func getChecks() []*analysis.Analyzer {
	var mychecks = []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		Multichecker,
		NoExistAnalyzer,
	}
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	return mychecks
}

func Check() {
	multichecker.Main(
		getChecks()...,
	)
}
