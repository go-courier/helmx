package kubetypes

type SecurityContext struct {
	Capabilities             *Capabilities   `yaml:"capabilities,omitempty" json:"capabilities,omitempty"`
	RunAsUser                *int64          `yaml:"runAsUser,omitempty" json:"runAsUser,omitempty"`
	RunAsGroup               *int64          `yaml:"runAsGroup,omitempty" json:"runAsGroup,omitempty"`
	RunAsNonRoot             *bool           `yaml:"runAsNonRoot,omitempty" json:"runAsNonRoot,omitempty"`
	ReadOnlyRootFilesystem   *bool           `yaml:"readOnlyRootFilesystem,omitempty" json:"readOnlyRootFilesystem,omitempty"`
	AllowPrivilegeEscalation *bool           `yaml:"allowPrivilegeEscalation,omitempty" json:"allowPrivilegeEscalation,omitempty"`
	ProcMount                *string         `yaml:"procMount,omitempty" json:"procMount,omitempty"`
	Privileged               *bool           `yaml:"privileged,omitempty" json:"privileged,omitempty"`
	SELinuxOptions           *SELinuxOptions `yaml:"seLinuxOptions,omitempty" json:"seLinuxOptions,omitempty"`
}

type SELinuxOptions struct {
	// User is a SELinux user label that applies to the container.
	// +optional
	User string `yaml:"user,omitempty" json:"user,omitempty"`
	// Role is a SELinux role label that applies to the container.
	// +optional
	Role string `yaml:"role,omitempty" json:"role,omitempty"`
	// Type is a SELinux type label that applies to the container.
	// +optional
	Type string `yaml:"type,omitempty" json:"type,omitempty"`
	// Level is SELinux level label that applies to the container.
	// +optional
	Level string `yaml:"level,omitempty" json:"level,omitempty"`
}

// Capability represent POSIX capabilities type
type Capability string

// Adds and removes POSIX capabilities from running containers.
type Capabilities struct {
	// Added capabilities
	// +optional
	Add []Capability `yaml:"add,omitempty" json:"add,omitempty"`
	// Removed capabilities
	// +optional
	Drop []Capability `yaml:"drop,omitempty" json:"drop,omitempty"`
}
