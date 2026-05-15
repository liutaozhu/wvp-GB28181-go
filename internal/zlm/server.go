package zlm

import (
	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"

	"go.uber.org/zap"
)

// Server manages MediaServer nodes
type Server struct {
	log    *zap.Logger
	client *Client
}

func NewServer(log *zap.Logger, client *Client) *Server {
	return &Server{log: log, client: client}
}

// AutoConfig auto-configures ZLMediaKit via HTTP API
func (s *Server) AutoConfig() error {
	config, err := s.client.GetServerConfig()
	if err != nil {
		s.log.Warn("failed to get ZLM server config", zap.Error(err))
		return err
	}

	s.log.Info("ZLMediaKit auto-configured",
		zap.String("secret", config.API.Secret),
		zap.String("snapRoot", config.API.SnapRoot),
	)
	return nil
}

// GetMediaServerFromDB retrieves media server from database
func (s *Server) GetMediaServerFromDB(id string) (*model.MediaServer, error) {
	var ms model.MediaServer
	err := database.DB.Where("id = ?", id).First(&ms).Error
	return &ms, err
}

// ListMediaServers lists all media servers
func (s *Server) ListMediaServers() ([]model.MediaServer, error) {
	var servers []model.MediaServer
	err := database.DB.Find(&servers).Error
	return servers, err
}

// SaveMediaServer saves media server to database
func (s *Server) SaveMediaServer(ms *model.MediaServer) error {
	return database.DB.Save(ms).Error
}

// DeleteMediaServer deletes media server from database
func (s *Server) DeleteMediaServer(id string) error {
	return database.DB.Where("id = ?", id).Delete(&model.MediaServer{}).Error
}
