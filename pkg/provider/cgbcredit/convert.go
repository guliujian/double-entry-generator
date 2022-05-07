package cgbcredit

import (
	"fmt"
	"math"

	"github.com/deb-sig/double-entry-generator/pkg/ir"
)

func (c *CGBCredit) convertToIR() *ir.IR {
	i := ir.New()
	for _, o := range c.Orders {
		irO := ir.Order{
			Peer:    o.CardNo,
			Item:    o.Description,
			PayTime: o.TransDate,
			Money:   math.Abs(o.Money),
			TxType:  convertType(o.Money),
		}
		irO.Metadata = getMetadata(o)
		i.Orders = append(i.Orders, irO)
	}
	return i
}

func getMetadata(o Order) map[string]string {
	data := map[string]string{
		"source": "广发信用卡",
	}
	if o.TransCurrency != "人民币" {
		data["transcurrency"] = o.TransCurrency
		data["moneyoriginal"] = fmt.Sprintf("%f", o.MoneyOriginal)
	}
	data["tranfertype"] = string(convertType(o.Money))
	return data
}

func convertType(money float64) ir.TxType {
	if money > 0 {
		return ir.TxTypeRecv
	}
	return ir.TxTypeSend
}
