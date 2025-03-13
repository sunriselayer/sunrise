package types

import (
	"fmt"
)

func SelfDelegateProxyAccountModuleName(owner string) string {
	return fmt.Sprintf("%s/proxy/%s", ModuleName, owner)
}
