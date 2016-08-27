// Copyright (c) 2016. See AUTHORS file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hodor

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

// HTMLTemplateEngine #TODO
type HTMLTemplateEngine interface {
	EnableRecompiling(enable bool)
	CompileTemplates(path string) error
	Render(name string, data interface{}, w http.ResponseWriter)
}

// ============================================================================
// 					Default html template engine (html/template)
// ============================================================================

type goTemplateEngine struct {
	templatePath       string
	templates          *template.Template
	recompilingEnabled bool
}

// NewDefaultTemplateEngine #TODO
func NewDefaultTemplateEngine() HTMLTemplateEngine {
	return &goTemplateEngine{recompilingEnabled: false}
}

func (eng *goTemplateEngine) CompileTemplates(templatePath string) error {
	eng.templatePath = templatePath

	// collect templates
	var tmplList []string
	var err error
	filepath.Walk(templatePath, func(path string, info os.FileInfo, ínternalError error) error {
		if ínternalError != nil {
			err = ínternalError
		} else {
			if !info.IsDir() {
				tmplList = append(tmplList, path)
			}
		}
		return nil
	})

	// compile templates
	templates, err := template.ParseFiles(tmplList...)
	if err != nil {
		panic("Error while compiling template files: " + err.Error())
	}
	eng.templates = templates

	return err
}

func (eng *goTemplateEngine) Render(name string, data interface{}, w http.ResponseWriter) {
	// recompile before render if enabled
	if eng.recompilingEnabled {
		eng.CompileTemplates(eng.templatePath)
	}
	// render
	eng.templates.ExecuteTemplate(w, name, data)
}

func (eng *goTemplateEngine) EnableRecompiling(enable bool) {
	eng.recompilingEnabled = enable
}
