package winservices

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	fi "github.com/miroslav-matejovsky/winii/fileinfo"
)

type ServiceDetails struct {
	Name             string
	DisplayName      string
	Description      string
	PathToExecutable string
	StartupType      string
	ServiceStatus    string
	ServiceType      string
	ErrorControl     string
	Dependencies     []string
	ServiceStartName string
	DelayedAutoStart bool
	Executable       ServiceExecutable
	Recovery         ServiceRecovery
}

type ServiceExecutable struct {
	ExecutableFile ExecutableFile
	ConfigFiles    []ServiceConfigFile
}

type ExecutableFile struct {
	Path           string
	Version        string
	ProductVersion string
	CreationTime   time.Time
	LastAccessTime time.Time
	LastWriteTime  time.Time
}

type ServiceRecovery struct {
	Command                 string
	FirstFailure            string
	FirstFailureAfter       time.Duration
	SecondFailure           string
	SecondFailureAfter      time.Duration
	SubsequentFailures      string
	SubsequentFailuresAfter time.Duration
	MoreThan3Actions        bool
	ResetFailCountAfter     time.Duration
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
	service, err := s.mgr.OpenService(name)
	if err != nil {
		return nil, fmt.Errorf("could not open service: %w", err)
	}
	defer func() { _ = service.Close() }()

	// Get service configuration
	config, err := service.Config()
	if err != nil {
		return nil, fmt.Errorf("could not get service configuration: %w", err)
	}

	// Get current status
	status, err := service.Query()
	if err != nil {
		return nil, fmt.Errorf("could not query service: %w", err)
	}

	// Get service recovery options
	recoveryActions, err := service.RecoveryActions()
	if err != nil {
		return nil, fmt.Errorf("could not get service recovery actions: %w", err)
	}

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
		return nil, fmt.Errorf("could not get service reset period: %w", err)
	}

	recoveryCommand, err := service.RecoveryCommand()
	if err != nil {
		return nil, fmt.Errorf("could not get service recovery command: %w", err)
	}

	executable, err := findServiceExecutable(config.BinaryPathName)
	if err != nil {
		return nil, err
	}

	wf, err := fi.NewWinFileInfo(executable)
	if err != nil {
		return nil, fmt.Errorf("failed to get win file info for %v: %v", name, err)
	}
	versions, err := wf.GetVersions()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for %v: %v", name, err)
	}
	fileTime, err := wf.GetFileTimestamps()
	if err != nil {
		return nil, fmt.Errorf("failed to get file time for %v: %v", name, err)
	}

	executableDir := filepath.Dir(executable)
	configFiles, err := collectServiceConfigFiles(executableDir)
	if err != nil {
		slog.Warn("not all config files collected", "service", name, "error", err)
	}

	return &ServiceDetails{
		Name:             name,
		DisplayName:      config.DisplayName,
		Description:      config.Description,
		PathToExecutable: config.BinaryPathName,
		StartupType:      startTypeToString(config.StartType),
		ServiceStatus:    stateToString(status.State),
		ServiceType:      serviceTypeToString(config.ServiceType),
		ErrorControl:     errorControlToString(config.ErrorControl),
		Dependencies:     config.Dependencies,
		ServiceStartName: config.ServiceStartName,
		DelayedAutoStart: config.DelayedAutoStart,
		Recovery: ServiceRecovery{
			Command:                 recoveryCommand,
			FirstFailure:            recoveryDetails[0],
			FirstFailureAfter:       recoveryDelays[0],
			SecondFailure:           recoveryDetails[1],
			SecondFailureAfter:      recoveryDelays[1],
			SubsequentFailures:      recoveryDetails[2],
			SubsequentFailuresAfter: recoveryDelays[2],
			MoreThan3Actions:        moreThan3RecoveryActions,
			ResetFailCountAfter:     time.Duration(resetSeconds) * time.Second,
		},
		Executable: ServiceExecutable{
			ExecutableFile: ExecutableFile{
				Path:           executable,
				Version:        versions.FileVersion.String(),
				ProductVersion: versions.ProductVersion.String(),
				CreationTime:   fileTime.CreationTime,
				LastAccessTime: fileTime.LastAccessTime,
				LastWriteTime:  fileTime.LastWriteTime,
			},
			ConfigFiles: configFiles,
		},
	}, nil
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
