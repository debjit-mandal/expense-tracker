package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Expense struct {
	Name     string
	Price    float64
	Category string
	Date     time.Time
}

type Budget struct {
	Limit    float64
	Expenses []Expense
}

func main() {
	budget := Budget{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your budget limit: ")
	limitInput, _ := reader.ReadString('\n')
	limitInput = trimNewLine(limitInput)
	limit, _ := strconv.ParseFloat(limitInput, 64)
	budget.Limit = limit

	for {
		fmt.Println("\n1. Add an expense")
		fmt.Println("2. View expenses")
		fmt.Println("3. Filter expenses by category")
		fmt.Println("4. View remaining budget")
		fmt.Println("5. Generate financial report")
		fmt.Println("6. Export expenses to CSV")
		fmt.Println("7. Import expenses from CSV")
		fmt.Println("8. Exit")

		fmt.Print("\nEnter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = trimNewLine(choice)

		switch choice {
		case "1":
			addExpense(reader, &budget)
		case "2":
			viewExpenses(budget)
		case "3":
			filterExpensesByCategory(reader, budget)
		case "4":
			viewRemainingBudget(budget)
		case "5":
			generateFinancialReport(budget)
		case "6":
			exportExpensesToCSV(budget)
		case "7":
			importExpensesFromCSV(reader, &budget)
		case "8":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func addExpense(reader *bufio.Reader, budget *Budget) {
	expense := Expense{}

	fmt.Print("\nEnter expense name: ")
	name, _ := reader.ReadString('\n')
	expense.Name = trimNewLine(name)

	fmt.Print("Enter expense amount: ")
	priceInput, _ := reader.ReadString('\n')
	priceInput = trimNewLine(priceInput)
	price, _ := strconv.ParseFloat(priceInput, 64)
	expense.Price = price

	fmt.Print("Enter expense category: ")
	category, _ := reader.ReadString('\n')
	expense.Category = trimNewLine(category)

	expense.Date = time.Now()

	budget.Expenses = append(budget.Expenses, expense)
	fmt.Println("Expense added successfully!")
}

func viewExpenses(budget Budget) {
	if len(budget.Expenses) == 0 {
		fmt.Println("No expenses found.")
		return
	}

	fmt.Println("\nExpenses:")
	for _, expense := range budget.Expenses {
		fmt.Printf("- %s: $%.2f (%s) [%s]\n", expense.Name, expense.Price, expense.Category, expense.Date.Format("2006-01-02 15:04:05"))
	}
}

func filterExpensesByCategory(reader *bufio.Reader, budget Budget) {
	if len(budget.Expenses) == 0 {
		fmt.Println("No expenses found.")
		return
	}

	fmt.Print("Enter category to filter expenses: ")
	category, _ := reader.ReadString('\n')
	category = trimNewLine(category)

	filteredExpenses := make([]Expense, 0)

	for _, expense := range budget.Expenses {
		if expense.Category == category {
			filteredExpenses = append(filteredExpenses, expense)
		}
	}

	if len(filteredExpenses) == 0 {
		fmt.Println("No expenses found in the specified category.")
		return
	}

	fmt.Println("\nFiltered Expenses:")
	for _, expense := range filteredExpenses {
		fmt.Printf("- %s: $%.2f (%s) [%s]\n", expense.Name, expense.Price, expense.Category, expense.Date.Format("2006-01-02 15:04:05"))
	}
}

func viewRemainingBudget(budget Budget) {
	totalExpenses := 0.0
	for _, expense := range budget.Expenses {
		totalExpenses += expense.Price
	}

	remainingBudget := budget.Limit - totalExpenses
	fmt.Printf("\nRemaining Budget: $%.2f\n", remainingBudget)
}

func generateFinancialReport(budget Budget) {
	if len(budget.Expenses) == 0 {
		fmt.Println("No expenses found.")
		return
	}

	file, _ := os.Create("financial_report.txt")
	defer file.Close()

	totalExpenses := 0.0
	for _, expense := range budget.Expenses {
		totalExpenses += expense.Price
	}

	remainingBudget := budget.Limit - totalExpenses

	fmt.Fprintf(file, "Financial Report\n\n")
	fmt.Fprintf(file, "Budget Limit: $%.2f\n", budget.Limit)
	fmt.Fprintf(file, "Total Expenses: $%.2f\n", totalExpenses)
	fmt.Fprintf(file, "Remaining Budget: $%.2f\n\n", remainingBudget)

	fmt.Fprintf(file, "Expenses:\n")
	for _, expense := range budget.Expenses {
		fmt.Fprintf(file, "- %s: $%.2f (%s) [%s]\n", expense.Name, expense.Price, expense.Category, expense.Date.Format("2006-01-02 15:04:05"))
	}

	fmt.Println("Financial report generated successfully!")
}

func exportExpensesToCSV(budget Budget) {
	if len(budget.Expenses) == 0 {
		fmt.Println("No expenses found.")
		return
	}

	file, _ := os.Create("expenses.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Name", "Price", "Category", "Date"}
	writer.Write(headers)

	for _, expense := range budget.Expenses {
		record := []string{
			expense.Name,
			strconv.FormatFloat(expense.Price, 'f', 2, 64),
			expense.Category,
			expense.Date.Format("2006-01-02 15:04:05"),
		}
		writer.Write(record)
	}

	fmt.Println("Expenses exported to CSV successfully!")
}

func importExpensesFromCSV(reader *bufio.Reader, budget *Budget) {
	fmt.Print("Enter the path of the CSV file: ")
	filePath, _ := reader.ReadString('\n')
	filePath = trimNewLine(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	for _, record := range records {
		expense := Expense{
			Name:     record[0],
			Price:    parseFloat(record[1]),
			Category: record[2],
		}

		date, err := time.Parse("2006-01-02 15:04:05", record[3])
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}
		expense.Date = date

		budget.Expenses = append(budget.Expenses, expense)
	}

	fmt.Println("Expenses imported from CSV successfully!")
}

func parseFloat(value string) float64 {
	value = strings.ReplaceAll(value, ",", "")
	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println("Error parsing float value:", err)
		return 0.0
	}
	return parsedValue
}

func trimNewLine(str string) string {
	return strings.TrimSuffix(str, "\n")
}
