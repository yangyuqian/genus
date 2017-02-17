package generate

import (
	"github.com/pkg/errors"
	cli "gopkg.in/urfave/cli.v1"
)

func doGenerate(ctx *cli.Context) (err error) {
	switch framework := ctx.String("framework"); framework {
	case "sqlboiler":
		return doGenerateSqlboiler(ctx)
	case "gorm":
		return doGenerateORM(ctx)
	default:
		return errors.Errorf("Unkown framework <%s>", framework)
	}

	return
}
