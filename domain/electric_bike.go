package domain

// ElectricBike - электровелосипед
type ElectricBike struct {
	baseVehicle
	energy float64
}

func NewElectricBike(name string, simplicity int, energy float64) *ElectricBike {
	e := &ElectricBike{}
	e.name = name
	e.simplicity = simplicity
	e.energy = energy
	return e
}

func (e *ElectricBike) Kind() string {
	return "Электровелосипед"
}

func (e *ElectricBike) DailyEnergy() float64 {
	return e.energy
}

var _ EnergyConsumer = (*ElectricBike)(nil)

func init() {
	RegisterVehicleKind("electric-bike", func(name string, simplicity int, energy float64) Vehicle {
		return NewElectricBike(name, simplicity, energy)
	})
}
