package components

import (
	"fmt"

	"github.com/dcos/dcos-checks/constants"
)

type diagnosticsResponse struct {
	Units []struct {
		ID          string `json:"id"`
		Health      int    `json:"health"`
		Output      string `json:"output"`
		Description string `json:"description"`
		Help        string `json:"help"`
		Name        string `json:"name"`
	} `json:"units"`
}

func (d *diagnosticsResponse) checkHealth(args []string) ([]string, int) {
	var errorList []string
	skipcomp := make(map[string]bool)
	if len(args) != 0 {
		// We want to skip some components
		m := make(map[string]bool)
		for _, unit := range d.Units {
			// convert list to map
			m[unit.ID] = true
		}
		for _, unit := range args {
			if !m[unit] {
				errorList = append(errorList, fmt.Sprintf("component %s does not exist", unit))
				return errorList, constants.StatusFailure
			}
			skipcomp[unit] = true
		}
	}

	for _, unit := range d.Units {
		if (unit.Health != constants.StatusOK) && !skipcomp[unit.ID] {
			errorList = append(errorList, fmt.Sprintf("component %s has health status %d", unit.Name, unit.Health))
		}
	}
	retCode := constants.StatusOK
	if len(errorList) > 0 {
		retCode = constants.StatusFailure
	}
	return errorList, retCode
}
