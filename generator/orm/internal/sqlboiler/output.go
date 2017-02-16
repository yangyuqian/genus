package sqlboiler

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"

	"github.com/pkg/errors"
)

var (
	// templateByteBuffer is re-used by all template construction to avoid
	// allocating more memory than is needed. This will later be a problem for
	// concurrency, address it then.
	templateByteBuffer = &bytes.Buffer{}

	rgxRemoveNumberedPrefix = regexp.MustCompile(`[0-9]+_`)
	rgxSyntaxError          = regexp.MustCompile(`(\d+):\d+: `)

	testHarnessWriteFile = ioutil.WriteFile
)

// generateOutput builds the file output and sends it to outHandler for saving
func generateOutput(state *State, data *templateData) error {
	return executeTemplates(executeTemplateData{
		state:                state,
		data:                 data,
		templates:            state.Templates,
		importSet:            defaultTemplateImports,
		combineImportsOnType: true,
		fileSuffix:           ".go",
	})
}

// generateTestOutput builds the test file output and sends it to outHandler for saving
func generateTestOutput(state *State, data *templateData) error {
	return executeTemplates(executeTemplateData{
		state:                state,
		data:                 data,
		templates:            state.TestTemplates,
		importSet:            defaultTestTemplateImports,
		combineImportsOnType: false,
		fileSuffix:           "_test.go",
	})
}

// generateSingletonOutput processes the templates that should only be run
// one time.
func generateSingletonOutput(state *State, data *templateData) error {
	return executeSingletonTemplates(executeTemplateData{
		state:          state,
		data:           data,
		templates:      state.SingletonTemplates,
		importNamedSet: defaultSingletonTemplateImports,
		fileSuffix:     ".go",
	})
}

// generateSingletonTestOutput processes the templates that should only be run
// one time.
func generateSingletonTestOutput(state *State, data *templateData) error {
	return executeSingletonTemplates(executeTemplateData{
		state:          state,
		data:           data,
		templates:      state.SingletonTestTemplates,
		importNamedSet: defaultSingletonTestTemplateImports,
		fileSuffix:     ".go",
	})
}

type executeTemplateData struct {
	state *State
	data  *templateData

	templates *templateList

	importSet      imports
	importNamedSet map[string]imports

	combineImportsOnType bool

	fileSuffix string
}

func executeTemplates(e executeTemplateData) error {
	if e.data.Table.IsJoinTable {
		return nil
	}

	out := templateByteBuffer
	out.Reset()

	var imps imports
	imps.standard = e.importSet.standard
	imps.thirdParty = e.importSet.thirdParty
	if e.combineImportsOnType {
		imps = combineTypeImports(imps, importsBasedOnType, e.data.Table.Columns)
	}

	writePackageName(out, e.state.Config.PkgName)
	writeImports(out, imps)

	for _, tplName := range e.templates.Templates() {
		if err := executeTemplate(out, e.templates.Template, tplName, e.data); err != nil {
			return err
		}
	}

	fName := e.data.Table.Name + e.fileSuffix
	if err := writeFile(e.state.Config.OutFolder, fName, out); err != nil {
		return err
	}

	return nil
}

func executeSingletonTemplates(e executeTemplateData) error {
	if e.data.Table.IsJoinTable {
		return nil
	}

	out := templateByteBuffer
	for _, tplName := range e.templates.Templates() {
		out.Reset()

		fName := tplName
		ext := filepath.Ext(fName)
		fName = rgxRemoveNumberedPrefix.ReplaceAllString(fName[:len(fName)-len(ext)], "")

		imps := imports{
			standard:   e.importNamedSet[fName].standard,
			thirdParty: e.importNamedSet[fName].thirdParty,
		}

		writePackageName(out, e.state.Config.PkgName)
		writeImports(out, imps)

		if err := executeTemplate(out, e.templates.Template, tplName, e.data); err != nil {
			return err
		}

		if err := writeFile(e.state.Config.OutFolder, fName+e.fileSuffix, out); err != nil {
			return err
		}
	}

	return nil
}

func generateTestMainOutput(state *State, data *templateData) error {
	if state.TestMainTemplate == nil {
		return errors.New("No TestMain template located for generation")
	}

	out := templateByteBuffer
	out.Reset()

	var imps imports
	imps.standard = defaultTestMainImports[state.Config.DriverName].standard
	imps.thirdParty = defaultTestMainImports[state.Config.DriverName].thirdParty

	writePackageName(out, state.Config.PkgName)
	writeImports(out, imps)

	if err := executeTemplate(out, state.TestMainTemplate, state.TestMainTemplate.Name(), data); err != nil {
		return err
	}

	if err := writeFile(state.Config.OutFolder, "main_test.go", out); err != nil {
		return err
	}

	return nil
}

// writePackageName writes the package name correctly, ignores errors
// since it's to the concrete buffer type which produces none
func writePackageName(out *bytes.Buffer, pkgName string) {
	_, _ = fmt.Fprintf(out, "package %s\n\n", pkgName)
}

// writeImports writes the package imports correctly, ignores errors
// since it's to the concrete buffer type which produces none
func writeImports(out *bytes.Buffer, imps imports) {
	if impStr := buildImportString(imps); len(impStr) > 0 {
		_, _ = fmt.Fprintf(out, "%s\n", impStr)
	}
}

// writeFile writes to the given folder and filename, formatting the buffer
// given.
func writeFile(outFolder string, fileName string, input *bytes.Buffer) (err error) {
	byt, err := formatBuffer(input)
	if err != nil {
		return err
	}

	path := filepath.Join(outFolder, fileName)
	if err = testHarnessWriteFile(path, byt, 0666); err != nil {
		return errors.Wrapf(err, "failed to write output file %s", path)
	}

	return nil
}

// executeTemplate takes a template and returns the output of the template
// execution.
func executeTemplate(buf *bytes.Buffer, t *template.Template, name string, data *templateData) error {
	if err := t.ExecuteTemplate(buf, name, data); err != nil {
		return errors.Wrapf(err, "failed to execute template: %s", name)
	}
	return nil
}

func formatBuffer(buf *bytes.Buffer) ([]byte, error) {
	output, err := format.Source(buf.Bytes())

	if err == nil {
		return output, nil
	}

	matches := rgxSyntaxError.FindStringSubmatch(err.Error())
	if matches == nil {
		return nil, errors.Wrap(err, "failed to format template")
	}

	lineNum, _ := strconv.Atoi(matches[1])
	scanner := bufio.NewScanner(buf)
	errBuf := &bytes.Buffer{}
	line := 1
	for ; scanner.Scan(); line++ {
		if delta := line - lineNum; delta < -50 || delta > 50 {
			continue
		}

		if line == lineNum {
			errBuf.WriteString(">>>> ")
		} else {
			fmt.Fprintf(errBuf, "% 4d ", line)
		}
		errBuf.Write(scanner.Bytes())
		errBuf.WriteByte('\n')
	}

	return nil, errors.Wrapf(err, "failed to format template\n\n%s\n", errBuf.Bytes())
}
