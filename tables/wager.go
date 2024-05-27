package tables

type Wager struct {
	Hand      Hand // The hand associated with the wager
	AmountBet int  // The amount of the initial bet
	DoubleBet int  // The amount of the double bet
	AmountWon int  // The amount won from the wager
}

func NewWager() *Wager {
	return new(Wager)
}

func (w *Wager) SplitWager(split *Wager) {
	split.AmountBet = w.AmountBet
	split.Hand.Draw(w.Hand.SplitPair())
}

func (w *Wager) Reset() {
	w.Hand.Reset()
	w.AmountBet = 0
	w.DoubleBet = 0
	w.AmountWon = 0
}

func (w *Wager) Bet(b int) {
	if b%2 != 0 {
		panic("All bets must be in multiples of 2.")
	}
	w.AmountBet = b
}

func (w *Wager) Double() {
	w.DoubleBet = w.AmountBet
}

func (w *Wager) Blackjack() bool {
	return w.Hand.Blackjack()
}

func (w *Wager) WonBlackjack() {
	w.AmountWon += (w.AmountBet * 3) / 2
}

func (w *Wager) Won() {
	w.AmountWon += w.AmountBet + w.DoubleBet
}

func (w *Wager) Lost() {
	w.AmountWon -= (w.AmountBet + w.DoubleBet)
}
