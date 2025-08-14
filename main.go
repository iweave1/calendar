package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// formatMonthCalendar formats the calendar for the given date
func formatMonthCalendar(date time.Time, isYearly bool) []string {
	year, month, day := date.Year(), date.Month(), date.Day()

	// Get the first day of the month
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, date.Location())

	// Determine the number of days in the month
	daysInMonth := firstDay.AddDate(0, 1, -1).Day()

	// Get the weekday the month starts on
	startWeekday := firstDay.Weekday()

	line := "Su Mo Tu We Th Fr Sa" // Separator line
	lineLength := len(line)        // Length of the separator

	// Calculate padding
	monthYear := month.String() + " " + strconv.Itoa(year)
	padding := (lineLength - len(monthYear)) / 2
	formattedMonth := fmt.Sprintf("%*s%s%*s", padding, "", monthYear, padding, "")

	// Print centered month with separator line
	var output []string
	output = append(output, formattedMonth)
	output = append(output, line)
	output = append(output, strings.Repeat("-", lineLength)) // Add the separator line

	// Build the first week line with correct padding
	firstWeek := ""
	for i := 0; i < int(startWeekday); i++ {
		firstWeek += "   " // Three spaces for alignment
	}

	weekLines := 0
	for i := 1; i <= daysInMonth; i++ {
		if i == day && !isYearly {
			firstWeek += fmt.Sprintf("\x1b[31m%2d\x1b[0m ", i) // Highlight current day in red
		} else {
			firstWeek += fmt.Sprintf("%2d ", i)
		}

		if (int(startWeekday)+i)%7 == 0 {
			weekLine := strings.TrimRight(firstWeek, " ")
			if len(weekLine) < len(line) {
				weekLine += strings.Repeat(" ", len(line)-len(weekLine))
			}
			output = append(output, weekLine)
			firstWeek = ""
			weekLines++
		}
	}
	if strings.TrimSpace(firstWeek) != "" {
		weekLine := strings.TrimRight(firstWeek, " ")
		if len(weekLine) < len(line) {
			weekLine += strings.Repeat(" ", len(line)-len(weekLine))
		}
		output = append(output, weekLine)
		weekLines++
	}

	// Pad with blank week lines if less than 6
	for weekLines < 6 {
		output = append(output, strings.Repeat(" ", len(line)))
		weekLines++
	}
	return output
}

func main() {
	// Use current date
	currentDate := time.Now()
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if len(arg) == 10 {
			// Try parsing as full date YYYY-MM-DD
			if inputDate, err := time.Parse("2006-01-02", arg); err == nil {
				currentDate = inputDate
				fmt.Println(strings.Join(formatMonthCalendar(currentDate, false), "\n"))
			} else {
				fmt.Println("Invalid date format. Please use YYYY-MM-DD or just YYYY.")
				return
			}
		} else if len(arg) == 4 {
			// Try parsing as year YYYY
			if inputYear, err := time.Parse("2006", arg); err == nil {
				currentDate = time.Date(inputYear.Year(), time.January, 1, 0, 0, 0, 0, inputYear.Location())
				var calendarSlice [][]string
				// Generate calendar for the next 12 months
				for i := 0; i < 12; i++ {
					calendarSlice = append(calendarSlice, formatMonthCalendar(currentDate, true))
					currentDate = currentDate.AddDate(0, 1, 0) // Move to the next month
				}
			} else {
				fmt.Println("Invalid date format. Please use YYYY-MM-DD.")
				return
			}
			return
		}
	} else {
		fmt.Println(strings.Join(formatMonthCalendar(currentDate, false), "\n"))
	}
}

func PrintCalendarGrid(months [][]string) {
	const cols = 4
	const rows = 3

	for row := 0; row < rows; row++ {
		// Find the max number of lines in this row of months
		maxLines := 0
		headerWidths := make([]int, cols)
		for col := 0; col < cols; col++ {
			idx := row*cols + col
			if idx < len(months) {
				if len(months[idx]) > maxLines {
					maxLines = len(months[idx])
				}
				headerWidths[col] = len(months[idx][0])
			}
		}
		// Pad each month in this row to maxLines using its header width
		padded := make([][]string, cols)
		for col := 0; col < cols; col++ {
			idx := row*cols + col
			if idx < len(months) {
				m := make([]string, len(months[idx]))
				for i, line := range months[idx] {
					if len(line) < headerWidths[col] {
						m[i] = line + strings.Repeat(" ", headerWidths[col]-len(line))
					} else {
						m[i] = line
					}
				}
				for len(m) < maxLines {
					m = append(m, strings.Repeat(" ", headerWidths[col]))
				}
				padded[col] = m
			} else {
				padded[col] = make([]string, maxLines)
				for i := 0; i < maxLines; i++ {
					padded[col][i] = strings.Repeat(" ", 20) // fallback width
				}
			}
		}
		// Print each line of the 4 months in this row
		for line := 0; line < maxLines; line++ {
			for col := 0; col < cols; col++ {
				fmt.Print(padded[col][line])
				if col < cols-1 {
					fmt.Print("   ") // space between months
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
