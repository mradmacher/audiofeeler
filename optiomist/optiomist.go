package optiomist

// Describes status of an option.
type OptionStatus int

// Option has nil value.
const OptionNil = -1

// Option is undefined.
const OptionNone = 0

// Option has some value.
const OptionSome = 1

type Option[T any] struct {
	value  T
	status OptionStatus
}

// Defines behaviour of an Option.
type Optionable interface {
	IsSome() bool
	IsNone() bool
	IsNil() bool
	AnyValue() any
}

// It converts provided value of type T to an Option[T].
// If the second parameter is true, it returns a "some" Option.
// Otherwise it returns "none" Option.
func Optiomize[T any](value T, some bool) Option[T] {
	if some {
		return Some[T](value)
	} else {
		return None[T]()
	}
}

// It compares two options of same comparable type.
// For "some" variant, two options are equal if their values are equal.
// Otherwise thery are equal if both are "none" or both are "nil".
func IsEql[T comparable](opt1, opt2 Option[T]) bool {
	if opt1.status == opt2.status {
		if opt1.status == OptionSome  {
			return opt1.value == opt2.value
		} else {
			return true
		}
	}

	return false
}
// Creates an option with a value of type T.
func Some[T any](value T) Option[T] {
	return Option[T]{value, OptionSome}
}

// Creates an option with an undefined value of type T.
func None[T any]() Option[T] {
	return Option[T]{status: OptionNone}
}

// Creates an option with a nil value of type T.
func Nil[T any]() Option[T] {
	return Option[T]{status: OptionNil}
}

// Does the option have defined value?
func (opt Option[T]) IsSome() bool {
	return opt.status == OptionSome || opt.status == OptionNil
}

// Is the option undefined?
func (opt Option[T]) IsNone() bool {
	return opt.status == OptionNone
}

// Does the option have nil value?
func (opt Option[T]) IsNil() bool {
	return opt.status == OptionNil
}

// Returns the value of an option without specified type.
func (opt Option[T]) AnyValue() any {
	return opt.value
}

// Returns the value of an option with the option's type.
func (opt Option[T]) Value() T {
	return opt.value
}
