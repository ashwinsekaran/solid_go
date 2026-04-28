package uc

type Store interface {
	GetById(id int) (string, bool)
	Save(id int, name string)
}

type GetUc func(id int) (string, bool)
type SaveUc func(id int, name string)

func MakeGetUc(store Store) GetUc {
	return func(id int) (string, bool) {
		return store.GetById(id)
	}
}

func MakeSaveUc(store Store) SaveUc {
	return func(id int, name string) {
		store.Save(id, name)
	}
}
