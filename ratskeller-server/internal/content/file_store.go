package content

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
)

type FileStore struct {
	configDir string
	mu        sync.RWMutex
}

func NewFileStore(configDir string) *FileStore {
	return &FileStore{
		configDir: configDir,
	}
}

func (s *FileStore) GetSite(c *gin.Context) {
	s.readJSON(c, "site.json")
}

func (s *FileStore) UpdateSite(c *gin.Context) {
	s.writeJSON(c, "site.json")
}

func (s *FileStore) GetOpeningHours(c *gin.Context) {
	s.readJSON(c, "opening_hours.json")
}

func (s *FileStore) UpdateOpeningHours(c *gin.Context) {
	s.writeJSON(c, "opening_hours.json")
}

func (s *FileStore) readJSON(c *gin.Context, filename string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.configDir, filename)

	data, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "config_read_failed",
		})
		return
	}

	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "config_parse_failed",
		})
		return
	}

	c.JSON(http.StatusOK, payload)
}

func (s *FileStore) writeJSON(c *gin.Context, filename string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var payload map[string]any
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_json_payload",
		})
		return
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "config_marshal_failed",
		})
		return
	}

	path := filepath.Join(s.configDir, filename)
	tmpPath := path + ".tmp"

	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "config_write_failed",
		})
		return
	}

	if err := os.Rename(tmpPath, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "config_replace_failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "updated",
	})
}
