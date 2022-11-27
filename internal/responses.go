package internal

type LastBlockResp struct {
	Result string `json:"result"`
}

type TransactionsBlockResp struct {
	Result  struct {
		Transactions    []struct {
			From    string  `json:"from"`
			To      string  `json:"to"`
			Value   string  `json:"value"`
		} `json:"transactions"`
	} `json:"result"`
}