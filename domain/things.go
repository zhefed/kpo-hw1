package domain

// Thing - вещь на балансе
type Thing interface {
	Inventoried
	Category() string
}

// Helmet - шлем. Энергию не потребляет
type Helmet struct {
	item
	size string
}

func NewHelmet(name string, size string) *Helmet {
	h := &Helmet{}
	h.name = name
	h.size = size
	return h
}

func (h *Helmet) Category() string {
	return "Шлем (размер " + h.size + ")"
}

// DockStation - док-станция
type DockStation struct {
	item
	energy float64
}

func NewDockStation(name string, energy float64) *DockStation {
	d := &DockStation{}
	d.name = name
	d.energy = energy
	return d
}

func (d *DockStation) Category() string {
	return "Док-станция"
}

func (d *DockStation) DailyEnergy() float64 {
	return d.energy
}

// ChargingCabinet - зарядный шкаф
type ChargingCabinet struct {
	item
	energy float64
}

func NewChargingCabinet(name string, energy float64) *ChargingCabinet {
	c := &ChargingCabinet{}
	c.name = name
	c.energy = energy
	return c
}

func (c *ChargingCabinet) Category() string {
	return "Зарядный шкаф"
}

func (c *ChargingCabinet) DailyEnergy() float64 {
	return c.energy
}
