package indexmap

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexMap(t *testing.T) {
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	ok := imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))
	assert.True(t, ok)

	persons := GenPersons()
	InsertData(imap, persons)

	for i, person := range persons {
		assert.Equal(t,
			persons[i], imap.Get(persons[i].ID))

		assert.Equal(t,
			person, imap.Get(person.ID))

		assert.Equal(t,
			person, imap.GetBy(NameIndex, person.Name))

		assert.Nil(t, imap.GetBy(InvalidIndex, person.Name))

		result := imap.GetAllBy(NameIndex, person.Name)
		assert.Equal(t, 1, len(result))
		assert.Contains(t, result, person)

		assert.Nil(t, imap.getAllBy(InvalidIndex, person.Name))
	}

	// Add index after inserting data
	ok = imap.AddIndex(CityIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.City}
	}))
	assert.True(t, ok)

	for _, person := range persons {
		assert.Equal(t,
			person, imap.GetBy(NameIndex, person.Name))

		result := imap.GetAllBy(CityIndex, person.City)
		assert.Contains(t, result, person)
	}

	// Remove
	imap.Remove(persons[0].ID)
	assert.Nil(t, imap.Get(persons[0].ID))
	assert.Nil(t, imap.GetBy(NameIndex, persons[0].Name))
	assert.Empty(t, imap.GetAllBy(NameIndex, persons[0].Name))

	imap.RemoveBy(CityIndex, "San Francisco")
	assert.Empty(t, imap.GetAllBy(CityIndex, "San Francisco"))
	assert.Equal(t, 1, len(imap.GetAllBy(CityIndex, "Shanghai")))

	// Update
	imap.Clear()
	InsertData(imap, persons)
	imap.Update(persons[0].ID, func(value *Person) (*Person, bool) {
		value.Name = "Tracer"
		return value, true
	})
	assert.Equal(t, "Tracer", imap.Get(persons[0].ID).Name)

	count := len(imap.GetAllBy(CityIndex, "Shanghai"))
	imap.UpdateBy(CityIndex, "Shanghai", func(value *Person) (*Person, bool) {
		value.City = "Beijing"
		return value, true
	})
	assert.Empty(t, imap.GetAllBy(CityIndex, "Shanghai"))
	assert.Equal(t, count, len(imap.GetAllBy(CityIndex, "Beijing")))

	// Collect
	keys, values := imap.Collect()
	assert.Equal(t, imap.Len(), len(keys))
	assert.Equal(t, imap.Len(), len(values))
	for i := range keys {
		assert.Equal(t, values[i], imap.Get(keys[i]))
	}

	// CollectBy
	ks, vs := imap.CollectBy(CityIndex)
	for i := range ks {
		exp := vs[i]
		sort.SliceStable(exp, func(k, j int) bool {
			return exp[k].ID < exp[j].ID
		})

		act := imap.GetAllBy(CityIndex, ks[i])
		sort.SliceStable(act, func(k, j int) bool {
			return act[k].ID < act[j].ID
		})

		assert.Equal(t, exp, act)
	}

	// Range
	count = 0
	imap.Range(func(key int64, value *Person) bool {
		count++
		assert.Equal(t, value, imap.Get(key))
		return true
	})
	assert.Equal(t, imap.Len(), count)

	// RangeBy
	count = 0
	imap.RangeBy(CityIndex, func(key any, vals []*Person) bool {
		count++
		exp := imap.indexes[CityIndex].inner[key].Collect()
		sort.SliceStable(exp, func(i, j int) bool {
			return exp[i].ID < exp[j].ID
		})
		sort.SliceStable(vals, func(i, j int) bool {
			return vals[i].ID < vals[j].ID
		})
		assert.Equal(t, exp, vals)
		return true
	})
	assert.Equal(t, len(imap.indexes[CityIndex].inner), count)
}

func TestAddExistedIndex(t *testing.T) {
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	ok := imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))
	assert.True(t, ok)

	ok = imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.City}
	}))
	assert.False(t, ok)
}

func TestIndexMap_PrimaryKey(t *testing.T) {
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))
	p := GenPersons()[0]
	assert.Equal(t, p.ID, imap.PrimaryKey(p))
}

func BenchmarkInsertOnlyPrimaryInt(b *testing.B) {
	n := len(names)
	rand.Seed(123)
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))
	for i := 0; i < b.N; i++ {
		pi := int64(i)
		imap.Insert(&Person{pi, names[i%n], rand.Intn(106), "city", nil})
		r := imap.Get(pi)
		assert.Equal(b, pi, r.ID)
	}
}

func BenchmarkParallelInsertMonlyPrimaryInt(b *testing.B) {
	n := int64(len(names))
	rand.Seed(123)
	var i int64
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))
	b.RunParallel(func(pb *testing.PB) {
		age := rand.Intn(106)
		for pb.Next() {
			atomic.AddInt64(&i, 1)
			pi := i
			imap.Insert(&Person{pi, names[i%n], age, "city", nil})
			r := imap.Get(pi)
			assert.Equal(b, pi, r.ID)
		}
	})
}

func BenchmarkNativeMap(b *testing.B) {
	n := len(names)
	rand.Seed(123)
	imap := make(map[int64]*Person)
	for i := 0; i < b.N; i++ {
		pi := int64(i)
		imap[pi] = (&Person{pi, names[i%n], 10, "city", nil})
		r := imap[pi]
		assert.Equal(b, pi, r.ID)
	}
}

func BenchmarkNativeSyncMap(b *testing.B) {
	n := len(names)
	rand.Seed(123)
	var imap sync.Map
	age := rand.Intn(106)
	for i := 0; i < b.N; i++ {
		pi := int64(i)
		imap.Store(pi, &Person{pi, names[i%n], age, "city", nil})
		r, _ := imap.Load(pi)
		assert.Equal(b, pi, (r.(*Person)).ID)
	}
}
func BenchmarkParallelNativeSyncMap(b *testing.B) {
	n := int64(len(names))
	rand.Seed(123)
	var i int64
	var imap sync.Map
	b.RunParallel(func(pb *testing.PB) {
		age := rand.Intn(106)
		for pb.Next() {
			atomic.AddInt64(&i, 1)
			pi := i
			imap.Store(pi, &Person{pi, names[i%n], age, "city", nil})
			r, _ := imap.Load(pi)
			assert.Equal(b, pi, (r.(*Person)).ID)
		}
	})
}

func FuzzAddSecondaryIndex(f *testing.F) {
	var i int64
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))
	imap.AddIndex("name", NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))
	f.Add("John", "Doh", 34)
	f.Fuzz(func(t *testing.T, first string, city string, age int) {
		atomic.AddInt64(&i, 1)
		uniqName := fmt.Sprintf("%s-%d", first, i)
		pi := i
		imap.Insert(&Person{pi, uniqName, age, city, nil})
		ret := imap.GetAllBy("name", uniqName)
		assert.Equal(t, pi, ret[0].ID)
		assert.Equal(t, uniqName, ret[0].Name)
	})
}
