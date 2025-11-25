package route

import (
	"database/sql"
	"uas-pbe-praksem5/app/model"
	"uas-pbe-praksem5/app/repository"
	"uas-pbe-praksem5/app/service"
	"uas-pbe-praksem5/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAll(app *fiber.App, db *sql.DB) {
	// repos
	userRepo := repository.NewUserRepo(db)
	refreshRepo := repository.NewRefreshTokenRepo(db)
	studentRepo := repository.NewStudentRepo(db)
	lecturerRepo := repository.NewLecturerRepo(db)
	achRefRepo := repository.NewAchievementRefRepo(db)

	// services
	userSvc := service.NewUserService(userRepo, db)
	authSvc := service.NewAuthService(userRepo, refreshRepo)
	studentSvc := service.NewStudentService(studentRepo)
	lecturerSvc := service.NewLecturerService(lecturerRepo)
	achRefSvc := service.NewAchievementRefService(achRefRepo)

	v1 := app.Group("/api/v1")

	// auth
	v1.Post("/auth/login", func(c *fiber.Ctx) error {
		var req model.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		access, refresh, user, err := authSvc.Login(req)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if access == "" {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}
		return c.JSON(fiber.Map{"access_token": access, "refresh_token": refresh, "user": user})
	})

	v1.Post("/auth/refresh", func(c *fiber.Ctx) error {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		access, refresh, err := authSvc.Refresh(body.RefreshToken)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"access_token": access, "refresh_token": refresh})
	})

	v1.Post("/auth/logout", func(c *fiber.Ctx) error {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if err := authSvc.Logout(body.RefreshToken); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "logged out"})
	})

	v1.Get("/auth/profile", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		uid := c.Locals("user_id").(string)
		user, err := userSvc.GetByID(uid)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if user == nil {
			return c.Status(404).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(user)
	})

	// users (admin)
	users := v1.Group("/users", middleware.AuthRequired(), middleware.AdminOnly())
	users.Get("/", func(c *fiber.Ctx) error {
		res, err := userSvc.ListUsers()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(res)
	})
	users.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		u, err := userSvc.GetByID(id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if u == nil {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		}
		return c.JSON(u)
	})
	users.Post("/", func(c *fiber.Ctx) error {
		var req model.CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if err := userSvc.CreateUser(req); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(201).JSON(fiber.Map{"message": "user created"})
	})
	users.Put("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var req model.CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if err := userSvc.UpdateUser(id, req); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "user updated"})
	})
	users.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := userSvc.DeleteUser(id); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "user deleted"})
	})
	users.Put("/:id/role", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var b struct {
			RoleID string `json:"role_id"`
		}
		if err := c.BodyParser(&b); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if err := userSvc.UpdateRole(id, b.RoleID); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "role updated"})
	})

	// students (admin)
	stu := v1.Group("/students", middleware.AuthRequired(), middleware.AdminOnly())
	stu.Post("/", func(c *fiber.Ctx) error {
		var s model.Student
		if err := c.BodyParser(&s); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		// simple create using repo directly via service or expand service
		if err := studentSvc.Repo.Create(s); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(201).JSON(fiber.Map{"message": "student created"})
	})

	// lecturers (admin)
	lec := v1.Group("/lecturers", middleware.AuthRequired(), middleware.AdminOnly())
	lec.Post("/", func(c *fiber.Ctx) error {
		var l model.Lecturer
		if err := c.BodyParser(&l); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if err := lecturerSvc.Repo.Create(l); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(201).JSON(fiber.Map{"message": "lecturer created"})
	})

	// achievement refs (student create -> later link to mongo)
	ach := v1.Group("/achievements", middleware.AuthRequired())
	ach.Post("/", func(c *fiber.Ctx) error {
		b := struct {
			StudentID string `json:"student_id"`
			MongoID   string `json:"mongo_id"`
		}{}
		if err := c.BodyParser(&b); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if err := achRefSvc.CreateRef(b.StudentID, b.MongoID); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(201).JSON(fiber.Map{"message": "achievement reference created"})
	})
	ach.Post("/:id/submit", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := achRefSvc.UpdateStatus(id, "submitted", nil, nil); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "submitted"})
	})
	ach.Post("/:id/verify", middleware.RoleIs("Dosen Wali"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		verifier := c.Locals("user_id").(string)
		if err := achRefSvc.UpdateStatus(id, "verified", &verifier, nil); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "verified"})
	})
	ach.Post("/:id/reject", middleware.RoleIs("Dosen Wali"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		body := struct {
			Note string `json:"note"`
		}{}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if err := achRefSvc.UpdateStatus(id, "rejected", nil, &body.Note); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "rejected"})
	})
}
