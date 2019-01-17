package libraries

"github.com/marzhall/zanth/libraries/worldParts"

type GameEntities struct {
	Entities         []Entity
    Players          []worldParts.Player
	emap             map[Entity]int
	NumberOfEntities int
	openSpots        chan int
}

var AllGameEntities GameEntities

func (gEnts *GameEntities) Add(ent Entity) {
    select {
    case i := <-gEnts.openSpots:
        gEnts.Entities[i] = ent
        gEnts.emap[ent] = i
    default:
        gEnts.NumberOfEntities += 1
        gEnts.Entities[gEnts.NumberOfEntities] = ent
        gEnts.emap[ent] = gEnts.NumberOfEntities
    }
}

func (gEnts *GameEntities) Remove(ent Entity) {
    entIndex := gEnts.emap[ent]
    gEnts.Entities[entIndex] = nil
    delete(gEnts.emap, ent)
    gEnts.openSpots <- entIndex
}

func (gEnts *GameEntities) Close() {
	close(gEnts.openSpots)
}

func init() {
	entities := make([]Entity, 10000)
	emap := make(map[Entity]int)
	openSpots := make(chan int)

	AllGameEntities = GameEntities{entities, emap, 0, openSpots}
}
