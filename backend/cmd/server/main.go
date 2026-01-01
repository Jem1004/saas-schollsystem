package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/school-management/backend/internal/config"
	"github.com/school-management/backend/internal/domain/models"
	"github.com/school-management/backend/internal/middleware"
	"github.com/school-management/backend/internal/modules/attendance"
	"github.com/school-management/backend/internal/modules/auth"
	"github.com/school-management/backend/internal/modules/bk"
	"github.com/school-management/backend/internal/modules/device"
	"github.com/school-management/backend/internal/modules/displaytoken"
	"github.com/school-management/backend/internal/modules/grade"
	"github.com/school-management/backend/internal/modules/homeroom"
	"github.com/school-management/backend/internal/modules/notification"
	"github.com/school-management/backend/internal/modules/parent"
	"github.com/school-management/backend/internal/modules/publicdisplay"
	"github.com/school-management/backend/internal/modules/realtime"
	"github.com/school-management/backend/internal/modules/schedule"
	"github.com/school-management/backend/internal/modules/school"
	"github.com/school-management/backend/internal/modules/settings"
	"github.com/school-management/backend/internal/modules/student"
	"github.com/school-management/backend/internal/modules/tenant"
	"github.com/school-management/backend/internal/shared/database"
	"github.com/school-management/backend/internal/shared/fcm"
	"github.com/school-management/backend/internal/shared/redis"
)

