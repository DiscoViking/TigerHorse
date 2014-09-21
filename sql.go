// Storage controller for a mySQL database.
package main

import (
	"database/sql"
	"time"
)

type Controller struct {
	db *sql.DB
}

// Opens a connection to the mySQL database and confirms success.
func New(source string) (*Controller, error) {
	db, err := sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Controller{db}, nil
}

// Adds a person to the database by name.
func (c *Controller) AddPerson(name string) error {
	_, err := c.db.Exec("INSERT INTO people (name) VALUES (?)", name)

	return err
}

// Adds a transaction to the database.
func (c *Controller) AddTransaction(t Transaction) error {
	var err error

	// We need to make multiple insertions atomically, so start a transaction.
	tx, err := c.db.Begin()

	// Make sure we abort the transaction if we exit with an error.
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert transaction.
	r, err := tx.Exec(
		"INSERT INTO transactions (amount,buyer,guests,time,note)"+
			"VALUES (?,?,?,?,?)",
		t.Value, t.Buyer, t.Guests, t.Time, t.Note)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	// Insert links to beneficiaries.
	for _, p := range t.Involved {
		_, err := c.db.Exec("INSERT INTO link VALUES (?,?)", id, p)
		if err != nil {
			return err
		}
	}

	// Commit the transaction to apply all changes atomically.
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Gets all the transactions from the database.
func (c *Controller) GetTransactions() ([]Transaction, error) {
	// Firstly, get the actual transactions.
	rows, err := c.db.Query("SELECT id,amount,buyer,guests,note,time FROM transactions")
	if err != nil {
		return nil, err
	}

	var txs []Transaction
	for rows.Next() {
		var id int64
		var amount int64
		var buyer int64
		var guests int64
		var note string
		var time time.Time
		err := rows.Scan(&id, &amount, &buyer, &guests, &note, &time)
		if err != nil {
			return nil, err
		}

		txs = append(txs, Transaction{
			Id:       id,
			Value:    amount,
			Buyer:    buyer,
			Involved: []int64{},
			Guests:   0,
			Note:     note,
			Time:     time,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	// Now collect all the information about beneficiaries for each transaction.
	for _, tx := range txs {
		rows, err := c.db.Query("SELECT person FROM link WHERE transaction=?", tx.Id)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var id int64
			err := rows.Scan(&id)
			if err != nil {
				return nil, err
			}

			tx.Involved = append(tx.Involved, id)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if err := rows.Close(); err != nil {
			return nil, err
		}
	}

	return txs, nil
}

// Gets all the people from the database.
func (c *Controller) GetPeople() ([]Person, error) {
	people := []Person{}

	rows, err := c.db.Query("SELECT id,name FROM people")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		people = append(people, Person{id, name})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	return people, nil
}
