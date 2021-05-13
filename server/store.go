// Utils for interacting with KVStore
package main

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")
)

// KVEnsure makes sure the initial value for a key is set to the value provided, if it does not already exists
// Returns whether the value was set and error
func (p *Plugin) KVEnsure(key string, newValue []byte) (bool, error) {
	_, loadErr := p.KVLoad(key)
	switch loadErr {
	case nil:
		// value already set
		return false, nil
	case ErrNotFound:
		break
	default:
		return false, loadErr
	}

	ok, err := p.API.KVCompareAndSet(key, nil, newValue)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (p *Plugin) KVLoad(key string) ([]byte, error) {
	data, appErr := p.API.KVGet(key)
	if appErr != nil {
		return nil, errors.WithMessage(appErr, "failed plugin KVGet")
	}
	if data == nil {
		return nil, ErrNotFound
	}
	return data, nil
}
