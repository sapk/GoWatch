package admin

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
	"github.com/sapk/GoWatch/modules/tools"
)

var (
	startTime = time.Now()
)

type sysStatus struct {
	Uptime       string
	NumGoroutine int

	// General statistics.
	MemAllocated string // bytes allocated and still in use
	MemTotal     string // bytes allocated (even if freed)
	MemSys       string // bytes obtained from system (sum of XxxSys below)
	Lookups      uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    string // bytes allocated and still in use
	HeapSys      string // bytes obtained from system
	HeapIdle     string // bytes in idle spans
	HeapInuse    string // bytes in non-idle span
	HeapReleased string // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  string // bootstrap stacks
	StackSys    string
	MSpanInuse  string // mspan structures
	MSpanSys    string
	MCacheInuse string // mcache structures
	MCacheSys   string
	BuckHashSys string // profiling bucket hash table
	GCSys       string // GC metadata
	OtherSys    string // other system allocations

	// Garbage collector statistics.
	NextGC       string // next run in HeapAlloc time (bytes)
	LastGC       string // last run in absolute time (ns)
	PauseTotalNs string
	PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC        uint32
}

func getSystemStatus() sysStatus {
	var sys sysStatus
	sys.Uptime = tools.TimeSinceHuman(startTime)

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	sys.NumGoroutine = runtime.NumGoroutine()

	sys.MemAllocated = tools.FileSize(int64(m.Alloc))
	sys.MemTotal = tools.FileSize(int64(m.TotalAlloc))
	sys.MemSys = tools.FileSize(int64(m.Sys))
	sys.Lookups = m.Lookups
	sys.MemMallocs = m.Mallocs
	sys.MemFrees = m.Frees

	sys.HeapAlloc = tools.FileSize(int64(m.HeapAlloc))
	sys.HeapSys = tools.FileSize(int64(m.HeapSys))
	sys.HeapIdle = tools.FileSize(int64(m.HeapIdle))
	sys.HeapInuse = tools.FileSize(int64(m.HeapInuse))
	sys.HeapReleased = tools.FileSize(int64(m.HeapReleased))
	sys.HeapObjects = m.HeapObjects

	sys.StackInuse = tools.FileSize(int64(m.StackInuse))
	sys.StackSys = tools.FileSize(int64(m.StackSys))
	sys.MSpanInuse = tools.FileSize(int64(m.MSpanInuse))
	sys.MSpanSys = tools.FileSize(int64(m.MSpanSys))
	sys.MCacheInuse = tools.FileSize(int64(m.MCacheInuse))
	sys.MCacheSys = tools.FileSize(int64(m.MCacheSys))
	sys.BuckHashSys = tools.FileSize(int64(m.BuckHashSys))
	sys.GCSys = tools.FileSize(int64(m.GCSys))
	sys.OtherSys = tools.FileSize(int64(m.OtherSys))

	sys.NextGC = tools.FileSize(int64(m.NextGC))
	sys.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	sys.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	sys.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	sys.NumGC = m.NumGC

	return sys
}

// Dashboard genrate the home admin page
func Dashboard(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	if err := auth.VerificationAuth(ctx, sess, []string{"admin.dashboard"}); err != nil {
		return
	}
	fillGlobalPage(ctx, db, "admin_dashboard")

	ctx.Data["SysStatus"] = getSystemStatus()
	ctx.HTML(200, "admin/dashboard")
}

func fillGlobalPage(ctx *macaron.Context, db *db.Db, page string) {
	ctx.Data["page_admin"] = true
	if page != "" {
		ctx.Data[page] = true
	}
	ctx.Data["users_count"] = db.NbUsers()
	ctx.Data["equipements_count"] = db.NbEquipements()
}
