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

package games

import (
	"github.com/playbymail/empyr/cmd/server/pkg/orders"
	"log"
)

func (g *Game) Apply(turn []*orders.Order) error {
	var debug bool
	var o interface{}
	for _, order := range turn {
		switch {
		case order.Accept != nil:
			o = order.Accept
		case order.AddOn != nil:
			o = order.AddOn
		case order.AfterManeuverEnergyWeaponFire != nil:
			o = order.AfterManeuverEnergyWeaponFire
		case order.AfterManeuverMissileFire != nil:
			o = order.AfterManeuverMissileFire
		case order.AssembleFactory != nil:
			o = order.AssembleFactory
		case order.AssembleFactoryGroup != nil:
			o = order.AssembleFactoryGroup
		case order.AssembleItem != nil:
			o = order.AssembleItem
		case order.AssembleMine != nil:
			o = order.AssembleMine
		case order.AssembleMineGroup != nil:
			o = order.AssembleMineGroup
		case order.AutoReturnFire != nil:
			o = order.AutoReturnFire
		case order.BuildChange != nil:
			o = order.BuildChange
		case order.Close != nil:
			o = order.Close
		case order.CloseProximityTargeting != nil:
			o = order.CloseProximityTargeting
		case order.CombineFactoryGroup != nil:
			o = order.CombineFactoryGroup
		case order.ControlPlanet != nil:
			o = order.ControlPlanet
		case order.Debug != nil:
			// Debug toggles the debug flag for the game
			debug = order.Debug.On
		case order.DefensiveSupport != nil:
			o = order.DefensiveSupport
		case order.DefineCargoHold != nil:
			o = order.DefineCargoHold
		case order.Disassemble != nil:
			o = order.Disassemble
		case order.Disband != nil:
			o = order.Disband
		case order.Dock != nil:
			o = order.Dock
		case order.Dodge != nil:
			// Target must allocate percentage of Speed to avoid hostile weapons fire.
			// Setting the percentage to zero disables dodging.
			ship, ok := g.Ships.id[order.Dodge.ShipID]
			if !ok {
				continue // ignore errors with orders
			}
			ship.Dodge = order.Dodge.Percentage / 100
			log.Printf("[orders] dodge %q %6.2f%%\n", ship.Name, ship.Dodge)
		case order.Draft != nil:
			o = order.Draft
		case order.ExpendCommittedBufferResearchPoints != nil:
			o = order.ExpendCommittedBufferResearchPoints
		case order.ExpendPrototype != nil:
			o = order.ExpendPrototype
		case order.ExpendResearchPointsOnly != nil:
			o = order.ExpendResearchPointsOnly
		case order.FactoryGroupChange != nil:
			o = order.FactoryGroupChange
		case order.Give != nil:
			o = order.Give
		case order.HomePortChange != nil:
			o = order.HomePortChange
		case order.Invade != nil:
			o = order.Invade
		case order.Jump != nil:
			o = order.Jump
		case order.Junk != nil:
			o = order.Junk
		case order.LaunchRobotProbe != nil:
			o = order.LaunchRobotProbe
		case order.LoadCargo != nil:
			o = order.LoadCargo
		case order.Merge != nil:
			o = order.Merge
		case order.Message != nil:
			o = order.Message
		case order.MineChange != nil:
			o = order.MineChange
		case order.MineShutDown != nil:
			o = order.MineShutDown
		case order.MineStartUp != nil:
			o = order.MineStartUp
		case order.Move != nil:
			o = order.Move
		case order.Name != nil:
			o = order.Name
		case order.Note != nil:
			o = order.Note
		case order.OffensiveSupport != nil:
			o = order.OffensiveSupport
		case order.Pay != nil:
			o = order.Pay
		case order.PermissionToColonize != nil:
			o = order.PermissionToColonize
		case order.PickUpItem != nil:
			o = order.PickUpItem
		case order.PickUpPopulation != nil:
			o = order.PickUpPopulation
		case order.PreManeuverEnergyWeaponFire != nil:
			o = order.PreManeuverEnergyWeaponFire
		case order.PreManeuverMissileFire != nil:
			o = order.PreManeuverMissileFire
		case order.Probe != nil:
			o = order.Probe
		case order.ProbeOrbit != nil:
			o = order.ProbeOrbit
		case order.ProbeSystem != nil:
			o = order.ProbeSystem
		case order.Ration != nil:
			o = order.Ration
		case order.Run != nil:
			o = order.Run
		case order.Scrap != nil:
			o = order.Scrap
		case order.SetUp != nil:
			o = order.SetUp
		case order.ShutDown != nil:
			o = order.ShutDown
		case order.StartUp != nil:
			o = order.StartUp
		case order.Survey != nil:
			o = order.Survey
		case order.TacticalManeuver != nil:
			o = order.TacticalManeuver
		case order.Transfer != nil:
			o = order.Transfer
		case order.UncontrolPlanet != nil:
			o = order.UncontrolPlanet
		case order.Undock != nil:
			o = order.Undock
		case order.UnloadCargo != nil:
			o = order.UnloadCargo
		case order.Withdraw != nil:
			o = order.Withdraw
		default:
			panic("!implemented")
		}
	}
	if debug {
		log.Printf("[debug] %v\n", o)
	}

	return nil
}
