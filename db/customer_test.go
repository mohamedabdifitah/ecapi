package db

import (
	"fmt"
	"testing"
)

func TestSave(t *testing.T) {
	C := Customer{
		Email:      "example@example.com",
		FamilyName: "familyname",
		GivenName:  "uniqueGivenName",
		Password:   "@#$%^&*()",
	}
	res, err := C.Save()
	t.Error(err)
	fmt.Println(res)

}
func TestGetById(t *testing.T) {
	customer := Customer{}
	customer.GetById()
}
func TestGetAll(t *testing.T) {

}
