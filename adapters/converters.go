package adapters

import "errors"

func normalizePressure(pressure float64, unit string) (float64, error) {
	unitNotFoundError := errors.New("Unit not found")

	unitTable := make(map[string]float64)
	unitTable["mmHg"] = 1013.25 / 760

	rate, u_ok := unitTable[unit]
	if u_ok == false {
		return float64(0), unitNotFoundError
	}

	result := pressure * rate

	return result, nil
}
