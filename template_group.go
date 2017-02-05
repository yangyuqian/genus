package genus

type Imports map[string]string

// Represents single Go package
type TemplateGroup struct {
	Package        string
	Imports        Imports
	Templates      []*Template
	SkipFixImports bool
}
