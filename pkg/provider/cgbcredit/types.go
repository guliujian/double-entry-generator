package cgbcredit

import "time"

type Statistics struct {
	UserID          string    `json:"user_id,omitempty"`
	Username        string    `json:"username,omitempty"`
	ParsedItems     int       `json:"parsed_items,omitempty"`
	Start           time.Time `json:"start,omitempty"`
	End             time.Time `json:"end,omitempty"`
	TotalInRecords  int       `json:"total_in_records,omitempty"`
	TotalInMoney    float64   `json:"total_in_money,omitempty"`
	TotalOutRecords int       `json:"total_out_records,omitempty"`
	TotalOutMoney   float64   `json:"total_out_money,omitempty"`
}

type Order struct {
	TransDate     time.Time `json:"transDate,omitempty"`     //交易日
	Description   string    `json:"description,omitempty"`   //交易摘要
	Money         float64   `json:"money,omitempty"`         // 入账金额
	MoneyCurrency string    `json:"moneyCurrency,omitempty"` //入账币种
	CardNo        string    `json:"cardNo,omitempty"`        //卡号后四位
	MoneyOriginal float64   `json:"moneyOriginal,omitempty"` //交易金额
	TransCurrency string    `json:"curreny,omitempty"`       //交易币种
}
