package render

import (
	"Quiklink_BE/internal/models"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var pathToTemplates = "./templates"

// TemplateRenderer renders templates using html/template
func TemplateRenderer(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	tc, _ = CreateTemplateCache()

	//get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		return errors.New("Could not get template from cache")
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	//get all the files name *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	//range through all pages
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
