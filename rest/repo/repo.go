package repo

type Repo struct {
	Store map[int]string
}

func NewRepo() *Repo {
	return &Repo{
		Store: make(map[int]string),
	}
}

func (r *Repo) GetById(id int) (string, bool) {
	name, ok := r.Store[id]
	return name, ok
}

func (r *Repo) Save(id int, name string) {
	r.Store[id] = name
}
