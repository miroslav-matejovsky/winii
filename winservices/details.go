package winservices

import (
	"fmt"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	fi "github.com/miroslav-matejovsky/winii/fileinfo"
)

type ServiceDetails struct {
	ServiceName string // Name of the service, always populated

	CollectionError error // Error encountered during collection, if any

	DisplayName string // DisplayName is the human-readable name of the service.

	Description string // Description provides a brief description of the service.

	PathToExecutable string // PathToExecutable is the full path to the service's executable.

	StartupType string // StartupType indicates how the service is started (e.g., "auto", "manual").

	ServiceStatus string // ServiceStatus shows the current status of the service (e.g., "running", "stopped").

	ServiceType string // ServiceType describes the type of service (e.g., "win32_own_process").

	ErrorControl string // ErrorControl specifies the action taken on failure (e.g., "normal", "severe").

	Dependencies []string // Dependencies lists the names of services this service depends on.

	ServiceStartName string // ServiceStartName is the account name under which the service runs.

	DelayedAutoStart bool // DelayedAutoStart indicates if the service starts with a delay on boot.

	Executable ServiceExecutable // Executable contains details about the service's executable and config files.

	Recovery ServiceRecovery // Recovery defines the recovery actions for service failures.
}

type ServiceExecutable struct {
	ExecutableFile ExecutableFile      // ExecutableFile holds information about the main executable file.
	ConfigFiles    []ServiceConfigFile // ConfigFiles lists configuration files associated with the service.
}

type ExecutableFile struct {
	Path           string    // Path is the full path to the executable file.
	Version        string    // Version is the file version of the executable.
	ProductVersion string    // ProductVersion is the product version of the executable.
	CreationTime   time.Time // CreationTime is when the file was created.
	LastAccessTime time.Time // LastAccessTime is when the file was last accessed.
	LastWriteTime  time.Time // LastWriteTime is when the file was last modified.
}

type ServiceRecovery struct {
	Command                 string        // Command is the command to run on failure.
	FirstFailure            string        // FirstFailure is the action for the first failure.
	FirstFailureAfter       time.Duration // FirstFailureAfter is the delay before the first failure action.
	SecondFailure           string        // SecondFailure is the action for the second failure.
	SecondFailureAfter      time.Duration // SecondFailureAfter is the delay before the second failure action.
	SubsequentFailures      string        // SubsequentFailures is the action for subsequent failures.
	SubsequentFailuresAfter time.Duration // SubsequentFailuresAfter is the delay before subsequent failure actions.
	MoreThan3Actions        bool          // MoreThan3Actions indicates if more than 3 actions are configured.
	ResetFailCountAfter     time.Duration // ResetFailCountAfter is the time after which the failure count resets.
}

// GetServiceDetails retrieves detailed information about a Windows service by its name.
func (s *WinSvcManager) GetServiceDetails(name string) (*ServiceDetails, error) {
	if err := s.Connect(); err != nil {
		return nil, err
	}
	// Check if the service exists
	exists, err := s.serviceExists(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrServiceNotFound
	}
	details := &ServiceDetails{
		ServiceName: name,
	}
	service, err := s.mgr.OpenService(name)
	if err != nil {
		details.CollectionError = err
		return details, nil
	}
	defer func() { _ = service.Close() }()

	// Get service configuration
	config, err := service.Config()
	if err != nil {
		details.CollectionError = err
	} else {
		details.DisplayName = config.DisplayName
		details.Description = config.Description
		details.PathToExecutable = config.BinaryPathName
		details.StartupType = startTypeToString(config.StartType)
		details.ServiceType = serviceTypeToString(config.ServiceType)
		details.ErrorControl = errorControlToString(config.ErrorControl)
		details.Dependencies = config.Dependencies
		details.ServiceStartName = config.ServiceStartName
		details.DelayedAutoStart = config.DelayedAutoStart
	}

	// Get current status
	status, err := service.Query()
	if err != nil {
		if details.CollectionError == nil {
			details.CollectionError = err
		}
	} else {
		details.ServiceStatus = stateToString(status.State)
	}

	// Get service recovery options
	recoveryActions, err := service.RecoveryActions()
	if err != nil {
		if details.CollectionError == nil {
			details.CollectionError = err
		}
	} else {
		recoveryDetails := make([]string, 3)
		recoveryDelays := make([]time.Duration, 3)
		moreThan3RecoveryActions := false
		for i, action := range recoveryActions {
			if i >= 3 {
				moreThan3RecoveryActions = true
				break
			}
			recoveryDetails[i] = recoverActionToString(action.Type)
			recoveryDelays[i] = action.Delay
		}

		resetSeconds, err := service.ResetPeriod()
		if err != nil {
			if details.CollectionError == nil {
				details.CollectionError = err
			}
		} else {
			details.Recovery.ResetFailCountAfter = time.Duration(resetSeconds) * time.Second
		}

		recoveryCommand, err := service.RecoveryCommand()
		if err != nil {
			if details.CollectionError == nil {
				details.CollectionError = err
			}
		} else {
			details.Recovery.Command = recoveryCommand
		}

		details.Recovery.FirstFailure = recoveryDetails[0]
		details.Recovery.FirstFailureAfter = recoveryDelays[0]
		details.Recovery.SecondFailure = recoveryDetails[1]
		details.Recovery.SecondFailureAfter = recoveryDelays[1]
		details.Recovery.SubsequentFailures = recoveryDetails[2]
		details.Recovery.SubsequentFailuresAfter = recoveryDelays[2]
		details.Recovery.MoreThan3Actions = moreThan3RecoveryActions
	}

	// Collect executable information if PathToExecutable is available
	if details.PathToExecutable != "" {
		executable, err := findServiceExecutable(details.PathToExecutable)
		if err != nil {
			if details.CollectionError == nil {
				details.CollectionError = err
			}
		} else {
			details.Executable.ExecutableFile.Path = executable
			wf, err := fi.NewWinFileInfo(executable)
			if err != nil {
				if details.CollectionError == nil {
					details.CollectionError = err
				}
			} else {
				versions, err := wf.GetVersions()
				if err != nil {
					if details.CollectionError == nil {
						details.CollectionError = err
					}
				} else {
					details.Executable.ExecutableFile.Version = versions.FileVersion.String()
					details.Executable.ExecutableFile.ProductVersion = versions.ProductVersion.String()
				}
				fileTime, err := wf.GetFileTimestamps()
				if err != nil {
					if details.CollectionError == nil {
						details.CollectionError = err
					}
				} else {
					details.Executable.ExecutableFile.CreationTime = fileTime.CreationTime
					details.Executable.ExecutableFile.LastAccessTime = fileTime.LastAccessTime
					details.Executable.ExecutableFile.LastWriteTime = fileTime.LastWriteTime
				}
			}
			executableDir := filepath.Dir(executable)
			configFiles, err := collectServiceConfigFiles(executableDir)
			if err != nil {
				if details.CollectionError == nil {
					details.CollectionError = err
				}
			} else {
				details.Executable.ConfigFiles = configFiles
			}
		}
	}

	return details, nil
}

