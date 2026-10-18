package console

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"hw1/domain"
	"hw1/service"
)

// App - консольный интерфейс учёта парка. Только ввод и вывод,
// вся логика находится в service.Park.
type App struct {
	park    *service.Park
	scanner *bufio.Scanner
}

func NewApp(park *service.Park, input io.Reader) *App {
	return &App{
		park:    park,
		scanner: bufio.NewScanner(input),
	}
}

func (a *App) Run() {
	fmt.Println("=== ВышКат: учёт парка. «Доедем до дедлайна» ===")
	for {
		a.printMenu()
		choice, ok := a.readLine("Выберите пункт: ")
		if !ok {
			fmt.Println()
			fmt.Println("Ввод закончился, выходим.")
			return
		}
		fmt.Println()

		switch choice {
		case "1":
			a.addVehicle()
		case "2":
			a.addThing()
		case "3":
			a.showReport()
		case "4":
			a.showBeginners()
		case "5":
			a.showInventory()
		case "6":
			a.runDemo()
		case "0":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Нет такого пункта, попробуйте ещё раз.")
		}
		fmt.Println()
	}
}

func (a *App) printMenu() {
	fmt.Println("-----------------------------------------")
	fmt.Println("1. Принять транспорт (с техосмотром)")
	fmt.Println("2. Поставить вещь на баланс")
	fmt.Println("3. Отчёт: количество и энергопотребление")
	fmt.Println("4. Транспорт для новичков")
	fmt.Println("5. Все инвентарные номера")
	fmt.Println("6. Демо-сценарий (заполнить парк примерами)")
	fmt.Println("0. Выход")
}

// readLine печатает подсказку и читает строку. ok == false, если ввод закончился.
func (a *App) readLine(prompt string) (string, bool) {
	fmt.Print(prompt)
	if !a.scanner.Scan() {
		return "", false
	}
	return strings.TrimSpace(a.scanner.Text()), true
}

// readInt читает целое число, пока пользователь не введёт корректное.
func (a *App) readInt(prompt string) (int, bool) {
	for {
		text, ok := a.readLine(prompt)
		if !ok {
			return 0, false
		}
		number, err := strconv.Atoi(text)
		if err == nil {
			return number, true
		}
		fmt.Println("Нужно ввести целое число.")
	}
}

// readFloat читает дробное число, можно писать и через точку, и через запятую.
func (a *App) readFloat(prompt string) (float64, bool) {
	for {
		text, ok := a.readLine(prompt)
		if !ok {
			return 0, false
		}
		text = strings.ReplaceAll(text, ",", ".")
		number, err := strconv.ParseFloat(text, 64)
		if err == nil {
			return number, true
		}
		fmt.Println("Нужно ввести число, например 0.5")
	}
}

func (a *App) addVehicle() {
	kinds := domain.VehicleKinds()
	fmt.Println("Виды транспорта:")
	for i, kind := range kinds {
		fmt.Printf("  %d. %s\n", i+1, kind.HumanName)
	}
	index, ok := a.readInt("Номер вида: ")
	if !ok {
		return
	}
	if index < 1 || index > len(kinds) {
		fmt.Println("Такого вида транспорта нет.")
		return
	}
	kind := kinds[index-1]

	name, ok := a.readLine("Наименование (модель): ")
	if !ok {
		return
	}
	simplicity, ok := a.readInt("Простота для новичка (1-10): ")
	if !ok {
		return
	}
	energy := 0.0
	if kind.UsesEnergy {
		energy, ok = a.readFloat("Суточный расход энергии, кВт·ч: ")
		if !ok {
			return
		}
	}

	a.acceptVehicle(kind.Create(name, simplicity, energy))
}

// acceptVehicle отправляет транспорт в парк и печатает результат техосмотра.
func (a *App) acceptVehicle(v domain.Vehicle) {
	err := a.park.AddVehicle(v)
	if err != nil {
		fmt.Printf("ОТКЛОНЁН: %s «%s» - %s\n", v.Kind(), v.Name(), err)
		return
	}
	fmt.Printf("ПРИНЯТ: %s «%s», инв. № %s\n", v.Kind(), v.Name(), formatNumber(v.InventoryNumber()))
}

func (a *App) addThing() {
	fmt.Println("Виды вещей:")
	fmt.Println("  1. Шлем")
	fmt.Println("  2. Док-станция")
	fmt.Println("  3. Зарядный шкаф")
	choice, ok := a.readLine("Номер вида: ")
	if !ok {
		return
	}
	if !slices.Contains([]string{"1", "2", "3"}, choice) {
		fmt.Println("Такого вида вещей нет.")
		return
	}

	name, ok := a.readLine("Наименование: ")
	if !ok {
		return
	}

	var thing domain.Thing
	switch choice {
	case "1":
		size, ok := a.readLine("Размер (S/M/L): ")
		if !ok {
			return
		}
		thing = domain.NewHelmet(name, size)
	case "2":
		energy, ok := a.readFloat("Суточный расход энергии, кВт·ч: ")
		if !ok {
			return
		}
		thing = domain.NewDockStation(name, energy)
	case "3":
		energy, ok := a.readFloat("Суточный расход энергии, кВт·ч: ")
		if !ok {
			return
		}
		thing = domain.NewChargingCabinet(name, energy)
	}

	a.acceptThing(thing)
}

