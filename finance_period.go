package churchtools

import (
	"encoding/json"
	"fmt"
)

type FinancePeriod struct {
	ID          int
	StartDate   string
	EndDate     string
	IsClosed    bool
	ClientID    int
	Client      *FinanceClient `json:"-"`
	Permissions Permissions
	// TODO donationReceiptsCreated
}

type FinanceGetPeriodsResult struct {
	Data []FinancePeriod
}

var (
	FinancePeriods   []FinancePeriod
	FinancePeriodMap = make(map[int]*FinancePeriod)
)

func (conn *Connector) FinanceGetPeriods() ([]FinancePeriod, error) {
	content, err := conn.Get("finance/accountingperiods", true)
	if err != nil {
		return nil, err
	}

	var result FinanceGetPeriodsResult
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	var ok bool
	for _, period := range result.Data {
		period.Client, ok = FinanceClientMap[period.ClientID]
		if !ok {
			return nil, fmt.Errorf("invalid clientId %d in FinancePeriod %d", period.ClientID, period.ID)
		}
		FinancePeriods = append(FinancePeriods, period)
		FinancePeriodMap[period.ID] = &period
	}

	return FinancePeriods, nil
}
