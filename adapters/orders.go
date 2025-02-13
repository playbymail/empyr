// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package adapters

import (
	"fmt"
	"github.com/playbymail/empyr/ec"
	"github.com/playbymail/empyr/models/coordinates"
	mo "github.com/playbymail/empyr/models/orders"
	"github.com/playbymail/empyr/models/units"
	po "github.com/playbymail/empyr/parsers/orders"
)

func CoordToEngineCoord(in po.Coordinates) coordinates.Coordinates {
	return coordinates.Coordinates{
		X:      in.X,
		Y:      in.Y,
		Z:      in.Z,
		System: in.System,
		Orbit:  in.Orbit,
	}
}

func ItemsToEngineItems(in []*po.TransferDetail) (out []*mo.TransferDetail) {
	for _, item := range in {
		out = append(out, &mo.TransferDetail{
			Unit:     UnitToEngineUnit(item.Unit),
			Quantity: item.Quantity,
		})
	}
	return out
}

func UnitToEngineUnit(in po.Unit) units.Unit {
	return units.Unit{
		Name:      in.Name,
		TechLevel: in.TechLevel,
	}
}

func OrdersToEngineOrders(in []any) (out []mo.Order) {
	for _, o := range in {
		switch order := o.(type) {
		case *po.Abandon:
			out = append(out, &ec.Abandon{
				Line:     order.Line,
				Location: CoordToEngineCoord(order.Location),
			})
		case *po.AssembleFactoryGroup:
			out = append(out, &ec.AssembleFactoryGroup{
				Line:        order.Line,
				Id:          order.Id,
				Quantity:    order.Quantity,
				Unit:        UnitToEngineUnit(order.Unit),
				Manufacture: UnitToEngineUnit(order.Manufacture),
			})
		case *po.AssembleMineGroup:
			out = append(out, &ec.AssembleMineGroup{
				Line:      order.Line,
				Id:        order.Id,
				DepositId: order.DepositId,
				Quantity:  order.Quantity,
				Unit:      UnitToEngineUnit(order.Unit),
			})
		case *po.AssembleUnit:
			out = append(out, &ec.AssembleUnit{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
				Unit:     UnitToEngineUnit(order.Unit),
			})
		case *po.Bombard:
			out = append(out, &ec.Bombard{
				Line:         order.Line,
				Id:           order.Id,
				PctCommitted: order.PctCommitted,
				TargetId:     order.TargetId,
			})
		case *po.Buy:
			out = append(out, &ec.Buy{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
				Unit:     UnitToEngineUnit(order.Unit),
				Bid:      order.Bid,
			})
		case *po.CheckRebels:
			out = append(out, &ec.CheckRebels{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
			})
		case *po.Claim:
			out = append(out, &ec.Claim{
				Line:     order.Line,
				Id:       order.Id,
				Location: CoordToEngineCoord(order.Location),
			})
		case *po.ConvertRebels:
			out = append(out, &ec.ConvertRebels{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
			})
		case *po.CounterAgents:
			out = append(out, &ec.CounterAgents{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
			})
		case *po.Discharge:
			out = append(out, &ec.Discharge{
				Line:       order.Line,
				Id:         order.Id,
				Quantity:   order.Quantity,
				Profession: order.Profession,
			})
		case *po.Draft:
			out = append(out, &ec.Draft{
				Line:       order.Line,
				Id:         order.Id,
				Quantity:   order.Quantity,
				Profession: order.Profession,
			})
		case *po.ExpandFactoryGroup:
			out = append(out, &ec.ExpandFactoryGroup{
				Line:         order.Line,
				Id:           order.Id,
				FactoryGroup: order.FactoryGroup,
				Quantity:     order.Quantity,
				Unit:         UnitToEngineUnit(order.Unit),
			})
		case *po.ExpandMineGroup:
			out = append(out, &ec.ExpandMineGroup{
				Line:      order.Line,
				Id:        order.Id,
				MineGroup: order.MineGroup,
				Quantity:  order.Quantity,
				Unit:      UnitToEngineUnit(order.Unit),
			})
		case *po.Grant:
			out = append(out, &ec.Grant{
				Line:     order.Line,
				Location: CoordToEngineCoord(order.Location),
				Kind:     order.Kind,
				TargetId: order.TargetId,
			})
		case *po.InciteRebels:
			out = append(out, &ec.InciteRebels{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
				TargetId: order.TargetId,
			})
		case *po.Invade:
			out = append(out, &ec.Invade{
				Line:         order.Line,
				Id:           order.Id,
				PctCommitted: order.PctCommitted,
				TargetId:     order.TargetId,
			})
		case *po.Jump:
			out = append(out, &ec.Jump{
				Line:     order.Line,
				Id:       order.Id,
				Location: CoordToEngineCoord(order.Location),
			})
		case *po.Move:
			out = append(out, &ec.Move{
				Line:  order.Line,
				Id:    order.Id,
				Orbit: order.Orbit,
			})
		case *po.Name:
			out = append(out, &ec.Name{
				Line:     order.Line,
				Location: CoordToEngineCoord(order.Location),
				Name:     order.Name,
			})
		case *po.NameUnit:
			out = append(out, &ec.NameUnit{
				Line: order.Line,
				Id:   order.Id,
				Name: order.Name,
			})
		case *po.News:
			out = append(out, &ec.News{
				Line:      order.Line,
				Location:  CoordToEngineCoord(order.Location),
				Article:   order.Article,
				Signature: order.Signature,
			})
		case *po.PayAll:
			out = append(out, &ec.PayAll{
				Line:       order.Line,
				Profession: order.Profession,
				Rate:       order.Rate,
			})
		case *po.PayLocal:
			out = append(out, &ec.PayLocal{
				Line:       order.Line,
				Id:         order.Id,
				Profession: order.Profession,
				Rate:       order.Rate,
			})
		case *po.Probe:
			out = append(out, &ec.Probe{
				Line:  order.Line,
				Id:    order.Id,
				Orbit: order.Orbit,
			})
		case *po.ProbeSystem:
			out = append(out, &ec.ProbeSystem{
				Line:     order.Line,
				Id:       order.Id,
				Location: CoordToEngineCoord(order.Location),
			})
		case *po.Raid:
			out = append(out, &ec.Raid{
				Line:         order.Line,
				Id:           order.Id,
				PctCommitted: order.PctCommitted,
				TargetId:     order.TargetId,
				TargetUnit:   UnitToEngineUnit(order.TargetUnit),
			})
		case *po.RationAll:
			out = append(out, &ec.RationAll{
				Line: order.Line,
				Rate: order.Rate,
			})
		case *po.RationLocal:
			out = append(out, &ec.RationLocal{
				Line: order.Line,
				Id:   order.Id,
				Rate: order.Rate,
			})
		case *po.RecycleFactoryGroup:
			out = append(out, &ec.RecycleFactoryGroup{
				Line:         order.Line,
				Id:           order.Id,
				FactoryGroup: order.FactoryGroup,
				Quantity:     order.Quantity,
				Unit:         UnitToEngineUnit(order.Unit),
			})
		case *po.RecycleMineGroup:
			out = append(out, &ec.RecycleMineGroup{
				Line:      order.Line,
				Id:        order.Id,
				MineGroup: order.MineGroup,
				Quantity:  order.Quantity,
				Unit:      UnitToEngineUnit(order.Unit),
			})
		case *po.RecycleUnit:
			out = append(out, &ec.RecycleUnit{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
				Unit:     UnitToEngineUnit(order.Unit),
			})
		case *po.RetoolFactoryGroup:
			out = append(out, &ec.RetoolFactoryGroup{
				Line:         order.Line,
				Id:           order.Id,
				FactoryGroup: order.FactoryGroup,
				Unit:         UnitToEngineUnit(order.Unit),
			})
		case *po.Revoke:
			out = append(out, &ec.Revoke{
				Line:     order.Line,
				Location: CoordToEngineCoord(order.Location),
				Kind:     order.Kind,
				TargetId: order.TargetId,
			})
		case *po.ScrapFactoryGroup:
			out = append(out, &ec.ScrapFactoryGroup{
				Line:         order.Line,
				Id:           order.Id,
				FactoryGroup: order.FactoryGroup,
				Quantity:     order.Quantity,
				Unit:         UnitToEngineUnit(order.Unit),
			})
		case *po.ScrapMineGroup:
			out = append(out, &ec.ScrapMineGroup{
				Line:      order.Line,
				Id:        order.Id,
				MineGroup: order.MineGroup,
				Quantity:  order.Quantity,
				Unit:      UnitToEngineUnit(order.Unit),
			})
		case *po.ScrapUnit:
			out = append(out, &ec.ScrapUnit{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
				Unit:     UnitToEngineUnit(order.Unit),
			})
		case *po.Secret:
			out = append(out, &ec.Secret{
				Line:   order.Line,
				Handle: order.Handle,
				Game:   order.Game,
				Turn:   order.Turn,
				Token:  order.Token,
			})
		case *po.Sell:
			out = append(out, &ec.Sell{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
				Unit:     UnitToEngineUnit(order.Unit),
				Ask:      order.Ask,
			})
		case *po.Setup:
			out = append(out, &ec.Setup{
				Line:     order.Line,
				Id:       order.Id,
				Location: CoordToEngineCoord(order.Location),
				Kind:     order.Kind,
				Action:   order.Action,
				Items:    ItemsToEngineItems(order.Items),
			})
		case *po.StealSecrets:
			out = append(out, &ec.StealSecrets{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
				TargetId: order.TargetId,
			})
		case *po.StoreFactoryGroup:
			out = append(out, &ec.StoreFactoryGroup{
				Line:         order.Line,
				Id:           order.Id,
				FactoryGroup: order.FactoryGroup,
				Quantity:     order.Quantity,
				Unit:         UnitToEngineUnit(order.Unit),
			})
		case *po.StoreMineGroup:
			out = append(out, &ec.StoreMineGroup{
				Line:      order.Line,
				Id:        order.Id,
				MineGroup: order.MineGroup,
				Quantity:  order.Quantity,
				Unit:      UnitToEngineUnit(order.Unit),
			})
		case *po.StoreUnit:
			out = append(out, &ec.StoreUnit{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
				Unit:     UnitToEngineUnit(order.Unit),
			})
		case *po.SupportAttack:
			out = append(out, &ec.SupportAttack{
				Line:         order.Line,
				Id:           order.Id,
				PctCommitted: order.PctCommitted,
				SupportId:    order.SupportId,
				TargetId:     order.TargetId,
			})
		case *po.SupportDefend:
			out = append(out, &ec.SupportDefend{
				Line:         order.Line,
				Id:           order.Id,
				SupportId:    order.SupportId,
				PctCommitted: order.PctCommitted,
			})
		case *po.SuppressAgents:
			out = append(out, &ec.SuppressAgents{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
				TargetId: order.TargetId,
			})
		case *po.Survey:
			out = append(out, &ec.Survey{
				Line:  order.Line,
				Id:    order.Id,
				Orbit: order.Orbit,
			})
		case *po.SurveySystem:
			out = append(out, &ec.SurveySystem{
				Line:     order.Line,
				Id:       order.Id,
				Location: CoordToEngineCoord(order.Location),
			})
		case *po.Transfer:
			out = append(out, &ec.Transfer{
				Line:     order.Line,
				Id:       order.Id,
				Quantity: order.Quantity,
				Unit:     UnitToEngineUnit(order.Unit),
				TargetId: order.TargetId,
			})
		case *po.Unknown:
			// ignore unknown orders
		default:
			panic(fmt.Sprintf("unknown type %T", o))
		}
	}
	return out
}
