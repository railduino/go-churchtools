package churchtools

import (
	"encoding/json"
	"fmt"
)

type FinanceClientMeta struct {
	CreatedPerson  MetaPerson `json:"createdPerson"`
	CreatedDate    string     `json:"createdDate"`
	ModifiedPerson MetaPerson `json:"modifiedPerson"`
	ModifiedDate   string     `json:"modifiedDate"`
}

type FinanceClient struct {
	ID      int               `json:"id"`
	Name    string            `json:"name"`
	SortKey int               `json:"sortKey"`
	Meta    FinanceClientMeta `json:"meta"`
}

type FinanceGetClientsResult struct {
	Data []FinanceClient `json:"data"`
}

func (conn *Connector) FinanceGetClients() ([]FinanceClient, error) {
	result, err := conn.Get("finance/clients", true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("FinanceGetClients '%s'\n", string(result))

	var clients FinanceGetClientsResult
	if err := json.Unmarshal(result, &clients); err != nil {
		return nil, err
	}

	return clients.Data, nil
}
