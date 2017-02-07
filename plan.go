package genus

type Plan interface {
	Render() error
}

type PlanItem struct {
	Package         string
	PlanType        string
	Suffix          string
	TemplateDir     string
	BaseDir         string
	BasePackage     string
	RelativePackage string
	Filename        string
	TemplateNames   []string
	SkipExists      bool
	SkipFormat      bool
	SkipFixImports  bool
	Merge           bool
	Imports         Imports
	Data            []interface{}
}
