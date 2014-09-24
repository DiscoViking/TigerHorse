package main

// A dummy storage to use for testing, just keeps lists of stuff.

type dummyStorage struct {
	people  []*Person
	txs     []*Transaction
	pId     int64
	tId     int64
	nextErr error
}

func newDummyStorage() Storage {
	return &dummyStorage{
		people:  []*Person{},
		txs:     []*Transaction{},
		pId:     1,
		tId:     1,
		nextErr: nil,
	}
}

func (s *dummyStorage) AddPerson(name string) error {
	// Error if set.
	if s.nextErr != nil {
		err := s.nextErr
		s.nextErr = nil
		return err
	}

	// Add the person.
	s.people = append(s.people, &Person{s.pId, name})

	// Increment person Id ready for next.
	s.pId++

	return nil
}

func (s *dummyStorage) AddTransaction(t *Transaction) error {
	// Error if set.
	if s.nextErr != nil {
		err := s.nextErr
		s.nextErr = nil
		return err
	}

	// Add the transaction.
	// We have to add a copy since we're dealing with pointers.
	t.Id = s.tId
	s.txs = append(s.txs, &Transaction{
		t.Id,
		t.Value,
		t.Buyer,
		t.Involved, // Should really make a copy of this too, I'll do it when it bites me.
		t.Guests,
		t.Note,
		t.Time,
	})

	// Increment transaction Id ready for next.
	s.tId++

	return nil
}

func (s *dummyStorage) GetPeople() ([]*Person, error) {
	// Error if set.
	if s.nextErr != nil {
		err := s.nextErr
		s.nextErr = nil
		return nil, err
	}

	// Return the people list.
	return s.people, nil
}

func (s *dummyStorage) GetTransactions() ([]*Transaction, error) {
	// Error if set.
	if s.nextErr != nil {
		err := s.nextErr
		s.nextErr = nil
		return nil, err
	}

	// Return the people list.
	return s.txs, nil
}

func (s *dummyStorage) Init() error {
	// Error if set.
	if s.nextErr != nil {
		err := s.nextErr
		s.nextErr = nil
		return err
	}
	return nil
}
func (s *dummyStorage) Clear() error {
	// Error if set.
	if s.nextErr != nil {
		err := s.nextErr
		s.nextErr = nil
		return err
	}

	s.people = []*Person{}
	s.txs = []*Transaction{}
	s.pId = 1
	s.tId = 1

	return nil
}
