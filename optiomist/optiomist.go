package optiomist

type OptionStatus int
const OptionNil = -1
const OptionNone = 0
const OptionSome = 1

type Option[T any] struct {
    value T
    status OptionStatus
}

type Optionable interface {
    IsSome() bool
    IsNone() bool
    IsNil() bool
    Value() any
}

func Optiomize[T any](value T, valid bool) Option[T] {
    if valid {
        return Some[T](value)
    } else {
        return None[T]()
    }
}

func Some[T any](value T) Option[T] {
    return Option[T]{value, OptionSome}
}
func None[T any]() Option[T] {
    return Option[T]{ status: OptionNone }
}
func Nil[T any]() Option[T] {
    return Option[T]{ status: OptionNil }
}
func (opt Option[T]) IsSome() bool {
    return opt.status == OptionSome || opt.status == OptionNil
}
func (opt Option[T]) IsNone() bool {
    return opt.status == OptionNone
}
func (opt Option[T]) IsNil() bool {
    return opt.status == OptionNil
}
func (opt Option[T]) Value() any {
    return opt.value
}
func (opt Option[T]) TypedValue() T {
    return opt.value
}
