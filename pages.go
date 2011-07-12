
//package pages implements a set of types and methods to help manage multiple pages across an web application.
//
//Apps are created, and parse directories for .tpl files.  Each .tpl file becomes a Page, 
// and is able to be executed and rendered to an io.Writer
// 
//A formatter map is also applied to templates as they are parsed, and enables the execution of nested templates.
// if a directory contains index.html, and a subdirectory 'bits' contains a file named header.tpl
// then within index.html, you can execute a subtemplate by calling {@|embed:bits/header.tpl}
// this passes the entire data context to the sub template.  
//

package pages

import (
	"io"
	"log"
	"template"
	"os"
	"path"
	"path/filepath"
)

const(
	templateExtension = ".tpl"
	embedCommand = "embed"
)

// debug template

var debugTplString string = "{page}<div><h3>Debug</h3><h4>Headers</h4>" +
	"{.section headers}{.repeated section @}{Key} = {Value}<br />{.end}" +
	"{.or}No Headers{.end}<br />" +
	"<h4>Url Params</h4>{.section params}{.repeated section @}{Key} = {Value}<br />" +
	"{.or}No params{.end}{.end}" + 
	"</div>"

var debugTpl *template.Template





//A page contains a parsed template, the filename, and a map of the data
type Page struct {
	Template *template.Template
	Filename string
	Data map[string]interface{}
}


// use the page's internal data map and execute a page that has already been parsed.
func (t *Page) Execute(wr io.Writer) {	
	if t == nil {
		log.Println("specified page is nil")
		return
	}
	
	if err := t.Template.Execute(wr, t.Data); err != nil {
		log.Printf("error executing template: %s", err)
	}
}

// parse a file, and store the templates in Page.Template
func (t *Page) ParseFile(formatterMap map[string]func(io.Writer, string, ...interface{})) {
	var err os.Error
	t.Template, err = template.ParseFile(t.Filename, formatterMap)
	if err != nil {
		log.Fatal("Cannot Parse " + t.Filename + "; got " + err.String())
	}
}



// an app contains a map of pages, as well as a formattermap.
// add a directory containing .tpl files to an app, and then perform 
type App struct {
	Pages map[string]*Page
	FormatterMap map[string]func(io.Writer, string, ...interface{})
}


//execute the template stored under id.
func (app *App) Execute(id string, wr io.Writer, data map[string]interface{}) {
	page, ok := app.Pages[id]
	if ok {
		page.Data = data
		page.Execute(wr)
	} else {
		log.Fatal("Template " + id + " Not Found")
	}
}

//create a new App
func NewApp() (*App) {
	app := new(App)
	app.Pages = make(map[string]*Page)
	app.FormatterMap = template.FormatterMap{}
	return app
}

//add a page to the template map, and parse the file as a template
func (app *App) addPage(id, filename string) (page *Page, err os.Error) {
	p := &Page{Filename: filename}
	app.Pages[id] = p
	return p, nil
}

// parse all the templates
func (app *App) parseAllPages() {
	for k,v := range(app.Pages) {
		log.Printf("Parsing %s %v\n", k, app.FormatterMap)
		v.ParseFile(app.FormatterMap)
		log.Printf("Finished Parsing %s\n",k)
	}
}


//add a directory of templates to the template map
//templates are added with their filename as their id
func(app *App) AddDirectory(dirname string) (err os.Error) {
	v := &TemplateVisitor{templateBase: dirname, app: app}
	filepath.Walk(dirname, v, nil)
	
	for k,_ := range(app.Pages) {
		app.FormatterMap[embedCommand + ":" + k] = func(wr io.Writer, str string, data ...interface{}) {
			filename := str[len(embedCommand)+1:]
			if page,ok := app.Pages[filename]; ok {
				err := page.Template.Execute(wr, data)
				if err != nil {
					log.Println("Failed to execute " + filename)
				}
			}
		}
	}

	app.parseAllPages()
	return nil	
}

// a visitor that carries the app and template base and catches all properly named template files
type TemplateVisitor struct {
	app *App
	templateBase string
}

// visit all directories
func (v *TemplateVisitor) VisitDir(p string, f *os.FileInfo) bool {
	return true
}
// for each file matching the proper extension, we want to add it to the app's template set
func (v *TemplateVisitor) VisitFile(p string, f *os.FileInfo) {
	if path.Ext(f.Name) == templateExtension {
		_, err := v.app.addPage(p[len(v.templateBase):], p)
		if err != nil {
			log.Fatal("Cannot parse Template: " + err.String())
		}
		log.Printf("Parsed Template %s, %s\n",p[len(v.templateBase):], p)
	}
}
