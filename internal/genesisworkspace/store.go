package genesisworkspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type fileState struct {
	ByIdempotency map[string]Workspace `json:"by_idempotency"`
}

type Store struct {
	path string
	mu   sync.Mutex
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func (s *Store) Find(owner, idempotency string) (Workspace, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, err := s.load()
	if err != nil {
		return Workspace{}, false, err
	}
	item, ok := state.ByIdempotency[storeKey(owner, idempotency)]
	return item, ok, nil
}

func (s *Store) Save(owner, idempotency string, workspace Workspace) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, err := s.load()
	if err != nil {
		return err
	}
	state.ByIdempotency[storeKey(owner, idempotency)] = workspace
	return s.write(state)
}

func (s *Store) List(owner string) ([]Workspace, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, err := s.load()
	if err != nil {
		return nil, err
	}
	items := make([]Workspace, 0)
	for _, item := range state.ByIdempotency {
		if item.OwnerPlayerID == owner {
			items = append(items, item)
		}
	}
	return items, nil
}

func (s *Store) load() (fileState, error) {
	state := fileState{ByIdempotency: map[string]Workspace{}}
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return state, err
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return state, fmt.Errorf("decode genesis workspace store: %w", err)
	}
	if state.ByIdempotency == nil {
		state.ByIdempotency = map[string]Workspace{}
	}
	return state, nil
}

func (s *Store) write(state fileState) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(s.path), ".genesis-workspaces-*.json")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, s.path)
}

func storeKey(owner, idempotency string) string {
	return owner + "\x00" + idempotency
}
