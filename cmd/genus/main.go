package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"github.com/yangyuqian/genus"
	"github.com/yangyuqian/genus/cmd/genus/schema"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	specPath := os.Args[1]

	if err := validateSpec(specPath); err != nil {
		log.Printf("Validate specification error %+v", err)
	}

	planner, err := genus.NewPackagePlanner(specPath)
	if err != nil {
		log.Printf("Can not initialize planner %v", err)
		os.Exit(1)
	}

	if planErr := planner.Plan(); planErr != nil {
		log.Printf("Can not warmup planner %v", planErr)
		os.Exit(1)
	}

	if perfErr := planner.Perform(); perfErr != nil {
		log.Printf("Can not peform planner %v", perfErr)
		os.Exit(1)
	}
}

func validateSpec(specPath string) (err error) {
	if !filepath.IsAbs(specPath) {
		specPath, err = filepath.Abs(specPath)
		if err != nil {
			return err
		}
	}

	schemaLoader := gojsonschema.NewBytesLoader(schema.MustAsset("cmd/genus/schema/plan_schema.json"))

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
