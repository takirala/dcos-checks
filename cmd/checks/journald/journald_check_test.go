package journald

import (
	"context"
	"os/user"
	"strconv"
	"testing"

	"github.com/dcos/dcos-checks/constants"
	"github.com/pkg/errors"
)

func mockCheckDirFn(err error) checkDirectoryFn {
	return func(s string, u uint32, m map[string]uint32) error {
		return err
	}
}

func TestJournalCheckFailure(t *testing.T) {
	e := errors.New("my error")
	c, err := newMockJournalCheck(e)
	if err != nil {
		t.Fatal(err)
	}

	output, code, err := c.Run(context.TODO(), nil)
	if output != "" {
		t.Fatalf("expected empty output. Got %s", output)
	}

	if code != constants.StatusUnknown {
		t.Fatalf("expected code ...Got %d", code)
	}

	if err != e {
		t.Fatalf("expect error %s, got %s", e, err)
	}
}

func TestJournalCheckSuccess(t *testing.T) {
	c, err := newMockJournalCheck(nil)
	out, code, err := c.Run(context.TODO(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if code != constants.StatusOK {
		t.Fatalf("Expect non 0 code. Got %d", code)
	}

	if out == "" {
		t.Fatal("Expect non empty output")
	}
}

func newMockJournalCheck(e error) (*journalCheck, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return nil, err
	}

	c := &journalCheck{
		checkDirFn: mockCheckDirFn(e),
		Path:       "/tmp",
		lookupGroup: grp{
			id: uint32(gid),
		},

		checkBits: map[string]uint32{"test": 1},
	}

	return c, nil
}
