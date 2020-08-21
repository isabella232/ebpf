package link

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/DataDog/ebpf"
	"github.com/DataDog/ebpf/asm"
	"github.com/DataDog/ebpf/internal/btf"
)

func TestIter(t *testing.T) {
	prog, err := ebpf.NewProgram(&ebpf.ProgramSpec{
		Type:       ebpf.Tracing,
		AttachType: ebpf.AttachTraceIter,
		AttachTo:   "bpf_map",
		Instructions: asm.Instructions{
			asm.Mov.Imm(asm.R0, 0),
			asm.Return(),
		},
		License: "MIT",
	})
	if errors.Is(err, btf.ErrNotFound) {
		t.Skip("Kernel doesn't support iter:", err)
	}
	if err != nil {
		t.Fatal("Can't load program:", err)
	}
	defer prog.Close()

	it, err := AttachIter(IterOptions{
		Program: prog,
	})
	if err != nil {
		t.Fatal("Can't create iter:", err)
	}

	file, err := it.Open()
	if err != nil {
		t.Fatal("Can't open iter instance:", err)
	}
	defer file.Close()

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(contents) != 0 {
		t.Error("Non-empty output from no-op iterator:", string(contents))
	}

	testLink(t, it, testLinkOptions{
		prog: prog,
		loadPinned: func(s string) (Link, error) {
			return LoadPinnedIter(s)
		},
	})
}
