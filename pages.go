package pages

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"template"
	"os"
	)


type App struct {
	Pages map[string]*Page
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
var formatterMap = template.FormatterMap{
	"embed": EmbedFormatter,
}


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
func (t *Page) ParseFile() {
	var err os.Error
	t.Template, err = template.ParseFile(t.Filename, formatterMap)
	if err != nil {
		log.Fatal("Cannot Parse " + t.Filename + "; got " + err.String())
	}
}





//execute the template stored under id.
func (app *App) Execute(id string, wr io.Writer, data map[string]interface{}) {
	app.Pages[id].Data = data
	app.Pages[id].Execute(wr)
}

//create a new App
func NewApp() (*App) {
	app := new(App)
	app.Pages = make(map[string]*Page)
	return app
}

//add a page to the template map, and parse the file as a template
func (app *App) AddPage(id, filename string) (page *Page, err os.Error) {
	p := new(Page)
	p.Filename = filename
	p.ParseFile()
	app.Pages[id] = p
	return p, nil
}

//add a directory of templates to the template map
//templates are added with their filename as their id
func(app *App) AddDirectory(dirname string) (err os.Error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	for _, file := range(files) {
		log.Println(file.Name)

		_, err := app.AddPage(file.Name, dirname + "/" + file.Name)
		if err != nil {
			return err
		}
	}



	return nil
	
}



