package cards

type Wager struct {
	Hand         Hand  // The hand associated with the wager
	AmountBet    int64 // The amount of the initial bet
	AmountWon    int64 // The amount won from the wager
	DoubleBet    int64 // The amount of the double bet
	DoubleWon    int64 // The amount won from the wager
	InsuranceBet int64
	InsuranceWon int64
}

const (
	MinimumBet = 2
	MaximumBet = 98
)

func NewWager() *Wager {
	return new(Wager)
}

func (w *Wager) SplitWager(split *Wager) {
	split.Reset()
	split.AmountBet = w.AmountBet
	split.Hand.Draw(w.Hand.SplitPair())
}

func (w *Wager) Reset() {
	w.Hand.Reset()
	w.AmountBet = 0
	w.AmountWon = 0
	w.DoubleBet = 0
	w.DoubleWon = 0
	w.InsuranceBet = 0
	w.InsuranceWon = 0
}

func (w *Wager) Bet(bet int64) {
	w.AmountBet = (ClampInt(bet, MinimumBet, MaximumBet) + 1) / 2 * 2
}

func (w *Wager) Double() {
	w.DoubleBet = w.AmountBet
}

func (w *Wager) Blackjack() bool {
	return w.Hand.Blackjack()
}

func (w *Wager) WonBlackjack(pays, bet int64) {
	w.AmountWon = int64((w.AmountBet * pays) / bet)
}

func (w *Wager) Won() {
	w.AmountWon = w.AmountBet
	w.DoubleWon = w.DoubleBet
}

func (w *Wager) Lost() {
	w.AmountWon = -w.AmountBet
	w.DoubleWon = -w.DoubleBet
}

func (w *Wager) Push() {
}

func (w *Wager) WonInsurance() {
	w.InsuranceWon = w.InsuranceBet * 2
}

func (w *Wager) LostInsurance() {
	w.InsuranceWon = -w.InsuranceBet
}

func ClampInt(value, min, max int64) int64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

