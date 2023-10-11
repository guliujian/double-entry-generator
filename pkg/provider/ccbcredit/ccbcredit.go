package ccbcredit

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/PuerkitoBio/goquery"
	"github.com/deb-sig/double-entry-generator/pkg/ir"
)

type CcbCredit struct {
	Statistics Statistics `json:"statistics,omitempty"`
	LineNum    int        `json:"line_num,omitempty"`
	Orders     []Order    `json:"orders,omitempty"`
}

func New() *CcbCredit {
	return &CcbCredit{
		Statistics: Statistics{},
		LineNum:    0,
		Orders:     make([]Order, 0),
	}
}

func (c *CcbCredit) Translate(filename string) (*ir.IR, error) {
	log.SetPrefix("[Provider-CcbCredit]")
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	email, err := parsemail.Parse(bufio.NewReader(file))
	if err != nil {
		return nil, err
	}
	var doc goquery.Document
	decoded := base64.NewDecoder(base64.StdEncoding, strings.NewReader(email.HTMLBody))
	b, err := ioutil.ReadAll(decoded)
	if err != nil {
		if _, ok := err.(base64.CorruptInputError); ok {
			docs, err := goquery.NewDocumentFromReader(strings.NewReader(email.HTMLBody))
			if err != nil {
				return nil, err
			}
			doc = *docs
		} else {
			return nil, err
		}
	} else {
		docs, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		doc = *docs
	}
	var begin int
	doc.Find("tr").Find("tbody").Find("tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.TrimSpace(s.Text()) == "【交易明细】" {
			begin = i
			return false
		}
		return true
	})
	len := doc.Find("tr").Find("tbody").Find("tr").Length()
	var rs []*Order
	var bill *Order
	var exit bool
	doc.Find("tr").Find("tbody").Find("tr").Slice(begin+4, len).EachWithBreak(func(i int, s *goquery.Selection) bool {
		c.LineNum++
		s.Find("td").EachWithBreak(func(i int, s *goquery.Selection) bool {
			value := s.Text()
			if value == "*** 结束 The End ***" {
				exit = true
				return false
			}
			switch i {
			case 0:
				bill = &Order{}
				tm, _ := time.Parse("2006-01-02", strings.TrimSpace(value))
				bill.TransDate = tm
				rs = append(rs, bill)
			case 1:
				tm, _ := time.Parse("2006-01-02", strings.TrimSpace(value))
				bill.PostDate = tm
			case 2:
				bill.CardNo = strings.TrimSpace(value)
			case 3:
				bill.Description = strings.TrimSpace(strings.ReplaceAll(value, "\u00a0", ""))
			case 4:
				bill.TransCurrency = strings.TrimSpace(value)
			case 5:
				value = strings.TrimSpace(value)
				value = strings.ReplaceAll(value, ",", "")
				amount, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Printf("parse original amount failed: %s", err)
					return true
				}
				bill.MoneyOriginal = amount
			case 6:
				bill.MoneyCurrency = strings.TrimSpace(value)
			case 7:
				value = strings.TrimSpace(value)
				value = strings.ReplaceAll(value, ",", "")
				amount, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Printf("parse original amount failed: %s", err)
					return true
				}
				bill.Money = amount
			}
			// fmt.Println(i, value)
			return true
		})
		return !exit
	})

	for _, v := range rs {
		c.Orders = append(c.Orders, *v)
	}
	log.Printf("Finished to parse the file %s", filename)
	// log.Printf("data1: %v", c.Orders)
	return c.convertToIR(), nil
}
