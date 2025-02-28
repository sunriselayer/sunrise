// GetWithdraw returns the withdraw for the given id
func (k Keeper) GetWithdraw(ctx context.Context, id uint64) (withdraw types.Withdraw, found bool, err error) {
	has, err := k.Withdraws.Has(ctx, id)
	if err != nil {
		return withdraw, false, err
	}
	if !has {
		return withdraw, false, nil
	}
	val, err := k.Withdraws.Get(ctx, id)
	if err != nil {
		return withdraw, false, err
	}
	return val, true, nil
}

// SetWithdraw sets the withdraw
func (k Keeper) SetWithdraw(ctx context.Context, withdraw types.Withdraw) error {
	return k.Withdraws.Set(ctx, withdraw.Id, withdraw)
}

// DeleteWithdraw removes the withdraw
func (k Keeper) DeleteWithdraw(ctx context.Context, id uint64) error {
	return k.Withdraws.Remove(ctx, id)
}

// GetAllWithdraws returns all withdraws
func (k Keeper) GetAllWithdraws(ctx context.Context) (list []types.Withdraw, err error) {
	err = k.Withdraws.Walk(ctx, nil, func(key uint64, value types.Withdraw) (bool, error) {
		list = append(list, value)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
} 