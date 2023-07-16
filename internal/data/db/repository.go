package db

import "errors"

var (
	ErrNotFound = errors.New("record not found")
)

type Repository[Entity any, ID any] struct {
}

func (r *Repository[Entity, ID]) Create(entity *Entity) error {
	_, err := DB.Insert(entity)
	return err
}

func (r *Repository[Entity, ID]) FindOneByID(id ID) (*Entity, error) {
	entity := new(Entity)
	has, err := DB.ID(id).Get(entity)
	if !has {
		return nil, ErrNotFound
	}
	return entity, err
}

func (r *Repository[Entity, ID]) FindOne(filter *Entity) (*Entity, error) {
	entity := new(Entity)
	err := DB.Find(entity, filter)
	return entity, err
}

func (r *Repository[Entity, ID]) UpdateByID(id ID, entity *Entity) error {
	aff, err := DB.ID(id).Update(entity)
	if aff == 0 {
		return ErrNotFound
	}
	return err
}

func (r *Repository[Entity, ID]) Update(filter *Entity, entity *Entity) (affected int, err error) {
	aff, err := DB.Update(entity, filter)
	if aff == 0 {
		return 0, ErrNotFound
	}
	return int(aff), err
}

func (r *Repository[Entity, ID]) DeleteByID(id ID) error {
	aff, err := DB.ID(id).Delete()
	if aff == 0 {
		return ErrNotFound
	}
	return err
}

func (r *Repository[Entity, ID]) Delete(filter *Entity) (affected int, err error) {
	aff, err := DB.Delete(filter)
	if aff == 0 {
		return 0, ErrNotFound
	}
	return int(aff), err
}

func (r *Repository[Entity, ID]) Upsert(id ID, entity *Entity) error {
	aff, err := DB.ID(id).Update(entity)
	if aff == 1 || err != nil {
		return err
	}

	_, err = DB.Insert(entity)
	return err
}
