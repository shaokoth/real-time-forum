package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	MaxUploadSize = 10 << 20 // 10MB
	UploadsDir    = "./uploads"
)

func init() {
	// Create the uploads directory if it doesn't exist
	if _, err := os.Stat(UploadsDir); os.IsNotExist(err) {
		os.Mkdir(UploadsDir, 0o755) // Read/write for owner, read for others
	}
}

func HandleImageUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, "Image file too big", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create unique filename
	fileExtension := strings.ToLower(filepath.Ext(handler.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	if !allowedExtensions[fileExtension] {
		http.Error(w, "Invalid file type. Only JPG, JPEG, PNG, GIF are allowed.", http.StatusBadRequest)
		return
	}
	newFilename := fmt.Sprintf("%s-%d%s", uuid.New().String(), time.Now().UnixNano(), fileExtension)
	filepath := UploadsDir + newFilename

	// Save file
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Could not copy file", http.StatusInternalServerError)
		return
	}

	imageURL := "/uploads/" + newFilename

	// Respond with success and path
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"image_url": imageURL,
	})
}
