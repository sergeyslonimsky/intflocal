package excluded

import "example.com/extiface"

// Service uses external interfaces that will be suppressed by excludePackages or excludeTypes settings.
// No "want" comments — zero diagnostics are expected when the proper settings are applied.
type Service struct {
	dep     extiface.MyInterface
	another extiface.AnotherInterface
}
