package tmpl

import (
	"bytes"
	"encoding/json"
	"github.com/go-courier/reflectx"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
	"text/template"
)

var HelperFuncs = template.FuncMap{
	"default":      valueDefault,
	"spaces":       spaces,
	"toJson":       toJson,
	"toYamlIndent": toYamlIndent,
	"quote":        strconv.Quote,
	"join":         strings.Join,
	"repeat":       strings.Repeat,
	"trimSpace":    strings.TrimSpace,
}

func toJson(v interface{}) string { // nolint
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func toYamlIndent(v interface{}, ident string) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	if bytes.HasPrefix(data, []byte{'{', '}'}) || bytes.HasPrefix(data, []byte{'[', ']'}) {
		return ""
	}
	return indent(ident, string(data))
}

func valueDefault(d interface{}, given ...interface{}) interface{} {
	if reflectx.IsEmptyValue(given) || reflectx.IsEmptyValue(given[0]) {
		return d
	}
	return given[0]
}

func indent(ident string, v string) string {
	return ident + strings.Replace(strings.TrimSpace(v), "\n", "\n"+ident, -1)
}

func spaces(spaces int) string {
	return strings.Repeat(" ", spaces)
}
