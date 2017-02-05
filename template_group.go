package genus

// "log"

type Imports map[string]string

// Represents single Go package
type TemplateGroup struct {
	Package        string
	Imports        Imports
	Templates      []*Template
	SkipFixImports bool
}

// // TODO: V1 doesn't support windows
// func (tg *TemplateGroup) ensureGOOS() (err error) {
// 	if runtime.GOOS == "windows" {
// 		return errors.New("Windows is not supported for now")
// 	}
//
// 	return
// }
//
// // Operations must be performed under gopath
// func (tg *TemplateGroup) ensureGopath() (err error) {
// 	gopath := os.Getenv("GOPATH")
// 	strings.Split(gopath, ":")
//
// 	return
// }
