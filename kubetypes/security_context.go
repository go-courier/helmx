package kubetypes

type SecurityContext struct {
	Capabilities             *Capabilities   `yaml:"capabilities,omitempty"`
	RunAsUser                *int64          `yaml:"runAsUser,omitempty"`
	RunAsGroup               *int64          `yaml:"runAsGroup,omitempty"`
	RunAsNonRoot             *bool           `yaml:"runAsNonRoot,omitempty"`
	ReadOnlyRootFilesystem   *bool           `yaml:"readOnlyRootFilesystem,omitempty"`
	AllowPrivilegeEscalation *bool           `yaml:"allowPrivilegeEscalation,omitempty"`
	ProcMount                *string         `yaml:"procMount,omitempty"`
	Privileged               *bool           `yaml:"privileged,omitempty"`
	SELinuxOptions           *SELinuxOptions `yaml:"seLinuxOptions,omitempty"`
}

type SELinuxOptions struct {
	// User is a SELinux user label that applies to the container.
	// +optional
	User string `yaml:"user,omitempty"`
	// Role is a SELinux role label that applies to the container.
	// +optional
	Role string `yaml:"role,omitempty"`
	// Type is a SELinux type label that applies to the container.
	// +optional
	Type string `yaml:"type,omitempty"`
	// Level is SELinux level label that applies to the container.
	// +optional
	Level string `yaml:"level,omitempty"`
}

// Capability represent POSIX capabilities type
type Capability string

// Adds and removes POSIX capabilities from running containers.
type Capabilities struct {
	// Added capabilities
	// +optional
	Add []Capability `yaml:"add,omitempty"`
	// Removed capabilities
	// +optional
	Drop []Capability `yaml:"drop,omitempty"`
}
