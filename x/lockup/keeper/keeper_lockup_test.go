// This file is a unit test for the keeper's lockup logic.
// It covers the following functions:
// - TrackDelegation
// - TrackUndelegation
// - CheckUnbondingEntriesMature
// - AddRewardsToLockupAccount
package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func TestKeeper_TrackDelegation(t *testing.T) {
	owner := sdk.AccAddress("owner")

	testCases := []struct {
		name               string
		initialDelLocking  math.Int
		initialDelFree     math.Int
		balance            math.Int
		locked             math.Int
		amount             math.Int
		expectedDelLocking math.Int
		expectedDelFree    math.Int
		expectedErr        string
	}{
		{
			name:               "success - delegate from locked funds",
			initialDelLocking:  math.ZeroInt(),
			initialDelFree:     math.ZeroInt(),
			balance:            math.NewInt(1000),
			locked:             math.NewInt(1000),
			amount:             math.NewInt(100),
			expectedDelLocking: math.NewInt(100),
			expectedDelFree:    math.ZeroInt(),
		},
		{
			name:               "success - delegate from free funds",
			initialDelLocking:  math.NewInt(500),
			initialDelFree:     math.ZeroInt(),
			balance:            math.NewInt(1000),
			locked:             math.NewInt(500),
			amount:             math.NewInt(100),
			expectedDelLocking: math.NewInt(500),
			expectedDelFree:    math.NewInt(100),
		},
		{
			name:               "success - delegate from both locked and free",
			initialDelLocking:  math.NewInt(500),
			initialDelFree:     math.ZeroInt(),
			balance:            math.NewInt(1000),
			locked:             math.NewInt(600),
			amount:             math.NewInt(200),
			expectedDelLocking: math.NewInt(600), // 500 + 100
			expectedDelFree:    math.NewInt(100), // 0 + 100
		},
		{
			name:        "fail - zero amount",
			balance:     math.NewInt(1000),
			locked:      math.NewInt(1000),
			amount:      math.ZeroInt(),
			expectedErr: "delegation amount cannot be zero",
		},
		{
			name:        "fail - insufficient total balance",
			balance:     math.NewInt(99),
			locked:      math.NewInt(99),
			amount:      math.NewInt(100),
			expectedErr: "total balance is less than delegation amount",
		},
		{
			name:              "fail - insufficient free funds",
			initialDelLocking: math.NewInt(950),
			initialDelFree:    math.ZeroInt(),
			balance:           math.NewInt(1000),
			locked:            math.NewInt(950),
			amount:            math.NewInt(100),
			expectedErr:       types.ErrInsufficientUnlockedFunds.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)

			lockup := types.LockupAccount{
				Owner:            owner.String(),
				Id:               1,
				DelegatedLocking: tc.initialDelLocking,
				DelegatedFree:    tc.initialDelFree,
			}
			err := f.keeper.SetLockupAccount(f.ctx, lockup)
			require.NoError(t, err)

			err = f.keeper.TrackDelegation(f.ctx, owner, 1, tc.balance, tc.locked, tc.amount)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				lockup, err := f.keeper.GetLockupAccount(f.ctx, owner, 1)
				require.NoError(t, err)
				require.True(t, tc.expectedDelLocking.Equal(lockup.DelegatedLocking), "expected del locking %s, got %s", tc.expectedDelLocking, lockup.DelegatedLocking)
				require.True(t, tc.expectedDelFree.Equal(lockup.DelegatedFree), "expected del free %s, got %s", tc.expectedDelFree, lockup.DelegatedFree)
			}
		})
	}
}

func TestKeeper_TrackUndelegation(t *testing.T) {
	owner := sdk.AccAddress("owner")

	testCases := []struct {
		name               string
		initialDelLocking  math.Int
		initialDelFree     math.Int
		amount             math.Int
		expectedDelLocking math.Int
		expectedDelFree    math.Int
		expectedErr        string
	}{
		{
			name:               "success - undelegate from free",
			initialDelLocking:  math.NewInt(500),
			initialDelFree:     math.NewInt(100),
			amount:             math.NewInt(50),
			expectedDelLocking: math.NewInt(500),
			expectedDelFree:    math.NewInt(50),
		},
		{
			name:               "success - undelegate from locked",
			initialDelLocking:  math.NewInt(500),
			initialDelFree:     math.ZeroInt(),
			amount:             math.NewInt(100),
			expectedDelLocking: math.NewInt(400),
			expectedDelFree:    math.ZeroInt(),
		},
		{
			name:               "success - undelegate from both",
			initialDelLocking:  math.NewInt(500),
			initialDelFree:     math.NewInt(100),
			amount:             math.NewInt(200),
			expectedDelLocking: math.NewInt(400),
			expectedDelFree:    math.ZeroInt(),
		},
		{
			name:        "fail - zero amount",
			amount:      math.ZeroInt(),
			expectedErr: "invalid coins",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)

			lockup := types.LockupAccount{
				Owner:            owner.String(),
				Id:               1,
				DelegatedLocking: tc.initialDelLocking,
				DelegatedFree:    tc.initialDelFree,
			}
			err := f.keeper.SetLockupAccount(f.ctx, lockup)
			require.NoError(t, err)

			err = f.keeper.TrackUndelegation(f.ctx, owner, 1, tc.amount)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				lockup, err := f.keeper.GetLockupAccount(f.ctx, owner, 1)
				require.NoError(t, err)
				require.True(t, tc.expectedDelLocking.Equal(lockup.DelegatedLocking), "expected del locking %s, got %s", tc.expectedDelLocking, lockup.DelegatedLocking)
				require.True(t, tc.expectedDelFree.Equal(lockup.DelegatedFree), "expected del free %s, got %s", tc.expectedDelFree, lockup.DelegatedFree)
			}
		})
	}
}

