package lockupdepinject

import (
	"cosmossdk.io/x/accounts/accountstd"

	lockup "github.com/sunriselayer/sunrise/x/accounts/non_voting_delegatable_lockup"
)

func ProvideAllLockupAccounts() []accountstd.DepinjectAccount {
	return []accountstd.DepinjectAccount{
		// ProvidePeriodicLockingAccount(),
		ProvideContinuousLockingAccount(),
		// ProvidePermanentLockingAccount(),
		// ProvideDelayedLockingAccount(),
	}
}

func ProvideContinuousLockingAccount() accountstd.DepinjectAccount {
	return accountstd.DIAccount(lockup.CONTINUOUS_LOCKING_ACCOUNT, lockup.NewContinuousLockingAccount)
}

// func ProvidePeriodicLockingAccount() accountstd.DepinjectAccount {
// 	return accountstd.DIAccount(lockup.PERIODIC_LOCKING_ACCOUNT, lockup.NewPeriodicLockingAccount)
// }

// func ProvideDelayedLockingAccount() accountstd.DepinjectAccount {
// 	return accountstd.DIAccount(lockup.DELAYED_LOCKING_ACCOUNT, lockup.NewDelayedLockingAccount)
// }

// func ProvidePermanentLockingAccount() accountstd.DepinjectAccount {
// 	return accountstd.DIAccount(lockup.PERMANENT_LOCKING_ACCOUNT, lockup.NewPermanentLockingAccount)
// }
