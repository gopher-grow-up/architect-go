package hashmap

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	m := &HashMap{}
	m.Set("hello", "world")
	get, _ := m.GetStringKey("hello")
	fmt.Println(get.(string))
}
