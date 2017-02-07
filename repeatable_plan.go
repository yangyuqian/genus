package genus

import "errors"

// Plan for generations of repeatable files
type RepeatablePlan struct {
	Data []interface{}
	SingletonPlan
}

func (p *RepeatablePlan) Render() (err error) {
	if err = p.init(); err != nil {
		return err
	}

	if len(p.Data) <= 0 {
		return errors.New("Data of RepeatablePlan not set")
	}

	for _, dat := range p.Data {
		if err = p.TemplateGroup.Render(dat); err != nil {
			return err
		}
	}
	return
}
