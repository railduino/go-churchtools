package churchtools

import (
	"encoding/json"
	//"fmt"
)

type FinanceClient struct {
	ID         int
	Name       string
	Street     string
	PostalCode string
	City       string
	Phone      string
	SortKey    int
	Meta       MetaInfo
}

type FinanceGetClientsResult struct {
	Data []FinanceClient
}

var (
	FinanceClients   []FinanceClient
	FinanceClientMap = make(map[int]*FinanceClient)
)

func (conn *Connector) FinanceGetClients() ([]FinanceClient, error) {
	content, err := conn.Get("finance/clients", true)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("FinanceGetClients '%s'\n", string(content))

	var result FinanceGetClientsResult
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	for _, client := range result.Data {
		FinanceClients = append(FinanceClients, client)
		FinanceClientMap[client.ID] = &client
	}

	return FinanceClients, nil
}
