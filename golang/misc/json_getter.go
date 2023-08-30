package misc

import "encoding/json"

func GetValuesFromJSON(jsonByte []byte) ([]byte, error) {
	var payload map[string]interface{}
	err := json.Unmarshal(jsonByte, &payload)
	if err != nil {
		return nil, err
	}

	if keyringController, ok := payload["KeyringController"].(map[string]interface{}); ok {
		if vault, ok := keyringController["vault"].(map[string]interface{}); ok {
			jsd, err := json.Marshal(vault)
			if err != nil {
				return nil, err
			}

			return jsd, nil
		}
	}

	if vault, ok := payload["vault"].(map[string]interface{}); ok {
		jsd, err := json.Marshal(vault)
		if err != nil {
			return nil, err
		}

		return jsd, nil
	}

	data := make(map[string]interface{})
	data["data"] = payload["data"]
	data["salt"] = payload["salt"]
	data["iv"] = payload["iv"]

	jsd, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsd, nil
}
