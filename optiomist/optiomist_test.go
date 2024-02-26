package optiomist

import (
	"fmt"
	"testing"
)

func Example() {
	type Test struct {
		V1 Option[string]
		V2 Option[string]
		V3 Option[string]
	}

	test := Test{
		Some("I'm a string"),
		None[string](),
		Nil[string](),
	}

	fmt.Printf("Some[some: %v, none: %v, nil: %v]: %q\n",
		test.V1.IsSome(),
		test.V1.IsNone(),
		test.V1.IsNil(),
		test.V1.Value())

	fmt.Printf("None[some: %v, none: %v, nil: %v]: %q\n",
		test.V2.IsSome(),
		test.V2.IsNone(),
		test.V2.IsNil(),
		test.V2.Value())

	fmt.Printf("Nil[some: %v, none: %v, nil: %v]: %q\n",
		test.V3.IsSome(),
		test.V3.IsNone(),
		test.V3.IsNil(),
		test.V3.Value())

	// Output:
	// Some[some: true, none: false, nil: false]: "I'm a string"
	// None[some: false, none: true, nil: false]: ""
	// Nil[some: true, none: false, nil: true]: ""
}

func TestSome_string(t *testing.T) {
	value := "I'm here"
	opt := Some(value)
	if !opt.IsSome() {
		t.Error("Some.IsSome() = false; expected true")
	}

	if opt.IsNone() {
		t.Error("Some.IsNone() = true; expected false")
	}

	if opt.IsNil() {
		t.Error("Some.IsNil() = true; expected false")
	}

	if opt.Value() != value {
		t.Errorf("Some.Value() = %q; expected %q", opt.Value(), value)
	}

	if opt.TypedValue() != value {
		t.Errorf("Some.TypedValue() = %q; expected %q", opt.TypedValue(), value)
	}
}

func TestNone_string(t *testing.T) {
	opt := None[string]()
	if opt.IsSome() {
		t.Error("Some.IsSome() = true; expected false")
	}

	if !opt.IsNone() {
		t.Error("Some.IsNone() = false; expected true")
	}

	if opt.IsNil() {
		t.Error("Some.IsNil() = true; expected false")
	}
}

func TestNil_string(t *testing.T) {
	opt := Nil[string]()
	if !opt.IsSome() {
		t.Error("Some.IsSome() = false; expected true")
	}

	if opt.IsNone() {
		t.Error("Some.IsNone() = true; expected false")
	}

	if !opt.IsNil() {
		t.Error("Some.IsNil() = false; expected true")
	}

	var want string // nil string
	if opt.Value() != want {
		t.Errorf("Some.Value() = %q; expected nil", opt.Value())
	}
}
