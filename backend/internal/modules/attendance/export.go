package attendance

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// generateAttendanceExcel generates an Excel file from attendance records
// Requirements: 1.1 - Generate Excel file (.xlsx) containing attendance records
// Requirements: 1.4, 1.5 - Include student info and attendance details
func generateAttendanceExcel(records []ExportAttendanceRecord) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// Set sheet name
	sheetName := "Attendance"
	f.SetSheetName("Sheet1", sheetName)

	// Define headers
	// Requirements: 1.4 - Include student information (NIS, NISN, name, class)
	// Requirements: 1.5 - Include attendance details (date, check-in time, check-out time, status)
	headers := []string{
		"No",
		"NIS",
		"NISN",
		"Nama Siswa",
		"Kelas",
		"Tanggal",
		"Jam Masuk",
		"Jam Pulang",
		"Status",
		"Jadwal",
	}

	// Set header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create header style: %w", err)
	}

	// Write headers
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Set column widths
	columnWidths := map[string]float64{
		"A": 5,   // No
		"B": 15,  // NIS
		"C": 15,  // NISN
		"D": 30,  // Nama Siswa
		"E": 15,  // Kelas
		"F": 12,  // Tanggal
		"G": 12,  // Jam Masuk
		"H": 12,  // Jam Pulang
		"I": 12,  // Status
		"J": 20,  // Jadwal
	}
	for col, width := range columnWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	// Data style
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create data style: %w", err)
	}

	// Write data rows
	// Requirements: 1.6 - IF no attendance records match, return empty file with headers only
	for i, record := range records {
		row := i + 2 // Start from row 2 (after header)

		// Translate status to Indonesian
		status := translateStatus(record.Status)

		rowData := []interface{}{
			i + 1,
			record.StudentNIS,
			record.StudentNISN,
			record.StudentName,
			record.ClassName,
			record.Date,
			record.CheckInTime,
			record.CheckOutTime,
			status,
			record.ScheduleName,
		}

		for j, value := range rowData {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheetName, cell, value)
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}
	}

	// Write to buffer
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write Excel to buffer: %w", err)
	}

	return buffer.Bytes(), nil
}

// generateMonthlyRecapExcel generates an Excel file from monthly recap data
// Requirements: 2.5 - Export monthly recap to Excel with summary statistics
// Requirements: 8.2 - Include total_sick and total_excused columns per student
func generateMonthlyRecapExcel(recap *MonthlyRecapResponse) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// Set sheet name
	sheetName := "Rekap Bulanan"
	f.SetSheetName("Sheet1", sheetName)

	// Define headers - updated to include Sakit and Izin columns
	// Requirements: 8.2, 8.3, 8.4 - Include sick and excused with proper labels
	headers := []string{
		"No",
		"NIS",
		"NISN",
		"Nama Siswa",
		"Kelas",
		"Hadir",
		"Terlambat",
		"Sangat Terlambat",
		"Tidak Hadir",
		"Sakit",
		"Izin",
		"Persentase Kehadiran (%)",
	}

	// Set header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create header style: %w", err)
	}

	// Write title
	monthName := getMonthName(recap.Month)
	title := fmt.Sprintf("Rekap Kehadiran Bulan %s %d", monthName, recap.Year)
	if recap.ClassName != "" {
		title += fmt.Sprintf(" - Kelas %s", recap.ClassName)
	}
	f.SetCellValue(sheetName, "A1", title)
	f.MergeCell(sheetName, "A1", "L1") // Updated to L1 for 12 columns

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	f.SetCellStyle(sheetName, "A1", "L1", titleStyle)

	// Write total days info
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("Total Hari Sekolah: %d hari", recap.TotalDays))
	f.MergeCell(sheetName, "A2", "L2") // Updated to L2 for 12 columns

	// Write headers (row 4)
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 4)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Set column widths - updated to include Sakit and Izin columns
	columnWidths := map[string]float64{
		"A": 5,   // No
		"B": 15,  // NIS
		"C": 15,  // NISN
		"D": 30,  // Nama Siswa
		"E": 15,  // Kelas
		"F": 10,  // Hadir
		"G": 12,  // Terlambat
		"H": 18,  // Sangat Terlambat
		"I": 12,  // Tidak Hadir
		"J": 10,  // Sakit
		"K": 10,  // Izin
		"L": 25,  // Persentase
	}
	for col, width := range columnWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	// Data style
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create data style: %w", err)
	}

	// Write data rows - updated to include TotalSick and TotalExcused
	for i, student := range recap.StudentRecaps {
		row := i + 5 // Start from row 5 (after title, info, empty row, and header)

		rowData := []interface{}{
			i + 1,
			student.StudentNIS,
			student.StudentNISN,
			student.StudentName,
			student.ClassName,
			student.TotalPresent,
			student.TotalLate,
			student.TotalVeryLate,
			student.TotalAbsent,
			student.TotalSick,
			student.TotalExcused,
			fmt.Sprintf("%.2f%%", student.AttendancePercent),
		}

		for j, value := range rowData {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheetName, cell, value)
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}
	}

	// Write to buffer
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write Excel to buffer: %w", err)
	}

	return buffer.Bytes(), nil
}

// generateExportFilename generates the filename for attendance export
// Requirements: 1.7 - Name the file with format attendance_{school}_{date_range}.xlsx
func generateExportFilename(schoolName, startDate, endDate string) string {
	// Sanitize school name for filename
	sanitizedSchool := sanitizeFilename(schoolName)
	return fmt.Sprintf("attendance_%s_%s_to_%s.xlsx", sanitizedSchool, startDate, endDate)
}

// generateMonthlyRecapFilename generates the filename for monthly recap export
func generateMonthlyRecapFilename(schoolName string, month, year int) string {
	sanitizedSchool := sanitizeFilename(schoolName)
	monthName := getMonthName(month)
	return fmt.Sprintf("rekap_bulanan_%s_%s_%d.xlsx", sanitizedSchool, strings.ToLower(monthName), year)
}

// sanitizeFilename removes or replaces characters that are not safe for filenames
func sanitizeFilename(name string) string {
	// Replace spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")
	// Remove special characters
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return -1
	}, name)
	// Convert to lowercase
	return strings.ToLower(name)
}

// translateStatus translates attendance status to Indonesian
func translateStatus(status string) string {
	switch status {
	case "on_time":
		return "Tepat Waktu"
	case "late":
		return "Terlambat"
	case "very_late":
		return "Sangat Terlambat"
	case "absent":
		return "Tidak Hadir"
	case "sick":
		return "Sakit"
	case "excused":
		return "Izin"
	default:
		return status
	}
}

// getMonthName returns the Indonesian month name
func getMonthName(month int) string {
	months := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	if month >= 1 && month <= 12 {
		return months[month]
	}
	return ""
}
