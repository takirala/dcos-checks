package cmd

import (
	"context"
	"testing"
)

func TestDetectIPCheck_Run(t *testing.T) {
	mockCLICfg := &CLIConfigFlags{}

	check := DetectIPCheck{"./fixture/detect_ip.bad"}
	_, _, err := check.Run(context.TODO(), mockCLICfg)
	if err == nil {
		t.Fatal("expect error")
	}

	check = DetectIPCheck{"./fixture/detect_ip.good"}
	_, _, err = check.Run(context.TODO(), mockCLICfg)
	if err != nil {
		t.Fatal(err)
	}

	check = DetectIPCheck{"./fixture/detect_ip.empty"}
	_, _, err = check.Run(context.TODO(), mockCLICfg)
	if err == nil {
		t.Fatal("expect error")
	}
}
