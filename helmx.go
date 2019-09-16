package helmx

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-courier/helmx/spec"
	"github.com/go-courier/helmx/tmpl"
	"gopkg.in/yaml.v2"
)

func MustBytes(data []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return data
}

func NewHelmX() *HelmX {
	return &HelmX{
		TemplateMgr: tmpl.NewTemplateMgr(),
	}
}

type HelmX struct {
	spec.Spec
	*tmpl.TemplateMgr
}

func (hx *HelmX) FromYAML(data []byte) error {
	return yaml.Unmarshal(data, &hx.Spec)
}

func (hx *HelmX) ToYAML() ([]byte, error) {
	return yaml.Marshal(hx.Spec)
}

func MustReadOrFetch(fileOrURL string) []byte {
	data, err := ReadOrFetch(fileOrURL)
	if err != nil {
		panic(err)
	}
	return data
}

func ReadOrFetch(fileOrURL string) ([]byte, error) {
	if strings.HasPrefix(fileOrURL, "http") {
		resp, err := http.Get(fileOrURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}
	return ioutil.ReadFile(fileOrURL)
}
