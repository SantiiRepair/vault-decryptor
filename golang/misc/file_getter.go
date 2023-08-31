package misc

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	color "github.com/fatih/color"
)

func ExtractVaultFromFile(data string) []byte {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	var vaultBody string

	// Attempt 1: raw JSON
	err := json.Unmarshal([]byte(data), &vaultBody)
	if err == nil {
		formated := strings.Trim(strings.ReplaceAll(vaultBody, "\\", ""), "\"")
		vaultB, err := json.Marshal(formated)
		if err != nil {
			red.Printf("[ERROR]: %s", err)
			os.Exit(1)
		}
		return vaultB
	}

	// Attempt 2: pre-v3 cleartext
	matches := regexp.MustCompile(`{"wallet-seed":"([^"}]*)"`).FindStringSubmatch(data)
	if len(matches) > 0 {
		mnemonic := matches[1]
		vaultMatches := regexp.MustCompile(`"wallet":("{[ -~]*\\"version\\":2}")`).FindStringSubmatch(data)
		var vault interface{}
		if len(vaultMatches) > 0 {
			json.Unmarshal([]byte(vaultMatches[1]), &vault)
		}
		vaultM := map[string]interface{}{
			"data": map[string]interface{}{
				"mnemonic": mnemonic,
				"vault":    vault,
			},
		}

		str := fmt.Sprint(vaultM)
		formated := strings.Trim(strings.ReplaceAll(str, "\\", ""), "\"")
		vaultB, err := json.Marshal(formated)
		if err != nil {
			red.Printf("[ERROR]: %s", err)
			os.Exit(1)
		}

		return vaultB
	}

	// Attempt 3: chromium 000003.log file on Linux
	matches = regexp.MustCompile(`"KeyringController":{"vault":"{[^{}]*}"`).FindStringSubmatch(data)
	if len(matches) > 0 {
		vaultBody = matches[0][29:]
		var vault interface{}
		formated := strings.Trim(strings.ReplaceAll(vaultBody, "\\", ""), "\"")
		json.Unmarshal([]byte(formated), &vault)
		vaultB, err := json.Marshal(vault)
		if err != nil {
			red.Printf("[ERROR]: %s", err)
			os.Exit(1)
		}

		return vaultB
	}

	// Attempt 4: chromium 000005.ldb on Windows
	matchRegex := regexp.MustCompile(`Keyring[0-9][^\}]*(\{[^\{\}]*\\"\})`)
	captureRegex := regexp.MustCompile(`Keyring[0-9][^\}]*(\{[^\{\}]*\\"\})`)
	ivRegex := regexp.MustCompile(`\\"iv.{1,4}[^A-Za-z0-9+\/]{1,10}([A-Za-z0-9+\/]{10,40}=*)`)
	dataRegex := regexp.MustCompile(`\\"[^":,is]*\\":\\"([A-Za-z0-9+\/]*=*)`)
	saltRegex := regexp.MustCompile(`,\\"salt.{1,4}[^A-Za-z0-9+\/]{1,10}([A-Za-z0-9+\/]{10,100}=*)`)
	matcher := matchRegex.FindAllStringSubmatch(data, -1)
	var vaults []interface{}
	for _, m := range matcher {
		captures := captureRegex.FindStringSubmatch(m[0])
		if len(captures) > 1 {
			dataMatch := dataRegex.FindStringSubmatch(captures[1])
			ivMatch := ivRegex.FindStringSubmatch(captures[1])
			saltMatch := saltRegex.FindStringSubmatch(captures[1])
			if len(dataMatch) > 1 && len(ivMatch) > 1 && len(saltMatch) > 1 {
				vault := map[string]interface{}{
					"data": dataMatch[1],
					"iv":   ivMatch[1],
					"salt": saltMatch[1],
				}
				vaults = append(vaults, vault)
			}
		}
	}
	if len(vaults) == 0 {
		return nil
	}
	if len(vaults) > 1 {
		yellow.Println("[WARNING]: Found multiple vaults!", vaults)
	}

	str := fmt.Sprint(vaults[0])
	formated := strings.Trim(strings.ReplaceAll(str, "\\", ""), "\"")
	vaultB, err := json.Marshal(formated)
	if err != nil {
		red.Printf("[ERROR]: %s", err)
		os.Exit(1)
	}

	return vaultB
}
