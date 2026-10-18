package main

import (
	"fmt"
	"os"

	"go.uber.org/dig"

	"hw1/console"
	"hw1/service"
)

const maxDailyEnergy = 3.0

func main() {
	container, err := buildContainer()
	if err != nil {
		fmt.Println("Ошибка сборки зависимостей:", err)
		os.Exit(1)
	}

	// dig сам создаст App и всё, что ему нужно
	err = container.Invoke(func(app *console.App) {
		app.Run()
	})
	if err != nil {
		fmt.Println("Ошибка запуска:", err)
		os.Exit(1)
	}
}

func newServiceCenter() service.ServiceCenter {
	return service.NewBasicServiceCenter(maxDailyEnergy)
}

func newApp(park *service.Park) *console.App {
	return console.NewApp(park, os.Stdin)
}

func buildContainer() (*dig.Container, error) {
	container := dig.New()

	constructors := []interface{}{
		newServiceCenter,
		service.NewPark,
		newApp,
	}
	for _, constructor := range constructors {
		err := container.Provide(constructor)
		if err != nil {
			return nil, err
		}
	}
	return container, nil
}
