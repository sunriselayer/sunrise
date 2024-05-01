package namespace

import (
	tmrand "github.com/cometbft/cometbft/libs/rand"
	ns "github.com/sunriselayer/sunrise-app/pkg/namespace"
)

func RandomNamespace() ns.Namespace {
	for {
		id := RandomVerzionZeroID()
		namespace, err := ns.New(ns.NamespaceVersionZero, id)
		if err != nil {
			continue
		}
		return namespace
	}
}

func RandomVerzionZeroID() []byte {
	return append(ns.NamespaceVersionZeroPrefix, tmrand.Bytes(ns.NamespaceVersionZeroIDSize)...)
}
