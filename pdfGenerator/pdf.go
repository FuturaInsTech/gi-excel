package pdfGenerator

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// pdf requestpdf struct
type RequestPdf struct {
	body string
}

// new request to pdf function
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

// parsing template function
func (r *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()

	return nil
}

// generate pdf function

func (r *RequestPdf) GeneratePDF(writer io.Writer) (bool, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	page := wkhtmltopdf.NewPageReader(strings.NewReader(r.body))
	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)
	pdfg.SetOutput(writer)
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	return true, nil

}
