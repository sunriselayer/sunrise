package zkp

import (
	"reflect"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// Circuit defines a pre-image knowledge proof
// mimc(secret preImage) = public hash
type DemoCircuit struct {
	// struct tag on a variable is optional
	// default uses variable name and secret visibility.
	PreImage frontend.Variable
	Hash     frontend.Variable `gnark:",public"`
}

// Define declares the circuit's constraints
// Hash = mimc(PreImage)
func (circuit *DemoCircuit) Define(api frontend.API) error {
	// hash function
	mimc, _ := mimc.NewMiMC(api)

	// specify constraints
	// mimc(preImage) == hash
	mimc.Write(circuit.PreImage)
	api.AssertIsEqual(circuit.Hash, mimc.Sum())

	return nil
}

type CommitmentCircuit struct {
	Public []frontend.Variable `gnark:",public"`
	X      []frontend.Variable
}

func (c *CommitmentCircuit) Define(api frontend.API) error {
	commitment, err := api.(frontend.Committer).Commit(c.X...)
	if err != nil {
		return err
	}
	sum := frontend.Variable(0)
	for i, x := range c.X {
		sum = api.Add(sum, api.Mul(x, i+1))
	}
	for _, p := range c.Public {
		sum = api.Add(sum, p)
	}
	api.AssertIsDifferent(commitment, sum)
	return nil
}

func Hollow(c frontend.Circuit) frontend.Circuit {
	cV := reflect.ValueOf(c).Elem()
	t := reflect.TypeOf(c).Elem()
	res := reflect.New(t) // a new object of the same type as c
	resE := res.Elem()
	resC := res.Interface().(frontend.Circuit)

	frontendVar := reflect.TypeOf((*frontend.Variable)(nil)).Elem()

	for i := 0; i < t.NumField(); i++ {
		fieldT := t.Field(i).Type
		if fieldT.Kind() == reflect.Slice && fieldT.Elem().Implements(frontendVar) { // create empty slices for witness slices
			resE.Field(i).Set(reflect.ValueOf(make([]frontend.Variable, cV.Field(i).Len())))
		} else if fieldT != frontendVar { // copy non-witness variables
			resE.Field(i).Set(cV.Field(i))
		}
	}

	return resC
}
