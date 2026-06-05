package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/internal/service"
)

type GenId interface {
	GenerateId(c *gin.Context)
	//cfg *api.Config
}

type genIdHandler struct {
	genIdService service.GenId
	cfg          *config.Config
}

func NewGenId(genIdSvc service.GenId, cfg *config.Config) GenId {
	return &genIdHandler{
		genIdService: genIdSvc,
		cfg:          cfg,
	}
}

// GenerateId generates a new ID
// @Summary Generate a new ID
// @Description Generate a new ID
// @Tags genid
// @Success 200 {object} map[string]interface{}
// @Router /health-check [get]
func (s *genIdHandler) GenerateId(c *gin.Context) {
	service_name := s.cfg.ServiceName
	instance_id := s.cfg.InstanceId
	if len(instance_id) == 0 {
		instance_id = s.genIdService.GenerateId()
	}

	if instance_id == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Err"})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"message": "OK",
			"service_name": service_name,
			"instance_id":  instance_id,
		})
}
