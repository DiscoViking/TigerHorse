package main

// Gets data and performs calculations.

// Gets only the transactions of the given person.
func (p *Person) Transactions(s Storage) ([]*Transaction, error) {
	txs := []*Transaction{}
	allTxs, err := s.GetTransactions()
	if err != nil {
		return nil, err
	}

	for _, tx := range allTxs {
		involved := false
		if tx.Buyer == p.Id {
			involved = true
		}
		for _, id := range tx.Involved {
			if id == p.Id {
				involved = true
				break
			}
		}
		if involved {
			txs = append(txs, tx)
		}
	}

	return txs, nil
}

// Returns total money spent in pennies by person.
func (p *Person) Spent(s Storage) (int64, error) {
	txs, err := s.GetTransactions()
	if err != nil {
		return 0, err
	}

	var spent int64 = 0
	for _, tx := range txs {
		if tx.Buyer == p.Id {
			spent += tx.Value
		}
	}

	return spent, nil
}

// Returns the total value a person has gained from the system.
// To calculate this, we take every transaction the person benefited from,
// divide the value of the transaction by how many people took part in it,
// and sum these individual values.
// Return value is in pennies.
func (p *Person) Received(s Storage) (int64, error) {
	txs, err := s.GetTransactions()
	if err != nil {
		return 0, err
	}

	var value int64 = 0
	for _, tx := range txs {
		// Work out if we gained value from the transaction.
		involved := false
		for _, id := range tx.Involved {
			if id == p.Id {
				involved = true
				break
			}
		}
		if !involved {
			continue
		}

		// Total people who benefited from the transaction.
		n := int64(len(tx.Involved)) + tx.Guests

		// Increase total value.
		value += tx.Value / n
	}

	return value, nil
}