func (a *App) acceptThing(t domain.Thing) {
	err := a.park.AddThing(t)
	if err != nil {
		fmt.Printf("ОТКЛОНЕНО: %s «%s» - %s\n", t.Category(), t.Name(), err)
		return
	}
	fmt.Printf("ПОСТАВЛЕНО НА БАЛАНС: %s «%s», инв. № %s\n", t.Category(), t.Name(), formatNumber(t.InventoryNumber()))
}

func (a *App) showReport() {
	vehicles := a.park.Vehicles()
	things := a.park.Things()

	fmt.Println("=== Отчёт по парку ===")
	fmt.Printf("Транспорт: %d ед.\n", len(vehicles))
	for _, kc := range a.park.VehicleCountByKind() {
		fmt.Printf("  - %s: %d\n", kc.Kind, kc.Count)
	}

	fmt.Printf("Вещи: %d ед.\n", len(things))
	fmt.Printf("Всего на балансе: %d ед.\n", len(vehicles)+len(things))
	fmt.Printf("Суммарное суточное энергопотребление: %.2f кВт·ч\n", a.park.TotalEnergy())
}

func (a *App) showBeginners() {
	list := a.park.BeginnerVehicles()
	fmt.Printf("=== Транспорт для новичков (простота >= %d) ===\n", service.MinBeginnerSimplicity)
	if len(list) == 0 {
		fmt.Println("Подходящего транспорта нет.")
		return
	}
	for _, v := range list {
		fmt.Printf("%s  %-18s %-25s простота %d\n",
			formatNumber(v.InventoryNumber()), v.Kind(), v.Name(), v.Simplicity())
	}
}

func (a *App) showInventory() {
	fmt.Println("=== Инвентарные номера ===")
	all := a.park.AllInventory()
	if len(all) == 0 {
		fmt.Println("На балансе пока ничего нет.")
		return
	}
	for _, obj := range all {
		fmt.Printf("%s  %-28s %s\n", formatNumber(obj.InventoryNumber()), describe(obj), obj.Name())
	}
}

// describe возвращает вид объекта: для транспорта - вид, для вещи - категорию.
func describe(obj domain.Inventoried) string {
	if v, ok := obj.(domain.Vehicle); ok {
		return v.Kind()
	}
	if t, ok := obj.(domain.Thing); ok {
		return t.Category()
	}
	return "Объект"
}

func formatNumber(number int) string {
	return fmt.Sprintf("ВК-%04d", number)
}

// runDemo показывает все сценарии подряд: приём, отказ техосмотра,
// вещи, отчёт, список для новичков и инвентарные номера.
func (a *App) runDemo() {
	fmt.Println(">>> Шаг 1. Приём транспорта через техосмотр")
	a.acceptVehicle(domain.NewScooter("Ninebot Max G30", 7, 0.35))
	a.acceptVehicle(domain.NewScooter("Kugoo Kirin G1", 4, 0.6))
	a.acceptVehicle(domain.NewBicycle("Stels Navigator 300", 9))
	a.acceptVehicle(domain.NewElectricBike("Eltreco XT 600", 6, 0.8))
	a.acceptVehicle(domain.NewElectricBike("Самосбор с рынка", 5, 7.5)) // слишком большой расход
	a.acceptVehicle(domain.NewScooter("Прототип ФКН", 12, 0.3))         // простота вне 1-10
	a.acceptVehicle(domain.NewBicycle("", 8))                           // нет наименования
	fmt.Println()

	fmt.Println(">>> Шаг 2. Вещи на балансе")
	a.acceptThing(domain.NewHelmet("Шлем ВышКат", "M"))
	a.acceptThing(domain.NewHelmet("Шлем ВышКат", "L"))
	a.acceptThing(domain.NewDockStation("Док у Покровки", 1.2))
	a.acceptThing(domain.NewChargingCabinet("Шкаф в корпусе на Мясницкой", 2.5))
	fmt.Println()

	fmt.Println(">>> Шаг 3. Отчёт")
	a.showReport()
	fmt.Println()

	fmt.Println(">>> Шаг 4. Кого можно выдавать новичкам")
	a.showBeginners()
	fmt.Println()

	fmt.Println(">>> Шаг 5. Все инвентарные номера")
	a.showInventory()
}
