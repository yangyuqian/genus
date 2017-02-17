package spec

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

const (
	schemaPath = "cmd/genus/spec/plan_schema.json"
)

func ValidateSpec(specPath string) (err error) {
	if !filepath.IsAbs(specPath) {
		specPath, err = filepath.Abs(specPath)
		if err != nil {
			return err
		}
	}

	schemaLoader := gojsonschema.NewBytesLoader(MustAsset(schemaPath))

	specLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", specPath))

	res, err := gojsonschema.Validate(schemaLoader, specLoader)
	if err != nil {
		return err
	}

	if !res.Valid() {
		for _, e := range res.Errors() {
			log.Printf("Invalid '%s': %s", e.Field(), e.Description())
		}
		return errors.New("Specification validation failed")
	}

	log.Printf("Validate specification %s successfully", specPath)

	return
}
