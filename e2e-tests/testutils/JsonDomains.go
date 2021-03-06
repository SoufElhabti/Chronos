package testutils

import (
	"encoding/json"
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/ssaUtils"
	"go/token"
	"golang.org/x/tools/go/ssa"
	"strings"
)

type GuardedAccessJSON struct {
	ID          int
	Pos         string
	Value       string
	OpKind      domain.OpKind
	LocksetJson *LocksetJson
	State       *ContextJSON
}

func GuardedAccessToJSON(ga *domain.GuardedAccess) GuardedAccessJSON {
	prog := ssaUtils.GlobalProgram
	dumpJson := GuardedAccessJSON{}
	dumpJson.ID = ga.ID
	dumpJson.Pos = GetRelativePath(ga.Pos, prog)
	dumpJson.Value = GetRelativePath(ga.Value.Pos(), prog)
	dumpJson.OpKind = ga.OpKind
	dumpJson.LocksetJson = LocksetToJSON(ga.Lockset)
	dumpJson.State = ContextToJSON(ga.State)
	return dumpJson
}
func MarshalJSON(ga *domain.GuardedAccess) ([]byte, error) {
	dump, err := json.Marshal(GuardedAccessToJSON(ga))
	return dump, err
}

// Lockset name ans pos
type LocksetJson struct {
	ExistingLocks   map[string]string
	ExistingUnlocks map[string]string
}

func LocksetToJSON(ls *domain.Lockset) *LocksetJson {
	prog := ssaUtils.GlobalProgram

	dumpJson := &LocksetJson{ExistingLocks: make(map[string]string, 0), ExistingUnlocks: make(map[string]string, 0)}
	for lockInit, lock := range ls.ExistingLocks {
		lockInitPos := GetRelativePath(lockInit, prog)
		lockPos := GetRelativePath(lock.Pos(), prog)
		dumpJson.ExistingLocks[lockInitPos] = lockPos
	}
	for unlockInit, lock := range ls.ExistingUnlocks {
		unlockInitPos := GetRelativePath(unlockInit, prog)
		unlockPos := GetRelativePath(lock.Pos(), prog)
		dumpJson.ExistingUnlocks[unlockInitPos] = unlockPos
	}
	return dumpJson
}

// GetRelativePath - {abs_path}{pkg_name}/{some_file} -> {pkg_name}/{some_file}
func GetRelativePath(pos token.Pos, prog *ssa.Program) string {
	if !pos.IsValid() {
		return ""
	}
	position := prog.Fset.Position(pos)
	path := position.String()
	path1 := strings.Split(path, ssaUtils.GlobalPackageName)
	path2 := ssaUtils.GlobalPackageName + path1[1]
	return path2
}

type ContextJSON struct {
	GoroutineID int
	Clock       domain.VectorClock
}

func ContextToJSON(gs *domain.Context) *ContextJSON {
	dumpJson := ContextJSON{}
	dumpJson.GoroutineID = gs.GoroutineID
	dumpJson.Clock = gs.Clock
	return &dumpJson
}
