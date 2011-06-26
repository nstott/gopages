package pages

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"template"
	"os"
	)

type Page struct {
	Template *template.Template
	Filename string
	Data map[string]interface{}
}


//type FormatterMap map[string]func(io.Writer, string, ...interface{})
var formatterMap = template.FormatterMap{
	"embed": Embed,
}

var Pages = make(map[string]*Page)

var debugTplString string = "{page}<div><h3>Debug</h3><h4>Headers</h4>" +
	"{.section headers}{.repeated section @}{Key} = {Value}<br />{.end}" +
	"{.or}No Headers{.end}<br />" +
	"<h4>Url Params</h4>{.section params}{.repeated section @}{Key} = {Value}<br />" +
	"{.or}No params{.end}{.end}" + 
	"</div>"

var debugTpl *template.Template


/* 
 * Allows templates to embed other templates 
 */
func Embed(wr io.Writer, str string, data ...interface{}) {

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
 * methods to help with templates 
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

//
func NewPage(id, filename string)  (page *Page, err os.Error) {
	p := new (Page)
	p.Filename = filename
	p.ParseFile()
	Pages[id] = p
	return p, nil
}

//execute the template stored under id.
//we use this, so that we don't have to expose the page map
func Execute(id string, wr io.Writer, data map[string]interface{}) {
	Pages[id].Data = data
	Pages[id].Execute(wr)
}






