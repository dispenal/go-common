package common_utils

func ConvertInterfaceP(from any, to any) {
	data, err := Marshal(from)
	PanicIfAppError(err, "failed when marshal interface", 422)

	err = Unmarshal(data, to)
	PanicIfAppError(err, "failed when unmarshal interface", 422)
}

func ConvertInterfaceE(from any, to any) error {
	data, err := Marshal(from)
	if err != nil {
		return err
	}
	return Unmarshal(data, to)
}
