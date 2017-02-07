package command

import (
	"log"

	"github.com/yangyuqian/genus"
	"github.com/yangyuqian/genus/cmd/genus/command/generate"
	"github.com/yangyuqian/genus/cmd/genus/spec"
	cli "gopkg.in/urfave/cli.v1"
)

var GenerateCmd = cli.Command{
	Name:      "generate",
	ShortName: "g",
	Usage:     "Perform code generation with given Generation Plan Specification",
	Before: func(ctx *cli.Context) (err error) {
		specPath := ctx.String("spec")
		if specPath == "" {
			return
		}

		if err := spec.ValidateSpec(specPath); err != nil {
			log.Printf("Validate specification %s error %+v", specPath, err)
		}
		return
	},
	Action: doGenerate,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "spec", Usage: "Location of Generation Plan Specification"},
	},
	Subcommands: cli.Commands{generate.ORMCmd},
}

func doGenerate(ctx *cli.Context) (err error) {
	planner, err := genus.NewPackagePlanner(ctx.String("spec"))
	if err != nil {
		log.Printf("Can not initialize planner due to %v", err)
		return err
	}

	if planErr := planner.Plan(); planErr != nil {
		log.Printf("Can not warmup planner due to %v", planErr)
		return planErr
	}

	if perfErr := planner.Perform(); perfErr != nil {
		log.Printf("Can not peform planner due to %v", perfErr)
		return perfErr
	}
	return
}
