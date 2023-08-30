package misc

import "encoding/json"

func getDataFromJSON(jsonData string) (map[string]interface{}, error) {
	var payload map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &payload)
	if err != nil {
		return nil, err
	}

	if keyringController, ok := payload["KeyringController"].(map[string]interface{}); ok {
		if vault, ok := keyringController["vault"].(map[string]interface{}); ok {
			return vault, nil
		}
	}

	if vault, ok := payload["vault"].(map[string]interface{}); ok {
		return vault, nil
	}

	data := make(map[string]interface{})
	data["data"] = payload["data"]
	data["salt"] = payload["salt"]
	data["iv"] = payload["iv"]

	return data, nil
}
