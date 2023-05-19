// empyr - a reimagining of Vern Holford's Empyrean Challenge
// Copyright (C) 2023 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package orders

import (
	"encoding/json"
	"io"
	"sort"
)

// Decode extracts from request body.
func Decode(r io.ReadCloser) (orders []*Order, errors []error) {
	dec := json.NewDecoder(r)

	// reject any request with unknown verbs.
	dec.DisallowUnknownFields()

	if err := dec.Decode(&orders); err != nil {
		return nil, append(errors, err)
	}

	// loop through all of the orders and assign the sort priority
	// based on the order type. yeah, i think that this is gross.
	for _, order := range orders {
		switch {
		case order.Debug != nil:
			order.priority = 0

		// Combat Orders Stage - Prefire Segment
		case order.Dodge != nil:
			order.priority = 10001
		case order.Accept != nil:
			order.priority = 10002
		case order.AutoReturnFire != nil:
			order.priority = 10003
		case order.CloseProximityTargeting != nil:
			order.priority = 10004
		// Combat Orders Stage - Pre-Maneuver Fire Segment
		case order.PreManeuverEnergyWeaponFire != nil:
			order.priority = 10101
		case order.PreManeuverMissileFire != nil:
			order.priority = 10102
		// Combat Orders Stage - Allocate Damage
		// Combat Orders Stage - Maneuver Segment
		case order.Undock != nil:
			order.priority = 10301
		case order.Run != nil:
			order.priority = 10302
		case order.TacticalManeuver != nil:
			order.priority = 10303
		case order.Close != nil:
			order.priority = 10304
		case order.Dock != nil:
			order.priority = 10305
		// Combat Orders Stage - Allocate Damage
		// Combat Orders Stage - Post-Maneuver Fire Segment
		case order.AfterManeuverEnergyWeaponFire != nil:
			order.priority = 10501
		case order.AfterManeuverMissileFire != nil:
			order.priority = 10501
		// Combat Orders Stage - Allocate Damage
		// Combat Orders Stage - Ground Combat Segment
		case order.Withdraw != nil:
			order.priority = 10701
		case order.DefensiveSupport != nil:
			order.priority = 10702
		case order.Invade != nil:
			order.priority = 10703
		case order.OffensiveSupport != nil:
			order.priority = 10704
		// Combat Orders Stage - Cycle Ground Combat
		// Permissions Orders Stage
		case order.PermissionToColonize != nil:
			order.priority = 11001
		case order.HomePortChange != nil:
			order.priority = 11002
		// Permissions Orders Stage - Diplomacy
		// Disassembly Stage
		case order.Disassemble != nil:
			order.priority = 12001
		case order.Scrap != nil:
			order.priority = 12002
		case order.Junk != nil:
			order.priority = 12003
		case order.CombineFactoryGroup != nil:
			order.priority = 12004
		// Setup Stage
		case order.DefineCargoHold != nil:
			order.priority = 13001
		case order.SetUp != nil:
			order.priority = 13002
		case order.AddOn != nil:
			order.priority = 13003
		// Transfer Stage
		case order.UnloadCargo != nil:
			order.priority = 14001
		case order.Transfer != nil:
			order.priority = 14002
		case order.PickUpItem != nil:
			order.priority = 14003
		case order.PickUpPopulation != nil:
			order.priority = 14003
		case order.LoadCargo != nil:
			order.priority = 14005
		// Draft Orders Stage
		case order.Draft != nil:
			order.priority = 15001
		case order.Disband != nil:
			order.priority = 15002
		// Assembly Stage - Order Processing Segment
		case order.AssembleFactory != nil:
			order.priority = 16101
		case order.AssembleFactoryGroup != nil:
			order.priority = 16101
		case order.AssembleItem != nil:
			order.priority = 16101
		case order.AssembleMine != nil:
			order.priority = 16101
		case order.AssembleMineGroup != nil:
			order.priority = 16101
		case order.ExpendResearchPointsOnly != nil:
			order.priority = 16110
		case order.ExpendPrototype != nil:
			order.priority = 16111
		case order.FactoryGroupChange != nil:
			order.priority = 16112
		case order.BuildChange != nil:
			order.priority = 16113
		case order.MineChange != nil:
			order.priority = 16114
		case order.ShutDown != nil:
			order.priority = 16115
		case order.StartUp != nil:
			order.priority = 16116
		// Assembly Stage - Non Prototype TL Increases Segment
		case order.ExpendCommittedBufferResearchPoints != nil:
			order.priority = 16201
		// Build Change Stage (merged with Assembly Stage)
		// Surveys and Probes Stage
		case order.Probe != nil:
			order.priority = 17001
		case order.Survey != nil:
			order.priority = 17002
		case order.LaunchRobotProbe != nil:
			order.priority = 17003
		// Pay Change Stage
		case order.Pay != nil:
			order.priority = 18001
		case order.Ration != nil:
			order.priority = 18002
		// Naming Orders Stage
		case order.Name != nil:
			order.priority = 19001
		case order.Note != nil:
			order.priority = 19002
		case order.ControlPlanet != nil:
			order.priority = 19003
		case order.UncontrolPlanet != nil:
			order.priority = 19004
		case order.Message != nil:
			order.priority = 19005
		// Ship Travel Stage
		case order.Jump != nil:
			order.priority = 20001
		case order.Move != nil:
			order.priority = 20002
		// Probe Stage
		case order.ProbeOrbit != nil:
			order.priority = 21001
		case order.ProbeSystem != nil:
			order.priority = 21002
		// Give Stage
		case order.Give != nil:
			order.priority = 22001
		// Production Stage
		// Produce Out Put Stage
		// Send Out Put Stage

		// Unknown Stage Priority
		case order.Merge != nil:
			order.priority = 99999
		case order.MineShutDown != nil:
			order.priority = 99999
		case order.MineStartUp != nil:
			order.priority = 99999

		default:
			panic("assert(order.Type != unknown)")
		}
	}

	// must do a stable sort because some order types are sensitive
	// to the order of individual lines
	sort.SliceStable(orders, func(i, j int) bool {
		return orders[i].priority < orders[j].priority
	})

	return orders, errors
}
