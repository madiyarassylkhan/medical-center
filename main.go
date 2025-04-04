package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	impl "medical-center/internal/gorm"
	"medical-center/internal/handler"
	"medical-center/internal/models/appointment"
	"medical-center/internal/models/department"
	"medical-center/internal/models/doctor"
	"medical-center/internal/models/schedule"

	"medical-center/internal/service"
)

func main() {
	db, err := gorm.Open(postgres.Open("postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	db.AutoMigrate(
		&department.Department{},
		&doctor.Doctor{},
		&schedule.Schedule{},
		&appointment.Appointment{},
	)

	deptRepo := impl.NewDepartmentRepository(db)
	doctorRepo := impl.NewDoctorRepository(db)
	scheduleRepo := impl.NewScheduleRepository(db)
	appointmentRepo := impl.NewAppoinmentRepository(db)

	deptService := service.NewDepartmentService(deptRepo)
	doctorService := service.NewDoctorService(doctorRepo)
	scheduleService := service.NewScheduleService(scheduleRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo)

	deptHandler := handler.NewDepartmentHandler(deptService)
	doctorHandler := handler.NewDoctorHandler(doctorService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.POST("/departments", deptHandler.CreateDepartment)
		api.GET("/departments", deptHandler.GetAllDepartments)
		api.GET("/departments/:id", deptHandler.GetDepartment)
		api.PUT("/departments/:id", deptHandler.UpdateDepartment)
		//api.DELETE("/departments/:id", deptHandler.DeleteDepartment)
		api.GET("/departments/:id/slots", deptHandler.GetDepartmentSlots)

		api.POST("/doctors", doctorHandler.CreateDoctor)
		api.GET("/doctors", doctorHandler.GetAllDoctors)
		api.GET("/doctors/:id", doctorHandler.GetDoctor)
		api.PUT("/doctors/:id", doctorHandler.UpdateDoctor)
		//api.DELETE("/doctors/:id", doctorHandler.DeleteDoctor)
		api.PATCH("/doctors/:id/availability", doctorHandler.SetAvailability)

		api.POST("/schedules", scheduleHandler.CreateSlot)
		api.GET("/schedules/:id", scheduleHandler.GetSlot)
		api.GET("/schedules/doctor/:doctor_id", scheduleHandler.GetDoctorSlots)
		api.POST("/schedules/:id/book", scheduleHandler.BookSlot)
		api.GET("/schedules/available", scheduleHandler.GetAvailableSlots)

		api.POST("/appointments", appointmentHandler.CreateAppointment)
		api.GET("/appointments", appointmentHandler.GetAllAppointments)
		api.GET("/appointments/:id", appointmentHandler.GetAppointment)
		api.PUT("/appointments/:id", appointmentHandler.UpdateAppointment)
		//api.DELETE("/appointments/:id", appointmentHandler.DeleteAppointment)
		api.GET("/appointments/department/:department_id", appointmentHandler.GetAppointmentsByDepartment)
	}

	router.Run(":8080")
}
