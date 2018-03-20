package common

import (
	"net"
	"net/http"
	"testing"
)

func TestCLIConfigFlags_IP(t *testing.T) {
	DCOSConfig.Role = "master"
	DCOSConfig.DetectIP = "fixture/detect_ip"
	ip, err := DCOSConfig.IP(&http.Client{})
	if err != nil {
		t.Fatal(err)
	}

	expectedIP := net.IPv4(192, 168, 0, 1)
	if !ip.Equal(expectedIP) {
		t.Fatalf("expected ip %s. Got %s", expectedIP, ip)
	}
}
