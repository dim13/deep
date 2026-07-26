package deep

import (
	"testing"
	"time"
)

func TestCopyInt(t *testing.T) {
	src := 42
	dst := Copy(src)
	if !Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
}

func TestCopyBytes(t *testing.T) {
	src := []byte("test")
	dst := Copy(src)
	if !Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
	src[0] = 'A' // clobber original value
	if Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
}

func TestCopyString(t *testing.T) {
	src := "test"
	dst := Copy(src)
	if dst != src {
		t.Errorf("got %v, want %v", dst, src)
	}
}

type TestStruct struct {
	A, B int
	Any  any
	Ptr  *int
	Map  map[int]string
	Sub  *SubStruct
}

type SubStruct struct {
	String string
}

func TestCopyStruct(t *testing.T) {
	src := TestStruct{
		A:   42,
		Any: "X",
		Ptr: new(42),
		Map: map[int]string{
			0: "zero",
			1: "one",
			2: "two",
		},
		Sub: &SubStruct{String: "YYY"},
	}
	dst := Copy(src)
	if !Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
	src.Any = "clobber"
	src.Map[0] = "clobber"
	src.Sub.String = "clobber"
	if Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
}

func TestCopyPtrStruct(t *testing.T) {
	src := &SubStruct{String: "test"}
	dst := Copy(src)
	if src == dst {
		t.Error("expected different pointers")
	}
	if !Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
}

func TestCopyMap(t *testing.T) {
	src := map[int]string{
		0: "zero",
		1: "one",
		2: "two",
	}
	dst := Copy(src)
	if !Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
	src[0] = "clobber"
	if Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
}

func TestCopyArray(t *testing.T) {
	src := [5]int{1, 2, 3, 4, 5}
	dst := Copy(src)
	if !Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
	src[2] = 30
	if Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
}

func TestCopyArrayOfPointers(t *testing.T) {
	a, b := 1, 2
	src := [2]*int{&a, &b}
	dst := Copy(src)
	if !Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
	*src[0] = 100
	if Equal(dst, src) {
		t.Error("expected dst to be independent of src")
	}
}

func TestCopyFunc(t *testing.T) {
	src := func(x int) int { return x }
	dst := Copy(src)
	if src(5) != dst(5) {
		t.Errorf("got %v, want %v", dst(5), src(5))
	}
}

func TestCopySlice(t *testing.T) {
	src := []string{"test", "abc", "xyz", "uvw", "ijk"}
	dst := Copy(src)
	if !Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
	src[3] = "clobber"
	if Equal(dst, src) {
		t.Errorf("got %v, want %v", dst, src)
	}
}

func TestCopyTime(t *testing.T) {
	src := time.Now()
	dst := Copy(src)
	if !dst.Equal(src) {
		t.Errorf("got %v, want %v", dst, src)
	}
}

type unexportedFields struct {
	Pub  int
	priv int
}

func TestCopyUnexportedFields(t *testing.T) {
	src := unexportedFields{Pub: 1, priv: 2}
	dst := Copy(src)
	if !Equal(dst, src) {
		t.Errorf("got %+v, want %+v", dst, src)
	}
}
