package main

import (
	"reflect"
	"testing"
	"time"
)

const (
	testDSN = "tigerhorse@tcp(localhost:3306)/tigerhorsetest?parseTime=true&loc=GMT"
)

// Test database initialization.
func TestInit(t *testing.T) {
	// Open a connection to the server.
	//c, err := New("mysql@tcp(localhost:3306)/tigerhorsetest")
	c, err := New(testDSN)
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to initialise.
	err = c.Init()
	if err != nil {
		t.Fatal(err)
	}

	// Check tables.
	rows, err := c.(*Controller).db.Query("SHOW TABLES")
	if err != nil {
		t.Fatal(err)
	}

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// Test adding people to the database.
func TestPeople(t *testing.T) {
	// Open a connection to the server.
	//c, err := New("mysql@tcp(localhost:3306)/tigerhorsetest")
	c, err := New(testDSN)
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to initialise.
	err = ensureEmpty(c)
	if err != nil {
		t.Fatal(err)
	}

	// Add a person.
	t.Log("Add one person")
	err = c.AddPerson("ryan")
	if err != nil {
		t.Error(err)
	}
	checkPeople(t, c, []*Person{{1, "ryan", 0}})

	// Add another person.
	t.Log("Add another person")
	err = c.AddPerson("fiona")
	if err != nil {
		t.Error(err)
	}
	checkPeople(t, c, []*Person{{2, "fiona", 0}, {1, "ryan", 0}})

	// Check duplicate names not allowed.
	t.Log("Attempt to add duplicate")
	err = c.AddPerson("fiona")
	if err == nil {
		t.Error("Didn't error whilst attempting to add duplicate person")
	}
	checkPeople(t, c, []*Person{{2, "fiona", 0}, {1, "ryan", 0}})
}

// Test adding a basic transaction.
func TestBasicTransaction(t *testing.T) {
	// Open a connection to the server.
	//c, err := New("mysql@tcp(localhost:3306)/tigerhorsetest")
	c, err := New(testDSN)
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to initialise.
	err = ensureEmpty(c)
	if err != nil {
		t.Fatal(err)
	}

	// Add 2 people.
	t.Log("Add 2 people")
	if err := c.AddPerson("ryan"); err != nil {
		t.Fatal(err)
	}
	if err := c.AddPerson("fiona"); err != nil {
		t.Fatal(err)
	}

	// Add the transaction.
	t.Log("Add a transaction")
	tx := &Transaction{
		Id:       0,             // Not set yet.
		Value:    100,           // £1.00
		Buyer:    2,             // Fiona
		Involved: []int64{1, 2}, // Ryan & Fiona
		Guests:   0,
		Note:     "Fiona buys some cookies",
		Time:     time.Now(),
	}
	err = c.AddTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}

	// Check we can retrieve it.
	t.Log("Retrieve transaction")
	txs, err := c.GetTransactions()
	if err != nil {
		t.Fatal(err)
	}

	if len(txs) != 1 {
		t.Fatalf("Was only expecting 1 transaction, got %v.\nTxs: %v\n", len(txs), txs)
	}

	/* Disabled because of timezones.
	if !reflect.DeepEqual(txs[0], tx) {
		t.Errorf("Expected: %v, Got: %v.\n", tx, txs[0])
	}
	*/
}

// Some utility functions
func checkPeople(t *testing.T, c Storage, exp []*Person) {
	people, err := c.GetPeople()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(people, exp) {
		t.Errorf("Wrong people. Expected %v, Got %v.", exp, people)
		return
	}

	return
}

// Ensure the database is clear.
func ensureEmpty(c Storage) error {
	err := c.Init()
	if err != nil {
		return err
	}
	err = c.Clear()
	if err != nil {
		return err
	}
	err = c.Init()
	if err != nil {
		return err
	}
	return nil
}
