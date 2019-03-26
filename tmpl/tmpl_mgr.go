package tmpl

import (
	"github.com/go-courier/helmx/spec"
	"io"
	"text/template"
)

func MergeFuncMap(funcMaps ...template.FuncMap) template.FuncMap {
	funcMap := template.FuncMap{}

	for _, fm := range funcMaps {
		for k, fn := range fm {
			funcMap[k] = fn
		}
	}

	return funcMap
}

func NewTemplateMgr() *TemplateMgr {
	return &TemplateMgr{
		templates: map[string]*Tmpl{},
		funcMap:   MergeFuncMap(KubeFuncs, HelperFuncs),
	}
}

type Tmpl struct {
	Validates []func(s *spec.Spec) bool
	*template.Template
}

type TemplateMgr struct {
	funcMap   template.FuncMap
	templates map[string]*Tmpl
}

func (tplMgr *TemplateMgr) AddFunc(name string, fn interface{}) {
	tplMgr.funcMap[name] = fn
}

func (tplMgr *TemplateMgr) AddTemplate(name string, text string, validates ...func(s *spec.Spec) bool) {
	if err := tplMgr.addTemplate(name, text, validates...); err != nil {
		panic(err)
	}
}

func (tplMgr *TemplateMgr) addTemplate(name string, text string, validates ...func(s *spec.Spec) bool) error {
	tmpl, err := template.New(name).Funcs(tplMgr.funcMap).Parse(text)
	if err != nil {
		return err
	}
	tplMgr.templates[name] = &Tmpl{
		Validates: validates,
		Template:  tmpl,
	}
	return nil
}

func (tplMgr *TemplateMgr) ExecuteAll(writer io.Writer, s *spec.Spec) error {
	count := 0

	for name := range tplMgr.templates {
		if count != 0 {
			_, err := writer.Write([]byte(`

---

`))
			if err != nil {
				return err
			}
		}

		ok, err := tplMgr.execute(name, writer, s)
		if err != nil {
			return err
		}
		if ok {
			count ++
		}
	}
	return nil
}

func (tplMgr TemplateMgr) execute(name string, writer io.Writer, s *spec.Spec) (bool, error) {
	if tmpl, ok := tplMgr.templates[name]; ok {
		valid := true
		if len(tmpl.Validates) > 0 {
			for _, validate := range tmpl.Validates {
				v := validate(s)
				if !v {
					valid = v
					break
				}
			}
		}

		if !valid {
			return false, nil
		}
		err := tmpl.Execute(writer, s)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
