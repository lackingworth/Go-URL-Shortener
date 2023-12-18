package database

import "testing"

func TestCreateClient(t *testing.T) {
	address := "127.0.0.1:6379"
	r1 := CreateClient(1, address)
	r3 := CreateClient(3, address)
	
	if r1 == nil{
		t.Errorf("Redis not connected")
	}
	
	if r1 == nil{
		t.Errorf("Redis not connected")
	}
	
	if r1 == r3{
		t.Errorf("Different instances must be unique")
	}
}