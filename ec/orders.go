// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ec

import (
	"github.com/playbymail/empyr/models/coordinates"
	"github.com/playbymail/empyr/models/orders"
	"github.com/playbymail/empyr/models/units"
)

// Orders holds all of a player's orders for a single turn.
type Orders struct {
	Validated bool
	Handle    string
	Game      string
	Turn      int
	Secret    *Secret
	Orders    []orders.Order
	Error     error
}

type Abandon struct {
	Line     int
	Location coordinates.Coordinates // location to be abandoned
}

func (o *Abandon) Execute() error { panic("!") }

type AssembleFactoryGroup struct {
	Line        int
	Id          int        // id of unit being ordered
	Quantity    int        // number of units to assemble
	Unit        units.Unit // factory units to assemble
	Manufacture units.Unit // product unit to be manufactured
}

func (o *AssembleFactoryGroup) Execute() error { panic("!") }

type AssembleMineGroup struct {
	Line      int
	Id        int        // id of unit being ordered
	DepositId string     // deposit to assemble mines at
	Quantity  int        // number of units to assemble
	Unit      units.Unit // mine units to assemble
}

func (o *AssembleMineGroup) Execute() error { panic("!") }

type AssembleUnit struct {
	Line     int
	Id       int        // id of unit being ordered
	Quantity int        // number of units to assemble
	Unit     units.Unit // unit to assemble
}

func (o *AssembleUnit) Execute() error { panic("!") }

type Bombard struct {
	Line         int
	Id           int // id of unit being ordered
	PctCommitted int
	TargetId     int // id of unit being attacked
}

func (o *Bombard) Execute() error { panic("!") }

type Buy struct {
	Line     int
	Id       int        // id of unit being ordered
	Quantity int        // number of units to purchase
	Unit     units.Unit // unit to sell
	Bid      float64    // bid per unit
}

func (o *Buy) Execute() error { panic("!") }

type CheckRebels struct {
	Line     int
	Id       int // id of unit being ordered
	Quantity int // number of units to use
}

func (o *CheckRebels) Execute() error { panic("!") }

type Claim struct {
	Line     int
	Id       int                     // id of unit being ordered
	Location coordinates.Coordinates // location to be claimed
}

func (o *Claim) Execute() error { panic("!") }

type ConvertRebels struct {
	Line     int
	Id       int // id of unit being ordered
	Quantity int // number of units to use
}

func (o *ConvertRebels) Execute() error { panic("!") }

type CounterAgents struct {
	Line     int
	Id       int // id of unit being ordered
	Quantity int // number of units to use
}

func (o *CounterAgents) Execute() error { panic("!") }

type Discharge struct {
	Line       int
	Id         int    // id of unit being ordered
	Quantity   int    // number of units to use
	Profession string // profession to discharge from
}

func (o *Discharge) Execute() error { panic("!") }

type Draft struct {
	Line       int
	Id         int    // id of unit being ordered
	Quantity   int    // number of units to use
	Profession string // profession to draft into
}

func (o *Draft) Execute() error { panic("!") }

type ExpandFactoryGroup struct {
	Line         int
	Id           int        // id of unit being ordered
	FactoryGroup string     // factory group to expand
	Quantity     int        // number of units to assemble
	Unit         units.Unit // mine units to assemble
}

func (o *ExpandFactoryGroup) Execute() error { panic("!") }

type ExpandMineGroup struct {
	Line      int
	Id        int        // id of unit being ordered
	MineGroup string     // mine group to expand
	Quantity  int        // number of units to assemble
	Unit      units.Unit // mine units to assemble
}

func (o *ExpandMineGroup) Execute() error { panic("!") }

type Grant struct {
	Line     int
	Location coordinates.Coordinates // coordinates of system and orbit
	Kind     string                  // kind of grant
	TargetId int                     // nation to grant
}

func (o *Grant) Execute() error { panic("!") }

type InciteRebels struct {
	Line     int
	Id       int // id of unit being ordered
	Quantity int // number of units to use
	TargetId int // id of nation to target
}