func TestKeeper_CheckUnbondingEntriesMature(t *testing.T) {
	owner := sdk.AccAddress("owner")
	now := time.Now()

	testCases := []struct {
		name               string
		setup              func(f *fixture)
		checkTime          time.Time
		expectedEntries    int
		expectedDelLocking math.Int
	}{
		{
			name: "one entry matures",
			setup: func(f *fixture) {
				lockup := types.LockupAccount{
					Owner:            owner.String(),
					Id:               1,
					DelegatedLocking: math.NewInt(100),
					UnbondEntries: &types.UnbondingEntries{
						Entries: []*types.UnbondingEntry{
							{Amount: math.NewInt(50), EndTime: now.Add(-time.Hour).Unix()}, // matured
							{Amount: math.NewInt(50), EndTime: now.Add(time.Hour).Unix()},  // not matured
						},
					},
				}
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)
			},
			checkTime:          now,
			expectedEntries:    1,
			expectedDelLocking: math.NewInt(50), // 100 - 50
		},
		{
			name: "no entries mature",
			setup: func(f *fixture) {
				lockup := types.LockupAccount{
					Owner:            owner.String(),
					Id:               1,
					DelegatedLocking: math.NewInt(100),
					UnbondEntries: &types.UnbondingEntries{
						Entries: []*types.UnbondingEntry{
							{Amount: math.NewInt(50), EndTime: now.Add(time.Hour).Unix()},
						},
					},
				}
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)
			},
			checkTime:          now,
			expectedEntries:    1,
			expectedDelLocking: math.NewInt(100), // no change
		},
		{
			name: "no unbonding entries",
			setup: func(f *fixture) {
				lockup := types.LockupAccount{
					Owner:            owner.String(),
					Id:               1,
					DelegatedLocking: math.NewInt(100),
					UnbondEntries:    nil,
				}
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)
			},
			checkTime:          now,
			expectedEntries:    0,
			expectedDelLocking: math.NewInt(100), // no change
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.setup(f)

			sdkCtx := sdk.UnwrapSDKContext(f.ctx)
			ctx := sdkCtx.WithBlockTime(tc.checkTime)

			err := f.keeper.CheckUnbondingEntriesMature(ctx, owner, 1)
			require.NoError(t, err)

			lockup, err := f.keeper.GetLockupAccount(ctx, owner, 1)
			require.NoError(t, err)

			if lockup.UnbondEntries != nil {
				require.Len(t, lockup.UnbondEntries.Entries, tc.expectedEntries)
			} else {
				require.Equal(t, tc.expectedEntries, 0)
			}
			require.True(t, tc.expectedDelLocking.Equal(lockup.DelegatedLocking))
		})
	}
}

func TestKeeper_AddRewardsToLockupAccount(t *testing.T) {
	owner := sdk.AccAddress("owner")
	f := initFixture(t)

	// Setup initial account
	lockup := types.LockupAccount{
		Owner:             owner.String(),
		Id:                1,
		AdditionalLocking: math.NewInt(100),
	}
	err := f.keeper.SetLockupAccount(f.ctx, lockup)
	require.NoError(t, err)

	// Add rewards
	rewardAmount := math.NewInt(50)
	err = f.keeper.AddRewardsToLockupAccount(f.ctx, owner, 1, rewardAmount)
	require.NoError(t, err)

	// Verify
	updatedLockup, err := f.keeper.GetLockupAccount(f.ctx, owner, 1)
	require.NoError(t, err)

	expectedAmount := math.NewInt(150)
	require.True(t, expectedAmount.Equal(updatedLockup.AdditionalLocking))
}
