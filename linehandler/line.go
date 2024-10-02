package linehandler

import (
	"bufio"
	"io"
	"strings"

	"github.com/vizee/export"
	"github.com/vizee/export/linehandler/lines"
)

func Serve(r io.Reader, w io.Writer) error {
	lineProcessor := &lines.LineProcessor{
		Cmds: map[string]*lines.Command{
			"GET": {
				Args: 1,
				Handle: func(args []string) ([]byte, error) {
					v := export.Lookup(args[0])
					if v == nil {
						return []byte("NO SUCH KEY"), nil
					}
					return v.Get(), nil
				},
			},
			"SET": {
				Args: 2,
				Handle: func(args []string) ([]byte, error) {
					v := export.Lookup(args[0])
					if v == nil {
						return []byte("NO SUCH KEY"), nil
					}
					err := v.Set(args[1])
					if err != nil {
						return nil, err
					}
					return []byte("OK"), nil
				},
			},
			"LIST": {
				Args: 0,
				Handle: func(_ []string) ([]byte, error) {
					keys := export.AllVars()
					return []byte(strings.Join(keys, "\n")), nil
				},
			},
		},
		UpperName: true,
		AddReturn: true,
	}
	lineProcessor.Default = lineProcessor.Cmds["GET"]
	return lineProcessor.Execute(bufio.NewReaderSize(r, 256), w)
}
