package processors

import (
	"main/types"
)

type Processors struct {
	Runtime *types.Runtime
}

func (p *Processors) All() []types.Processor {
	return []types.Processor{
		p.structuredData,
		p.inference,
		p.validate,
	}
}
