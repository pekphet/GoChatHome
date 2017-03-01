package db

import (
	"fmt"
	"testing"
)

func TestGetEquip(t *testing.T) {
	fmt.Println(GetEquip(1))
	fmt.Println(GetEquip(2))
	fmt.Println(GetEquip(1))
	fmt.Println(GetEquip(2))
	fmt.Println(GetEquip(2))
	fmt.Println(GetEquip(2))
}
