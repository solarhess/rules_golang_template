package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const _Usage = `render A simple utility to render go templates to STDOUT
Usage: 
	render 
		--template=./template.tpl   -- Required. Supplies the template file to use. 
		--values=./values.json      -- Optional. Supplies values for the template. If omitted, the environment will be used.
		--output=./output.txt 		-- Optional. The file where the rendered template will be written. If omitted, stdout.

	values.json contains special directives to load json or text from files on the filesystem.
	Here's an example
	{
		"literal":"a literal string",   // sets the value of  "literal" = "a literal string" 
		"fileValue: {"__FILE__":"./file.txt"},  // sets the value of "fileValue" to the contents of ./file.txt
		"jsonValue: {"__JSON__":"./file.json"},  // sets jsonValue to the parsed json in ./file.json
	}

Go template rules are defined in detail here:
	https://golang.org/pkg/text/template/
`

func loadEnvironmentValues() (templateData map[string]interface{}) {
	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]interface{} {
		items := make(map[string]interface{})
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}

	templateData = getenvironment(os.Environ(), func(item string) (key, val string) {
		index := strings.Index(item, "=")
		if index > -1 {
			key = item[0:index]
			val = item[index+1:]
		} else {
			key = item
			val = ""
		}
		return
	})

	return
}

func loadJsonValues(dataFilename string) (templateData map[string]interface{}, err error) {

	items := make(map[string]interface{})
	templateData = make(map[string]interface{})

	data, errdf := ioutil.ReadFile(dataFilename)
	if errdf != nil {
		panic(errdf)
	}

	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}

	for k, v := range items {
		switch v.(type) {
		case map[string]interface{}:
			mapV := v.(map[string]interface{})

			textFilename, textExists := mapV["__FILE__"]
			jsonFilename, jsonExists := mapV["__JSON__"]

			var mapValue interface{} = nil

			if textExists {
				if _, fileCheckErr := os.Stat(textFilename.(string)); !os.IsNotExist(fileCheckErr) {
					fileData, fileErr := ioutil.ReadFile(textFilename.(string))
					if fileErr == nil {
						mapValue = string(fileData)
					} else {
						mapValue = fmt.Sprintf("<can't open file %s>", textFilename)
					}
				}
			}
			if jsonExists {
				if _, fileCheckErr := os.Stat(jsonFilename.(string)); !os.IsNotExist(fileCheckErr) {
					var loadMapErr error
					mapValue, loadMapErr = loadJsonValues(jsonFilename.(string))
					if loadMapErr != nil {
						mapValue = fmt.Sprintf("<can't load json %s>", jsonFilename)
					}
				}
			}
			if mapValue == nil {
				mapValue = v
			}
			templateData[k] = mapValue
		default:
			templateData[k] = v
		}
	}

	return
}

func loadValues(dataFilename string) map[string]interface{} {
	if dataFilename == "" {
		return loadEnvironmentValues()
	}

	values, err := loadJsonValues(dataFilename)
	if err != nil {
		panic(err)
	}
	return values

}

func main() {

	var dataFilename string
	var templateFilename string
	var outputFilename string
	var dataOutputFilename string

	valuesJsonFilenameFlag := flag.String("values", "", "The values file to use, formatted as JSON")
	templateFilenameFlag := flag.String("template", "", "The filename to read the template from")
	outputFilenameFlag := flag.String("output", "", "The filename to write the template output")
	dataOutputFilenameFlag := flag.String("data_output", "", "The filename to write the json data used to render the template")

	flag.Parse()

	if *outputFilenameFlag != "" {
		outputFilename = *outputFilenameFlag
	} else {
		outputFilename = ""
	}

	if *valuesJsonFilenameFlag != "" {
		dataFilename = *valuesJsonFilenameFlag
	} else {
		dataFilename = ""
	}

	if *templateFilenameFlag != "" {
		templateFilename = *templateFilenameFlag
	} else {
		templateFilename = ""
	}
	if *dataOutputFilenameFlag != "" {
		dataOutputFilename = *dataOutputFilenameFlag
	} else {
		dataOutputFilename = ""
	}

	outputStr, data, renderErr := Render(templateFilename, dataFilename)
	if renderErr != nil {
		panic(renderErr)
	}

	if outputFilename == "" {
		print(outputStr)
	} else {
		err := ioutil.WriteFile(outputFilename, []byte(outputStr), 0644)
		if err != nil {
			panic(err)
		}
	}

	if dataOutputFilename != "" {
		dataJson, _ := json.Marshal(data)
		err := ioutil.WriteFile(dataOutputFilename, dataJson, 0644)
		if err != nil {
			panic(err)
		}
	}

}

func Render(templateFilename string,
	dataFilename string) (output string, data map[string]interface{}, err error) {

	templateData := loadValues(dataFilename)

	templateBytes, errtp := ioutil.ReadFile(templateFilename)
	if errtp != nil {
		return "", nil, errtp
	}
	templateStr := string(templateBytes)

	tmpl, err := template.New("test").Parse(templateStr)

	stringOut := bytes.NewBufferString(output)
	tmpl.Execute(stringOut, templateData)

	return stringOut.String(), templateData, nil
}
