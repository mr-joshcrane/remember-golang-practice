package remember

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type Memory struct {
	input  []io.Reader
	output io.ReadWriter
}

type Option func(*Memory) error

func WithInput(args []string) Option {
	return func(m *Memory) error {
		var input = []io.Reader{}
		for _, str := range args {
			input = append(input, bytes.NewBufferString(str))
		}
		m.input = input
		return nil
	}
}

func WithOutput(output string) Option {
	return func(m *Memory) error {
		m.output = bytes.NewBufferString(output)
		return nil
	}
}

func DefaultFileStore() (*os.File, error) {
	_, err := os.Stat("store.txt")
	if err != nil {
		if err != os.ErrNotExist {
			return nil, err
		}
		_, err := os.Create("store.txt")
		if err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile("store.txt", os.O_APPEND|os.O_RDWR, 0644)
	return f, nil
}

func NewMemory(opts ...Option) (Memory, error) {
	file, err := DefaultFileStore()
	if err != nil {
		return Memory{}, err
	}
	m := Memory{
		[]io.Reader{os.Stdin},
		file,
	}
	for _, opt := range opts {
		err := opt(&m)
		if err != nil {
			return Memory{}, err
		}
	}
	return m, nil
}

func (m *Memory) Recall() string {
	data, err := io.ReadAll(m.output)
	if err != nil {
		panic(err.Error())
	}
	return string(data)
}

func (m *Memory) Memorise() error {
	var memorable []string
	for _, item := range m.input {
		word, err := io.ReadAll(item)
		if err != nil {
			return err
		}
		memorable = append(memorable, string(word))
	}
	data := strings.Join(memorable, " ")
	data = data + "\n"
	_, err := m.output.Write([]byte(data))
	return err
}

func RunReminder() string {
	m, err := NewMemory(
		WithInput(os.Args[1:]),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(m.input) == 0 {
		return m.Recall()
	}

	err = m.Memorise()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return "OK I've remembered that"
}
