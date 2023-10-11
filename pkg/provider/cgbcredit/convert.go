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
			Type:    convertType(o.Money),
		}
		if o.MoneyCurrency != "人民币" {
			irO.Currency = convertCurrencyName(o.MoneyCurrency)
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
	if o.TransCurrency != o.MoneyCurrency {
		data["transcurrency"] = o.TransCurrency
		data["moneyoriginal"] = fmt.Sprintf("%f", o.MoneyOriginal)
	}
	data["tranfertype"] = string(convertType(o.Money))
	return data
}

func convertType(money float64) ir.Type {
	if money > 0 {
		return ir.TypeRecv
	}
	return ir.TypeSend
}

// convert currency to CNY
var currencyMapping = map[string]string{
	"人民币": "CNY",
	"美元":  "USD",
	"欧元":  "EUR",
	"英镑":  "GBP",
	"日元":  "JPY",
	"韩元":  "KRW",
	"澳元":  "AUD",
	"加元":  "CAD",
	"港币":  "HKD",
}

func convertCurrencyName(c string) string {
	// 将中文货币名转换为英文大写
	englishCurrency := currencyMapping[c]

	return englishCurrency
}
