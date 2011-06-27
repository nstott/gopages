package pages

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"template"
	"os"
	"path"
	"path/filepath"
)

const(
	templateExtension = ".tpl"
)

type App struct {
	Pages map[string]*Page
	FormatterMap map[string]func(io.Writer, string, ...interface{})
}

type Page struct {
	Template *template.Template
	Filename string
	Data map[string]interface{}
}


// debug template
var debugTplString string = "{page}<div><h3>Debug</h3><h4>Headers</h4>" +
	"{.section headers}{.repeated section @}{Key} = {Value}<br />{.end}" +
	"{.or}No Headers{.end}<br />" +
	"<h4>Url Params</h4>{.section params}{.repeated section @}{Key} = {Value}<br />" +
	"{.or}No params{.end}{.end}" + 
	"</div>"

var debugTpl *template.Template


//formatters



/* 
 * Allows templates to embed other templates 
 */
func EmbedFormatter(wr io.Writer, str string, data ...interface{}) {

	var b []byte
	var ok bool
	if len(data) == 1 {
		b, ok = data[0].([]byte)
	}

    	if !ok {
    		var buf bytes.Buffer
    		fmt.Fprint(&buf, data...)
    		b = buf.Bytes()
	}
	fmt.Fprint(wr,string(b))
}

/* 
 * Page methods 
 */
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
func (app *App) AddPage(id, filename string) (page *Page, err os.Error) {
	p := &Page{Filename: filename}
	app.Pages[id] = p
	return p, nil
}

// parse all the templates
func (app *App) ParseAllPages() {
	for k,v := range(app.Pages) {
		log.Printf("Parsing %s", k)
		v.ParseFile(app.FormatterMap)				
	}
}


//add a directory of templates to the template map
//templates are added with their filename as their id
func(app *App) AddDirectory(dirname string) (err os.Error) {
	v := &TemplateVisitor{templateBase: dirname, app: app}
	filepath.Walk(dirname, v, nil)
	
	for k,_ := range(app.Pages) {
		app.FormatterMap["embed:" + k] = EmbedFormatter
	}

	app.ParseAllPages()
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
		_, err := v.app.AddPage(p[len(v.templateBase):], p)
		if err != nil {
			log.Fatal("Cannot parse Template: " + err.String())
		}
		log.Printf("Parsed Template %s, %s\n",p[len(v.templateBase):], p)
	}
}
