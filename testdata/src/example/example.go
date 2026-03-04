package example

import "example/extiface"

// localIface is a locally defined private interface — this is the expected pattern.
type localIface interface {
	Do()
}

// concreteType is a regular struct, not an interface.
type concreteType struct {
	value string
}

// Good demonstrates correct usage — local private interface and concrete types.
type Good struct {
	dep  localIface   // OK: local private interface
	name string       // OK: not an interface
	data concreteType // OK: concrete type
	err  error        // OK: builtin interface
}

// Bad demonstrates incorrect usage — imported external interfaces.
type Bad struct {
	dep     extiface.MyInterface      // want `struct field "dep" uses external interface "extiface.MyInterface"; define it locally as a private interface`
	another extiface.AnotherInterface // want `struct field "another" uses external interface "extiface.AnotherInterface"; define it locally as a private interface`
}

// EmbeddedBad demonstrates embedded external interface.
type EmbeddedBad struct {
	extiface.MyInterface // want `struct field "\(embedded\)" uses external interface "extiface.MyInterface"; define it locally as a private interface`
}

// PointerBad demonstrates pointer to external interface.
type PointerBad struct {
	dep *extiface.MyInterface // want `struct field "dep" uses external interface "extiface.MyInterface"; define it locally as a private interface`
}
