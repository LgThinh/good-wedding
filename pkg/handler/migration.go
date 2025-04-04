package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"net/http"
)

type MigrationHandler struct {
	db *gorm.DB
}

func NewMigrationHandler(db *gorm.DB) *MigrationHandler {
	return &MigrationHandler{db: db}
}

// MigratePublic godoc
// @Summary  Database migration
// @Tags     internal
// @Security Authorization
// @Security Access Token
// @Security Entity Key
// @Security User ID
// @Success  200
// @Router   /internal/migrate-public [post]
func (h *MigrationHandler) MigratePublic(ctx *gin.Context) {
	_ = h.db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	_ = h.db.Exec(`SET client_encoding = 'UTF8'`)

	sqlCommands := []string{
		`CREATE SCHEMA IF NOT EXISTS public`,
	}

	for _, sql := range sqlCommands {
		if err := h.db.Exec(sql).Error; err != nil {
			log.Printf(err.Error())
			h.db.Rollback()
			ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	h.db.Config.NamingStrategy = schema.NamingStrategy{
		TablePrefix: "public.",
	}

	models := []interface{}{
		//&model.Todo{},
	}
	if err := h.db.AutoMigrate(models...); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, "Migrate public schema success")
}

// MigrateLog godoc
// @Summary  Database log migration
// @Tags     internal
// @Security Authorization
// @Security Access Token
// @Security Entity Key
// @Security User ID
// @Success  200
// @Router   /internal/migrate-log [post]
func (h *MigrationHandler) MigrateLog(ctx *gin.Context) {
	_ = h.db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	_ = h.db.Exec(`SET client_encoding = 'UTF8'`)

	sqlCommands := []string{
		`CREATE SCHEMA IF NOT EXISTS log`,
	}

	for _, sql := range sqlCommands {
		if err := h.db.Exec(sql).Error; err != nil {
			log.Printf(err.Error())
			h.db.Rollback()
			ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	h.db.Config.NamingStrategy = schema.NamingStrategy{
		TablePrefix: "log.",
	}

	models := []interface{}{
		//&model.Todo{},
	}

	if err := h.db.AutoMigrate(models...); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "Migrate log schema success")
}
