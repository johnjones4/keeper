package processors

import (
	"main/core"
)

type Processors struct {
	Runtime *core.Runtime
}

func (p *Processors) All() []core.Processor {
	return []core.Processor{
		p.structuredData,
		p.inference,
		p.validate,
	}
}
