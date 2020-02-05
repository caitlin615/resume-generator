package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/caitlin615/resume-generator/pdf"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	var err error

	// Flags
	var resumePath = flag.String("resume", "example.yaml", "Path to the resume YAML file")
	flag.Parse()

	// Load resume
	r, err := loadResume(*resumePath)
	if err != nil {
		log.Fatal(err)
	}
	r.ResumeName = strings.TrimSuffix(*resumePath, filepath.Ext(*resumePath))

	// Templates extensions
	// There must exist a template named TemplatesPath+Extension for each extension
	var templatesExtensions = []Exporter{HTML, MD, TXT, JSON, XML}

	// Save resume using the templates
	for _, exporter := range templatesExtensions {
		var filename string
		filename, err = exporter.SaveAs(r)
		if err != nil {
			log.Printf("Unable to export as %s: %s", filename, err)
		} else {
			log.Printf("Exported to %s", (filename))
		}
	}
	// Save resume as PDF
	err = pdf.SaveHTMLAsPDF("output/" + r.ResumeName + ".html")
	if err != nil {
		log.Fatal(err)
	}
}

//
// Resume structure
//

// Resume format
type Resume struct {
	ResumeName string    `json:"-"`
	Name       string    `json:",omitempty"`
	Title      string    `json:",omitempty"`
	Contact    Contact   `json:",omitempty"`
	Summary    string    `json:",omitempty"`
	Sections   []Section `json:",omitempty"`
}

// Contact section
type Contact struct {
	Phone    string `json:",omitempty"`
	Address  string `json:",omitempty"`
	Email    string `json:",omitempty"`
	Webpage  Link   `json:",omitempty"`
	Linkedin Link   `json:",omitempty"`
	Github   Link   `json:",omitempty"`
}

// Link to URL
type Link struct {
	Name string `json:",omitempty"`
	URL  string `json:",omitempty"`
}

// Section of the resume
type Section struct {
	Name    string  `json:",omitempty"`
	Entries []Entry `json:",omitempty"`
}

// Entry of a section
type Entry struct {
	What        string   `json:",omitempty"`
	URL         string   `json:",omitempty"`
	Where       string   `json:",omitempty"`
	When        string   `json:",omitempty"`
	Location    string   `json:",omitempty"`
	Description string   `json:",omitempty"`
	Details     []string `json:",omitempty"`
}

//
// Load from YAML
//

func loadResume(yamlPath string) (*Resume, error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("Open input YAML (%s) failed\n%s", yamlPath, err)
	}
	resume := Resume{}
	err = yaml.Unmarshal(yamlFile, &resume)
	if err != nil {
		return nil, fmt.Errorf("Read input YAML (%s) failed\n%s", yamlPath, err)
	}
	return &resume, nil
}

//
// Export to file
//

var HTML = DefaultExporter{".html"}
var XML = DefaultExporter{".xml"}
var TXT = DefaultExporter{".txt"}
var MD = DefaultExporter{".md"}
var JSON = JSONExporter{}

type Exporter interface {
	SaveAs(r *Resume) (filename string, err error)
}

type DefaultExporter struct {
	Extension string
}

type JSONExporter struct{}

func (j JSONExporter) SaveAs(r *Resume) (outputPath string, err error) {
	outputPath = "output/" + r.ResumeName + ".json"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("Open output file (%s) failed\n%s", outputPath, err)
	}
	defer outputFile.Close()
	output, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Write output file (%s) failed\n%s", outputPath, err)
	}
	outputFile.Write(output)
	return
}

func (d DefaultExporter) SaveAs(r *Resume) (outputPath string, err error) {
	templatePath := "templates/tmpl" + d.Extension
	outputPath = "output/" + r.ResumeName + d.Extension
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("Parse template file (%s) failed\n%s", templatePath, err)
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("Open output file (%s) failed\n%s", outputPath, err)
	}
	defer outputFile.Close()
	err = tmpl.Execute(outputFile, *r)
	if err != nil {
		return "", fmt.Errorf("Execute template (%s) failed\n%s", templatePath, err)
	}
	return
}
