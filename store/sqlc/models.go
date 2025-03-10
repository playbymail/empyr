// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"time"
)

type Deposit struct {
	DepositID int64
	Quantity  int64
	YieldPct  int64
	EffTurn   int64
	EndTurn   int64
	Active    int64
}

type Deposits struct {
	ID           int64
	PlanetID     int64
	DepositNo    int64
	Kind         int64
	InitialQty   int64
	RemainingQty int64
	YieldPct     int64
}

type Empires struct {
	ID           int64
	GameID       int64
	EmpireNo     int64
	Name         string
	HomeSystemID int64
	HomeStarID   int64
	HomeOrbitID  int64
	HomePlanetID int64
	Handle       string
}

type Games struct {
	ID           int64
	Code         string
	Name         string
	DisplayName  string
	CurrentTurn  int64
	LastEmpireNo int64
	HomeSystemID int64
	HomeStarID   int64
	HomeOrbitID  int64
	HomePlanetID int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Orbits struct {
	ID       int64
	StarID   int64
	OrbitNo  int64
	Kind     int64
	Scarcity int64
}

type Planets struct {
	ID           int64
	OrbitID      int64
	Kind         int64
	Habitability int64
	Scarcity     int64
}

type SorcDetails struct {
	ID          int64
	SorcID      int64
	TurnNo      int64
	TechLevel   int64
	Name        string
	UemQty      int64
	UemPay      float64
	UskQty      int64
	UskPay      float64
	ProQty      int64
	ProPay      float64
	SldQty      int64
	SldPay      float64
	CnwQty      int64
	SpyQty      int64
	Rations     float64
	BirthRate   float64
	DeathRate   float64
	Sol         float64
	OrbitID     int64
	IsOnSurface int64
}

type SorcInfrastructure struct {
	SorcDetailID int64
	Kind         string
	TechLevel    int64
	Qty          int64
}

type SorcInventory struct {
	SorcDetailID int64
	Kind         string
	TechLevel    int64
	QtyAssembled int64
	QtyStored    int64
}

type SorcPopulation struct {
	SorcDetailID int64
	Kind         string
	Qty          int64
}

type SorcSuperstructure struct {
	SorcDetailID int64
	Kind         string
	TechLevel    int64
	Qty          int64
}

type Sorcs struct {
	ID       int64
	EmpireID int64
	Kind     int64
}

type Stars struct {
	ID       int64
	SystemID int64
	Sequence string
	Scarcity int64
}

type SystemDistances struct {
	FromSystemID int64
	ToSystemID   int64
	Distance     int64
}

type Systems struct {
	ID       int64
	GameID   int64
	X        int64
	Y        int64
	Z        int64
	Scarcity int64
}

type Units struct {
	Code          string
	Mass          int64
	IsOperational int64
}
