package main

// Tests for the calculation module.

import (
	"testing"
	"time"
)

func TestSpent(t *testing.T) {
	s := newDummyStorage()
	s.AddPerson("ryan")
	s.AddPerson("fiona")

	// Add 3 transactions.
	tx := &Transaction{
		Id:       0,
		Value:    100,
		Buyer:    1,
		Involved: []int64{1, 2},
		Guests:   0,
		Note:     "",
		Time:     time.Now(),
	}
	s.AddTransaction(tx)
	tx.Value = 250
	s.AddTransaction(tx)
	tx.Value = 0
	s.AddTransaction(tx)

	// Check the Spent value is correct.
	ryan := &Person{1, "ryan", 0}
	fiona := &Person{2, "fiona", 0}

	if got, _ := ryan.Spent(s); got != 350 {
		t.Errorf("Ryan spent %v but we calculated %v!\n", 350, got)
	}
	if got, _ := fiona.Spent(s); got != 0 {
		t.Errorf("Fiona spent %v but we calculated %v!\n", 0, got)
	}
}

// Test we correctly calculate received value.
func TestReceived(t *testing.T) {
	s := newDummyStorage()
	s.AddPerson("ryan")
	s.AddPerson("fiona")

	// Add 3 transactions.
	tx := &Transaction{
		Id:       0,
		Value:    100,
		Buyer:    1,
		Involved: []int64{1, 2},
		Guests:   0,
		Note:     "",
		Time:     time.Now(),
	}
	s.AddTransaction(tx)
	tx.Value = 250
	tx.Involved = []int64{2} // This is ryan giving money directly to fiona.
	s.AddTransaction(tx)

	// Check the received value is correct.
	ryan := &Person{1, "ryan", 0}
	fiona := &Person{2, "fiona", 0}

	if got, _ := ryan.Received(s); got != 50 {
		t.Errorf("Ryan received %v but we calculated %v!\n", 50, got)
	}
	if got, _ := fiona.Received(s); got != 300 {
		t.Errorf("Fiona received %v but we calculated %v!\n", 300, got)
	}
}
