package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "golang.org/x/crypto/pbkdf2"
)

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

var wallet = []byte(`{
	"data": "n89RkxtR7smzTno/uxOYNo6FIcLt526pEdpr4zGbCx5whIh67BJKKc0eKOheRjzPfaqtsMJooneY83f7mRWFWP2MYyG/27SoiafYUUQwn3O0WtA0h9PDeNXgmrCphcqevVas28ova8ERuhPC6by4R15f3kW7vFKySy4zYnaNrJlJQVvapi7LSeDztuR1qLLjAA3mYDp4y46qCVXBZmO2DFrUFFgdT8JHAZgh2Ar2hJGeEqrT/0S8zRKzZyqsPwTyPQMYbXr9k+kUz7AsqnxGyB/YbFS4fHyqEbYyFXk3KTS+JkJl8TsOHyhQJmUxaFLyCwGphx0xRH1icyONALmyUp1Um7irzUe5BpFalstxPUiaq0sbyWPsIeFBcqgze+ViFC25TS6+dryEOR6Ywm7vps9fez+DwqN0WD2TOpz9gaVqfzpPUiCK1/fyH4sfC1q4P3+qlfK9KWrmLbwW1NWAFcaqQ7QM1IYZSPS8cxu21FSXfqrghD99GYHvK/3xu5j5HUT/byWfW1tH+6jK3uZXHqa1lnKpwwI/kQwQIeJIOYJsQTMcwTSsKU4kMBzHNkRO1N/CDTPARQZ8dKphgI4opSUl8lwzPOAtZ0RdH9nX+SLQ2JTAD8axrd4UMelSTbfr2z0xsTSqzYU1eUC/dAb3Ih+YJYnDW/6qbOGLzPSiNYsm/uD3uCqTpYtdCFWK7ZVJkNybWj7+yOWUOGqL6mGJXdRIBPC9HGhV8qhbjTjR1yaVUbvsy7Fnw8ERP9Ct",
	"iv": "tXUhOuobmBraXu1HqHXc6w==",
	"salt": "6iAr4lY7y0SSeI/J0XyhEPYcgcOhWqB+8bBAUHUmS0A="
  }`)

func main() {
    var payload Payload
    json.Unmarshal(wallet, &payload)

    iv, _ := base64.StdEncoding.DecodeString(payload.Iv)
    salt, _ := base64.StdEncoding.DecodeString(payload.Salt)
    data, _ := base64.StdEncoding.DecodeString(payload.Data)

    password := "Elmetamask1"

    key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
    block, _ := aes.NewCipher(key)

    // Initialization Vector size is fixed to 16 bytes.
    gcm, _ := cipher.NewGCMWithNonceSize(block, len(iv))
    plaintext, err := gcm.Open(nil, iv, data, nil)

    if err != nil {
        panic(err)
    }

    var vault []Vault
    json.Unmarshal(plaintext, &vault)
    fmt.Println(string(vault[0].Data.Mnemonic))
}