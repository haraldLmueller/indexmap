package indexmap

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

type Person struct {
	ID   int64
	Name string
	Age  int
	City string
	Like []string
}

const (
	InvalidIndex = "invalid"
	NameIndex    = "name"
	CityIndex    = "city"
	LikeIndex    = "like"
)

func GenPersons() map[int64]*Person {
	return map[int64]*Person{
		0: {0, "Ashe", 38, "San Francisco", []string{"Bob", "Cassidy"}},
		1: {1, "Bob", 18, "San Francisco", nil},
		2: {2, "Cassidy", 40, "Shanghai", []string{"Bob", "Ashe"}},
		3: {3, "Harald", 40, "Nürnberg", []string{"Cassidy"}},
	}
}

var names = []string{"James", "Mary", "Robert", "Patricia", "John", "Jennifer", "Michael", "Linda", "David", "Elizabeth", "William", "Barbara", "Richard",
	"Susan", "Joseph", "Jessica", "Thomas", "Sarah", "Charles", "Karen", "Christopher", "Lisa", "Daniel", "Nancy", "Matthew", "Betty", "Anthony",
	"Margaret", "Mark", "Sandra", "Donald", "Ashley", "Steven", "Kimberly", "Paul", "Emily", "Andrew", "Donna", "Joshua", "Michelle", "Kenneth",
	"Carol", "Kevin", "Amanda", "Brian", "Dorothy", "George", "Melissa", "Timothy", "Deborah", "Ronald", "Stephanie", "Edward", "Rebecca", "Jason",
	"Sharon", "Jeffrey", "Laura", "Ryan", "Cynthia", "Jacob", "Kathleen", "Gary", "Amy", "Nicholas", "Angela", "Eric", "Shirley", "Jonathan",
	"Anna", "Stephen", "Brenda", "Larry", "Pamela", "Justin", "Emma", "Scott", "Nicole", "Brandon", "Helen", "Benjamin", "Samantha", "Samuel",
	"Katherine", "Gregory", "Christine", "Alexander", "Debra", "Frank", "Rachel", "Patrick", "Carolyn", "Raymond", "Janet", "Jack", "Catherine",
	"Dennis", "Maria", "Jerry", "Heather", "Tyler", "Diane", "Aaron", "Ruth", "Jose", "Julie", "Adam", "Olivia", "Nathan", "Joyce", "Henry",
	"Virginia", "Douglas", "Victoria", "Zachary", "Kelly", "Peter", "Lauren", "Kyle", "Christina", "Ethan", "Joan", "Walter", "Evelyn", "Noah",
	"Judith", "Jeremy", "Megan", "Christian", "Andrea", "Keith", "Cheryl", "Roger", "Hannah", "Terry", "Jacqueline", "Gerald", "Martha", "Harold",
	"Gloria", "Sean", "Teresa", "Austin", "Ann", "Carl", "Sara", "Arthur", "Madison", "Lawrence", "Frances", "Dylan", "Kathryn", "Jesse",
	"Janice", "Jordan", "Jean", "Bryan", "Abigail", "Billy", "Alice", "Joe", "Julia", "Bruce", "Judy", "Gabriel", "Sophia", "Logan", "Grace",
	"Albert", "Denise", "Willie", "Amber", "Alan", "Doris", "Juan", "Marilyn", "Wayne", "Danielle", "Elijah", "Beverly", "Randy", "Isabella",
	"Roy", "Theresa", "Vincent", "Diana", "Ralph", "Natalie", "Eugene", "Brittany", "Russell", "Charlotte", "Bobby", "Marie", "Mason", "Kayla",
	"Philip", "Alexis", "Louis", "Lori"}
var lastNames = []string{"Abraham", "Allan", "Alsop", "Anderson", "Arnold", "Avery", "Bailey", "Baker", "Ball", "Bell", "Berry", "Black", "Blake", "Bond",
	"Bower", "Brown", "Buckland", "Burgess", "Butler", "Cameron", "Campbell", "Carr", "Chapman", "Churchill", "Clark", "Clarkson", "Coleman", "Cornish",
	"Davidson", "Davies", "Dickens", "Dowd", "Duncan", "Dyer", "Edmunds", "Ellison", "Ferguson", "Fisher", "Forsyth", "Fraser", "Gibson", "Gill", "Glover",
	"Graham", "Grant", "Gray", "Greene", "Hamilton", "Hardacre", "Harris", "Hart", "Hemmings", "Henderson", "Hill", "Hodges", "Howard", "Hudson", "Hughes",
	"Hunter", "Ince", "Jackson", "James", "Johnston", "Jones", "Kelly", "Kerr", "King", "Knox", "Lambert", "Langdon", "Lawrence", "Lee", "Lewis", "Lyman",
	"MacDonald", "Mackay", "Mackenzie", "MacLeod", "Manning", "Marshall", "Martin", "Mathis", "May", "McDonald", "McLean", "McGrath", "Metcalfe", "Miller",
	"Mills", "Mitchell", "Morgan", "Morrison", "Murray", "Nash", "Newman", "Nolan", "North", "Ogden", "Oliver", "Paige", "Parr", "Parsons", "Paterson",
	"Payne", "Peake", "Peters", "Piper", "Poole", "Powell", "Pullman", "Quinn", "Rampling", "Randall", "Rees", "Reid", "Roberts", "Robertson", "Ross", "Russell",
	"Rutherford", "Sanderson", "Scott", "Sharp", "Short", "Simpson", "Skinner", "Slater", "Smith", "Springer", "Stewart", "Sutherland", "Taylor", "Terry",
	"Thomson", "Tucker", "Turner", "Underwood", "Vance", "Vaughan", "Walker", "Wallace", "Walsh", "Watson", "Welch", "White", "Wilkins", "Wilson",
	"Wright", "Young"}
