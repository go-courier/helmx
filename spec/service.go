package spec

import (
    "bytes"
    "errors"
    "io"
    "strings"

    "github.com/go-courier/helmx/kubetypes"
)

type Service struct {
    Pod                      `yaml:",inline"`
    kubetypes.DeploymentOpts `yaml:",inline"`

    Ports     []Port        `json:"ports,omitempty" yaml:"ports,omitempty"`
    Ingresses []IngressRule `json:"ingresses,omitempty" yaml:"ingresses,omitempty"`
    TLS       []IngressTLS  `json:"tls,omitempty" yaml:"tls,omitempty"`

    Headless bool `json:"headless,omitempty" yaml:"headless,omitempty"`
}

type Pod struct {
    Initials                []Container `json:"initials,omitempty" yaml:"initials,omitempty"`
    Container               `yaml:",inline"`
    kubetypes.PodOpts       `yaml:",inline"`
    ServiceAccountRoleRules []RoleRule `yaml:"serviceAccountRoleRules,omitempty" json:"serviceAccountRoleRules,omitempty"`
    Hosts                   []Hosts    `yaml:"hosts,omitempty" json:"hosts,omitempty"`
}

func ParseRoleRule(r string) (*RoleRule, error) {
    parts := strings.Split(r, "#")

    if len(parts) != 2 {
        return nil, errors.New("invalid role rule")
    }

    s := &RoleRule{}
    s.Verbs = strings.Split(parts[1], ",")

    appGroups := ""
    resources := ""
    resourceNames := ""

    appGroupsAndResources := strings.Split(parts[0], ".")

    if len(appGroupsAndResources) == 2 {
        appGroups = appGroupsAndResources[0]
        resources = appGroupsAndResources[1]
    } else {
        resources = appGroupsAndResources[0]
    }

    resourceAndNames := strings.Split(resources, "=")
    if len(resourceAndNames) == 2 {
        resources = resourceAndNames[0]
        resourceNames = resourceAndNames[1]
    } else {
        resources = resourceAndNames[0]
    }

    s.ApiGroups = strings.Split(appGroups, ",")
    s.Resources = strings.Split(resources, ",")
    s.ResourceNames = strings.Split(resourceNames, ",")

    if len(s.ResourceNames) == 1 && s.ResourceNames[0] == "" {
        s.ResourceNames = nil
    }

    return s, nil
}

// openapi:strfmt role-rule
type RoleRule struct {
    // apps,extensions.deployments#get,list,patch,create,update,delete
    // apps,extensions.deployments=a,b,c,d#update,get
    ApiGroups     []string
    Resources     []string
    ResourceNames []string

    Verbs []string
}

func (r RoleRule) String() string {
    buf := bytes.NewBuffer(nil)

    if len(r.ApiGroups) > 0 && r.ApiGroups[0] != "" {
        _, _ = io.WriteString(buf, strings.Join(r.ApiGroups, ","))
        _, _ = io.WriteString(buf, ".")
    }

    _, _ = io.WriteString(buf, strings.Join(r.Resources, ","))

    if len(r.ResourceNames) > 0 {
        _, _ = io.WriteString(buf, "=")
        _, _ = io.WriteString(buf, strings.Join(r.ResourceNames, ","))
    }

    _, _ = io.WriteString(buf, "#")
    _, _ = io.WriteString(buf, strings.Join(r.Verbs, ","))

    return buf.String()
}

func (r RoleRule) MarshalText() ([]byte, error) {
    return []byte(r.String()), nil
}

func (r *RoleRule) UnmarshalText(data []byte) error {
    ir, err := ParseRoleRule(string(data))
    if err != nil {
        return err
    }
    *r = *ir
    return nil
}
