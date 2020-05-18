package spec

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var versionRegexp = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)

func ParseVersion(s string) (*Version, error) {
	matched := versionRegexp.FindAllStringSubmatch(s, -1)

	if len(matched) == 0 || len(matched[0]) != 4 {
		return nil, errors.New(s + " is not an available version")
	}

	parts := matched[0]

	major, _ := strconv.ParseInt(parts[1], 10, 10)
	minor, _ := strconv.ParseInt(parts[2], 10, 10)
	patch, _ := strconv.ParseInt(parts[3], 10, 10)

	v := &Version{
		Major: int(major),
		Minor: int(minor),
		Patch: int(patch),
	}

	if s != parts[0] {
		ps := strings.Split(s, "-")

		switch len(ps) {
		case 3:
			v.Suffix = ps[2]
			v.Prefix = ps[0]
		case 2:
			if parts[0] == ps[0] {
				v.Suffix = ps[1]
			} else {
				v.Prefix = ps[0]
			}
		}
	}

	return v, nil
}

// openapi:strfmt version
type Version struct {
	Suffix string
	Prefix string
	Major  int
	Minor  int
	Patch  int
}

func (v Version) String() string {
	versions := []string{
		strconv.Itoa(v.Major),
		strconv.Itoa(v.Minor),
		strconv.Itoa(v.Patch),
	}

	version := strings.Join(versions, ".")

	if v.Prefix != "" {
		version = v.Prefix + "-" + version
	}

	if v.Suffix != "" {
		version = version + "-" + v.Suffix
	}

	return strings.ToLower(version)
}

func (v Version) IncrMajor() Version {
	v.Major = v.Major + 1
	v.Minor = 0
	v.Patch = 0
	return v
}

func (v Version) IncrMinor() Version {
	v.Minor = v.Minor + 1
	v.Patch = 0
	return v
}

func (v Version) IncrPatch() Version {
	v.Patch = v.Patch + 1
	return v
}

func (v Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *Version) UnmarshalText(data []byte) error {
	version, err := ParseVersion(string(data))
	if err != nil {
		return err
	}
	*v = *version
	return nil
}
