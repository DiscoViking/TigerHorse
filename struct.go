package main

// Types and methods for transactions.

import (
	"time"
)

// One monetary transaction.
type Transaction struct {
	Id       int64     // Unique identifier.
	Amount   int64     // Amount spent in pennies.
	Buyer    int64     // Id of the person who paid.
	Involved []int64   // Ids of all people who shared in this transaction.
	Guests   int64     // Number of non-spreadsheet people who partook.
	Note     string    // Short description of what it was for.
	Time     time.Time // Time this transaction was added to the system.
}

type Person struct {
	Id   int64  // Unique identifier.
	Name string // Name of person.
}
