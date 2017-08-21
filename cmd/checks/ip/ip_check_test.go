package ip

import (
	"context"
	"testing"

	"github.com/dcos/dcos-checks/common"
)

func TestDetectIPCheck_Run(t *testing.T) {
	mockCLICfg := &common.CLIConfigFlags{}

	check := detectIPCheck{"./fixture/detect_ip.bad"}
	_, _, err := check.Run(context.TODO(), mockCLICfg)
	if err == nil {
		t.Fatal("expect error")
	}

	check = detectIPCheck{"./fixture/detect_ip.good"}
	_, _, err = check.Run(context.TODO(), mockCLICfg)
	if err != nil {
		t.Fatal(err)
	}

	check = detectIPCheck{"./fixture/detect_ip.empty"}
	_, _, err = check.Run(context.TODO(), mockCLICfg)
	if err == nil {
		t.Fatal("expect error")
	}
}
