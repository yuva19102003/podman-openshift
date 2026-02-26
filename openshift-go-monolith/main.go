package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
)

var (
	startTime    = time.Now()
	requestCount int64
	writeCount   int64
	logger       *log.Logger
)

type AppInfo struct {
	AppName   string    `json:"app_name"`
	Env       string    `json:"environment"`
	DBUser    string    `json:"db_user"`
	Version   string    `json:"version"`
	Hostname  string    `json:"hostname"`
	Timestamp time.Time `json:"timestamp"`
}

type Stats struct {
	Uptime         string `json:"uptime"`
	TotalRequests  int64  `json:"total_requests"`
	WriteOps       int64  `json:"write_operations"`
	GoVersion      string `json:"go_version"`
	NumGoroutines  int    `json:"goroutines"`
	MemoryAllocMB  uint64 `json:"memory_alloc_mb"`
	ServerTime     string `json:"server_time"`
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	logger.Printf("[INFO] 📊 Request received: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	hostname, err := os.Hostname()
	if err != nil {
		logger.Printf("[WARN] ⚠️ Failed to get hostname: %v", err)
		hostname = "unknown"
	}
	
	info := AppInfo{
		AppName:   getEnvOrDefault("APP_NAME", "OpenShift Go Monolith"),
		Env:       getEnvOrDefault("APP_ENV", "development"),
		DBUser:    getEnvOrDefault("DB_USER", "not_configured"),
		Version:   "1.1.0",
		Hostname:  hostname,
		Timestamp: time.Now(),
	}

	logger.Printf("[INFO] 📤 Sending app info response: AppName=%s, Env=%s, Hostname=%s", 
		info.AppName, info.Env, info.Hostname)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		logger.Printf("[ERROR] 💥 Failed to encode JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	logger.Printf("[INFO] ✅ App info request completed successfully - hits different!")
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	atomic.AddInt64(&writeCount, 1)
	
	logger.Printf("[INFO] 📝 Write request received: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	
	// Create log directory if it doesn't exist
	logDir := "./data/log"
	logger.Printf("[DEBUG] 🔍 Ensuring log directory exists: %s", logDir)
	
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logger.Printf("[ERROR] 🚨 Failed to create log directory %s: %v", logDir, err)
		http.Error(w, fmt.Sprintf("Failed to create log directory: %v", err), http.StatusInternalServerError)
		return
	}
	logger.Printf("[DEBUG] ✅ Log directory ready: %s", logDir)

	// Create timestamped log file
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s-log.txt", timestamp)
	filepath := filepath.Join(logDir, filename)
	
	logger.Printf("[INFO] 📄 Creating log file: %s", filepath)

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logger.Printf("[ERROR] 💥 Failed to create log file %s: %v", filepath, err)
		http.Error(w, fmt.Sprintf("Failed to create log file: %v", err), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Write detailed log content with Gen Z vibes
	hostname, _ := os.Hostname()
	appName := getEnvOrDefault("APP_NAME", "OpenShift Go Monolith")
	env := getEnvOrDefault("APP_ENV", "development")
	
	logContent := fmt.Sprintf(`========================================
🚀 OpenShift Go Monolith - Volume Write Log
========================================

⏰ Timestamp:        %s
🔢 Operation Number: %d
📦 Application:      %s
🌍 Environment:      %s
🏠 Hostname:         %s
🌐 Client IP:        %s
🐹 Go Version:       %s
📊 Total Requests:   %d
⏱️  Uptime:           %s

========================================
📝 Log Entry Details
========================================

This log file was created as part of write operation #%d.
The application successfully wrote data to the persistent volume.
No cap, this is bussin fr fr! 💯

🖥️  System Information:
- Number of Goroutines: %d
- Memory Allocated: %d MB
- Status: Running smooth like butter 🧈

📡 Request Information:
- Method: %s
- Path: %s
- User Agent: %s
- Remote Address: %s

💭 Vibes: Immaculate ✨
🎯 Status: Mission accomplished, chief! 
🔥 Performance: Absolutely slaying rn

========================================
✅ End of Log - Stay hydrated! 💧
========================================
`,
		time.Now().Format(time.RFC3339),
		atomic.LoadInt64(&writeCount),
		appName,
		env,
		hostname,
		r.RemoteAddr,
		runtime.Version(),
		atomic.LoadInt64(&requestCount),
		time.Since(startTime).Round(time.Second).String(),
		atomic.LoadInt64(&writeCount),
		runtime.NumGoroutine(),
		getMemoryUsageMB(),
		r.Method,
		r.URL.Path,
		r.UserAgent(),
		r.RemoteAddr,
	)

	logger.Printf("[DEBUG] 💾 Writing %d bytes to log file", len(logContent))
	
	if _, err := f.WriteString(logContent); err != nil {
		logger.Printf("[ERROR] 😱 Failed to write content to log file %s: %v", filepath, err)
		http.Error(w, fmt.Sprintf("Failed to write log content: %v", err), http.StatusInternalServerError)
		return
	}

	logger.Printf("[INFO] 🎉 Successfully wrote log file: %s - it's giving main character energy!", filepath)

	response := fmt.Sprintf(`✓ Data written to volume successfully

📁 File: %s
🔢 Operation: #%d
⏰ Timestamp: %s
📏 Size: %d bytes

📂 Log directory: %s

💯 Status: Absolutely fire! No printer, just facts! 🔥`, 
		filename,
		atomic.LoadInt64(&writeCount),
		time.Now().Format(time.RFC3339),
		len(logContent),
		logDir)
	
	logger.Printf("[INFO] ✨ Write operation completed successfully - we're so back!")
	w.Write([]byte(response))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	logger.Printf("[INFO] ❤️ Health check request from %s - checking the vibes...", r.RemoteAddr)
	w.Write([]byte("OK"))
	logger.Printf("[DEBUG] 💚 Health check response sent - we're thriving!")
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	logger.Printf("[INFO] 📈 Stats request received: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	stats := Stats{
		Uptime:         time.Since(startTime).Round(time.Second).String(),
		TotalRequests:  atomic.LoadInt64(&requestCount),
		WriteOps:       atomic.LoadInt64(&writeCount),
		GoVersion:      runtime.Version(),
		NumGoroutines:  runtime.NumGoroutine(),
		MemoryAllocMB:  getMemoryUsageMB(),
		ServerTime:     time.Now().Format(time.RFC3339),
	}

	logger.Printf("[DEBUG] 📊 Stats collected: Uptime=%s, Requests=%d, WriteOps=%d, Memory=%dMB - looking good!", 
		stats.Uptime, stats.TotalRequests, stats.WriteOps, stats.MemoryAllocMB)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		logger.Printf("[ERROR] 😱 Failed to encode stats JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	logger.Printf("[INFO] ✨ Stats request completed successfully - data is immaculate!")
}

func getMemoryUsageMB() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc / 1024 / 1024
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Printf("[REQUEST] 🌐 %s %s from %s - User-Agent: %s", 
			r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		
		next.ServeHTTP(w, r)
		
		duration := time.Since(start)
		logger.Printf("[RESPONSE] ⚡ %s %s completed in %v - speedrun any%%", r.Method, r.URL.Path, duration)
	})
}

func initLogger() {
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	logger.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	logger.Println("[INIT] 🎯 Logger initialized with detailed output - let's get this bread!")
}

func main() {
	// Initialize logger first
	initLogger()
	
	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Printf("[WARN] ⚠️ No .env file found or error loading it: %v", err)
		logger.Println("[INFO] 📝 Using system environment variables or defaults")
	} else {
		logger.Println("[INFO] ✅ Successfully loaded .env file")
	}
	
	logger.Println("========================================")
	logger.Println("🚀 OpenShift Go Monolith Server")
	logger.Println("========================================")
	logger.Printf("[INIT] 💫 Version: 1.1.0")
	logger.Printf("[INIT] 🐹 Go Version: %s", runtime.Version())
	logger.Printf("[INIT] 💻 OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logger.Printf("[INIT] ⚡ CPUs: %d", runtime.NumCPU())
	logger.Printf("[INIT] ⏰ Started at: %s", time.Now().Format(time.RFC3339))
	
	// Log environment variables
	logger.Printf("[CONFIG] 📦 APP_NAME: %s", getEnvOrDefault("APP_NAME", "not set"))
	logger.Printf("[CONFIG] 🌍 APP_ENV: %s", getEnvOrDefault("APP_ENV", "not set"))
	logger.Printf("[CONFIG] 👤 DB_USER: %s", getEnvOrDefault("DB_USER", "not set"))
	
	hostname, err := os.Hostname()
	if err != nil {
		logger.Printf("[WARN] ⚠️ Failed to get hostname: %v", err)
	} else {
		logger.Printf("[CONFIG] 🏠 Hostname: %s", hostname)
	}
	
	// Check data directory
	dataDir := "./data/log"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		logger.Printf("[WARN] 📁 Data directory %s does not exist, will be created on first write", dataDir)
	} else {
		logger.Printf("[INFO] ✅ Data directory %s exists and is accessible", dataDir)
	}
	
	// Setup routes with logging middleware
	logger.Println("[INIT] 🔧 Registering HTTP handlers...")
	
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/api/info", infoHandler)
	mux.HandleFunc("/api/write", writeHandler)
	mux.HandleFunc("/api/stats", statsHandler)
	mux.HandleFunc("/health", healthHandler)
	
	logger.Println("[INIT] 🛣️ Routes registered:")
	logger.Println("[INIT]   📄 GET  /              - Static files")
	logger.Println("[INIT]   📊 GET  /api/info      - Application info")
	logger.Println("[INIT]   💾 POST /api/write     - Write volume data")
	logger.Println("[INIT]   📈 GET  /api/stats     - Application statistics")
	logger.Println("[INIT]   ❤️ GET  /health        - Health check")
	
	// Wrap with logging middleware
	handler := loggingMiddleware(mux)
	
	logger.Println("========================================")
	logger.Printf("[INIT] 🎧 Server listening on :8080")
	logger.Println("[INIT] ✨ Ready to accept connections - let's goooo!")
	logger.Println("========================================")
	
	if err := http.ListenAndServe(":8080", handler); err != nil {
		logger.Printf("[FATAL] 💀 Server failed to start: %v", err)
		os.Exit(1)
	}
}