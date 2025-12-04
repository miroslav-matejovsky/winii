// package winservices provides helpers for reading information about Windows services.
//
// The library exposes a small, focused API for:
//   - enumerating services
//   - checking whether a service exists
//   - reading service state and human-friendly state/type names
//   - collecting detailed service information (executable path, file version, config files,
//     recovery options, and more)
//
// The package is a thin helper around golang.org/x/sys/windows/svc and
// golang.org/x/sys/windows/svc/mgr and is intended to be used on Windows only.
//
// Example
//
//	mgr := winservices.NewWinSvcManager()
//	defer func() { _ = mgr.Disconnect() }()
//
//	exists, err := mgr.ServiceExists("MyService")
//	if err != nil {
//	    // handle error (service manager not available, permission denied, ...)
//	}
//	if exists {
//	    state, stateStr, _ := mgr.GetServiceState("MyService")
//	    fmt.Printf("%s is %s (%d)\n", "MyService", stateStr, state)
//	}
//
// For detailed information (including file metadata and possible config files next to
// the service executable) call `GetServiceDetails(name, includeFiles bool)`.
package winservices