func startTypeToString(startType uint32) string {
	switch startType {
	case mgr.StartAutomatic:
		return "Automatic"
	case mgr.StartDisabled:
		return "Disabled"
	case mgr.StartManual:
		return "Manual"
	default:
		return fmt.Sprintf("Unknown (%d)", startType)
	}
}

func stateToString(state svc.State) string {
	switch state {
	case svc.Stopped:
		return "Stopped"
	case svc.StartPending:
		return "Start Pending"
	case svc.StopPending:
		return "Stop Pending"
	case svc.Running:
		return "Running"
	case svc.ContinuePending:
		return "Continue Pending"
	case svc.PausePending:
		return "Pause Pending"
	case svc.Paused:
		return "Paused"
	default:
		return fmt.Sprintf("Unknown (%d)", state)
	}
}

func serviceTypeToString(serviceType uint32) string {
	switch serviceType {
	case windows.SERVICE_KERNEL_DRIVER: // 1
		return "Kernel Driver"
	case windows.SERVICE_FILE_SYSTEM_DRIVER: // 2
		return "File System Driver"
	case windows.SERVICE_ADAPTER: // 4
		return "Adapter"
	case windows.SERVICE_RECOGNIZER_DRIVER: // 8
		return "Recognizer Driver"
	case windows.SERVICE_WIN32_OWN_PROCESS: // 16
		return "Win32 Own Process"
	case windows.SERVICE_WIN32_SHARE_PROCESS: // 32
		return "Win32 Share Process"
	case windows.SERVICE_WIN32: // 16 | 32
		return "Win32"
	case windows.SERVICE_INTERACTIVE_PROCESS: // 256
		return "Interactive Process"
	case windows.SERVICE_DRIVER: // 1 | 2 | 8
		return "Driver"
	case windows.SERVICE_TYPE_ALL: // 16 | 32 | 4 | 1 | 2 | 8 | 256
		return "All"
	default:
		return fmt.Sprintf("Unknown (%d)", serviceType)
	}
}

func errorControlToString(errorControl uint32) string {
	switch errorControl {
	case mgr.ErrorNormal:
		return "Normal"
	case mgr.ErrorSevere:
		return "Severe"
	case mgr.ErrorCritical:
		return "Critical"
	case mgr.ErrorIgnore:
		return "Ignore"
	default:
		return fmt.Sprintf("Unknown (%d)", errorControl)
	}
}

func recoverActionToString(action int) string {
	switch action {
	case windows.SC_ACTION_NONE:
		return "None"
	case windows.SC_ACTION_RESTART:
		return "Restart"
	case windows.SC_ACTION_REBOOT:
		return "Reboot"
	case windows.SC_ACTION_RUN_COMMAND:
		return "Run Command"
	default:
		return fmt.Sprintf("Unknown (%d)", action)
	}
}
