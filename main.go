package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Student представляет студента
type Student struct {
	ID        int
	LastName  string
	FirstName string
	Grades    []int
}

// NewStudent создает нового студента
func NewStudent(id int, lastName, firstName string) *Student {
	return &Student{
		ID:        id,
		LastName:  lastName,
		FirstName: firstName,
		Grades:    make([]int, 0),
	}
}

// AddGrade добавляет оценку студенту
func (s *Student) AddGrade(grade int) error {
	if grade < 2 || grade > 5 {
		return fmt.Errorf("оценка должна быть от 2 до 5")
	}
	s.Grades = append(s.Grades, grade)
	return nil
}

// AverageScore вычисляет средний балл студента
func (s *Student) AverageScore() float64 {
	if len(s.Grades) == 0 {
		return 0
	}

	sum := 0
	for _, grade := range s.Grades {
		sum += grade
	}
	return float64(sum) / float64(len(s.Grades))
}

// GetStatus возвращает статус студента по среднему баллу
func (s *Student) GetStatus() string {
	avg := s.AverageScore()
	switch {
	case avg >= 4.5:
		return "отличник"
	case avg >= 4.0:
		return "хорошист"
	case avg >= 3.0:
		return "троечник"
	default:
		return "неуспевающий"
	}
}

// DisplayInfo выводит информацию о студенте
func (s *Student) DisplayInfo() {
	fmt.Printf("ID: %d, Студент: %s %s\n", s.ID, s.LastName, s.FirstName)
	fmt.Printf("  Оценки: %v\n", s.Grades)
	fmt.Printf("  Средний балл: %.2f\n", s.AverageScore())
	fmt.Printf("  Статус: %s\n", s.GetStatus())
}

// Journal представляет журнал группы
type Journal struct {
	Students map[int]*Student
	NextID   int
}

// NewJournal создает новый журнал
func NewJournal() *Journal {
	return &Journal{
		Students: make(map[int]*Student),
		NextID:   1,
	}
}

// AddStudent добавляет нового студента в журнал
func (j *Journal) AddStudent(lastName, firstName string) int {
	id := j.NextID
	j.Students[id] = NewStudent(id, lastName, firstName)
	j.NextID++
	return id
}

// FindStudentByID находит студента по ID
func (j *Journal) FindStudentByID(id int) (*Student, bool) {
	student, exists := j.Students[id]
	return student, exists
}

// AddGradeToStudent добавляет оценку студенту по ID
func (j *Journal) AddGradeToStudent(studentID, grade int) error {
	student, exists := j.FindStudentByID(studentID)
	if !exists {
		return fmt.Errorf("студент с ID %d не найден", studentID)
	}
	return student.AddGrade(grade)
}

// DisplayAllStudents выводит всех студентов
func (j *Journal) DisplayAllStudents() {
	if len(j.Students) == 0 {
		fmt.Println("В журнале нет студентов")
		return
	}

	// Сортируем студентов по ID для красивого вывода
	ids := make([]int, 0, len(j.Students))
	for id := range j.Students {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	fmt.Println("\n=== СПИСОК СТУДЕНТОВ ===")
	for _, id := range ids {
		j.Students[id].DisplayInfo()
		fmt.Println()
	}
}

// FilterByAverageScore фильтрует студентов по среднему баллу
func (j *Journal) FilterByAverageScore(minScore, maxScore float64) []*Student {
	var result []*Student

	for _, student := range j.Students {
		avg := student.AverageScore()
		if avg >= minScore && avg <= maxScore {
			result = append(result, student)
		}
	}

	// Сортируем по среднему баллу (по убыванию)
	sort.Slice(result, func(i, j int) bool {
		return result[i].AverageScore() > result[j].AverageScore()
	})

	return result
}

// GetGroupStatistics возвращает статистику по группе
func (j *Journal) GetGroupStatistics() {
	if len(j.Students) == 0 {
		fmt.Println("В журнале нет студентов")
		return
	}

	totalStudents := len(j.Students)
	var totalAverage, sumAllAverages float64

	// Считаем статистику по статусам
	statusCount := make(map[string]int)

	for _, student := range j.Students {
		avg := student.AverageScore()
		sumAllAverages += avg
		status := student.GetStatus()
		statusCount[status]++
	}

	totalAverage = sumAllAverages / float64(totalStudents)

	fmt.Println("\n=== СТАТИСТИКА ГРУППЫ ===")
	fmt.Printf("Общее количество студентов: %d\n", totalStudents)
	fmt.Printf("Средний балл группы: %.2f\n", totalAverage)
	fmt.Println("Распределение по статусам:")
	for status, count := range statusCount {
		percentage := float64(count) / float64(totalStudents) * 100
		fmt.Printf("  %s: %d (%.1f%%)\n", status, count, percentage)
	}
}

// DisplayMenu показывает главное меню
func DisplayMenu() {
	fmt.Println("\n=== ЖУРНАЛ ГРУППЫ ===")
	fmt.Println("1. Добавить студента")
	fmt.Println("2. Добавить оценку студенту")
	fmt.Println("3. Показать всех студентов")
	fmt.Println("4. Фильтр по среднему баллу")
	fmt.Println("5. Статистика группы")
	fmt.Println("6. Удалить студента")
	fmt.Println("7. Выход")
	fmt.Print("Выберите действие: ")
}

// InputStudent читает данные студента с консоли
func InputStudent() (string, string) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введите фамилию студента: ")
	scanner.Scan()
	lastName := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите имя студента: ")
	scanner.Scan()
	firstName := strings.TrimSpace(scanner.Text())

	return lastName, firstName
}

