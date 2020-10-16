package main

import (
	"flag"
	"fmt"
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/pointerAnalysis"
	"github.com/amit-davidson/Chronos/ssaUtils"
	"golang.org/x/tools/go/ssa"
	"os"
)

func main() {
	file := flag.String("file", "", "The file containing the entry point of the function - main.go")
	flag.Parse()
	if *file == "" {
		fmt.Printf("Please provide a file to load\n")
		os.Exit(1)
	}
	ssaProg, ssaPkg, err := ssaUtils.LoadPackage(*file)
	if err != nil {
		fmt.Printf("Failed loading with the following error:%s\n", err)
		os.Exit(1)
	}
	ssaUtils.SetGlobalProgram(ssaProg)

	entryFunc := ssaPkg.Func("main")
	entryCallCommon := ssa.CallCommon{Value: entryFunc}
	functionState := ssaUtils.HandleCallCommon(domain.NewEmptyContext(), &entryCallCommon, entryFunc.Pos())
	err = pointerAnalysis.Analysis(ssaPkg, ssaProg, functionState.GuardedAccesses)
	if err != nil {
		fmt.Printf("Error in analysis:%s\n", err)
		os.Exit(1)
	}

}
