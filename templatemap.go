/*
Package templatemap implements a map of names to html templates
and utility functions to load them from directories.
*/
package templatemap

import (
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Map map[string]*template.Template

// Loads the templates from a given directory recursively.
// Each template is added to the map as it's relative path
// from the directory being loaded.
// Each directory may also contain a "_base.tmpl" file,
// wich will be added as an associated template to the
// directories and its childrens templates.
func LoadDir(path string) (Map, error) {
	m := make(Map)
	err := m.loadDir(nil, path, "")
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m Map) loadDir(super *template.Template, path, name string) error {
	basePath := filepath.Join(path, "_base.tmpl")
	baseName := name + "_base.tmpl"

	base, err := loadTemplate(super, basePath, baseName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, info := range infos {
		tmplPath := filepath.Join(path, info.Name())
		tmplName := name + info.Name()

		if info.IsDir() {
			err = m.loadDir(base, tmplPath, tmplName+"/")
			if err != nil {
				return err
			}
			continue
		}

		if filepath.Ext(info.Name()) != ".tmpl" || info.Name() == "_base.tmpl" {
			continue
		}

		tmpl, err := loadTemplate(base, tmplPath, tmplName)
		if err != nil {
			return err
		}
		m[tmplName] = tmpl
	}
	return nil
}

func loadTemplate(super *template.Template, path, name string) (*template.Template, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	_, err = io.Copy(&b, file)
	if err != nil {
		return nil, err
	}
	var tmpl *template.Template
	if super == nil {
		tmpl = template.New(name)
	} else {
		superC, err := super.Clone()
		if err != nil {
			return nil, err
		}
		tmpl = superC.New(name)
	}
	return tmpl.Parse(b.String())
}
