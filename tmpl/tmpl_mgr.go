package tmpl

import (
	"io"
	"text/template"

	"github.com/go-courier/helmx/spec"
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
		templates: map[string]*template.Template{},
		funcMap:   MergeFuncMap(KubeFuncs, HelperFuncs),
	}
}

type TemplateMgr struct {
	funcMap   template.FuncMap
	templates map[string]*template.Template
}

func (tplMgr *TemplateMgr) AddFunc(name string, fn interface{}) {
	tplMgr.funcMap[name] = fn
}

func (tplMgr *TemplateMgr) AddTemplate(name string, text string) {
	if err := tplMgr.addTemplate(name, text); err != nil {
		panic(err)
	}
}

func (tplMgr *TemplateMgr) addTemplate(name string, text string) error {
	tmpl, err := template.New(name).Funcs(tplMgr.funcMap).Parse(text)
	if err != nil {
		return err
	}
	tplMgr.templates[name] = tmpl
	return nil
}

func (tplMgr *TemplateMgr) ExecuteAll(writer io.Writer, s *spec.Spec) error {
	for name := range tplMgr.templates {
		if err := tplMgr.execute(name, writer, s); err != nil {
			return err
		}
	}
	return nil
}

func (tplMgr TemplateMgr) execute(name string, writer io.Writer, s *spec.Spec) error {
	if tmpl, ok := tplMgr.templates[name]; ok {
		if err := tmpl.Execute(writer, s); err != nil {
			return err
		}
	}
	return nil
}
