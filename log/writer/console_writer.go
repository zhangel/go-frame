package writer

import (
	"fmt"
	"os"
)

type ConsoleWriter struct {
	f *os.File
}

func NewConsoleWriter(stdout bool) (*ConsoleWriter, error) {
	if stdout {
		return &ConsoleWriter{f: os.Stdout}, nil
	} else {
		return &ConsoleWriter{f: os.Stderr}, nil
	}
}

func (s *ConsoleWriter) Write(bytes []byte) error {
	fmt.Fprintln(s.f, string(bytes))
	return nil
}

func (s *ConsoleWriter) Close() error {
	return nil
}
