package export

import (
	"strconv"
	"sync"
)

var (
	varmu sync.RWMutex
	vars  = make(map[string]interface{})
)

// Var is an abstract type for exported variable.
type Var interface {
	Load() []byte
	Store(v string) []byte
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

type readOnlyFunc func() []byte

func (fn readOnlyFunc) Load() []byte { return fn() }

func (fn readOnlyFunc) Store(string) []byte { return nil }

type handlerFunc func(string) []byte

func (fn handlerFunc) Load() []byte { return fn("") }

func (fn handlerFunc) Store(v string) []byte { return fn(v) }

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
		return linefeed(any.Load())
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
		return linefeed(any.Store(args[1]))
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
