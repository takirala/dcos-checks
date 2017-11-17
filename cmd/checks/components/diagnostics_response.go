package components

import (
	"fmt"
	"strings"

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

func (d *diagnosticsResponse) checkHealth() ([]string, int) {
	var errorList []string
	for _, unit := range d.Units {
		if (unit.Health != constants.StatusOK) && !(strings.Contains(unit.ID, "dcos-checks")) {
			errorList = append(errorList, fmt.Sprintf("component %s has health status %d", unit.Name, unit.Health))
		}
	}
	retCode := constants.StatusOK
	if len(errorList) > 0 {
		retCode = constants.StatusFailure
	}
	return errorList, retCode
}
