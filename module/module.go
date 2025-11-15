package module

// IModule defines the contract for application modules in the clean architecture pattern
// Each module must implement these methods to be registered with the application
type IModule interface {
	// GetInfo returns basic information about the module
	GetInfo() *Module
	// RegisterHTTP registers the module's HTTP routes with the application
	RegisterHTTP()
	// RegisterTask registers the module's background tasks with the application
	RegisterTask()
	// RegisterSchedule registers the module's scheduled jobs with the application
	RegisterSchedule()
}

// Module contains basic information about a module
// including its name and routing prefix for API endpoints
type Module struct {
	Name   string // The display name of the module
	Prefix string // The URL prefix used for routing
}
