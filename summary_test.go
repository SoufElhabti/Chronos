package main

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/ssaUtils"
	"StaticRaceDetector/testutils"
	"StaticRaceDetector/utils"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/ssa"
	"testing"
)

var shouldUpdate = false

func TestGetFunctionSummary(t *testing.T) {
	var testCases = []struct {
		name     string
		testPath string
		resPath  string
	}{
		{name: "Lock", testPath: "testutils/Lock/prog1.go", resPath: "testutils/Lock/prog1_expected.json"},
		{name: "LockAndUnlock", testPath: "testutils/LockAndUnlock/prog1.go", resPath: "testutils/LockAndUnlock/prog1_expected.json"},
		{name: "LockAndUnlockIfBranch", testPath: "testutils/LockAndUnlockIfBranch/prog1.go", resPath: "testutils/LockAndUnlockIfBranch/prog1_expected.json"},
		{name: "DeferredLockAndUnlockIfBranch", testPath: "testutils/DeferredLockAndUnlockIfBranch/prog1.go", resPath: "testutils/DeferredLockAndUnlockIfBranch/prog1_expected.json"},
		{name: "NestedDeferWithLockAndUnlock", testPath: "testutils/NestedDeferWithLockAndUnlock/prog1.go", resPath: "testutils/NestedDeferWithLockAndUnlock/prog1_expected.json"},
		{name: "NestedDeferWithLockAndUnlockAndGoroutine", testPath: "testutils/NestedDeferWithLockAndUnlockAndGoroutine/prog1.go", resPath: "testutils/NestedDeferWithLockAndUnlockAndGoroutine/prog1_expected.json"},
		{name: "LockAndUnlockIfMap", testPath: "testutils/LockAndUnlockIfMap/prog1.go", resPath: "testutils/LockAndUnlockIfMap/prog1_expected.json"},
		{name: "NestedFunctions", testPath: "testutils/NestedFunctions/prog1.go", resPath: "testutils/NestedFunctions/prog1_expected.json"},
		{name: "DataRaceMap", testPath: "testutils/DataRaceMap/prog1.go", resPath: "testutils/DataRaceMap/prog1_expected.json"},
		{name: "DataRaceShadowedErr", testPath: "testutils/DataRaceShadowedErr/prog1.go", resPath: "testutils/DataRaceShadowedErr/prog1_expected.json"},
		{name: "DataRaceProperty", testPath: "testutils/DataRaceProperty/prog1.go", resPath: "testutils/DataRaceProperty/prog1_expected.json"},
		{name: "DataRaceWithOnlyAlloc", testPath: "testutils/DataRaceWithOnlyAlloc/prog1.go", resPath: "testutils/DataRaceWithOnlyAlloc/prog1_expected.json"},
		{name: "DataRaceWithSameFunction", testPath: "testutils/DataRaceWithSameFunction/prog1.go", resPath: "testutils/DataRaceWithSameFunction/prog1_expected.json"},
		{name: "StructMethod", testPath: "testutils/StructMethod/prog1.go", resPath: "testutils/StructMethod/prog1_expected.json"},
		//{name: "DataRaceIceCreamMaker", testPath: "testutils/DataRaceIceCreamMaker/prog1.go", resPath: "testutils/DataRaceIceCreamMaker/prog1_expected.json"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			domain.GoroutineCounter.Reset()
			domain.GuardedAccessCounter.Reset()

			ssaProg, ssaPkg, err := ssaUtils.LoadPackage(tc.testPath)
			require.NoError(t, err)
			//ssaUtils.SetGlobalProgram(ssaProg)

			entryFunc := ssaPkg.Func("main")
			entryCallCommon := ssa.CallCommon{Value: entryFunc}
			guardedAccesses := make([]*domain.GuardedAccess, 0)
			goroutineState := ssaUtils.GetSummary(&guardedAccesses, domain.NewEmptyGoroutineState(), &entryCallCommon)
			lsRet := goroutineState.Lockset
			testresult := testutils.TestResult{Lockset: lsRet, GuardedAccess: guardedAccesses}
			dump, err := json.MarshalIndent(testresult, "", "\t")
			require.NoError(t, err)
			if shouldUpdate {
				utils.UpdateFile(t, tc.resPath, dump)
			}
			testutils.CompareResult(t, tc.resPath, lsRet, guardedAccesses)
			Analysis(ssaPkg, ssaProg, guardedAccesses)
		})
	}
}
