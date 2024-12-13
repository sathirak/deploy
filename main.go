package main

import (
    "encoding/json"
    "html/template"
    "log"
    "net/http"
    "time"
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/disk"
    "github.com/shirou/gopsutil/v3/mem"
    "runtime"
)

type Response struct {
    Message string `json:"message"`
    Status  bool   `json:"status"`
}

type LogRequest struct {
    Timestamp   string `json:"timestamp"`
    Method      string `json:"method"`
    Path        string `json:"path"`
    RemoteAddr  string `json:"remote_addr"`
    UserAgent   string `json:"user_agent"`
    StatusCode  int    `json:"status_code"`
}

type SystemStats struct {
    CPUUsage    float64 `json:"cpu_usage"`
    MemoryTotal uint64  `json:"memory_total"`
    MemoryUsed  uint64  `json:"memory_used"`
    DiskTotal   uint64  `json:"disk_total"`
    DiskFree    uint64  `json:"disk_free"`
    GoVersion   string  `json:"go_version"`
    NumCPU      int     `json:"num_cpu"`
}

func getSystemStats() (SystemStats, error) {
    v, _ := mem.VirtualMemory()
    c, _ := cpu.Percent(0, false)
    d, _ := disk.Usage("/")

    return SystemStats{
        CPUUsage:    c[0],
        MemoryTotal: v.Total,
        MemoryUsed:  v.Used,
        DiskTotal:   d.Total,
        DiskFree:    d.Free,
        GoVersion:   runtime.Version(),
        NumCPU:      runtime.NumCPU(),
    }, nil
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        logEntry := LogRequest{
            Timestamp:   time.Now().Format(time.RFC3339),
            Method:      r.Method,
            Path:        r.URL.Path,
            RemoteAddr:  r.RemoteAddr,
            UserAgent:   r.UserAgent(),
            StatusCode:  http.StatusOK,
        }
        
        logJSON, _ := json.Marshal(logEntry)
        log.Printf("%s\n", logJSON)
        
        next(w, r)
    }
}

var templateFuncs = template.FuncMap{
    "div": func(a, b float64) float64 {
        if b == 0 {
            return 0
        }
        return a / b * 100 // Multiply by 100 to get percentage
    },
    "sub": func(a, b uint64) uint64 {
        return a - b
    },
    "float64": func(u uint64) float64 {
        return float64(u)
    },
}

func main() {
    http.HandleFunc("/", loggingMiddleware(handleHome))
    http.HandleFunc("/json", loggingMiddleware(handleJSON))

    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    stats, _ := getSystemStats()
    tmpl := template.Must(template.New("index.html").Funcs(templateFuncs).ParseFiles("templates/index.html"))
    tmpl.Execute(w, stats)
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
    stats, _ := getSystemStats()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}