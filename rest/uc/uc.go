package uc

// Store is the repository interface the use-case layer depends on.
// The use-case never imports the repo package directly — only this interface (DIP).
type Store interface {
	GetById(id int) (string, bool)
	Save(id int, name string)
}

// GetUc and SaveUc are function types, not structs.
// Using functions as use-cases keeps them lightweight and easy to test in isolation.
type GetUc func(id int) (string, bool)
type SaveUc func(id int, name string)

// MakeGetUc is a factory that injects the store into the returned closure.
// The closure captures store once; subsequent calls don't need to pass it explicitly.
func MakeGetUc(store Store) GetUc {
	return func(id int) (string, bool) {
		return store.GetById(id)
	}
}

// MakeSaveUc mirrors MakeGetUc for the write path.
func MakeSaveUc(store Store) SaveUc {
	return func(id int, name string) {
		store.Save(id, name)
	}
}
