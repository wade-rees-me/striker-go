package cards

type Dealer struct {
	Hand	  Hand
	HitSoft17 bool
}

func NewDealer(hitSoft17 bool) *Dealer {
	d := new(Dealer)
	d.HitSoft17 = hitSoft17
	return d
}

func (d *Dealer) Reset() {
	d.Hand.Reset()
}

func (d *Dealer) Stand() bool {
	if d.HitSoft17 && d.Hand.Soft17() {
		return false
	}
	return d.Hand.Total() >= 17
}

func (d *Dealer) Draw(c *Card) *Card {
	return d.Hand.Draw(c)
}
