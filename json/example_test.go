package json_test

import (
	"encoding/json"
	"fmt"
)

type Query struct {
	Name *string `json:"name"`
	Age  *string `json:"age"`
}

func Example_marshal() {
	bytes, err := json.Marshal(Query{})
	fmt.Println(string(bytes), err)
	bytes, err = json.Marshal(&Query{})
	fmt.Println(string(bytes), err)
	// Output:
	// {"name":null,"age":null} <nil>
	// {"name":null,"age":null} <nil>
}

func Example_unmarshal1() {
	str := `{"name":"Name"}`
	obj := Query{}
	err := json.Unmarshal([]byte(str), &obj)
	fmt.Println(*obj.Name, obj.Age, err)
	// Output:
	// Name <nil> <nil>
}

func Example_unmarshal2() {
	str := `{}`
	obj := Query{}
	err := json.Unmarshal([]byte(str), &obj)
	fmt.Println(obj.Name, obj.Age, err)
	// Output:
	// <nil> <nil> <nil>
}
