package adapters

import "errors"

// Converts data to various units: mmHg -> bar, F -> C
func convertUnits(data float64, unit string) (float64, error) {
	unitNotFoundError := errors.New("Unit not found")

	dataTable := make(map[string]float64)
	dataTable["mmHg"] = data * 1013.25 / 760
	dataTable["kph"] = data / 3.6
	dataTable["F"] = (data - 32) * 5 / 9

	result, u_ok := dataTable[unit]
	if u_ok == false {
		return float64(0), unitNotFoundError
	}

	return result, nil
}
