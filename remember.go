package remember

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Memory struct {
	input []io.Reader
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

func NewMemory(opts ...Option) (Memory, error) {
	_, err := os.Stat("store.txt")
	if err != nil {
		os.Create("store.txt")
	}
	defaultFile, err := os.OpenFile("store.txt", os.O_APPEND|os.O_RDWR, 0644)
	m := Memory{
		[]io.Reader{os.Stdin},
		defaultFile,
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
	data, err := ioutil.ReadAll(m.output)
	if err != nil {
		panic("error opening data store")
	}
	return string(data)
}

func (m *Memory) Memorise() error {
	var memorable []string
	for _, item := range m.input {
		word, err := ioutil.ReadAll(item)
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

func Reminder() string {
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
