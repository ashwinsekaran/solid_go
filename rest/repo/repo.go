package repo

// Repo is an in-memory store backed by a plain map.
// In a real service this would be replaced by a database-backed implementation
// that satisfies the same uc.Store interface — no other layer needs to change.
type Repo struct {
	Store map[int]string
}

func NewRepo() *Repo {
	return &Repo{Store: make(map[int]string)}
}

// GetById looks up a record by integer ID.
// The bool return follows Go's comma-ok idiom — callers distinguish "not found" from errors.
func (r *Repo) GetById(id int) (string, bool) {
	name, ok := r.Store[id]
	return name, ok
}

// Save writes or overwrites a record. Map access is not concurrent-safe;
// add a sync.RWMutex here if multiple goroutines share the same Repo.
func (r *Repo) Save(id int, name string) {
	r.Store[id] = name
}
