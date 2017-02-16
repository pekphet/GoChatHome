package chathome

import (
	"fmt"
	"github.com/tyrchen/goutil/uniq"
)

const (
	NAME_PREFIX = "TEMP "
)

func getUniqName() string {
	return fmt.Sprintf("%s%d", NAME_PREFIX, uniq.GetUniq())
}
