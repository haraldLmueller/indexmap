# IndexMap
We often created a map with $ID \to Object$ to seek data, but this limits us to seek the data with only ID. to seek data with any field without SQL in database, IndexMap is the data structure you can reach this.

## Get Started
First, to create a IndexMap with primary index:
```golang
type Person struct {
	ID   int64
	Name string
	Age  int
	City string
	Like []string
}

persons := indexmap.NewIndexMap(indexmap.NewPrimaryIndex(func(value *Person) int64 {
    return value.ID
}))
```

Now it's just like the common map type, but then you can add index to seek person with the other field:
```golang
persons.AddIndex("name", indexmap.NewSecondaryIndex(func(value *Person) []any {
    return []any{value.Name}
}))
```
You have to provide the way to extract keys for the inserted object, all keys must be comparable.

The insertion updates indexes automatically:
```golang
ashe := &Person{
    ID:   1,
    Name: "Ashe",
    Age:  39,
    City: "San Francisco",
    Like: []string{"Bob", "Cassidy"},
}
bob := &Person{
    ID:   2,
    Name: "Bob",
    Age:  18,
    City: "San Francisco",
}
cassidy := &Person{
    ID:   3,
    Name: "Cassidy",
    Age:  40,
    City: "Shanghai",
    Like: []string{"Ashe", "Bob"},
}

persons.Insert(ashe)
persons.Insert(bob)
persons.Insert(cassidy)
```

Adding index after inserting data also works:
```golang
persons.AddIndex("city", indexmap.NewSecondaryIndex(func(value *Person) []any {
    return []any{value.City}
}))

// Like is a "contain" index
persons.AddIndex("like", indexmap.NewSecondaryIndex(func(value *Person) []any {
    like := make([]any, 0, len(value.Like))
    for i := range value.Like {
        like = append(like, value.Like[i])
    }
    return like
}))
```

And seek data with primary index or the added index:
```golang
fmt.Println("Search with ID or Name:")
fmt.Printf("%+v\n", persons.Get(ashe.ID))
fmt.Printf("%+v\n", persons.GetBy("name", ashe.Name))

fmt.Println("\nSearch persons come from San Francisco:")
for _, person := range persons.GetAllBy("city", "San Francisco") {
    fmt.Printf("%+v\n", person)
}

fmt.Println("\nSearch persons like Bob")
for _, person := range persons.GetAllBy("like", "Bob") {
    fmt.Printf("%+v\n", person)
}
```

which outputs:
```
Search with ID or Name:
&{ID:1 Name:Ashe Age:39 City:San Francisco Like:[Bob Cassidy]}
&{ID:1 Name:Ashe Age:39 City:San Francisco Like:[Bob Cassidy]}

Search persons come from San Francisco:
&{ID:1 Name:Ashe Age:39 City:San Francisco Like:[Bob Cassidy]}
&{ID:2 Name:Bob Age:18 City:San Francisco Like:[]}

Search persons like Bob
&{ID:3 Name:Cassidy Age:40 City:Shanghai Like:[Ashe Bob]}
&{ID:1 Name:Ashe Age:39 City:San Francisco Like:[Bob Cassidy]}
```

## Document
[Reference](https://pkg.go.dev/github.com/yah01/indexmap)

### Update Value
Inserting the different values with the same key works like the normal map type, the last one overwrites the others, but for a inserted value, modifing it outside may confuse the index, modify an internal value with `Update()/UpdateBy()`:
```
DO NOT:
person := persons.GetBy("name", "Ashe")
person.City = "Shanghai"
persons.Insert(person)

// Modify the internal value with Update()/UpdateBy()
persons.UpdateBy("name", "Ashe", func(value *Person) (*Person, bool) {
    if value.City == "Shanghai" {
        return value, false
    }
    value.City = "Shanghai"
    return value, true
})
```

### One-To-Many/Many-To-Many Index
It's OK to create an index that's not one-to-one, The `GetBy()` method returns one of the object if many ones exist, `GetAllBy()` return a slice with all matched objects. For the example of many-to-many index, refer [contain_index_example](./examples/contain_index/main.go)

## Performance
Let $n$ be the number of elements inserted, $m$ be the number of indexes:
| Operation | Complexity |
| --------- | ---------- |
| Get       | $O(1)$     |
| GetBy     | $O(1)$     |
| Insert    | $O(m)$     |
| Remove    | $O(m)$     |
| AddIndex  | $O(n)$     |

The more indexes, the slower the write operations.