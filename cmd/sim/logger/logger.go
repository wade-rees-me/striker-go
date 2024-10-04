package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"log"
	"github.com/google/uuid"
)

type Logger struct {
	simulator       string
	directory       string
	subdirectory    string
	guid            string
	simulationFile  *os.File
	debugFile       *os.File
	mu              sync.Mutex
}

// NewLogger creates a new Logger
func NewLogger(simulator string, debugFlag bool) *Logger {
	logger := &Logger{
		simulator: simulator,
		directory: getSimulationDirectory(),
	}

	logger.getSubdirectory()

	fullPath := filepath.Join(logger.directory, logger.subdirectory)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			log.Fatalf("Failed to create directory: %s", err)
		}
	}

	// Open the simulation file
	logger.simulationFile = logger.openFile("simulation")
	if debugFlag {
		logger.debugFile = logger.openFile("debug")
	}

	return logger
}

// Destructor for Logger
func (logger *Logger) Close() {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if logger.simulationFile != nil {
		logger.simulationFile.Close()
	}

	if logger.debugFile != nil {
		logger.debugFile.Close()
	}
}

// Log to simulation file
func (logger *Logger) Simulation(message string) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if logger.simulationFile != nil {
		logger.simulationFile.WriteString(message)
		fmt.Print(message) // Also print to console
	}

	if logger.debugFile != nil {
		logger.debugFile.WriteString(message)
	}
}

// Log to debug file
func (logger *Logger) Debug(message string) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if logger.debugFile != nil {
		logger.debugFile.WriteString(message)
	}
}

// Insert to a separate insert log
func (logger *Logger) Insert(message string) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	insertFileName := filepath.Join(logger.directory, logger.subdirectory, fmt.Sprintf("%s_%s_insert.txt", logger.simulator, logger.guid))
	insertFile, err := os.Create(insertFileName)
	if err != nil {
		log.Fatalf("Failed to create insert file: %s", err)
	}
	defer insertFile.Close()

	insertFile.WriteString(message)
}

// Create subdirectory based on date and generate a GUID
func (logger *Logger) getSubdirectory() {
	// Get current time and format it as "YYYY_MM_DD"
	now := time.Now()
	logger.subdirectory = fmt.Sprintf("%04d_%02d_%02d", now.Year(), now.Month(), now.Day())

	// Generate a UUID for GUID
	logger.guid = uuid.New().String()
}

// Helper function to open files
func (logger *Logger) openFile(fileType string) *os.File {
	fileName := filepath.Join(logger.directory, logger.subdirectory, fmt.Sprintf("%s_%s_%s.txt", logger.simulator, logger.guid, fileType))
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	return file
}

// Get the simulation directory path
func getSimulationDirectory() string {
	// This function should return the path where you want to store the logs.
	// For simplicity, we use the current directory here.
	return "/tmp/simulations"
}

/*
func main() {
	// Example usage of the logger
	logger := NewLogger("striker-go", true)
	defer logger.Close()

	logger.Simulation("Simulation log message.\n")
	logger.Debug("Debug log message.\n")
	logger.Insert("Insert log message.\n")
}
*/
