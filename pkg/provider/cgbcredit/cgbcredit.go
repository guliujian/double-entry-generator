package cgbcredit

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/deb-sig/double-entry-generator/pkg/ir"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type CGBCredit struct {
	Statistics Statistics `json:"statistics,omitempty"`
	LineNum    int        `json:"line_num,omitempty"`
	Orders     []Order    `json:"orders,omitempty"`
}

func New() *CGBCredit {
	return &CGBCredit{
		Statistics: Statistics{},
		LineNum:    0,
		Orders:     make([]Order, 0),
	}
}

func (c *CGBCredit) Translate(filename string) (*ir.IR, error) {
	log.SetPrefix("[Provider-CGBCredit]")
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(transform.NewReader(bufio.NewReader(csvFile), simplifiedchinese.GBK.NewDecoder()))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	reader.Comma = ','
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {

			return nil, err
		}
		c.LineNum++
		if c.LineNum <= 1 {
			continue
		}
		err = c.translateToOrders(line)
		if err != nil {
			return nil, fmt.Errorf("Failed to translate bill: line %d: %v",
				c.LineNum, err)
		}
	}
	log.Printf("Finished to parse the file %s", filename)
	return c.convertToIR(), nil
}

func (c *CGBCredit) translateToOrders(array []string) error {
	for idx, a := range array {
		a = strings.TrimSpace(a)
		a = strings.ReplaceAll(a, ",", "")
		a = strings.Trim(a, "\"")
		array[idx] = a
	}

	var err error
	var bill Order
	// log.Printf("111 , %#v", array)
	bill.TransDate, err = time.Parse("2006-01-02 -0700 CST", array[0]+" +0800 CST")
	if err != nil {
		return fmt.Errorf("parse create time %s error: %v", array[0], err)
	}

	bill.TransCurrency = array[2]
	money := array[3]
	money = strings.ReplaceAll(money, "+", "")
	// money = strings.ReplaceAll(money, ",", "")

	bill.MoneyOriginal, err = strconv.ParseFloat(money, 64)
	if err != nil {
		return fmt.Errorf("parse oringal money %s error: %v ", money, err)
	}
	bill.MoneyCurrency = array[4]
	money = array[5]
	money = strings.ReplaceAll(money, "+", "")
	// money = strings.ReplaceAll(money, ",", "")

	bill.Money, err = strconv.ParseFloat(money, 64)
	if err != nil {
		return fmt.Errorf("parse money %s error: %v ", money, err)
	}
	bill.CardNo = array[6]

	c.Orders = append(c.Orders, bill)
	return err
}
