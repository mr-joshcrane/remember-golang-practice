package remember_test

import (
	"remember"
	"testing"
)

func TestMemorise(t *testing.T) {
	t.Parallel()
	memory, err := remember.NewMemory(
		remember.WithInput([]string{"a", "thing", "to", "remember"}),
		remember.WithOutput(""),
	)
	if err != nil {
		t.Fatalf("could not instantiate memory object %q", err)
	}
	err = memory.Memorise()
	if err != nil {
		t.Fatal(err)
	}
	got := memory.Recall()
	want := "a thing to remember\n"
	if want != got {
		t.Fatalf("wanted %s but got %s", want, got)
	}
}

func TestMemoriseSuccessive(t *testing.T) {
	t.Parallel()
	memory, err := remember.NewMemory(
		remember.WithInput([]string{"another", "thing", "to", "remember!"}),
		remember.WithOutput("a thing to remember\n"),
	)
	if err != nil {
		t.Fatalf("could not instantiate memory object %q", err)
	}

	err = memory.Memorise()
	if err != nil {
		t.Fatal(err)
	}

	got := memory.Recall()
	want := "a thing to remember\nanother thing to remember!\n"
	if want != got {
		t.Fatalf("wanted %s but got %s", want, got)
	}
}
