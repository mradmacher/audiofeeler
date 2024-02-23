package audiofeeler

type OptionStatus int
const OptionNil = -1
const OptionNone = 0
const OptionSome = 1

type Optionable interface {
    string | uint | int
}

type Option[T Optionable] struct {
    value T
    status OptionStatus
}

func Some[T Optionable](value T) Option[T] {
    return Option[T]{value, OptionSome}
}
func None[T Optionable]() Option[T] {
    return Option[T]{ status: OptionNone }
}
func Nil[T Optionable]() Option[T] {
    return Option[T]{ status: OptionNil }
}
func (opt *Option[T]) IsSome() bool {
    return opt.status == OptionSome
}
func (opt *Option[T]) IsNone() bool {
    return opt.status == OptionNone
}
func (opt *Option[T]) IsNil() bool {
    return opt.status == OptionNil
}
func (opt *Option[T]) Value() T {
    return opt.value
}
