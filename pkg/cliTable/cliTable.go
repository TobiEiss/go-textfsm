package cliTable

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/TobiEiss/go-textfsm/pkg/ast"
	"github.com/TobiEiss/go-textfsm/pkg/models"
	"github.com/TobiEiss/go-textfsm/pkg/reader"
	"github.com/pkg/errors"
)

type row struct {
	template string
	attrs    map[string]string
}

type CliTable struct {
	templateDir string
	indexPath   string
	headers     []string
	entries     []row
}

func NewCliTable(templateDir string, indexFile string) CliTable {
	indexChan := make(chan string)

	// check if user entered the full path of indexFile
	if !strings.Contains(templateDir, indexFile) {
		indexFile = path.Join(templateDir, indexFile)
	}

	clitable := CliTable{}
	clitable.indexPath = indexFile
	clitable.templateDir = templateDir

	go reader.ReadLineByLine(indexFile, indexChan)
	clitable.parseIndex(indexChan)
	return clitable

}

func (table CliTable) CreateAST(attrs map[string]string) (models.AST, error) {
	tmpChan := make(chan string)
	templateName, err := table.FindTemplate(attrs)
	if err != nil {
		return models.AST{}, err
	}
	template := path.Join(table.templateDir, templateName)
	go reader.ReadLineByLine(template, tmpChan)
	return ast.CreateAST(tmpChan)

}

func (table *CliTable) FindTemplate(attrs map[string]string) (string, error) {
	for _, row := range table.entries {

		if compareAttrs(row.attrs, attrs) {

			return row.template, nil
		}
	}
	return "", errors.New("attrs does not match any template")

}

func (table *CliTable) parseIndex(indexChan chan string) {
	parsedHeaders := false
	for {
		line, ok := <-indexChan
		if !ok {
			break
		}

		// skip comments and empty lines
		if !strings.Contains(line, "#") && line != "" {
			lineSlice := sliceLine(line)

			// set headers, the headers are always the first line
			if !parsedHeaders {
				table.headers = lineSlice
				parsedHeaders = true
				continue
			}

			// create attrs then create a row and add it to table.entries
			tmpRow := row{}
			tmpAttr := make(map[string]string)
			tmpRow.template = lineSlice[0]

			for i, exp := range lineSlice[0:] {
				key := table.headers[i]

				// this is a command so we need to cmdVarCompletion
				if i == len(lineSlice)-1 {
					exp = cmdVarCompletion(exp)
				}

				tmpAttr[key] = strings.TrimSpace(exp)
			}

			tmpRow.attrs = tmpAttr

			table.entries = append(table.entries, tmpRow)

		}
	}

}

// compare if attrs are matching,
func compareAttrs(rowAttr map[string]string, userAttrs map[string]string) bool {
	for key, value := range userAttrs {
		rowRegex := rowAttr[key]
		re := regexp.MustCompile(rowRegex)

		if !re.MatchString(value) {

			return false
		}
	}
	return true

}

// variable length command completion,
// regexp inside [[]] is not supported,
func cmdVarCompletion(command string) string {
	if !strings.Contains(command, "[[") {
		return command
	}
	re, _ := regexp.Compile(`\[\[(\S*)]]`)

	matches := re.FindAllStringSubmatch(command, -1)

	for _, match := range matches {

		variable := match[1]
		variableExp := "\\[\\[" + variable + "]]"
		varRegex, _ := regexp.Compile(variableExp)
		command = varRegex.ReplaceAllString(command, expand(variable))

	}
	return command
}

// this convert "abc[[xyz]]" to "abc(x(y(z)?)?)?"
func expand(input string) string {
	if len(input) == 1 {
		return "(" + input + ")?"
	}
	return fmt.Sprintf("(%c", input[0]) + expand(input[1:]) + ")?"
}

func sliceLine(line string) []string {
	var result []string
	slice := strings.SplitN(line, ",", -1)

	for _, val := range slice {
		result = append(result, strings.TrimSpace(val))
	}
	return result

}