func (o *InciteRebels) Execute() error { panic("!") }

type Invade struct {
	Line         int
	Id           int // id of unit being ordered
	PctCommitted int
	TargetId     int // id of unit being attacked
}

func (o *Invade) Execute() error { panic("!") }

type Jump struct {
	Line     int
	Id       int                     // id of unit being ordered
	Location coordinates.Coordinates // coordinates to move to
}

func (o *Jump) Execute() error { panic("!") }

type Move struct {
	Line  int
	Id    int // id of unit being ordered
	Orbit int // orbit to move to
}

func (o *Move) Execute() error { panic("!") }

type Name struct {
	Line     int
	Location coordinates.Coordinates // coordinates of system or planet to name
	Name     string                  // new name for unit
}

func (o *Name) Execute() error { panic("!") }

type NameUnit struct {
	Line int
	Id   int    // id of unit being ordered
	Name string // new name for unit
}

func (o *NameUnit) Execute() error { panic("!") }

type News struct {
	Line      int
	Location  coordinates.Coordinates // location to send news to
	Article   string
	Signature string
}

func (o *News) Execute() error { panic("!") }

type PayAll struct {
	Line       int
	Profession string  // profession to change pay for
	Rate       float64 // new pay rate
}

func (o *PayAll) Execute() error { panic("!") }

type PayLocal struct {
	Line       int
	Id         int     // id of unit being ordered
	Profession string  // profession to change pay for
	Rate       float64 // new pay rate
}

func (o *PayLocal) Execute() error { panic("!") }

type Probe struct {
	Line  int
	Id    int // id of unit being ordered
	Orbit int // orbit to probe
}

func (o *Probe) Execute() error { panic("!") }

type ProbeSystem struct {
	Line     int
	Id       int                     // id of unit being ordered
	Location coordinates.Coordinates // location to probe
}

func (o *ProbeSystem) Execute() error { panic("!") }

type Raid struct {
	Line         int
	Id           int // id of unit being ordered
	PctCommitted int
	TargetId     int        // id of unit being raided
	TargetUnit   units.Unit // material to raid
}

func (o *Raid) Execute() error { panic("!") }

type RationAll struct {
	Line int
	Rate int // new ration percentage
}

func (o *RationAll) Execute() error { panic("!") }

type RationLocal struct {
	Line int
	Id   int // id of unit being ordered
	Rate int // new ration percentage
}

func (o *RationLocal) Execute() error { panic("!") }

type RecycleFactoryGroup struct {
	Line         int
	Id           int        // id of unit being ordered
	FactoryGroup string     // factory group to recycle units from
	Quantity     int        // number of units to recycle
	Unit         units.Unit // unit to recycle
}

func (o *RecycleFactoryGroup) Execute() error { panic("!") }

type RecycleMineGroup struct {
	Line      int
	Id        int        // id of unit being ordered
	MineGroup string     // mine group to recycle units from
	Quantity  int        // number of units to recycle
	Unit      units.Unit // unit to recycle
}

func (o *RecycleMineGroup) Execute() error { panic("!") }

type RecycleUnit struct {
	Line     int
	Id       int        // id of unit being ordered
	Quantity int        // number of units to recycle
	Unit     units.Unit // unit to recycle
}

func (o *RecycleUnit) Execute() error { panic("!") }

type RetoolFactoryGroup struct {
	Line         int
	Id           int        // id of unit being ordered
	FactoryGroup string     // factory group to retool
	Unit         units.Unit // new unit to manufacture
}

func (o *RetoolFactoryGroup) Execute() error { panic("!") }

type Revoke struct {
	Line     int
	Location coordinates.Coordinates // coordinates of system and orbit
	Kind     string                  // kind of grant
	TargetId int                     // nation to grant
}

func (o *Revoke) Execute() error { panic("!") }

