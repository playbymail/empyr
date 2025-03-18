// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

// TurnReport_t is the report for a turn. It contains data on the game,
// the current turn, the empire, and all of the ships and colonies that
// are controlled by the empire
type TurnReport_t struct {
	Heading *ReportHeading_t

	Colonies []*ColonyReport_t // list of colonies sorted by ID
	Ships    []*ShipReport_t   // list of ships sorted by ID
	Surveys  []*SurveyReport_t // list of surveys sorted by ID

	CreatedDate     string // date the report was created
	CreatedDateTime string // date and time the report was created
}

type ReportHeading_t struct {
	Game       string // code of the game, eg "A02"
	TurnNo     int64  // turn number, eg 1
	TurnCode   string // display for the turn number, eg "T00001"
	EmpireNo   int64  // empire number, eg 1
	EmpireCode string // display for the empire, eg "E001"
}

type ColonyReport_t struct {
	Id          int64
	IdCode      string // display for the colony, eg "CC-1"
	Kind        string // kind of colony, eg "Open Colony"
	Name        string // name of the colony, eg "Colony 1"
	Coordinates string // display for the system, eg "02/13/28A"
	OrbitNo     int64

	VitalStatistics    *ColonyStatisticsReport_t
	Census             *ColonyCensusReport_t
	Other              *ColonyOtherReport_t
	Transports         *ColonyTransportReport_t
	Inventory          []*ColonyInventoryLine_t
	ProductionConsumed []*ProductionConsumedLine_t
	ProductionCreated  []*ProductionCreatedLine_t
	FarmGroups         []*ColonyFarmGroupsReport_t
	MiningGroups       []*ColonyMiningGroupsReport_t
	FactoryGroups      []*ColonyFactoryGroupsReport_t
	Spies              []*ColonySpyReport_t
}

type ColonyStatisticsReport_t struct {
	TechLevel        int64  // tech level of the colony, eg 1
	StandardOfLiving string // standard of living, eg "0.250" (empty if no standard of living)
	Rations          string // rations rate, eg "0.250" (empty if no rations)
	BirthRate        string // birth rate, eg "0.000" (empty if no birth rate)
	DeathRate        string // death rate, eg "0.000" (empty if no death rate)
	FoodConsumed     string // food consumed, eg "1,000,000"
}

type ColonyCensusReport_t struct {
	Population      []*PopulationReport_t
	TotalEmployed   string // total employed population, eg "1,000,000"
	TotalPopulation string // total population, eg "1,000,000"
	TotalPay        string // total pay, eg "1,000,000"
}

type ColonyOtherReport_t struct {
	TotalMass       string // total mass, eg "1,000,000"
	TotalVolume     string // total space capacity, eg "1,000,000"
	AvailableVolume string // space available, eg "1,000,000"
}

type ProductionConsumedLine_t struct {
	Category  string // Farming, Mining, etc
	Fuel      string // used, eg "1,000,000"
	Gold      string // used, eg "1,000,000"
	Metals    string // used, eg "1,000,000"
	NonMetals string // used, eg "1,000,000"
}

type ProductionCreatedLine_t struct {
	Category     string // Farming, Mining, etc
	Farmed       string // used, eg "1,000,000"
	Mined        string // used, eg "1,000,000"
	Manufactured string // used, eg "1,000,000"
}

type ColonyTransportReport_t struct {
	Capacity  string // transport capacity, eg "1,000,000"
	Used      string // used, eg "1,000,000"
	Available string // transport available, eg "1,000,000"
}

type ColonyInventoryLine_t struct {
	Code            string // inventory code, eg "FOOD"
	NonAssemblyQty  string // quantity, eg "1,000,000"
	DisassembledQty string // quantity, eg "1,000,000"
	AssembledQty    string // quantity, eg "1,000,000"
	IsAssembled     bool   // is assembled, eg "true"
	IsStored        bool   // is stored, eg "true"
	IsOPU           bool   // is an assembly/operational unit
}

type ColonyFarmGroupsReport_t struct {
	GroupNo    string // group number, eg "01"
	TechLevel  int64  // tech level, eg 1
	NbrOfUnits string // number of units, eg "10,000,000"
}

type ColonyMiningGroupsReport_t struct {
	GroupNo      string // group number, eg "01"
	DepositNo    string // deposit number, eg "01"
	DepositQty   string // deposit quantity, eg "1,000,000"
	DepositKind  string // deposit kind, eg "FUEL"
	DepositYield string // deposit yield, eg "1%"
	Units        []*MiningGroupUnitReport_t
}

type MiningGroupUnitReport_t struct {
	TechLevel  int64  // tech level, eg 1
	NbrOfUnits string // quantity, eg "1,000,000"
}

type ColonyFactoryGroupsReport_t struct {
	GroupNo    string // group number, eg "01"
	Orders     string // orders, eg "AUT-1" (or "AUT-1*" if retooling)
	RetoolTurn string // retool turn, eg "1"
	Units      []*ColonyFactoryGroupReport_t
}

type ColonyFactoryGroupReport_t struct {
	TechLevel  int64  // tech level, eg 1
	NbrOfUnits string // quantity, eg "1,000,000"
	Pipeline   [3]*ColonyFactoryPipelineReport_t
}

type ColonyFactoryPipelineReport_t struct {
	Percentage string // pipeline percentage complete, eg "25%", "50%", "75%"
	Unit       string // unit being produced, eg "AUT-1"
	Qty        string // quantity in this pipeline, eg "1,000,000"
}

type ColonySpyReport_t struct {
	Group   string // spy group, eg "A"
	Qty     string // quantity, eg "1,000,000"
	Results []string
}

type ShipReport_t struct {
	Id     int64
	IdCode string // display for the colony, eg "SS-1"
}

type PopulationReport_t struct {
	Group       string // population code, eg "USK"
	Population  string // population quantity, eg "1,000,000"
	PctTotalPop string // percentage of total population, eg "95%"
	Employed    string // employed population, eg "1,000,000"
	PayRate     string // pay rate, eg "0.0625" (empty if no pay rate)
	TotalPay    string // total pay, eg "1,000,000"
	qty         int64
}

type SurveyReport_t struct {
	ID int64
}
