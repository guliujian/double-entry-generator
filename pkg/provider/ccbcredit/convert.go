package ccbcredit

import (
	"fmt"
	"math"
	"time"

	"github.com/deb-sig/double-entry-generator/pkg/ir"
)

func (c *CcbCredit) convertToIR() *ir.IR {
	i := ir.New()
	for _, o := range c.Orders {
		irO := ir.Order{
			Peer:    o.CardNo,
			Item:    o.Description,
			PayTime: getDate(o),
			Money:   math.Abs(o.Money),
			Type:    convertType(o.Money),
		}
		irO.Metadata = getMetadata(o)
		i.Orders = append(i.Orders, irO)
	}
	return i
}

func convertType(money float64) ir.Type {
	if money < 0 {
		return ir.TypeRecv
	}
	return ir.TypeSend
}

func getMetadata(o Order) map[string]string {
	data := map[string]string{
		"source": "建行信用卡",
	}
	if o.TransCurrency != "人民币" {
		data["transcurrency"] = o.TransCurrency
		data["moneyoriginal"] = fmt.Sprintf("%f", o.MoneyOriginal)
	}
	data["tranfertype"] = string(convertType(o.Money))
	return data
}

func getDate(o Order) time.Time {
	if o.TransDate.IsZero() {
		return o.PostDate
	}
	return o.TransDate
}
