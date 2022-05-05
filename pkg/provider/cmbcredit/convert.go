package cmbcredit

import (
	"math"

	"github.com/deb-sig/double-entry-generator/pkg/ir"
)

func (c *CmbCredit) convertToIR() *ir.IR {
	i := ir.New()
	for _, o := range c.Orders {
		irO := ir.Order{
			Peer:    o.CardNo,
			Item:    o.Description,
			PayTime: o.PostDate,
			Money:   math.Abs(o.Money),
			TxType:  convertType(o.Money),
		}
		irO.Metadata = getMetadata(o)
		i.Orders = append(i.Orders, irO)
	}
	return i
}

func convertType(money float64) ir.TxType {
	if money < 0 {
		return ir.TxTypeRecv
	}
	return ir.TxTypeSend
}

func getMetadata(o Order) map[string]string {
	return map[string]string{
		"source":  "招行信用卡",
		"country": o.Area,
	}
}
