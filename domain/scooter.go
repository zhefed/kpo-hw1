package domain

// Scooter - самокат
type Scooter struct {
	baseVehicle
	energy float64
}

func NewScooter(name string, simplicity int, energy float64) *Scooter {
	s := &Scooter{}
	s.name = name
	s.simplicity = simplicity
	s.energy = energy
	return s
}

func (s *Scooter) Kind() string {
	return "Самокат"
}

func (s *Scooter) DailyEnergy() float64 {
	return s.energy
}

var _ EnergyConsumer = (*Scooter)(nil)

func init() {
	RegisterVehicleKind("scooter", func(name string, simplicity int, energy float64) Vehicle {
		return NewScooter(name, simplicity, energy)
	})
}
