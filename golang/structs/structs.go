package structs

type Payload struct {
	Data string `json:"data"`
	Salt string `json:"salt"`
	Iv   string `json:"iv"`
}

type Vault struct {
	Type string `json:"type"`
	Data struct {
		Mnemonic         []byte `json:"mnemonic"`
		NumberOfAccounts int    `json:"numberOfAccounts"`
		HDPath           string `json:"hdPath"`
	} `json:"data"`
}