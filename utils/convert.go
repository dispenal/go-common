package common_utils

import "encoding/json"

func ConvertInterfaceP(from any, to any) {
	data, err := json.Marshal(from)
	PanicIfAppError(err, "failed when marshal interface", 422)

	err = json.Unmarshal(data, to)
	PanicIfAppError(err, "failed when unmarshal interface", 422)
}

func ConvertInterfaceE(from any, to any) error {
	data, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, to)
}
