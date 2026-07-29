package creative

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

var ErrJobRunning = errors.New("creative job is already running")

type JobStore struct {
	path string
	mu   sync.Mutex
}

func NewJobStore(path string) *JobStore { return &JobStore{path: path} }

func (s *JobStore) Save(job CreativeJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	jobs, err := s.load()
	if err != nil {
		return err
	}
	return s.saveLocked(jobs, job)
}
func (s *JobStore) Begin(job CreativeJob) (CreativeJob, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	jobs, err := s.load()
	if err != nil {
		return CreativeJob{}, false, err
	}
	for _, existing := range jobs {
		if existing.JobID == job.JobID {
			if existing.Status == "running" {
				return CreativeJob{}, false, ErrJobRunning
			}
			if existing.Status == "completed" {
				return existing, true, nil
			}
		}
	}
	return CreativeJob{}, false, s.saveLocked(jobs, job)
}
func (s *JobStore) saveLocked(jobs []CreativeJob, job CreativeJob) error {
	replaced := false
	for i := range jobs {
		if jobs[i].JobID == job.JobID {
			jobs[i] = job
			replaced = true
			break
		}
	}
	if !replaced {
		jobs = append(jobs, job)
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(s.path), ".creative-jobs-*.json")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err = tmp.Chmod(0o600); err == nil {
		_, err = tmp.Write(raw)
	}
	if err == nil {
		err = tmp.Sync()
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(name, s.path)
}
func (s *JobStore) List(tenantID, caseCode string) ([]CreativeJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	jobs, err := s.load()
	if err != nil {
		return nil, err
	}
	out := []CreativeJob{}
	for _, job := range jobs {
		if job.TenantID == tenantID && (caseCode == "" || job.CaseCode == caseCode) {
			out = append(out, job)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt > out[j].CreatedAt })
	return out, nil
}
func (s *JobStore) load() ([]CreativeJob, error) {
	raw, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return []CreativeJob{}, nil
	}
	if err != nil {
		return nil, err
	}
	var jobs []CreativeJob
	if err = json.Unmarshal(raw, &jobs); err != nil {
		return nil, fmt.Errorf("decode creative job store: %w", err)
	}
	return jobs, nil
}
