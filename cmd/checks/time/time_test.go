// +build linux

package time

import (
	"context"
	"syscall"
	"testing"

	"github.com/dcos/dcos-checks/constants"
	"github.com/pkg/errors"
)

func TestTimeCheckBadStatus(t *testing.T) {
	mockrunAdjtimex := func(t *syscall.Timex) (int, error) {
		t.Status = 0x0040
		return 0, nil
	}

	check := &timeCheck{
		runAdjtimex: mockrunAdjtimex,
	}

	msg, code, err := check.Run(context.TODO(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if code != constants.StatusFailure {
		t.Fatalf("expect status %d. Got %d", constants.StatusFailure, code)
	}

	expectedMsg := "Clock is out of sync / in unsync state. Must be synchronized for proper operation."
	if msg != expectedMsg {
		t.Fatalf("expect %s. Got %s", expectedMsg, msg)
	}
}

func TestTimeCheckClockStable(t *testing.T) {
	mockrunAdjtimex := func(t *syscall.Timex) (int, error) {
		t.Esterror = maxEstErrorUs + 1000
		return 0, nil
	}

	check := &timeCheck{
		runAdjtimex: mockrunAdjtimex,
	}

	msg, code, err := check.Run(context.TODO(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if code != constants.StatusFailure {
		t.Fatalf("expect status %d. Got %d", constants.StatusFailure, code)
	}

	expectedMsg := "Clock is less stable than allowed. Max estimated error exceeded by: 1ms"
	if msg != expectedMsg {
		t.Fatalf("expect %s. Got %s", expectedMsg, msg)
	}
}

func TestTimeCheckError(t *testing.T) {
	mockrunAdjtimex := func(t *syscall.Timex) (int, error) {
		return 1, errors.New("error")
	}

	check := &timeCheck{
		runAdjtimex: mockrunAdjtimex,
	}

	_, _, err := check.Run(context.TODO(), nil)
	if err == nil {
		t.Fatal("expect error. Got nil")
	}
}

func TestTimeCheck(t *testing.T) {
	mockrunAdjtimex := func(t *syscall.Timex) (int, error) {
		return 0, nil
	}

	check := &timeCheck{
		runAdjtimex: mockrunAdjtimex,
	}

	msg, code, err := check.Run(context.TODO(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if code != constants.StatusOK {
		t.Fatalf("expect code %d. Got %d", constants.StatusOK, code)
	}

	expectedMsg := "Clock is synced"
	if msg != expectedMsg {
		t.Fatalf("expect msg %s. Got %s", expectedMsg, msg)
	}
}
