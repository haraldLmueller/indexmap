package indexmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrimaryIndex(t *testing.T) {
	index := NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	})

	persons := GenPersons()

	for _, person := range persons {
		index.insert(person)
	}

	for _, person := range persons {
		assert.Equal(t,
			person, index.get(person.ID))
	}

	// Insert overwrite
	overwritePerson := &Person{1, "Tracer", 23, "London", []string{"Bob"}}
	index.insert(overwritePerson)
	assert.Equal(t, overwritePerson, index.get(overwritePerson.ID))

	// Remove get nil
	index.remove(overwritePerson.ID)
	assert.Nil(t, index.get(overwritePerson.ID))
}

func TestSecondaryIndex(t *testing.T) {
	// Many-to-One index
	index := NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name, value.City}
	})

	persons := GenPersons()

	for _, person := range persons {
		index.insert(person)
	}

	for _, person := range persons {
		result := index.get(person.Name)

		assert.Equal(t, 1, len(result))
		assert.Contains(t,
			result, person)

		result = index.get(person.City)
		assert.Contains(t,
			result, person)
	}

	// Insert makes One-to-Many, Many-to-Many
	ashe2 := &Person{4, "Ashe", 83, "Chengdu", nil}
	index.insert(ashe2)

	result := index.get(ashe2.Name)
	assert.Equal(t, 2, len(result))
	assert.Contains(t, result, ashe2)
	assert.Contains(t, result, persons[0])

	// Remove
	index.remove(ashe2)
	result = index.get(persons[0].Name)
	assert.Equal(t, 1, len(result))
	assert.Contains(t, result, persons[0])

	index.remove(persons[0])
	_, ok := index.inner[persons[0].Name]
	assert.False(t, ok)

}

func TestSecondaryContainsIndex(t *testing.T) {
	index := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	index.AddIndex("city", NewSecondaryIndex(func(value *Person) []any {
		return []any{value.City}
	}))

	// Like is a "contain" index
	index.AddIndex("like", NewSecondaryIndex(func(value *Person) []any {
		like := make([]any, 0, len(value.Like))
		for i := range value.Like {
			like = append(like, value.Like[i])
		}
		return like
	}))

	persons := GenPersons()

	for _, person := range persons {
		index.insert(person)
	}
	likes := index.getAllBy("like", "Bob")
	assert.Equal(t, 2, len(likes), "have to conain Ashe and Cassidy")
	likes = index.getAllBy("like", "Cassidy")
	assert.Equal(t, 2, len(likes), "have to conain Ashe and Harald")
	likes = index.getAllBy("like", "Ashe")
	assert.Equal(t, 1, len(likes), "Ashe have to like Cassidy")
	likes = index.getAllBy("like", "Harald")
	assert.Equal(t, 0, len(likes), "Harald have to like no likes")
}
