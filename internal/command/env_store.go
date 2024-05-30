package command

import (
	"strings"
	"sync"

	"golang.org/x/exp/slices"
)

// func envPairSize(k, v string) int {
// 	// +2 for the '=' and '\0' separators
// 	return len(k) + len(v) + 2
// }

type envStore struct {
	mu    sync.RWMutex
	items map[string]string
}

func newEnvStore() *envStore {
	return &envStore{items: make(map[string]string)}
}

func (s *envStore) Merge(envs ...string) (*envStore, error) {
	for _, env := range envs {
		if _, err := s.Set(splitEnv(env)); err != nil {
			return s, err
		}
	}
	return s, nil
}

func (s *envStore) Get(k string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.items[k]
	return v, ok
}

func (s *envStore) Set(k, v string) (*envStore, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// environSize := envPairSize(k, v)

	// for key, value := range s.items {
	// 	if key == k {
	// 		continue
	// 	}
	// 	environSize += envPairSize(key, value)
	// }

	// if environSize > MaxEnvironSizeInBytes {
	// 	return s, errors.New("could not set environment variable, environment size limit exceeded")
	// }

	s.items[k] = v

	return s, nil
}

func (s *envStore) Delete(k string) *envStore {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, k)
	return s
}

func (s *envStore) Items() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]string, 0, len(s.items))
	for k, v := range s.items {
		result = append(result, k+"="+v)
	}
	slices.Sort(result)
	return result
}

func diffEnvStores(initial, updated *envStore) (newOrUpdated, unchanged, deleted []string) {
	initial.mu.RLock()
	defer initial.mu.RUnlock()
	updated.mu.RLock()
	defer updated.mu.RUnlock()

	for k, v := range initial.items {
		uVal, ok := updated.items[k]
		if !ok {
			deleted = append(deleted, k)
		} else if v != uVal {
			newOrUpdated = append(newOrUpdated, k+"="+uVal)
		} else {
			unchanged = append(unchanged, k)
		}
	}
	for k, v := range updated.items {
		_, ok := initial.items[k]
		if ok {
			continue
		}
		newOrUpdated = append(newOrUpdated, k+"="+v)
	}
	return
}

func splitEnv(str string) (string, string) {
	parts := strings.SplitN(str, "=", 2)
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return parts[0], ""
	default:
		return parts[0], parts[1]
	}
}
