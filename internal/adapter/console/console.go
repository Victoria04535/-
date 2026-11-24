package console

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"pos/internal/application"
)

type Console struct {
	facade *application.POSFacade
	input  *bufio.Scanner
	output io.Writer
}

func NewConsole(
	facade *application.POSFacade,
	input io.Reader,
	output io.Writer,
) *Console {
	return &Console{
		facade: facade,
		input:  bufio.NewScanner(input),
		output: output,
	}
}

func (c *Console) readLine(prompt string) (string, error) {
	fmt.Fprint(c.output, prompt)

	if !c.input.Scan() {
		if err := c.input.Err(); err != nil {
			return "", err
		}

		return "", io.EOF
	}

	return strings.TrimSpace(c.input.Text()), nil
}

func (c *Console) readFloat(prompt string) (float64, error) {
	value, err := c.readLine(prompt)
	if err != nil {
		return 0, err
	}

	if value == "" {
		return 0, fmt.Errorf("пустой ввод: ожидается число")
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("некорректное число %q: введите число, например 8.5", value)
	}

	if math.IsNaN(parsed) || math.IsInf(parsed, 0) {
		return 0, fmt.Errorf("некорректное число %q: NaN и бесконечность не допускаются", value)
	}

	return parsed, nil
}

func (c *Console) Run() error {
	for {
		fmt.Fprintln(c.output)
		fmt.Fprintln(c.output, "========== ПОС ==========")
		fmt.Fprintln(c.output, "1. Добавить студента")
		fmt.Fprintln(c.output, "2. Список студентов")
		fmt.Fprintln(c.output, "3. Добавить критерий")
		fmt.Fprintln(c.output, "4. Список критериев")
		fmt.Fprintln(c.output, "5. Выставить оценку")
		fmt.Fprintln(c.output, "6. Просмотреть оценки студента")
		fmt.Fprintln(c.output, "7. Отменить оценку")
		fmt.Fprintln(c.output, "8. Рассчитать итог для конкретного студента")
		fmt.Fprintln(c.output, "9. Рассчитать итоги всех студентов")
		fmt.Fprintln(c.output, "10. Экспортировать итоги в CSV")
		fmt.Fprintln(c.output, "0. Выход")

		choice, err := c.readLine("Выберите действие: ")
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		switch choice {
		case "1":
			c.addStudent()
		case "2":
			c.listStudents()
		case "3":
			c.addCriterion()
		case "4":
			c.listCriteria()
		case "5":
			c.setEvaluation()
		case "6":
			c.listEvaluations()
		case "7":
			c.removeEvaluation()
		case "8":
			c.calculateTotal()
		case "9":
			c.listAllResults()
		case "10":
			c.exportResults()
		case "0":
			fmt.Fprintln(c.output, "Завершение работы.")
			return nil
		default:
			fmt.Fprintln(c.output, "Неизвестная команда.")
		}
	}
}