type ScrapFactoryGroup struct {
	Line         int
	Id           int        // id of unit being ordered
	FactoryGroup string     // factory group to scrap units from
	Quantity     int        // number of units to scrap
	Unit         units.Unit // unit to scrap
}

func (o *ScrapFactoryGroup) Execute() error { panic("!") }

type ScrapMineGroup struct {
	Line      int
	Id        int        // id of unit being ordered
	MineGroup string     // mine group to scrap units from
	Quantity  int        // number of units to scrap
	Unit      units.Unit // unit to scrap
}

func (o *ScrapMineGroup) Execute() error { panic("!") }

type ScrapUnit struct {
	Line     int
	Id       int        // id of unit being ordered
	Quantity int        // number of units to scrap
	Unit     units.Unit // unit to scrap
}

func (o *ScrapUnit) Execute() error { panic("!") }

type Secret struct {
	Line   int
	Handle string
	Game   string
	Turn   int
	Token  string
}

func (o *Secret) Execute() error { panic("!") }

type Sell struct {
	Line     int
	Id       int        // id of unit being ordered
	Quantity int        // number of units to sell
	Unit     units.Unit // unit to sell
	Ask      float64    // ask per unit
}

func (o *Sell) Execute() error { panic("!") }

type Setup struct {
	Line     int
	Id       int                     // id of unit establishing ship or colony
	Location coordinates.Coordinates // location being set up
	Kind     string                  // must be 'colony' or 'ship'
	Action   string                  // must be 'transfer'
	Items    []*orders.TransferDetail
}

func (o *Setup) Execute() error { panic("!") }

type StealSecrets struct {
	Line     int
	Id       int // id of unit being ordered
	Quantity int // number of units to use
	TargetId int // id of nation to target
}

func (o *StealSecrets) Execute() error { panic("!") }

type StoreFactoryGroup struct {
	Line         int
	Id           int        // id of unit being ordered
	FactoryGroup string     // factory group to store units from
	Quantity     int        // number of units to store
	Unit         units.Unit // unit to store
}

func (o *StoreFactoryGroup) Execute() error { panic("!") }

type StoreMineGroup struct {
	Line      int
	Id        int        // id of unit being ordered
	MineGroup string     // mine group to store units from
	Quantity  int        // number of units to store
	Unit      units.Unit // unit to store
}

func (o *StoreMineGroup) Execute() error { panic("!") }

type StoreUnit struct {
	Line     int
	Id       int        // id of unit being ordered
	Quantity int        // number of units to store
	Unit     units.Unit // unit to store
}

func (o *StoreUnit) Execute() error { panic("!") }

type SupportAttack struct {
	Line         int
	Id           int // id of unit being ordered
	PctCommitted int
	SupportId    int // id of unit being supported
	TargetId     int // id of unit being attacked
}

func (o *SupportAttack) Execute() error { panic("!") }

type SupportDefend struct {
	Line         int
	Id           int // id of unit being ordered
	SupportId    int // id of unit being supported
	PctCommitted int
}

func (o *SupportDefend) Execute() error { panic("!") }

type SuppressAgents struct {
	Line     int
	Id       int // id of unit being ordered
	Quantity int // number of units to use
	TargetId int // id of nation to target
}

func (o *SuppressAgents) Execute() error { panic("!") }

type Survey struct {
	Line  int
	Id    int // id of unit being ordered
	Orbit int // orbit to survey
}

func (o *Survey) Execute() error { panic("!") }

type SurveySystem struct {
	Line     int
	Id       int                     // id of unit being ordered
	Location coordinates.Coordinates // location to survey
}

func (o *SurveySystem) Execute() error { panic("!") }

type Transfer struct {
	Line     int
	Id       int        // id of unit being ordered
	Quantity int        // number of units to transfer
	Unit     units.Unit // unit to transfer
	TargetId int        // id of unit receiving units
}

func (o *Transfer) Execute() error { panic("!") }

// The Unknown order type captures unrecognized orders.
type Unknown struct {
	Line    int
	Command string
}

func (o *Unknown) Execute() error { panic("!") }
