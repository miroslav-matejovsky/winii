# winii

A Go library for gathering Windows inventory insights.

`winii` provides deep visibility into your Windows environment by collecting various details about the current system, including network configurations, Windows services, .NET runtimes, Visual C++ Redistributables, system information, and more.

## Features

- **System Insights**: Gather hardware and OS details (CPU, memory, disk, host info).
- **Runtime Insights**: Check installed .NET runtimes and Visual C++ Redistributables.
- **Network Insights**: Retrieve network interface and adapter information.
- **Windows Services**: Enumerate and detail Windows services, including executable paths, versions, and config files.

## Installation

```bash
go get github.com/miroslav-matejovsky/winii
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/miroslav-matejovsky/winii"
)

func main() {
    insights, err := winii.ProvideAllInventoryInsight()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("System: %+v\n", insights.System)
    fmt.Printf("Runtimes: %+v\n", insights.Runtimes)
    fmt.Printf("Network: %+v\n", insights.Network)
    fmt.Printf("Services: %d services found\n", len(insights.WindowsServices))
}
```

### Minimal Insights with Options

```go
options := winii.InsightsOptions{
 WinServiceOptions: winservices.InventoryInsightOption{},
}
return winii.ProvideInsights(options)

```

## API Reference

For detailed API documentation, see the [Go package documentation](https://pkg.go.dev/github.com/miroslav-matejovsky/winii).

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

See [LICENSE](LICENSE) for details.