// InputGrade читает оценку с консоли
func InputGrade() int {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Введите оценку (2-5): ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		grade, err := strconv.Atoi(input)
		if err != nil || grade < 2 || grade > 5 {
			fmt.Println("Ошибка: оценка должна быть числом от 2 до 5")
			continue
		}

		return grade
	}
}

func main() {
	journal := NewJournal()
	scanner := bufio.NewScanner(os.Stdin)

	// Добавим несколько тестовых студентов для демонстрации
	journal.AddStudent("Иванов", "Иван")
	journal.AddStudent("Петров", "Петр")
	journal.AddStudent("Сидорова", "Мария")

	// Добавим тестовые оценки
	journal.AddGradeToStudent(1, 5)
	journal.AddGradeToStudent(1, 4)
	journal.AddGradeToStudent(2, 3)
	journal.AddGradeToStudent(2, 4)
	journal.AddGradeToStudent(3, 5)
	journal.AddGradeToStudent(3, 5)

	for {
		DisplayMenu()

		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			// Добавление студента
			fmt.Println("\n--- ДОБАВЛЕНИЕ СТУДЕНТА ---")
			lastName, firstName := InputStudent()
			id := journal.AddStudent(lastName, firstName)
			fmt.Printf("Студент добавлен с ID: %d\n", id)

		case "2":
			// Добавление оценки
			fmt.Println("\n--- ДОБАВЛЕНИЕ ОЦЕНКИ ---")
			journal.DisplayAllStudents()

			fmt.Print("Введите ID студента: ")
			scanner.Scan()
			studentID, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil {
				fmt.Println("Ошибка: введите корректный ID")
				continue
			}

			grade := InputGrade()
			err = journal.AddGradeToStudent(studentID, grade)
			if err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			} else {
				fmt.Println("Оценка добавлена успешно")
			}

		case "3":
			// Показать всех студентов
			fmt.Println("\n--- ВСЕ СТУДЕНТЫ ---")
			journal.DisplayAllStudents()

		case "4":
			// Фильтр по среднему баллу
			fmt.Println("\n--- ФИЛЬТР ПО СРЕДНЕМУ БАЛЛУ ---")
			fmt.Println("1. Студенты со средним баллом ниже 4.0")
			fmt.Println("2. Студенты со средним баллом 4.0 и выше")
			fmt.Println("3. Произвольный диапазон")
			fmt.Print("Выберите тип фильтра: ")

			scanner.Scan()
			filterChoice := strings.TrimSpace(scanner.Text())

			var students []*Student

			switch filterChoice {
			case "1":
				students = journal.FilterByAverageScore(0, 3.99)
				fmt.Println("\n--- СТУДЕНТЫ СО СРЕДНИМ БАЛЛОМ НИЖЕ 4.0 ---")
			case "2":
				students = journal.FilterByAverageScore(4.0, 5.0)
				fmt.Println("\n--- СТУДЕНТЫ СО СРЕДНИМ БАЛЛОМ 4.0 И ВЫШЕ ---")
			case "3":
				fmt.Print("Введите минимальный балл: ")
				scanner.Scan()
				minScore, _ := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)

				fmt.Print("Введите максимальный балл: ")
				scanner.Scan()
				maxScore, _ := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)

				students = journal.FilterByAverageScore(minScore, maxScore)
				fmt.Printf("\n--- СТУДЕНТЫ СО СРЕДНИМ БАЛЛОМ ОТ %.1f ДО %.1f ---\n", minScore, maxScore)
			default:
				fmt.Println("Неверный выбор")
				continue
			}

			if len(students) == 0 {
				fmt.Println("Студенты не найдены")
			} else {
				for _, student := range students {
					student.DisplayInfo()
					fmt.Println()
				}
			}

		case "5":
			// Статистика группы
			journal.GetGroupStatistics()

		case "6":
			// Удаление студента
			fmt.Println("\n--- УДАЛЕНИЕ СТУДЕНТА ---")
			journal.DisplayAllStudents()

			fmt.Print("Введите ID студента для удаления: ")
			scanner.Scan()
			studentID, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil {
				fmt.Println("Ошибка: введите корректный ID")
				continue
			}

			if _, exists := journal.FindStudentByID(studentID); exists {
				delete(journal.Students, studentID)
				fmt.Println("Студент удален успешно")
			} else {
				fmt.Println("Студент с таким ID не найден")
			}

		case "7":
			fmt.Println("До свидания!")
			return

		default:
			fmt.Println("Ошибка: выберите действие от 1 до 7")
		}

		fmt.Print("\nНажмите Enter для продолжения...")
		scanner.Scan()
	}
}
