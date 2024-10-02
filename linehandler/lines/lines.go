package lines

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidLine   = errors.New("INVALID LINE")
	ErrNoSuchCommand = errors.New("NO SUCH COMMAND")
	ErrInvalidArgs   = errors.New("INVALID ARGUMENTS")
)

type Command struct {
	Args   int
	Handle func(args []string) ([]byte, error)
}

type LineProcessor struct {
	Cmds      map[string]*Command
	Default   *Command
	UpperName bool
	AddReturn bool
}

func (p *LineProcessor) parseCommand(s string) (*Command, []string, error) {
	name, remaining, ok := nextString(s)
	if !ok {
		return nil, nil, ErrInvalidLine
	}
	if p.UpperName {
		name = strings.ToUpper(name)
	}
	leadSpace := false
	cmd, ok := p.Cmds[name]
	if !ok {
		def := p.Default
		if def == nil {
			return nil, nil, ErrNoSuchCommand
		}
		cmd = def
		remaining = s
		leadSpace = true
	} else {
		leadSpace = strings.HasPrefix(remaining, " ")
	}

	args := make([]string, 0, cmd.Args)
	nargs := cmd.Args
	for nargs > 0 && remaining != "" {
		if !leadSpace {
			return nil, nil, ErrInvalidLine
		}
		arg, rem, ok := nextString(remaining)
		if !ok {
			return nil, nil, ErrInvalidLine
		}
		args = append(args, arg)
		remaining = rem
		leadSpace = strings.HasPrefix(rem, " ")
		nargs--
	}
	if len(remaining) > 0 || nargs > 0 {
		return nil, nil, ErrInvalidArgs
	}
	return cmd, args, nil
}

func (p *LineProcessor) ExecuteLine(s string) ([]byte, error) {
	cmd, args, err := p.parseCommand(s)
	if err != nil {
		return nil, err
	}
	return cmd.Handle(args)
}

func (p *LineProcessor) Execute(r io.Reader, w io.Writer) error {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReaderSize(r, 512)
	}
	for {
		lnbuf, err := br.ReadSlice('\n')
		if len(lnbuf) == 0 && err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		lnbuf = bytes.TrimSpace(lnbuf)
		if len(lnbuf) == 0 || lnbuf[0] == '#' {
			continue
		}
		out, err := p.ExecuteLine(string(lnbuf))
		if err != nil {
			if err == ErrInvalidLine || err == ErrInvalidArgs || err == ErrNoSuchCommand {
				_, err := io.WriteString(w, err.Error())
				if err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
				continue
			}
			return err
		}
		_, err = w.Write(out)
		if err != nil {
			return err
		}
		if p.AddReturn {
			_, _ = w.Write([]byte{'\n'})
		}
	}
	return nil
}

func nextString(s string) (string, string, bool) {
	i := 0
	for i < len(s) && unicode.IsSpace(rune(s[i])) {
		i++
	}
	s = s[i:]
	if len(s) == 0 {
		return "", "", false
	}

	if s[0] != '"' {
		space := strings.IndexByte(s, ' ')
		if space >= 0 {
			return s[:space], s[space:], true
		}
		return s, "", true
	}

	i = 1
	for i < len(s) && (s[i] != '"' || s[i-1] == '\\') {
		i++
	}
	if i == len(s) {
		return s, "", false
	}
	i++
	u, err := strconv.Unquote(s[0:i])
	if err != nil {
		return s, "", false
	}
	return u, s[i:], true
}
