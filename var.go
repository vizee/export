package export

import (
	"errors"
	"slices"
	"strconv"
	"sync"
)

var (
	ErrGetOnly = errors.New("get only")
)

var (
	varslock sync.RWMutex

	vars = make(map[string]Var)
)

type Var interface {
	Get() []byte
	Set(x string) error
}

func Register(name string, v Var) {
	varslock.Lock()
	vars[name] = v
	varslock.Unlock()
}

func Lookup(name string) Var {
	varslock.RLock()
	defer varslock.RUnlock()
	return vars[name]
}

func AllVars() []string {
	varslock.RLock()
	l := make([]string, 0, len(vars))
	for k := range vars {
		l = append(l, k)
	}
	varslock.RUnlock()
	slices.Sort(l)
	return l
}

type GetOnly func() []byte

func (fn GetOnly) Get() []byte {
	return fn()
}

func (fn GetOnly) Set(string) error {
	return ErrGetOnly
}

type SetOnly func(v string) error

func (fn SetOnly) Get() []byte {
	return []byte("set only")
}

func (fn SetOnly) Set(v string) error {
	return fn(v)
}

type boolVar bool

func (v *boolVar) Get() []byte {
	return strconv.AppendBool(make([]byte, 0, 5), bool(*v))
}

func (v *boolVar) Set(x string) error {
	t, err := strconv.ParseBool(x)
	if err != nil {
		return err
	}
	*v = boolVar(t)
	return nil
}

func Bool(name string, v *bool) {
	Register(name, (*boolVar)(v))
}

type intVar int

func (v *intVar) Get() []byte {
	return strconv.AppendInt(make([]byte, 0, 5), int64(*v), 10)
}

func (v *intVar) Set(x string) error {
	t, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return err
	}
	*v = intVar(t)
	return nil
}

func Int(name string, v *int) {
	Register(name, (*intVar)(v))
}

type int32Var int32

func (v *int32Var) Get() []byte {
	return strconv.AppendInt(make([]byte, 0, 5), int64(*v), 10)
}

func (v *int32Var) Set(x string) error {
	t, err := strconv.ParseInt(x, 10, 32)
	if err != nil {
		return err
	}
	*v = int32Var(t)
	return nil
}

func Int32(name string, v *int32) {
	Register(name, (*int32Var)(v))
}

type int64Var int64

func (v *int64Var) Get() []byte {
	return strconv.AppendInt(make([]byte, 0, 5), int64(*v), 10)
}

func (v *int64Var) Set(x string) error {
	t, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return err
	}
	*v = int64Var(t)
	return nil
}

func Int64(name string, v *int64) {
	Register(name, (*int64Var)(v))
}

type byteVar byte

func (v *byteVar) Get() []byte {
	return strconv.AppendUint(make([]byte, 0, 3), uint64(*v), 10)
}

func (v *byteVar) Set(x string) error {
	t, err := strconv.ParseUint(x, 10, 8)
	if err != nil {
		return err
	}
	*v = byteVar(t)
	return nil
}

func Byte(name string, v *byte) {
	Register(name, (*byteVar)(v))
}

type uint64Var uint64

func (v *uint64Var) Get() []byte {
	return strconv.AppendUint(make([]byte, 0, 5), uint64(*v), 10)
}

func (v *uint64Var) Set(x string) error {
	t, err := strconv.ParseUint(x, 10, 64)
	if err != nil {
		return err
	}
	*v = uint64Var(t)
	return nil
}

func Uint64(name string, v *uint64) {
	Register(name, (*uint64Var)(v))
}

type stringVar string

func (v *stringVar) Get() []byte {
	return []byte(*v)
}

func (v *stringVar) Set(x string) error {
	*v = stringVar(x)
	return nil
}

func String(name string, v *string) {
	Register(name, (*stringVar)(v))
}

type bytesVar []byte

func (v *bytesVar) Get() []byte {
	return *v
}

func (v *bytesVar) Set(x string) error {
	*v = bytesVar(x)
	return nil
}

func Bytes(name string, v *[]byte) {
	Register(name, (*bytesVar)(v))
}

type readOnly []byte

func (ro readOnly) Get() []byte {
	return ro
}

func (ro readOnly) Set(v string) error {
	return nil
}

func ReadOnly(name string, v []byte) {
	varslock.Lock()
	vars[name] = readOnly(v)
	varslock.Unlock()
}
