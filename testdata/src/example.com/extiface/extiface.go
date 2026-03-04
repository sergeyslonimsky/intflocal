package extiface

// MyInterface is an external interface used in tests.
type MyInterface interface {
	Do()
}

// AnotherInterface is another external interface.
type AnotherInterface interface {
	Process(data string) error
}