func main() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	// Initialize Redis connection
	redisClient, err := redis.Connect(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected successfully")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "School Management SaaS API",
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Tenant-ID",
		AllowCredentials: true,
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "school-management-api",
		})
	})

	// API routes group
	api := app.Group("/api/v1")

	// Initialize JWT Manager
	jwtManager := auth.NewJWTManager(cfg.JWT)

	// Initialize Auth Module
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, jwtManager)
	authHandler := auth.NewHandler(authService)

	// Register public auth routes (login, refresh, logout)
	authHandler.RegisterRoutes(api)

	// Protected routes group with auth middleware
	protected := api.Group("", middleware.AuthMiddleware(jwtManager))

	// Register protected auth routes (change-password, me)
	authHandler.RegisterProtectedRoutes(protected)

	// Initialize Tenant Module (Super Admin only)
	tenantRepo := tenant.NewRepository(db)
	tenantService := tenant.NewService(tenantRepo)
	tenantHandler := tenant.NewHandler(tenantService)

	// Initialize Device Module (Super Admin only)
	deviceRepo := device.NewRepository(db)
	deviceService := device.NewService(deviceRepo)
	deviceHandler := device.NewHandler(deviceService)

	// Initialize Pairing Module (for RFID card pairing)
	deviceStudentRepo := device.NewStudentRepository(db)
	pairingService := device.NewPairingService(deviceRepo, deviceStudentRepo)
	pairingHandler := device.NewPairingHandler(pairingService)

	// Super Admin routes - use specific path prefixes to avoid middleware conflicts
	// Schools management (Super Admin only)
	schoolsAdmin := protected.Group("/schools", middleware.SuperAdminOnly())
	tenantHandler.RegisterRoutesWithoutGroup(schoolsAdmin)

	// Devices management (Super Admin only)
	devicesAdmin := protected.Group("/devices", middleware.SuperAdminOnly())
	deviceHandler.RegisterRoutesWithoutGroup(devicesAdmin)

	// Public device routes (for ESP32 API key validation)
	deviceHandler.RegisterPublicRoutes(api)

	// Public pairing routes (for ESP32 RFID pairing)
	pairingHandler.RegisterPublicRoutes(api)

	// Tenant-scoped routes (for non-super_admin users)
	tenantScoped := protected.Group("", middleware.TenantMiddleware())

	// Initialize Settings Module FIRST to avoid route conflicts
	// Requirements: School Settings - attendance time, notification toggles, academic year
	settingsRepo := settings.NewRepository(db)
	settingsService := settings.NewService(settingsRepo)
	settingsHandler := settings.NewHandler(settingsService)

	// Settings routes for Admin Sekolah (full access to school settings)
	// Role check is done inside the handler
	settingsHandler.RegisterRoutes(tenantScoped)

	// Pairing routes for admin sekolah (start/cancel/status pairing)
	pairingRoutes := tenantScoped.Group("/pairing")
	pairingHandler.RegisterRoutesWithoutGroup(pairingRoutes)

	// Initialize School Module (Admin Sekolah)
	schoolRepo := school.NewRepository(db)
	schoolUserRepo := school.NewUserRepository(db)
	schoolService := school.NewService(schoolRepo, schoolUserRepo)
	schoolHandler := school.NewHandler(schoolService)

	// School routes for admin sekolah (classes, students, parents)
	adminSekolahRoutes := tenantScoped.Group("/school")
	schoolHandler.RegisterRoutes(adminSekolahRoutes)

	// Initialize Schedule Module (Admin Sekolah only)
	// Requirements: 3.1 - Multi-schedule support for different activities
	scheduleRepo := schedule.NewRepository(db)
	scheduleService := schedule.NewService(scheduleRepo)
	scheduleHandler := schedule.NewHandler(scheduleService)

	// Schedule routes for admin sekolah only
	scheduleRoutes := tenantScoped.Group("/schedules", middleware.AdminSekolahOnly())
	scheduleHandler.RegisterRoutesWithoutGroup(scheduleRoutes)

	// Initialize Display Token Module (Admin Sekolah only)
	// Requirements: 5.1, 6.1 - Display token management for public display access
	displayTokenRepo := displaytoken.NewRepository(db)
	displayTokenService := displaytoken.NewService(displayTokenRepo)
	displayTokenHandler := displaytoken.NewHandler(displayTokenService)

	// Display token routes for admin sekolah only
	displayTokenRoutes := tenantScoped.Group("/display-tokens", middleware.AdminSekolahOnly())
	displayTokenHandler.RegisterRoutesWithoutGroup(displayTokenRoutes)

	// Initialize Attendance Module
	attendanceRepo := attendance.NewRepository(db)
	attendancePolicy := attendance.NewAttendancePolicy(db)
	attendanceService := attendance.NewService(attendanceRepo, deviceService, attendancePolicy)
	attendanceHandler := attendance.NewHandler(attendanceService, attendanceRepo)

	// Initialize Real-Time Module
	// Requirements: 4.1, 4.2, 4.3 - Real-time attendance dashboard with WebSocket
	realtimeHub := realtime.NewHub()
	go realtimeHub.Run() // Start the hub in a goroutine
	realtimeRepo := realtime.NewRepository(db)
	realtimeService := realtime.NewService(realtimeRepo, realtimeHub)
	realtimeHandler := realtime.NewHandler(realtimeService, jwtManager)

	// Connect attendance service to real-time broadcaster
	// Requirements: 4.2 - Broadcast attendance updates in real-time
	attendanceService.SetRealtimeBroadcaster(realtimeService)

	// Register WebSocket routes (before other routes to handle upgrade properly)
	realtimeHandler.RegisterWebSocketRoutes(app)

	// Public attendance routes (for ESP32 RFID devices)
	attendanceHandler.RegisterPublicRoutes(api)

	// Attendance routes for admin sekolah and wali kelas
	attendanceRoutes := tenantScoped.Group("/attendance")
	attendanceHandler.RegisterRoutesWithoutGroup(attendanceRoutes)

	// Real-time routes for authenticated users
	// Requirements: 4.1 - Live attendance dashboard
	realtimeRoutes := tenantScoped.Group("/realtime")
	realtimeHandler.RegisterRoutes(realtimeRoutes)

	// Initialize Public Display Module
	// Requirements: 5.3, 5.4, 5.5, 5.6 - Public display for LCD screens without login
	publicDisplayRepo := publicdisplay.NewRepository(db)
	publicDisplayService := publicdisplay.NewService(publicDisplayRepo, displayTokenService)
	publicDisplayHandler := publicdisplay.NewHandler(publicDisplayService, realtimeHub, realtimeRepo)

	// Public display routes (no auth required - uses display token)
	publicDisplayHandler.RegisterPublicRoutes(api)

	// Initialize BK Module
	bkRepo := bk.NewRepository(db)
	bkService := bk.NewService(bkRepo)
	bkHandler := bk.NewHandler(bkService)

	// BK routes for Guru BK (full access)
	// Requirements: 6.1-6.5, 7.1-7.5, 8.1-8.5, 9.1-9.5
	bkRoutes := tenantScoped.Group("/bk")
	bkHandler.RegisterRoutesWithoutGroup(bkRoutes)

	// Initialize Grade Module
	// Requirements: 10.1, 10.2, 10.4, 10.5
	gradeRepo := grade.NewRepository(db)
	gradeService := grade.NewService(gradeRepo, db)
	gradeHandler := grade.NewHandler(gradeService)

	// Grade routes for Wali Kelas (full access to their class)
	gradeRoutes := tenantScoped.Group("/grades")
	gradeHandler.RegisterRoutesWithoutGroup(gradeRoutes)

	// Initialize Homeroom Module
	// Requirements: 11.1, 11.3, 11.4, 11.5
	homeroomRepo := homeroom.NewRepository(db)
	homeroomService := homeroom.NewService(homeroomRepo, db)
	homeroomHandler := homeroom.NewHandler(homeroomService)

	// Homeroom routes for Wali Kelas (full access to their class)
	homeroomRoutes := tenantScoped.Group("/homeroom")
	homeroomHandler.RegisterRoutesWithoutGroup(homeroomRoutes)

	// Initialize FCM Client
	// Requirements: 13.1, 13.2 - Firebase Cloud Messaging integration
	fcmClient, err := fcm.NewClient(cfg.FCM)
	if err != nil {
		log.Printf("Warning: Failed to initialize FCM client: %v", err)
		fcmClient = &fcm.Client{} // Use empty client
	}
	if fcmClient.IsInitialized() {
		log.Println("FCM client initialized successfully")
	} else {
		log.Println("FCM client not configured, push notifications disabled")
	}

	// Initialize Notification Module
	// Requirements: 17.1, 17.2, 17.3, 17.4, 17.5 - Notification system with queue and FCM
	notificationRepo := notification.NewRepository(db)
	notificationService := notification.NewService(notificationRepo, redisClient)
	notificationHandler := notification.NewHandler(notificationService)

	// Notification routes (accessible by all authenticated users)
	notificationHandler.RegisterRoutes(protected)

	// Initialize Parent Module
	// Requirements: 12.2, 14.4, 15.1, 15.2 - Parent data access for linked children
	parentRepo := parent.NewRepository(db)
	parentService := parent.NewService(parentRepo)
	parentHandler := parent.NewHandler(parentService)

	// Parent routes (accessible by parents only)
	parentRoutes := protected.Group("", middleware.RoleMiddleware(
		models.RoleParent,
	))
	parentHandler.RegisterRoutes(parentRoutes)

	// Initialize Student Module
	// Requirements: 16.1-16.5 - Student self-monitoring
	studentRepo := student.NewRepository(db)
	studentService := student.NewService(studentRepo)
	studentHandler := student.NewHandler(studentService)

	// Student routes (accessible by students only)
	studentRoutes := protected.Group("", middleware.RoleMiddleware(
		models.RoleStudent,
	))
	studentHandler.RegisterRoutes(studentRoutes)

	// Initialize and start Notification Worker
	// Requirements: 17.1, 17.2, 17.5 - Background queue processing with retry
	notificationWorker := notification.NewWorker(redisClient, fcmClient, notificationRepo)
	notificationWorker.Start()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		// Stop notification worker
		notificationWorker.Stop()

		if err := app.Shutdown(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	// Start server
	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// customErrorHandler handles all errors in a consistent format
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	})
}
