package churchtools

import (
	"encoding/json"
	"fmt"
)

type FinanceAccount struct {
	ID                      int
	Number                  string
	Name                    string
	AccountGroupID          int
	AccountingPeriodID      int
	IsDonationAccount       bool
	IsOpeningBalanceAccount bool
	Balance                 int
	Meta                    MetaInfo
	Permissions             Permissions
	// TODO TaxRateId
}

type FinanceGetAccountsResult struct {
	Data []FinanceAccount
}

func (conn *Connector) FinanceGetAccounts(period_id int) ([]FinanceAccount, error) {
	var accounts []FinanceAccount

	url := fmt.Sprintf("finance/accounts?accounting_period_id=%d", period_id)
	content, err := conn.Get(url, true)
	if err != nil {
		return nil, err
	}

	var result FinanceGetAccountsResult
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	for _, account := range result.Data {
		accounts = append(accounts, account)
	}

	return accounts, nil
}
