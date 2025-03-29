package second_faza

type Transaction struct {
	id int
}

func NewTransaction(id int) *Transaction {
	return &Transaction{
		id: id,
	}
}

func (tx *Transaction) Equals(other *Transaction) bool {
	if other == nil {
		return false
	}

	return tx.id == other.id
}

func (tx *Transaction) HashCode() int {
	return tx.id
}
