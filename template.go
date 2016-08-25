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
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

// HTMLTemplateEngine #TODO
type HTMLTemplateEngine interface {
	CompileTemplates(path string) error
	Render(name string, data interface{}, w http.ResponseWriter)
}

// ============================================================================
// 					Default html template engine (html/template)
// ============================================================================

type goTemplateEngine struct {
	templates *template.Template
}

// NewDefaultTemplateEngine #TODO
func NewDefaultTemplateEngine() HTMLTemplateEngine {
	return &goTemplateEngine{}
}

func (eng *goTemplateEngine) CompileTemplates(pathToTemplates string) error {
	// collect templates
	var tmplList []string
	var err error
	filepath.Walk(pathToTemplates, func(path string, info os.FileInfo, ínternalError error) error {
		fmt.Println(path)
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
	eng.templates.ExecuteTemplate(w, name, data)
}
