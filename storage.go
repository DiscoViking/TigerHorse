// This file provides an interface to interact directly with the underlying storage.
package main

type Storage interface {
	AddPerson(name string) error
	AddTransaction(t Transaction) error
	GetTransactions() ([]Transaction, error)
	GetPeople() ([]Person, error)
}
