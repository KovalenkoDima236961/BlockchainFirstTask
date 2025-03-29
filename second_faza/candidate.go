package second_faza

type Candidate struct {
	tx     *Transaction
	sender int
}

func NewCandidate(tx *Transaction, sender int) *Candidate {
	return &Candidate{tx: tx, sender: sender}
}
