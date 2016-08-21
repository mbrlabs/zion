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
