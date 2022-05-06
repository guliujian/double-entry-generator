package cmbcredit

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
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type CmbCredit struct {
	Statistics Statistics `json:"statistics,omitempty"`
	LineNum    int        `json:"line_num,omitempty"`
	Orders     []Order    `json:"orders,omitempty"`
}

func New() *CmbCredit {
	return &CmbCredit{
		Statistics: Statistics{},
		LineNum:    0,
		Orders:     make([]Order, 0),
	}
}

func (c *CmbCredit) Translate(filename string) (*ir.IR, error) {
	log.SetPrefix("[Provider-CmbCredit]")
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	email, err := parsemail.Parse(bufio.NewReader(file))
	if err != nil {
		return nil, err
	}
	decoded := base64.NewDecoder(base64.StdEncoding, strings.NewReader(email.HTMLBody))
	b, err := ioutil.ReadAll(decoded)
	if err != nil {
		return nil, err
	}
	reader := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder())
	// log.Println(email.HTMLBody)
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(d))
	if err != nil {
		return nil, err
	}
	var startDate, endDate time.Time
	var year int
	doc.Find("span[id=\"fixBand6\"]").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			if _, ok := s.Attr("valign"); !ok {
				return
			}
			value := s.Text()
			switch (i - 2) % 7 {
			case 1:
				stringSlice := strings.Split(value, "-")
				startDate, _ = time.Parse("2006/01/02", strings.TrimSpace(stringSlice[0]))
				endDate, _ = time.Parse("2006/01/02", strings.TrimSpace(stringSlice[1]))
				year = startDate.Year()
			}

		})
	})
	var rs []*Order
	var bill *Order
	doc.Find("span[id$=\"fixBand15\"]").Each(func(i int, s *goquery.Selection) {
		// var bill Order
		c.LineNum++
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			if _, ok := s.Attr("valign"); !ok {
				return
			}

			value := s.Text()
			switch (i - 2) % 7 {
			case 0:
				bill = &Order{}
				rs = append(rs, bill)
			case 1:
				tm, _ := time.Parse("20060102", strconv.Itoa(year)+value)
				if (tm.After(startDate) || tm.Equal(startDate)) && (tm.Before(endDate) || tm.Equal(endDate)) {
					bill.PostDate = tm
				} else if tm.Before(startDate) {
					tm, _ = time.Parse("20060102", strconv.Itoa(year+1)+value)
					if (tm.After(startDate) || tm.Equal(startDate)) && (tm.Before(endDate) || tm.Equal(endDate)) {
						bill.PostDate = tm
					}
				}
			case 2:
				bill.Description = strings.ReplaceAll(value, "\u00a0", "")
			case 3:
				value = strings.TrimPrefix(value, "￥")
				value = strings.TrimSpace(value)
				value = strings.ReplaceAll(value, ",", "")
				value = strings.ReplaceAll(value, "\u00a0", "")
				amount, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Printf("parse rmb amount failed: %s", err)
					return
				}
				bill.Money = amount
			case 4:
				bill.CardNo = strings.TrimSpace(value)
			case 5:
				bill.Area = strings.TrimSpace(value)
			case 6:
				value = strings.ReplaceAll(value, ",", "")
				amount, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Printf("parse original trans amount failed: %s", err)
					return
				}
				bill.MoneyOriginal = amount
			}

		})
	})
	for _, v := range rs {
		c.Orders = append(c.Orders, *v)
	}
	log.Printf("Finished to parse the file %s", filename)
	// log.Printf("data1: %v", c.Orders)
	return c.convertToIR(), nil
}
