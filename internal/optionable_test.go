package audiofeeler

import "testing"

func TestSomeString(t *testing.T) {
    var opt Option[string]
    value := "I'm here"
    opt = Some(value)
    if !opt.IsSome() {
        t.Errorf("Some.IsSome() = false; expected true")
    }

    if opt.IsNone() {
        t.Errorf("Some.IsNone() = true; expected false")
    }

    if opt.IsNil() {
        t.Errorf("Some.IsNil() = true; expected false")
    }

    if opt.Value() != value {
        t.Errorf("Some.Value() = %q; expected %q", opt.Value(), value)
    }
}

func TestNoneString(t *testing.T) {
    var opt Option[string]
    opt = None[string]()
    if opt.IsSome() {
        t.Errorf("Some.IsSome() = true; expected false")
    }

    if !opt.IsNone() {
        t.Errorf("Some.IsNone() = false; expected true")
    }

    if opt.IsNil() {
        t.Errorf("Some.IsNil() = true; expected false")
    }

    if opt.Value() != "" {
        t.Errorf("Some.Value() = %q; expected %q", opt.Value(), "")
    }
}

func TestNilString(t *testing.T) {
    var opt Option[string]
    opt = Nil[string]()
    if opt.IsSome() {
        t.Errorf("Some.IsSome() = true; expected false")
    }

    if opt.IsNone() {
        t.Errorf("Some.IsNone() = true; expected false")
    }

    if !opt.IsNil() {
        t.Errorf("Some.IsNil() = false; expected true")
    }

    if opt.Value() != "" {
        t.Errorf("Some.Value() = %q; expected %q", opt.Value(), "")
    }
}
