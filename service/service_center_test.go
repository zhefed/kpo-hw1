package service_test

import (
	"testing"

	"hw1/domain"
	"hw1/service"
)

const testMaxDailyEnergy = 3.0

func TestInspect_AcceptsValidVehicle(t *testing.T) {
	tests := []struct {
		name    string
		vehicle domain.Vehicle
	}{
		{"минимальная простота", domain.NewScooter("Ninebot", 1, 0.5)},
		{"максимальная простота", domain.NewScooter("Ninebot", 10, 0.5)},
		{"нулевой расход", domain.NewElectricBike("Eltreco", 6, 0)},
		{"расход ровно на пределе", domain.NewElectricBike("Eltreco", 6, testMaxDailyEnergy)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			center := service.NewBasicServiceCenter(testMaxDailyEnergy)

			ok, reason := center.Inspect(tt.vehicle)

			if !ok {
				t.Fatalf("ожидали, что техосмотр пройден, отказ: %q", reason)
			}
		})
	}
}

func TestInspect_RejectsInvalidVehicle(t *testing.T) {
	tests := []struct {
		name    string
		vehicle domain.Vehicle
	}{
		{"пустое имя", domain.NewBicycle("", 5)},
		{"простота меньше 1", domain.NewBicycle("Stels", 0)},
		{"простота больше 10", domain.NewBicycle("Stels", 11)},
		{"отрицательный расход", domain.NewScooter("Ninebot", 5, -0.01)},
		{"расход больше предела", domain.NewScooter("Ninebot", 5, 3.01)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			center := service.NewBasicServiceCenter(testMaxDailyEnergy)

			ok, reason := center.Inspect(tt.vehicle)

			if ok {
				t.Fatal("ожидали отказ техосмотра")
			}
			if reason == "" {
				t.Fatal("у отказа должна быть причина")
			}
		})
	}
}
