package engine

import (
	"reflect"
)

type Signal[T any] []func(T)

func (s *Signal[T]) Listen(f func(T)) {
	*s = append(*s, f)
}

func (s *Signal[T]) Forward(fwd *Signal[T]) {
	s.Listen(fwd.Emit)
}

func (s *Signal[T]) Emit(arg T) {
	for _, f := range *s {
		f(arg)
	}
}

// Listens to signalName on controlPtr; if available
func Listen[T any](controlPtr any, signalName string, f func(T)) {
	signal := getSignal[T](controlPtr, signalName)
	if signal != nil {
		signal.Listen(f)
	}
}

func ForwardSignal[T any](controlPtr any, signalName string, fwd *Signal[T]) {
	signal := getSignal[T](controlPtr, signalName)
	if signal != nil {
		signal.Listen(fwd.Emit)
	}
}

// uses reflection to find the signal by name
func getSignal[T any](controlPtr any, signalName string) *Signal[T] {
	field := reflect.ValueOf(controlPtr).Elem().FieldByName(signalName)
	if !field.IsValid() || !field.CanAddr() {
		return nil
	}
	return field.Addr().Interface().(*Signal[T])
}
