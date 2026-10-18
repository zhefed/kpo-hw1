package service_test

import (
	"math"
	"reflect"
	"testing"

	"hw1/domain"
	"hw1/service"
)

func TestTotalEnergy_EmptyParkIsZero(t *testing.T) {
	park := acceptingPark()

	total := park.TotalEnergy()

	if total != 0 {
		t.Fatalf("ожидали 0, получили %v", total)
	}
}

func TestTotalEnergy_SumsOnlyEnergyConsumers(t *testing.T) {
	park := acceptingPark()
	_ = park.AddVehicle(domain.NewScooter("Ninebot", 7, 0.5))
	_ = park.AddVehicle(domain.NewElectricBike("Eltreco", 6, 1.0))
	_ = park.AddVehicle(domain.NewBicycle("Stels", 9))
	_ = park.AddThing(domain.NewHelmet("Шлем", "M"))
	_ = park.AddThing(domain.NewDockStation("Док", 1.2))
	_ = park.AddThing(domain.NewChargingCabinet("Шкаф", 2.5))

	total := park.TotalEnergy()

	if math.Abs(total-5.2) > 1e-9 {
		t.Fatalf("ожидали 5.2, получили %v", total)
	}
}

func TestBeginnerVehicles_SimplicityThreshold(t *testing.T) {
	park := acceptingPark()
	tooHard := domain.NewScooter("Kugoo", 5, 0.6)
	justEnough := domain.NewElectricBike("Eltreco", 6, 0.8)
	_ = park.AddVehicle(tooHard)
	_ = park.AddVehicle(justEnough)

	got := park.BeginnerVehicles()

	want := []domain.Vehicle{justEnough}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ожидали %v, получили %v", want, got)
	}
}

func TestVehicleCountByKind_CountsInOrderOfAppearance(t *testing.T) {
	park := acceptingPark()
	_ = park.AddVehicle(domain.NewScooter("Ninebot", 7, 0.35))
	_ = park.AddVehicle(domain.NewBicycle("Stels", 9))
	_ = park.AddVehicle(domain.NewScooter("Kugoo", 4, 0.6))

	got := park.VehicleCountByKind()

	want := []service.KindCount{
		{Kind: "Самокат", Count: 2},
		{Kind: "Велосипед", Count: 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ожидали %v, получили %v", want, got)
	}
}
