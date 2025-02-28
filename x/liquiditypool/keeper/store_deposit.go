// GetDeposit returns the deposit for the given id
func (k Keeper) GetDeposit(ctx context.Context, id uint64) (deposit types.Deposit, found bool, err error) {
	has, err := k.Deposits.Has(ctx, id)
	if err != nil {
		return deposit, false, err
	}
	if !has {
		return deposit, false, nil
	}
	val, err := k.Deposits.Get(ctx, id)
	if err != nil {
		return deposit, false, err
	}
	return val, true, nil
}

// SetDeposit sets the deposit
func (k Keeper) SetDeposit(ctx context.Context, deposit types.Deposit) error {
	return k.Deposits.Set(ctx, deposit.Id, deposit)
}

// DeleteDeposit removes the deposit
func (k Keeper) DeleteDeposit(ctx context.Context, id uint64) error {
	return k.Deposits.Remove(ctx, id)
}

// GetAllDeposits returns all deposits
func (k Keeper) GetAllDeposits(ctx context.Context) (list []types.Deposit, err error) {
	err = k.Deposits.Walk(ctx, nil, func(key uint64, value types.Deposit) (bool, error) {
		list = append(list, value)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
} 