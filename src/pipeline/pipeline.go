package pipeline

import "fmt"

const (
	constPktTypeNew = iota
	constPktTypeExist
)

const (
	constStageInvalid = -1
	constStage0       = iota
	constStage1
	constStage2
	constStage3
	constStageMax
)

type Pkt struct {
	id      uint32
	pktType uint32
}

type Cache struct {
	pkt   *Pkt
	stage int
}

type Pipeline struct {
	caches  []*Cache
	current uint32
	pkts    []*Pkt
	count	int
}

func (pl *Pipeline) stage0(cache *Cache) {
	cache.stage = constStage1
}

func (pl *Pipeline) stage1(cache *Cache) {
	if cache.pkt.pktType == constPktTypeNew {
		cache.stage = constStage2
	} else {
		cache.stage = constStage3
	}
}

func (pl *Pipeline) stage2(cache *Cache) {
	cache.stage = constStage3
}

func (pl *Pipeline) stage3(cache *Cache) {
	cache.stage = constStageInvalid
}

func (pl *Pipeline) doStage(cache *Cache) {
	pl.count++
	fmt.Printf("ID[%d], ", pl.count)
	current := pl.current
	for i:=0; i<len(pl.caches); i++ {
		cache := pl.caches[current]
		if cache.stage == constStageInvalid {
			break
		}

		fmt.Printf("pkt%d.%d[%d] ", cache.pkt.id, cache.pkt.pktType, cache.stage)
		current = pl.nextCacheIndex(current)
	}
	fmt.Println()
	switch cache.stage {
	case constStage0:
		pl.stage0(cache)
	case constStage1:
		pl.stage1(cache)
	case constStage2:
		pl.stage2(cache)
	case constStage3:
		pl.stage3(cache)
	default:
		fmt.Println("Invalid stage: ", cache.stage)
	}
}

func (pl *Pipeline) nextCacheIndex(current uint32) uint32 {
	cLen := len(pl.caches)
	current++

	if current == uint32(cLen) {
		current = 0
	}

	return current

}

func (pl *Pipeline) next() {
	len := len(pl.caches)
	pl.current++

	if pl.current == (uint32(len)) {
		pl.current = 0
	}
}

func (pl *Pipeline) pipeRun(pkt *Pkt, end bool, index int) {
	if end {
		for i := constStage0; i < constStageMax; i++ {
			first := pl.current
			cache := pl.caches[pl.current]
			if cache.stage != constStageInvalid {
				pl.doStage(cache)
			}
			pl.next()
			for ;first!=pl.current;{
				cache := pl.caches[pl.current]
				if int(cache.stage) != constStageInvalid {
					pl.doStage(cache)
				}
				pl.next()
			}
		}

		return
	}

	pl.caches[pl.current].pkt = pkt
	pl.caches[pl.current].stage = constStage0
	for {
		cache := pl.caches[pl.current]
		if cache.stage != constStageInvalid {
			pl.doStage(cache)
			if cache.stage == constStageInvalid {
				return
			}
		} else {
			return
		}

		pl.next()
	}
}

func (pl *Pipeline) Run() {
	for i, pkt := range pl.pkts {
		pl.pipeRun(pkt, false, i)
	}

	pl.pipeRun(nil, true, len(pl.pkts))
}

func NewPipeline(cacheSize int, pkts []*Pkt) *Pipeline {
	pl := &Pipeline{}
	pl.pkts = pkts
	pl.caches = make([]*Cache, cacheSize)

	for i:=0; i<cacheSize; i++ {
		pl.caches[i] = &Cache{nil, constStageInvalid}
	}

	pl.current = 0
	pl.count = 0

	return pl
}

func NewPkts() []*Pkt {
	pkt := make([]*Pkt, 5)

	pkt[0] = &Pkt{0, constPktTypeNew}
	pkt[1] = &Pkt{1, constPktTypeExist}
	pkt[2] = &Pkt{2, constPktTypeNew}
	pkt[3] = &Pkt{3, constPktTypeExist}
	pkt[4] = &Pkt{4, constPktTypeNew}

	return pkt
}
