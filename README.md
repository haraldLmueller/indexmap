# IndexMap

This is a fork from github.com/yah01/indexmap. So many thanks to yah01 from Shanghai.  
Although the base from yah01 is working very well, it wasn't a full thread safe version.
The base version was extended by Harald Mueller.

A map (hash table) is often created with $ID \to Object$ to search for data (structured in tables). The map-type but is limited to search for data using only an ID. In order to search data using any field without a SQL database. The IndexMap data structure can achieve it

## Installation
Get the IndexMap package:
```shell
go get -u "github.com/haraldLmueller/indexmap"
```

Import the IndexMap package:
```golang
import "github.com/haraldLmueller/indexmap"
```

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

Now, it works just like the common map type with the possibility of adding an index to search for a person using another field:
```golang
persons.AddIndex("name", indexmap.NewSecondaryIndex(func(value *Person) []any {
    return []any{value.Name}
}))
```
It must provide a way to extract keys for the inserted objects, all keys must be comparable.

The insertion, updates all indexes automatically:
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

And search for data using the primary index or an added index:
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
```golang
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
[API Reference](https://pkg.go.dev/github.com/haraldLmueller/indexmap)

### Update Value
Inserting different values using the same key, works like the normal map type. The last one overwrites the value, but for an inserted value modifying it from the outside may confuse the index. It must modify an internal value using `Update()/UpdateBy()`:
```golang
// DO NOT:
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

### Serialize & Deserialize
An IndexMap can be serialized to JSON, the result is the same as serializing a normal map type. It doesn't contain the index information, resulting in an unrecoverable map (indexes cannot be recovered):
```golang
// Serialize
imapData, err := json.Marshal(imap)

// Deserialize
// You have to create an IndexMap with primary index,
// it's acceptable to add secondary index after deserializing
imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
    return value.ID
}))
err := json.Unmarshal(imapData, &imap)
```

### Iterate
As well as sync.Map, IndexMap can iterate using the `Range()` method:
```golang
imap.Range(func(key int64, value *Person) bool {
    fmt.Printf("key=%v, value=%+v\n", key, value)
    return true
})
```
There is an ordered version as well. To use that, you have to set
a compare function first, i.e. :
```golang
imap.SetCmpFn(func(value1, Value2 *Person) int {
		return cmp.Compare(value1.Age, Value2.Age)
	})
```
With that RangeOrdered ranges the values in ascending age order.

Additionally, a useful method to get all keys and values:
```golang
keys, values := imap.Collect()
```

## Performance
Let $n$ be the number of elements inserted, $m$ be the number of indexes:
| Operation | Complexity |
| --------- | ---------- |
| Get       | $O(1)$     |
| GetBy     | $O(1)$     |
| Insert    | $O(m)$     |
| Update    | $O(m)$     |
| Remove    | $O(m)$     |
| AddIndex  | $O(n)$     |

The more indexes, the slower the write operations.

### Benchmarks

### Version 1.2.0

```
goos: linux
goarch: amd64
pkg: github.com/haraldLmueller/indexmap
cpu: Intel(R) Celeron(R) N4100 CPU @ 1.10GHz
BenchmarkInsertOnlyPrimaryInt-4                       	  813411	      1755 ns/op	     151 B/op	       3 allocs/op
BenchmarkParallelInsertOnlyPrimaryInt-4               	  984680	      2012 ns/op	     183 B/op	       3 allocs/op
BenchmarkNativeMap-4                                  	 1000000	      1675 ns/op	     181 B/op	       3 allocs/op
BenchmarkNativeSyncMap-4                              	  792746	      2461 ns/op	     225 B/op	       7 allocs/op
BenchmarkParallelNativeSyncMap-4                      	  796434	      2591 ns/op	     227 B/op	       7 allocs/op
BenchmarkUpdateNoIndexedValue-4                       	  507723	      2270 ns/op	     129 B/op	       8 allocs/op
BenchmarkUpdatePrimaryIndexedValue-4                  	  549063	      2111 ns/op	     129 B/op	       8 allocs/op
BenchmarkUpdateOneSecondaryIndexedValue-4             	  438856	      2356 ns/op	     129 B/op	       8 allocs/op
BenchmarkUpdateTwoSecondaryIndexedValue-4             	  500124	      2311 ns/op	     129 B/op	       8 allocs/op
BenchmarkUpdatePrimaryAndTwoSecondaryIndexedValue-4   	  547125	      2197 ns/op	     130 B/op	       8 allocs/op
```