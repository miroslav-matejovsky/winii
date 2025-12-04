package inventory

import (
	"github.com/miroslav-matejovsky/winii/runtimes"
	"github.com/miroslav-matejovsky/winii/winservices"
)

type Result struct {
	Runtimes        runtimes.Insights
	WindowsServices []winservices.ServiceDetails
}
