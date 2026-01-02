package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/school-management/backend/internal/config"
	"github.com/school-management/backend/internal/domain/models"
	"github.com/school-management/backend/internal/shared/database"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Hash password helper
	hashPassword := func(password string) string {
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hash)
	}

	log.Println("Starting seed data...")

	// ============================================
	// 1. Create Super Admin
	// ============================================
	superAdmin := models.User{
		Username:     "superadmin",
		PasswordHash: hashPassword("admin123"),
		Email:        "superadmin@schoolsaas.com",
		Role:         models.RoleSuperAdmin,
		IsActive:     true,
		MustResetPwd: false, // Set false for testing
	}
	if err := db.FirstOrCreate(&superAdmin, models.User{Username: "superadmin"}).Error; err != nil {
		log.Fatalf("Failed to create super admin: %v", err)
	}
	log.Printf("✓ Super Admin created: %s (password: admin123)", superAdmin.Username)

	// ============================================
	// 2. Create Schools (Tenants)
	// ============================================
	schools := []models.School{
		{
			Name:     "SMP Negeri 1 Jakarta",
			Address:  "Jl. Merdeka No. 1, Jakarta Pusat",
			Phone:    "021-1234567",
			Email:    "smpn1jakarta@edu.id",
			IsActive: true,
		},
		{
			Name:     "SMP Negeri 2 Bandung",
			Address:  "Jl. Asia Afrika No. 10, Bandung",
			Phone:    "022-7654321",
			Email:    "smpn2bandung@edu.id",
			IsActive: true,
		},
	}

	for i := range schools {
		if err := db.FirstOrCreate(&schools[i], models.School{Name: schools[i].Name}).Error; err != nil {
			log.Fatalf("Failed to create school: %v", err)
		}
		log.Printf("✓ School created: %s (ID: %d)", schools[i].Name, schools[i].ID)
	}

	school1 := schools[0]

	// ============================================
	// 3. Create School Settings
	// ============================================
	settings := models.SchoolSettings{
		SchoolID:                     school1.ID,
		AttendanceStartTime:          "07:00",
		AttendanceEndTime:            "07:30",
		AttendanceLateThreshold:      30,
		AttendanceVeryLateThreshold:  60,
		EnableAttendanceNotification: true,
		EnableGradeNotification:      true,
		EnableBKNotification:         true,
		EnableHomeroomNotification:   true,
		AcademicYear:                 "2024/2025",
		Semester:                     1,
	}
	if err := db.FirstOrCreate(&settings, models.SchoolSettings{SchoolID: school1.ID}).Error; err != nil {
		log.Fatalf("Failed to create school settings: %v", err)
	}
	log.Printf("✓ School settings created for: %s", school1.Name)

	// ============================================
	// 4. Create Admin Sekolah
	// ============================================
	adminSekolah := models.User{
		SchoolID:     &school1.ID,
		Username:     "admin.smpn1",
		PasswordHash: hashPassword("admin123"),
		Email:        "admin@smpn1jakarta.edu.id",
		Role:         models.RoleAdminSekolah,
		IsActive:     true,
		MustResetPwd: false,
	}
	if err := db.FirstOrCreate(&adminSekolah, models.User{Username: "admin.smpn1"}).Error; err != nil {
		log.Fatalf("Failed to create admin sekolah: %v", err)
	}
	log.Printf("✓ Admin Sekolah created: %s (password: admin123)", adminSekolah.Username)

	// ============================================
	// 5. Create Classes
	// ============================================
	classes := []models.Class{
		{SchoolID: school1.ID, Name: "7A", Grade: 7, Year: "2024/2025"},
		{SchoolID: school1.ID, Name: "7B", Grade: 7, Year: "2024/2025"},
		{SchoolID: school1.ID, Name: "8A", Grade: 8, Year: "2024/2025"},
		{SchoolID: school1.ID, Name: "8B", Grade: 8, Year: "2024/2025"},
		{SchoolID: school1.ID, Name: "9A", Grade: 9, Year: "2024/2025"},
		{SchoolID: school1.ID, Name: "9B", Grade: 9, Year: "2024/2025"},
	}

	for i := range classes {
		if err := db.FirstOrCreate(&classes[i], models.Class{SchoolID: school1.ID, Name: classes[i].Name}).Error; err != nil {
			log.Fatalf("Failed to create class: %v", err)
		}
	}
	log.Printf("✓ Created %d classes", len(classes))

	class7A := classes[0]
	class7B := classes[1]

	// ============================================
	// 6. Create Wali Kelas (Homeroom Teachers)
	// ============================================
	waliKelas7A := models.User{
		SchoolID:     &school1.ID,
		Username:     "walikelas.7a",
		PasswordHash: hashPassword("guru123"),
		Email:        "walikelas7a@smpn1jakarta.edu.id",
		Role:         models.RoleWaliKelas,
		IsActive:     true,
		MustResetPwd: false,
	}
	if err := db.FirstOrCreate(&waliKelas7A, models.User{Username: "walikelas.7a"}).Error; err != nil {
		log.Fatalf("Failed to create wali kelas: %v", err)
	}
	log.Printf("✓ Wali Kelas 7A created: %s (password: guru123)", waliKelas7A.Username)

	// Update class with homeroom teacher
	db.Model(&class7A).Update("homeroom_teacher_id", waliKelas7A.ID)

	waliKelas7B := models.User{
		SchoolID:     &school1.ID,
		Username:     "walikelas.7b",
		PasswordHash: hashPassword("guru123"),
		Email:        "walikelas7b@smpn1jakarta.edu.id",
		Role:         models.RoleWaliKelas,
		IsActive:     true,
		MustResetPwd: false,
	}
	if err := db.FirstOrCreate(&waliKelas7B, models.User{Username: "walikelas.7b"}).Error; err != nil {
		log.Fatalf("Failed to create wali kelas: %v", err)
	}
	log.Printf("✓ Wali Kelas 7B created: %s (password: guru123)", waliKelas7B.Username)
	db.Model(&class7B).Update("homeroom_teacher_id", waliKelas7B.ID)

	// ============================================
	// 7. Create Guru BK
	// ============================================
	guruBK := models.User{
		SchoolID:     &school1.ID,
		Username:     "gurubk.smpn1",
		PasswordHash: hashPassword("guru123"),
		Email:        "gurubk@smpn1jakarta.edu.id",
		Role:         models.RoleGuruBK,
		IsActive:     true,
		MustResetPwd: false,
	}
	if err := db.FirstOrCreate(&guruBK, models.User{Username: "gurubk.smpn1"}).Error; err != nil {
		log.Fatalf("Failed to create guru BK: %v", err)
	}
	log.Printf("✓ Guru BK created: %s (password: guru123)", guruBK.Username)

	// ============================================
	// 8. Create Students
	// ============================================
	class7AID := class7A.ID
	class7BID := class7B.ID
	students := []models.Student{
		{SchoolID: school1.ID, ClassID: &class7AID, NIS: "2024001", NISN: "0012345001", Name: "Ahmad Rizki", RFIDCode: "RFID001", IsActive: true},
		{SchoolID: school1.ID, ClassID: &class7AID, NIS: "2024002", NISN: "0012345002", Name: "Budi Santoso", RFIDCode: "RFID002", IsActive: true},
		{SchoolID: school1.ID, ClassID: &class7AID, NIS: "2024003", NISN: "0012345003", Name: "Citra Dewi", RFIDCode: "RFID003", IsActive: true},
		{SchoolID: school1.ID, ClassID: &class7AID, NIS: "2024004", NISN: "0012345004", Name: "Dian Permata", RFIDCode: "RFID004", IsActive: true},
		{SchoolID: school1.ID, ClassID: &class7AID, NIS: "2024005", NISN: "0012345005", Name: "Eko Prasetyo", RFIDCode: "RFID005", IsActive: true},
		{SchoolID: school1.ID, ClassID: &class7BID, NIS: "2024006", NISN: "0012345006", Name: "Fitri Handayani", RFIDCode: "RFID006", IsActive: true},
		{SchoolID: school1.ID, ClassID: &class7BID, NIS: "2024007", NISN: "0012345007", Name: "Galih Pratama", RFIDCode: "RFID007", IsActive: true},
		{SchoolID: school1.ID, ClassID: &class7BID, NIS: "2024008", NISN: "0012345008", Name: "Hana Safitri", RFIDCode: "RFID008", IsActive: true},
	}

	for i := range students {
		if err := db.FirstOrCreate(&students[i], models.Student{NISN: students[i].NISN}).Error; err != nil {
			log.Fatalf("Failed to create student: %v", err)
		}
	}
	log.Printf("✓ Created %d students", len(students))


	// ============================================
	// 9. Create Parents and Link to Students
	// ============================================
	// Create parent users first
	parentUser1 := models.User{
		SchoolID:     &school1.ID,
		Username:     "parent.ahmad",
		PasswordHash: hashPassword("parent123"),
		Email:        "parent.ahmad@gmail.com",
		Role:         models.RoleParent,
		IsActive:     true,
		MustResetPwd: false,
	}
	if err := db.FirstOrCreate(&parentUser1, models.User{Username: "parent.ahmad"}).Error; err != nil {
		log.Fatalf("Failed to create parent user: %v", err)
	}

	parent1 := models.Parent{
		SchoolID: school1.ID,
		UserID:   parentUser1.ID,
		Name:     "Pak Rizki (Ayah Ahmad)",
		Phone:    "081234567890",
	}
	if err := db.FirstOrCreate(&parent1, models.Parent{UserID: parentUser1.ID}).Error; err != nil {
		log.Fatalf("Failed to create parent: %v", err)
	}

	// Link parent to student
	db.Exec("INSERT INTO student_parents (student_id, parent_id) VALUES (?, ?) ON CONFLICT DO NOTHING", students[0].ID, parent1.ID)
	log.Printf("✓ Parent created: %s (password: parent123)", parentUser1.Username)

	parentUser2 := models.User{
		SchoolID:     &school1.ID,
		Username:     "parent.budi",
		PasswordHash: hashPassword("parent123"),
		Email:        "parent.budi@gmail.com",
		Role:         models.RoleParent,
		IsActive:     true,
		MustResetPwd: false,
	}
	if err := db.FirstOrCreate(&parentUser2, models.User{Username: "parent.budi"}).Error; err != nil {
		log.Fatalf("Failed to create parent user: %v", err)
	}

	parent2 := models.Parent{
		SchoolID: school1.ID,
		UserID:   parentUser2.ID,
		Name:     "Bu Santoso (Ibu Budi)",
		Phone:    "081234567891",
	}
	if err := db.FirstOrCreate(&parent2, models.Parent{UserID: parentUser2.ID}).Error; err != nil {
		log.Fatalf("Failed to create parent: %v", err)
	}
	db.Exec("INSERT INTO student_parents (student_id, parent_id) VALUES (?, ?) ON CONFLICT DO NOTHING", students[1].ID, parent2.ID)
	log.Printf("✓ Parent created: %s (password: parent123)", parentUser2.Username)

	// ============================================
	// 10. Create Device (ESP32 RFID)
	// ============================================
	device := models.Device{
		SchoolID:    school1.ID,
		DeviceCode:  "ESP32-SMPN1-001",
		APIKey:      "sk_test_device_api_key_12345",
		Description: "RFID Reader - Gerbang Utama",
		IsActive:    true,
	}
	if err := db.FirstOrCreate(&device, models.Device{DeviceCode: "ESP32-SMPN1-001"}).Error; err != nil {
		log.Fatalf("Failed to create device: %v", err)
	}
	log.Printf("✓ Device created: %s (API Key: %s)", device.DeviceCode, device.APIKey)

	// ============================================
	// 11. Create Sample Attendance Records
	// ============================================
	today := time.Now()
	checkInTime := time.Date(today.Year(), today.Month(), today.Day(), 6, 55, 0, 0, time.Local)
	checkOutTime := time.Date(today.Year(), today.Month(), today.Day(), 14, 30, 0, 0, time.Local)

	for i := 0; i < 5; i++ {
		attendance := models.Attendance{
			StudentID:    students[i].ID,
			Date:         today,
			CheckInTime:  &checkInTime,
			CheckOutTime: &checkOutTime,
			Status:       "on_time",
			Method:       "rfid",
		}
		db.FirstOrCreate(&attendance, models.Attendance{StudentID: students[i].ID, Date: today})
	}
	log.Printf("✓ Created attendance records for today")

	// ============================================
	// 12. Create Sample Violations
	// ============================================
	violation := models.Violation{
		StudentID:   students[1].ID,
		Category:    "Keterlambatan",
		Level:       "ringan",
		Description: "Terlambat masuk kelas 15 menit",
		CreatedBy:   guruBK.ID,
	}
	db.FirstOrCreate(&violation, models.Violation{StudentID: students[1].ID, Category: "Keterlambatan"})
	log.Printf("✓ Created sample violation")

	// ============================================
	// 13. Create Sample Achievements
	// ============================================
	achievement := models.Achievement{
		StudentID:   students[0].ID,
		Title:       "Juara 1 Lomba Matematika",
		Point:       100,
		Description: "Juara 1 Olimpiade Matematika Tingkat Kota",
		CreatedBy:   guruBK.ID,
	}
	db.FirstOrCreate(&achievement, models.Achievement{StudentID: students[0].ID, Title: "Juara 1 Lomba Matematika"})
	log.Printf("✓ Created sample achievement")

	// ============================================
	// 14. Create Sample Grades
	// ============================================
	grades := []models.Grade{
		{StudentID: students[0].ID, Title: "Ulangan Harian Matematika", Score: 85, Description: "Bab 1: Bilangan Bulat", CreatedBy: waliKelas7A.ID},
		{StudentID: students[0].ID, Title: "Ulangan Harian IPA", Score: 90, Description: "Bab 1: Klasifikasi Makhluk Hidup", CreatedBy: waliKelas7A.ID},
		{StudentID: students[1].ID, Title: "Ulangan Harian Matematika", Score: 78, Description: "Bab 1: Bilangan Bulat", CreatedBy: waliKelas7A.ID},
		{StudentID: students[2].ID, Title: "Ulangan Harian Matematika", Score: 92, Description: "Bab 1: Bilangan Bulat", CreatedBy: waliKelas7A.ID},
	}
	for _, g := range grades {
		db.FirstOrCreate(&g, models.Grade{StudentID: g.StudentID, Title: g.Title})
	}
	log.Printf("✓ Created %d sample grades", len(grades))

	// ============================================
	// 15. Create Sample Homeroom Notes
	// ============================================
	note := models.HomeroomNote{
		StudentID: students[0].ID,
		TeacherID: waliKelas7A.ID,
		Content:   "Ahmad menunjukkan perkembangan yang baik dalam pelajaran. Aktif bertanya di kelas dan membantu teman-temannya.",
	}
	db.FirstOrCreate(&note, models.HomeroomNote{StudentID: students[0].ID, TeacherID: waliKelas7A.ID})
	log.Printf("✓ Created sample homeroom note")

	// ============================================
	// Summary
	// ============================================
	fmt.Println("\n" + repeatStr("=", 60))
	fmt.Println("SEED DATA COMPLETED SUCCESSFULLY!")
	fmt.Println(repeatStr("=", 60))
	fmt.Println("\nLogin Credentials:")
	fmt.Println(repeatStr("-", 60))
	fmt.Printf("%-20s %-25s %s\n", "Role", "Username", "Password")
	fmt.Println(repeatStr("-", 60))
	fmt.Printf("%-20s %-25s %s\n", "Super Admin", "superadmin", "admin123")
	fmt.Printf("%-20s %-25s %s\n", "Admin Sekolah", "admin.smpn1", "admin123")
	fmt.Printf("%-20s %-25s %s\n", "Wali Kelas 7A", "walikelas.7a", "guru123")
	fmt.Printf("%-20s %-25s %s\n", "Wali Kelas 7B", "walikelas.7b", "guru123")
	fmt.Printf("%-20s %-25s %s\n", "Guru BK", "gurubk.smpn1", "guru123")
	fmt.Printf("%-20s %-25s %s\n", "Parent (Ahmad)", "parent.ahmad", "parent123")
	fmt.Printf("%-20s %-25s %s\n", "Parent (Budi)", "parent.budi", "parent123")
	fmt.Println(repeatStr("-", 60))
	fmt.Println("\nDevice API Key for ESP32:")
	fmt.Printf("Device Code: %s\n", device.DeviceCode)
	fmt.Printf("API Key: %s\n", device.APIKey)
	fmt.Println(repeatStr("=", 60))
}

// Helper function for string repeat
func repeatStr(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
