package demo

import "slices"

// Sortable lets you add different kind of events in one slice
// and sort based on start tick
type Sortable interface {
	GetTick() Tick
}

func Sort(t []Sortable) {
	slices.SortFunc(t, func(a, b Sortable) int {
		return int(a.GetTick()) - int(b.GetTick())
	})
}

func (b *BombDrop) GetTick() Tick {
	return b.DropTick
}

func (b *BombPlant) GetTick() Tick {
	return b.Start
}

func (b *BombDefuse) GetTick() Tick {
	return b.Start
}

func (h *HostageCarry) GetTick() Tick {
	return h.Start
}

func (k Kill) GetTick() Tick {
	return k.Tick
}

func (a Assist) GetTick() Tick {
	return a.Tick
}

func (d Death) GetTick() Tick {
	return d.Tick
}

func (f *Flash) GetTick() Tick {
	return f.Start
}

func (s *Smoke) GetTick() Tick {
	return s.Start
}

func (h *He) GetTick() Tick {
	return h.Start
}

func (d *Decoy) GetTick() Tick {
	return d.Start
}

func (i *Incendiary) GetTick() Tick {
	return i.Start
}

func (m *Molotov) GetTick() Tick {
	return m.Start
}

func (i *ItemPurchase) GetTick() Tick {
	return i.Tick
}

func (i *ItemRefund) GetTick() Tick {
	return i.Tick
}

func (i *ItemDrop) GetTick() Tick {
	return i.Tick
}

func (i *ItemPickup) GetTick() Tick {
	return i.Tick
}

func (s *Spot) GetTick() Tick {
	return s.Start
}

func (s *SpottedBy) GetTick() Tick {
	return s.Start
}

func (r Reload) GetTick() Tick {
	return r.Tick
}

func (s Shot) GetTick() Tick {
	return s.Tick
}

func (m Message) GetTick() Tick {
	return m.Tick
}

func (c Chicken) GetTick() Tick {
	return c.Tick
}
