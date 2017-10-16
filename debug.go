package export

import (
	"bytes"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"time"
)

func readOSArgs() []byte {
	var reply []byte
	for _, arg := range os.Args {
		if reply != nil {
			reply = append(reply, ' ')
		}
		hasspace := false
		for i := 0; i < len(arg); i++ {
			if arg[i] == ' ' {
				hasspace = true
				break
			}
		}
		if hasspace {
			reply = append(reply, '"')
			reply = append(reply, arg...)
			reply = append(reply, '"')
		} else {
			reply = append(reply, arg...)
		}
	}
	return reply
}

func readNumGoroutine() []byte {
	return strconv.AppendInt(nil, int64(runtime.NumGoroutine()), 10)
}

func readRuntimeStack() []byte {
	n := 1 << 14
	var buf []byte
	for n <= 64<<20 {
		buf = make([]byte, n)
		n = runtime.Stack(buf, true)
		if n < len(buf) {
			break
		}
		n += n
	}
	return buf
}

func readMemStats() []byte {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return structToJSON(make([]byte, 0, 128), nil, reflect.ValueOf(&mem).Elem())
}

func handleGCStats(v string) []byte {
	if v == "gc" {
		runtime.GC()
	}
	var st debug.GCStats
	debug.ReadGCStats(&st)
	return structToJSON(make([]byte, 0, 64), nil, reflect.ValueOf(&st).Elem())
}

func handleGOMAXPROCS(v string) []byte {
	n := 0
	if v != "" {
		n, _ = strconv.Atoi(v)
	}
	return strconv.AppendInt(nil, int64(runtime.GOMAXPROCS(n)), 10)
}

func handleSetMaxThreads(v string) []byte {
	n := 0
	if v != "" {
		n, _ = strconv.Atoi(v)
	}
	return strconv.AppendInt(nil, int64(debug.SetMaxThreads(n)), 10)
}

func handleSetGCPercent(v string) []byte {
	n := 0
	if v != "" {
		n, _ = strconv.Atoi(v)
	}
	return strconv.AppendInt(nil, int64(debug.SetGCPercent(n)), 10)
}

func listProfiles() []byte {
	var reply []byte
	for _, p := range pprof.Profiles() {
		reply = append(reply, p.Name()...)
		reply = append(reply, ' ')
		reply = strconv.AppendInt(reply, int64(p.Count()), 10)
		reply = append(reply, '\n')
	}
	return reply
}

func fetchProfile(args []string) []byte {
	if len(args) == 0 {
		return nil
	}
	debug := 1
	if len(args) >= 2 {
		debug, _ = strconv.Atoi(args[1])
	}
	p := pprof.Lookup(args[0])
	if p == nil {
		return nil
	}
	var buf bytes.Buffer
	p.WriteTo(&buf, debug)
	return buf.Bytes()
}

func onCmdPprof(args []string) []byte {
	if len(args) == 0 {
		return listProfiles()
	}
	return fetchProfile(args)
}

func onCmdTrace(args []string) []byte {
	if len(args) == 0 || args[0] != "start" {
		return []byte("start [seconds:1]\n")
	}
	var sec int64
	if len(args) == 2 {
		sec, _ = strconv.ParseInt(args[1], 10, 64)
	}
	if sec <= 0 {
		sec = 1
	}

	var buf bytes.Buffer
	err := trace.Start(&buf)
	if err != nil {
		return append([]byte(err.Error()), '\n')
	}
	time.Sleep(time.Duration(sec) * time.Second)
	trace.Stop()
	return buf.Bytes()
}

func onCmdProfile(args []string) []byte {
	if len(args) == 0 || args[0] != "start" {
		return []byte("start [seconds:30]\n")
	}
	var sec int64
	if len(args) == 2 {
		sec, _ = strconv.ParseInt(args[1], 10, 64)
	}
	if sec <= 0 {
		sec = 30
	}
	var buf bytes.Buffer
	err := pprof.StartCPUProfile(&buf)
	if err != nil {
		return append([]byte(err.Error()), '\n')
	}
	time.Sleep(time.Duration(sec) * time.Second)
	pprof.StopCPUProfile()
	return buf.Bytes()
}

func init() {
	vars["os.Args"] = GetOnlyFunc(readOSArgs)
	vars["runtime.GOOS"] = runtime.GOOS
	vars["runtime.GOARCH"] = runtime.GOARCH
	vars["runtime.Version"] = runtime.Version()
	vars["runtime.NumGoroutine"] = GetOnlyFunc(readNumGoroutine)
	vars["runtime.Stack"] = GetOnlyFunc(readRuntimeStack)
	vars["runtime.MemStats"] = GetOnlyFunc(readMemStats)
	vars["runtime.GOMAXPROCS"] = VarFunc(handleGOMAXPROCS)
	vars["debug.GCStats"] = VarFunc(handleGCStats)
	vars["debug.SetMaxThreads"] = VarFunc(handleSetMaxThreads)
	vars["debug.SetGCPercent"] = VarFunc(handleSetGCPercent)

	cmds["PPROF"] = onCmdPprof
	cmds["TRACE"] = onCmdTrace
	cmds["PROFILE"] = onCmdProfile
}
