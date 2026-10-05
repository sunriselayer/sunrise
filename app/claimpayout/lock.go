// lock.go serializes payout record updates so two operators cannot consume one nonce.
package claimpayout

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func withLock(path string, fn func() error) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("withLock: create %s: %w", directory, err)
	}
	lock, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("withLock: open %s.lock: %w", path, err)
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("withLock: lock %s: %w", path, err)
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
	return fn()
}