var cities = []string{"Bladensburg", "Brambleton", "Edenburg", "Dubois", "Cotopaxi", "Sperryville", "Alleghenyville", "Westboro", "Tonopah", "Fowlerville",
	"Venice", "Wanship", "Diaperville", "Haring", "Morriston", "Kenvil", "Dahlen", "Canby", "Basye", "Marienthal", "Sutton", "Elwood",
	"Tilleda", "Crenshaw", "Loveland", "Canoochee", "Newkirk", "National", "Chesterfield", "Draper", "Turah", "Hall", "Dragoon", "Summertown", "Sims",
	"Guthrie", "Vivian", "Tuttle", "Ladera", "Drummond", "Ezel", "Marne", "Lookingglass", "Shasta", "Vandiver", "Sharon", "Glendale", "Loomis",
	"Statenville", "Gouglersville", "Sehili", "Catherine", "Whitmer", "Grimsley", "Salix", "Kersey", "Springdale", "Thermal", "Witmer", "Virgie",
	"Wakulla", "Indio", "Unionville", "Loretto", "Sabillasville", "Gracey", "Blodgett", "Aguila", "Harleigh", "Avalon", "Fairview",
	"Esmont", "Cascades", "Cleary", "Reno", "Holtville", "Lumberton", "Keller", "Caspar", "Biddle", "Dexter", "Whitehaven", "Fidelis", "Drytown",
	"Dorneyville", "Rivereno", "Independence", "Bodega", "Wanamie", "Townsend", "Caron", "Guilford", "Gallina", "Manila", "Itmann", "Whitewater",
	"Templeton", "Jessie", "Sena", "Charco", "Jamestown", "Imperial", "Vincent", "Nelson", "Abrams", "Glasgow", "Lynn", "Sugartown", "Navarre",
	"Marion", "Sanders", "Spelter", "Santel", "Outlook", "Ypsilanti", "Dotsero", "Mathews", "Loyalhanna", "Libertytown", "Terlingua", "Hackneyville",
	"Driftwood", "Stockdale", "Bynum", "Harrison", "Morningside", "Churchill", "Gambrills", "Brule", "Fairhaven", "Hinsdale", "Babb", "Buxton",
	"Biehle", "Catharine", "Dunbar", "Klagetoh", "Blandburg", "Roberts", "Romeville", "Hachita", "Leming", "Saranap", "Elliott", "Ronco", "Rossmore",
	"Bowie", "Roderfield", "Devon", "Trucksville", "Ribera", "Watchtower", "Orason", "Haena", "Fruitdale", "Riceville", "Urbana", "Moscow",
	"Fulford", "Cassel", "Shawmut", "Corinne", "Edmund", "Naomi", "Clara", "Duryea", "Chloride", "Axis", "Villarreal", "Talpa", "Rodman", "Goochland",
	"Deercroft", "Jacksonburg", "Kanauga", "Springville", "Concho", "Matheny", "Temperanceville", "Salunga", "Elfrida", "Stollings", "Lindisfarne",
	"Kimmell", "Fillmore", "Belmont", "Mansfield", "Fairforest", "Finzel", "Shelby", "Brenton", "Fairlee", "Brownlee", "Yettem", "Richmond", "Jeff",
	"Umapine", "Cuylerville", "Carbonville", "Alamo"}

func InsertData[K comparable, V any](imap *IndexMap[K, V], data map[K]*V) {
	for _, v := range data {
		imap.Insert(v)
	}
}

func createRandomPerson(id int64, myRand *rand.Rand) *Person {
	p := Person{ID: int64(id),
		Name: names[myRand.Intn(200)],
		City: cities[myRand.Intn(200)],
		Age:  myRand.Intn(103),
	}
	var friends []string
	for range myRand.Intn(3) {
		friends = append(friends, names[myRand.Intn(200)])
	}
	p.Like = friends
	return &p
}

func InsertRandomData(imap *IndexMap[int64, Person], n int) {
	myRand := rand.New(rand.NewSource(123))
	for i := range n {
		p := createRandomPerson(int64(i), myRand)
		imap.Insert(p)
	}
}

func CreateTestMap(n int) *IndexMap[int64, Person] {
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	_ = imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))

	_ = imap.AddIndex(CityIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.City}
	}))

	InsertRandomData(imap, n)
	return imap
}

func TestCreateJsons(t *testing.T) {
	myRand := rand.New(rand.NewSource(123))
	for i := range 100 {
		p := createRandomPerson(int64(i), myRand)
		jsonByte, _ := json.Marshal(p)
		os.WriteFile("jsonfiles/person_"+strconv.Itoa(i)+".json", jsonByte, 0666)
	}
}