func (c *Console) addStudent() {
	id, err := c.readLine("ID студента: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	name, err := c.readLine("Имя студента: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	if err := c.facade.AddStudent(id, name); err != nil {
		fmt.Fprintln(c.output, "Ошибка:", err)
		return
	}

	fmt.Fprintln(c.output, "Студент успешно добавлен.")
}

func (c *Console) listStudents() {
	students, err := c.facade.ListStudents()
	if err != nil {
		fmt.Fprintln(c.output, "Ошибка:", err)
		return
	}

	if len(students) == 0 {
		fmt.Fprintln(c.output, "Список студентов пуст.")
		return
	}

	fmt.Fprintln(c.output, "\nСписок студентов:")

	for _, student := range students {
		fmt.Fprintf(
			c.output,
			"ID: %s | Имя: %s\n",
			student.ID(),
			student.Name(),
		)
	}
}

func (c *Console) addCriterion() {
	id, err := c.readLine("ID критерия: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	name, err := c.readLine("Название критерия: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	maxScore, err := c.readFloat("Максимальный балл: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	weight, err := c.readFloat("Вес критерия: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	if err := c.facade.AddCriterion(
		id,
		name,
		maxScore,
		weight,
	); err != nil {
		fmt.Fprintln(c.output, "Ошибка:", err)
		return
	}

	fmt.Fprintln(c.output, "Критерий успешно добавлен.")
}

func (c *Console) listCriteria() {
	criteria, err := c.facade.ListCriteria()
	if err != nil {
		fmt.Fprintln(c.output, "Ошибка:", err)
		return
	}

	if len(criteria) == 0 {
		fmt.Fprintln(c.output, "Список критериев пуст.")
		return
	}

	fmt.Fprintln(c.output, "\nСписок критериев:")

	for _, criterion := range criteria {
		fmt.Fprintf(
			c.output,
			"ID: %s | Название: %s | Максимум: %.2f | Вес: %.2f\n",
			criterion.ID(),
			criterion.Name(),
			criterion.MaxScore(),
			criterion.Weight(),
		)
	}
}

func (c *Console) setEvaluation() {
	studentID, err := c.readLine("ID студента: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	criterionID, err := c.readLine("ID критерия: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	score, err := c.readFloat("Выставленный балл: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	comment, err := c.readLine("Комментарий: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	if err := c.facade.SetEvaluation(
		studentID,
		criterionID,
		score,
		comment,
	); err != nil {
		fmt.Fprintln(c.output, "Ошибка:", err)
		return
	}

	fmt.Fprintln(c.output, "Оценка успешно сохранена.")
}

func (c *Console) listEvaluations() {
	studentID, err := c.readLine("ID студента: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	evaluations, err := c.facade.ListEvaluations(studentID)
	if err != nil {
		fmt.Fprintln(c.output, "Ошибка:", err)
		return
	}

	if len(evaluations) == 0 {
		fmt.Fprintln(c.output, "У студента пока нет оценок.")
		return
	}

	fmt.Fprintln(c.output, "\nОценки студента:")

	for _, evaluation := range evaluations {
		fmt.Fprintf(
			c.output,
			"Критерий: %s | Балл: %.2f / %.2f | Комментарий: %s | Дата: %s\n",
			evaluation.CriterionID(),
			evaluation.Score().Value(),
			evaluation.Score().Max(),
			evaluation.Comment(),
			evaluation.CreatedAt().Format("2006-01-02 15:04:05"),
		)
	}
}

func (c *Console) calculateTotal() {
	studentID, err := c.readLine("ID студента: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	total, err := c.facade.CalculateTotal(studentID)
	if err != nil {
		fmt.Fprintln(c.output, "Ошибка:", err)
		return
	}

	fmt.Fprintf(
		c.output,
		"Итоговая оценка студента %s: %.2f / 10\n",
		studentID,
		total,
	)
}

func (c *Console) exportResults() {
	path, err := c.readLine(
		"Путь для CSV (по умолчанию results.csv): ",
	)
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	if path == "" {
		path = "results.csv"
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintln(c.output, "Ошибка создания файла:", err)
		return
	}

	defer file.Close()

	if err := c.facade.ExportResults(file); err != nil {
		fmt.Fprintln(c.output, "Ошибка экспорта:", err)
		return
	}

	fmt.Fprintln(c.output, "Результаты экспортированы:", path)
}

func (c *Console) removeEvaluation() {
	studentID, err := c.readLine("ID студента: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	criterionID, err := c.readLine("ID критерия: ")
	if err != nil {
		fmt.Fprintln(c.output, err)
		return
	}

	if err := c.facade.RemoveEvaluation(
		studentID,
		criterionID,
	); err != nil {
		fmt.Fprintln(c.output, "Ошибка:", err)
		return
	}

	fmt.Fprintln(c.output, "Оценка успешно отменена.")
}

func (c *Console) listAllResults() {
	results, err := c.facade.BuildResults()
	if err != nil {
		fmt.Fprintln(c.output, "Ошибка:", err)
		return
	}

	if len(results) == 0 {
		fmt.Fprintln(c.output, "Список студентов пуст.")
		return
	}

	fmt.Fprintln(c.output, "\nИтоги студентов:")

	for _, result := range results {
		fmt.Fprintf(
			c.output,
			"ID: %s | Имя: %s | Итог: %.2f / 10\n",
			result.Student.ID(),
			result.Student.Name(),
			result.Total,
		)
	}
}
