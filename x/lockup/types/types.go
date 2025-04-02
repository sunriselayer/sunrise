package types

import (
	"fmt"
)

func LockupAccountModule(owner string) string {
	return fmt.Sprintf("%s/%s", ModuleName, owner)
}
