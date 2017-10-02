package export

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"sync"
)

var (
	cmdmu sync.RWMutex
	cmds  = make(map[string]func([]string) []byte)
)

// Register a command handler or delete command if fn is nil.
func Register(cmd string, fn func([]string) []byte) {
	cmd = strings.ToUpper(cmd)
	cmdmu.Lock()
	if fn != nil {
		cmds[cmd] = fn
	} else {
		delete(cmds, cmd)
	}
	cmdmu.Unlock()
}

// Serve processes texted command from r, and writes the result to w.
func Serve(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)
	for {
		ln, err := br.ReadSlice('\n')
		if err != nil {
			return err
		}
		ln = bytes.TrimSpace(ln)
		tokens := strings.Split(string(ln), " ")
		cmdmu.RLock()
		fn := cmds[strings.ToUpper(tokens[0])]
		cmdmu.RUnlock()
		var reply []byte
		if fn != nil {
			reply = fn(tokens[1:])
		} else if len(ln) > 0 {
			cmdmu.RLock()
			for cmd := range cmds {
				reply = append(reply, cmd...)
				reply = append(reply, '\n')
			}
			cmdmu.RLock()
		}
		if len(reply) == 0 {
			continue
		}
		p := 0
		for p < len(reply) {
			n, err := w.Write(reply[p:])
			if err != nil {
				return err
			}
			p += n
		}
	}
}
