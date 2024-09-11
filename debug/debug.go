package debug

import (
	"bytes"
	"encoding/json"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"time"

	"github.com/vizee/export"
)

type debugFunc func(x string) []byte

func (f debugFunc) Get() []byte {
	return f("")
}

func (f debugFunc) Set(x string) error {
	f(x)
	return nil
}

func readOSArgs() []byte {
	data, _ := json.Marshal(os.Args)
	return data
}

func readNumGoroutine() []byte {
	return strconv.AppendInt(nil, int64(runtime.NumGoroutine()), 10)
}

func readRuntimeStack() []byte {
	const (
		initBuffer = 32 * 1024
		maxBuffer  = 4 * 1024 * 1024
	)
	n := initBuffer
	var buf []byte
	for n <= maxBuffer {
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
	data, err := json.Marshal(&mem)
	if err != nil {
		panic(err)
	}
	return data
}

func debugMemoryLimit(x string) []byte {
	n := int64(-1)
	if x != "" {
		var err error
		n, err = strconv.ParseInt(x, 10, 64)
		if err != nil {
			return []byte(err.Error())
		}
	}
	ret := debug.SetMemoryLimit(n)
	return strconv.AppendInt(nil, int64(ret), 10)
}

func debugGOMAXPROCS(x string) []byte {
	n := 0
	if x != "" {
		var err error
		n, err = strconv.Atoi(x)
		if err != nil {
			return []byte(err.Error())
		}
	}
	ret := runtime.GOMAXPROCS(n)
	return strconv.AppendInt(nil, int64(ret), 10)
}

func setMaxThreads(x string) error {
	n, err := strconv.Atoi(x)
	if err != nil {
		return err
	}
	debug.SetMaxThreads(n)
	return nil
}

func readGCStats(gc bool) []byte {
	if gc {
		runtime.GC()
	}
	var st debug.GCStats
	debug.ReadGCStats(&st)
	data, err := json.Marshal(&st)
	if err != nil {
		panic(err)
	}
	return data
}

func setGCPercent(x string) error {
	n, err := strconv.Atoi(x)
	if err != nil {
		return err
	}
	debug.SetGCPercent(n)
	return nil
}

func readProfiles() []byte {
	profiles := pprof.Profiles()
	list := make([]byte, 0, len(profiles))
	for i, p := range profiles {
		if i > 0 {
			list = append(list, '\n')
		}
		list = append(list, p.Name()...)
	}
	return list
}

func fetchProfile(name string, debug int) []byte {
	p := pprof.Lookup(name)
	if p == nil {
		return nil
	}
	var buf bytes.Buffer
	p.WriteTo(&buf, debug)
	return buf.Bytes()
}

func runTracing(dur time.Duration) []byte {
	var buf bytes.Buffer
	err := trace.Start(&buf)
	if err != nil {
		return append([]byte(err.Error()), '\n')
	}
	time.Sleep(dur)
	trace.Stop()
	return buf.Bytes()
}

func runCPUProfile(dur time.Duration) []byte {
	var buf bytes.Buffer
	err := pprof.StartCPUProfile(&buf)
	if err != nil {
		return append([]byte(err.Error()), '\n')
	}
	time.Sleep(dur)
	pprof.StopCPUProfile()
	return buf.Bytes()
}

func registerProfiles() {
	var debug int
	export.Register("debug.pprof.debug", debugFunc(func(v string) []byte {
		if v == "" {
			return []byte(strconv.Itoa(debug))
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			return []byte(err.Error())
		}
		debug = n
		return nil
	}))

	profiles := pprof.Profiles()
	for _, profile := range profiles {
		name := profile.Name()
		export.Register("debug.pprof.profile."+name, export.GetOnly(func() []byte {
			return fetchProfile(name, debug)
		}))
	}
}

func registerTrace() {
	var runDur time.Duration
	export.Register("debug.trace.duration", debugFunc(func(v string) []byte {
		if v == "" {
			return []byte(runDur.String())
		}
		d, err := time.ParseDuration(v)
		if err != nil {
			return []byte(err.Error())
		}
		runDur = d
		return nil
	}))
	export.Register("debug.trace.cpuprofile", export.GetOnly(func() []byte {
		return runCPUProfile(runDur)
	}))
	export.Register("debug.trace.tracing", export.GetOnly(func() []byte {
		return runTracing(runDur)
	}))
}

func init() {
	export.Register("debug.os.args", export.GetOnly(readOSArgs))
	export.Register("debug.runtime.goroutines", export.GetOnly(readNumGoroutine))
	export.Register("debug.runtime.stack", export.GetOnly(readRuntimeStack))
	export.Register("debug.runtime.memstats", export.GetOnly(readMemStats))
	export.Register("debug.runtime.memlimit", debugFunc(debugMemoryLimit))
	export.Register("debug.runtime.gomaxprocs", debugFunc(debugGOMAXPROCS))
	export.Register("debug.runtime.maxthreads", export.SetOnly(setMaxThreads))
	export.Register("debug.gc.stats", export.GetOnly(func() []byte {
		return readGCStats(false)
	}))
	export.Register("debug.gc.stats1", export.GetOnly(func() []byte {
		return readGCStats(true)
	}))
	export.Register("debug.gc.percent", export.SetOnly(setGCPercent))
	export.Register("debug.pprof.profiles", export.GetOnly(readProfiles))
	registerProfiles()
	registerTrace()
}
