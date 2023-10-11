package cgbcredit

import (
	"log"
	"strings"

	"github.com/deb-sig/double-entry-generator/pkg/config"
	"github.com/deb-sig/double-entry-generator/pkg/ir"
	"github.com/deb-sig/double-entry-generator/pkg/util"
)

type CGBCredit struct {
}

func (c CGBCredit) GetAllCandidateAccounts(cfg *config.Config) map[string]bool {
	uniqMap := make(map[string]bool)
	if cfg.Ccbcredit == nil || len(cfg.Ccbcredit.Rules) == 0 {
		return uniqMap
	}
	for _, r := range cfg.Ccbcredit.Rules {
		if r.MethodAccount != nil {
			uniqMap[*r.MethodAccount] = true
		}
		if r.TargetAccount != nil {
			uniqMap[*r.TargetAccount] = true
		}

	}
	uniqMap[cfg.DefaultPlusAccount] = true
	uniqMap[cfg.DefaultMinusAccount] = true
	return uniqMap
}

func (c CGBCredit) GetAccountsAndTags(o *ir.Order, cfg *config.Config, target, provider string) (bool, string, string, map[ir.Account]string, []string) {
	ignore := false
	if cfg.CGBCredit == nil || len(cfg.CGBCredit.Rules) == 0 {
		return ignore, cfg.DefaultMinusAccount, cfg.DefaultPlusAccount, nil, nil
	}

	resMinus := cfg.DefaultMinusAccount
	resPlus := cfg.DefaultPlusAccount
	var extraAccounts map[ir.Account]string
	var err error
	var tags = make([]string, 0)
	for _, r := range cfg.CGBCredit.Rules {
		match := true
		sep := ","
		if r.Separator != nil {
			sep = *r.Separator
		}

		matchFunc := util.SplitFindContains
		if r.FullMatch {
			matchFunc = util.SplitFindEquals
		}
		if r.Peer != nil {
			match = matchFunc(*r.Peer, o.Peer, sep, match)
		}
		if r.Item != nil {
			match = matchFunc(*r.Item, o.Item, sep, match)
		}
		if r.Time != nil {
			match, err = util.SplitFindTimeInterval(*r.Time, o.PayTime, match)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
		if match {
			// Support multiple matches, like one rule matches the
			// minus accout, the other rule matches the plus account.
			if r.TargetAccount != nil {
				if o.Type == ir.TypeRecv {
					resMinus = *r.TargetAccount
				} else {
					resPlus = *r.TargetAccount
				}
			}
			if r.MethodAccount != nil {
				if o.Type == ir.TypeRecv {
					resPlus = *r.MethodAccount
				} else {
					resMinus = *r.MethodAccount
				}
			}
			if r.PnlAccount != nil {
				extraAccounts = map[ir.Account]string{
					ir.PnlAccount: *r.PnlAccount,
				}
			}
			if r.Tags != nil {
				tags = strings.Split(*r.Tags, sep)
			}
			if r.DropDuplicate {
				resMinus = ""
				resPlus = ""
				extraAccounts = nil
			}

		}
	}
	if strings.HasPrefix(o.Item, "退款-") && ir.TypeRecv != o.Type {
		return ignore, resPlus, resMinus, extraAccounts, tags
	}
	return ignore, resMinus, resPlus, extraAccounts, tags
}