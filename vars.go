package export

import (
	"strconv"
	"sync"
	"sync/atomic"
)

var (
	varmu sync.RWMutex
	vars  = make(map[string]interface{})
)

// Var is an abstract type for exported variable.
type Var interface {
	Get() []byte
	Set(v string) []byte
}

// Any export a Var by name.
func Any(name string, v Var) {
	if v == nil {
		return
	}
	varmu.Lock()
	vars[name] = v
	varmu.Unlock()
}

// Bool export a boolean variable by name.
func Bool(name string, v *bool) {
	if v == nil {
		return
	}
	varmu.Lock()
	vars[name] = v
	varmu.Unlock()
}

// Int export an integer variable by name.
func Int(name string, v *int) {
	if v == nil {
		return
	}
	varmu.Lock()
	vars[name] = v
	varmu.Unlock()
}

// Int32 export an 32-bit integer variable by name.
func Int32(name string, v *int32) {
	if v == nil {
		return
	}
	varmu.Lock()
	vars[name] = v
	varmu.Unlock()
}

// Int64 export an 64-bit integer variable by name.
func Int64(name string, v *int64) {
	if v == nil {
		return
	}
	varmu.Lock()
	vars[name] = v
	varmu.Unlock()
}

// String export a string variable by name.
func String(name string, v *string) {
	if v == nil {
		return
	}
	varmu.Lock()
	vars[name] = v
	varmu.Unlock()
}

// ReadOnly export a readonly string variable by name.
func ReadOnly(name string, v string) {
	varmu.Lock()
	vars[name] = v
	varmu.Unlock()
}

// VarFunc type use a function as export.Var.
type VarFunc func(string) []byte

// Get calls function with an empty string.
func (fn VarFunc) Get() []byte { return fn("") }

// Set calls function with SET argument.
func (fn VarFunc) Set(v string) []byte { return fn(v) }

// GetOnlyFunc type use a function as export.Var but no Set.
type GetOnlyFunc func() []byte

// Get calls function.
func (fn GetOnlyFunc) Get() []byte { return fn() }

// Set always return nil.
func (fn GetOnlyFunc) Set(string) []byte { return nil }

func linefeed(reply []byte) []byte {
	if reply != nil {
		reply = append(reply, '\n')
	}
	return reply
}

func onCmdKeys(args []string) []byte {
	if len(args) != 0 {
		return nil
	}
	var reply []byte
	varmu.RLock()
	for k := range vars {
		reply = strconv.AppendQuote(reply, k)
		reply = append(reply, '\n')
	}
	varmu.RUnlock()
	return reply
}

func onCmdGet(args []string) []byte {
	if len(args) != 1 {
		return nil
	}
	varmu.RLock()
	v := vars[args[0]]
	varmu.RUnlock()
	if any, ok := v.(Var); ok {
		return linefeed(any.Get())
	}
	var (
		buf   [32]byte
		reply []byte
	)
	switch x := v.(type) {
	case *bool:
		reply = strconv.AppendBool(buf[:0], *x)
	case *int:
		reply = strconv.AppendInt(buf[:0], int64(*x), 10)
	case *int32:
		reply = strconv.AppendInt(buf[:0], int64(*x), 10)
	case *int64:
		reply = strconv.AppendInt(buf[:0], *x, 10)
	case *string:
		reply = strconv.AppendQuote(buf[:0], *x)
	case string:
		reply = strconv.AppendQuote(buf[:0], x)
	}
	return linefeed(reply)
}

func onCmdSet(args []string) []byte {
	if len(args) != 2 {
		return nil
	}
	varmu.RLock()
	v := vars[args[0]]
	varmu.RUnlock()
	if v == nil {
		return nil
	}
	if any, ok := v.(Var); ok {
		return linefeed(any.Set(args[1]))
	}
	var (
		buf   [32]byte
		reply []byte
	)
	switch x := v.(type) {
	case *bool:
		reply = strconv.AppendBool(buf[:0], *x)
		*x, _ = strconv.ParseBool(args[1])
	case *int:
		reply = strconv.AppendInt(buf[:0], int64(*x), 10)
		t, _ := strconv.ParseInt(args[1], 10, 64)
		*x = int(t)
	case *int32:
		t, _ := strconv.ParseInt(args[1], 10, 64)
		reply = strconv.AppendInt(buf[:0], int64(atomic.SwapInt32(x, int32(t))), 10)
	case *int64:
		t, _ := strconv.ParseInt(args[1], 10, 64)
		reply = strconv.AppendInt(buf[:0], atomic.SwapInt64(x, t), 10)
	case *string:
		reply = strconv.AppendQuote(buf[:0], *x)
		if s := args[1]; len(s) >= 2 && s[0] == s[len(s)-1] && s[0] == '"' {
			*x, _ = strconv.Unquote(args[1])
		} else {
			*x = s
		}
	}
	return linefeed(reply)
}

func init() {
	cmds["KEYS"] = onCmdKeys
	cmds["GET"] = onCmdGet
	cmds["SET"] = onCmdSet
}
