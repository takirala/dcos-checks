package journald

import (
	"os/user"
	"strconv"
)

type (
	grp struct {
		id   uint32
		name string
	}
)

func (g grp) gid() (uint32, error) {
	if g.name != "" {
		group, err := user.LookupGroup(g.name)
		if err != nil {
			return 0, err
		}

		gid, err := strconv.ParseUint(group.Gid, 10, 32)
		if err != nil {
			return 0, err
		}

		return uint32(gid), nil
	}

	return g.id, nil
}
