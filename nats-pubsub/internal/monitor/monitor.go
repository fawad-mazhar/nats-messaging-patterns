package monitor

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "os"
    "time"
)

// Monitor fetches and displays information from the NATS monitoring interface
type Monitor struct {
    baseURL string
    client  *http.Client
}

// ServerInfo holds basic NATS server information
type ServerInfo struct {
    ServerID    string  `json:"server_id"`
    Version     string  `json:"version"`
    GoVersion   string  `json:"go"`
    Host        string  `json:"host"`
    Port        int     `json:"port"`
    HTTPPort    int     `json:"http_port"`
    Uptime      string  `json:"uptime"`
    Mem         int64   `json:"mem"`
    Cores       int     `json:"cores"`
    CPU         float64 `json:"cpu"`
    Connections int     `json:"connections"`
}

// JetStreamResponse represents the top-level response from the JetStream API
type JetStreamResponse struct {
    Memory         int64           `json:"memory"`
    Storage        int64           `json:"storage"`
    Streams        int             `json:"streams"`
    Consumers      int             `json:"consumers"`
    Messages       int64           `json:"messages"`
    Bytes          int64           `json:"bytes"`
    AccountDetails []AccountDetail `json:"account_details"`
}

// AccountDetail represents account-specific information
type AccountDetail struct {
    Name         string         `json:"name"`
    StreamDetail []StreamDetail `json:"stream_detail"`
}

// StreamDetail represents detailed stream information
type StreamDetail struct {
    Name    string      `json:"name"`
    Created string      `json:"created"`
    State   StreamState `json:"state"`
}

// StreamState represents the current state of a stream
type StreamState struct {
    Messages      int64     `json:"messages"`
    Bytes         int64     `json:"bytes"`
    FirstSeq      int64     `json:"first_seq"`
    LastSeq       int64     `json:"last_seq"`
    FirstTS       time.Time `json:"first_ts"`
    LastTS        time.Time `json:"last_ts"`
    NumSubjects   int       `json:"num_subjects"`
    ConsumerCount int       `json:"consumer_count"`
}

// NewMonitor creates a new monitor instance
func NewMonitor(baseURL string) *Monitor {
    // Use environment variable if set
    if envURL := os.Getenv("APP_NATS_MONITOR_URL"); envURL != "" {
        baseURL = envURL
    }

    if baseURL == "" {
        baseURL = "http://localhost:8222"
    }

    // Ensure the URL is properly formatted
    if _, err := url.Parse(baseURL); err != nil {
        log.Printf("Warning: Invalid monitor URL %q, falling back to default", baseURL)
        baseURL = "http://localhost:8222"
    }

    log.Printf("Using NATS monitor URL: %s", baseURL)
    
    return &Monitor{
        baseURL: baseURL,
        client: &http.Client{
            Timeout: 5 * time.Second,
        },
    }
}

// Run starts the monitoring process with the specified interval
func (m *Monitor) Run(ctx context.Context, interval time.Duration) error {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    // Initial delay to allow NATS server to start
    time.Sleep(2 * time.Second)

    if err := m.logInfo(); err != nil {
        log.Printf("Error fetching initial information: %v", err)
    }

    for {
        select {
        case <-ticker.C:
            if err := m.logInfo(); err != nil {
                log.Printf("Error fetching information: %v", err)
            }
        case <-ctx.Done():
            return nil
        }
    }
}

func (m *Monitor) logInfo() error {
    // Fetch server info
    serverInfo, err := m.getServerInfo()
    if err != nil {
        return fmt.Errorf("error fetching server info: %w", err)
    }
    
    log.Printf("=== NATS Server Info ===")
    log.Printf("Server ID: %s", serverInfo.ServerID)
    log.Printf("Version: %s", serverInfo.Version)
    log.Printf("Uptime: %s", serverInfo.Uptime)
    log.Printf("Connections: %d", serverInfo.Connections)
    log.Printf("Memory: %d MB", serverInfo.Mem/(1024*1024))
    log.Printf("CPU: %.2f%%", serverInfo.CPU)
    
    // Fetch JetStream info
    jsInfo, err := m.getJetStreamInfo()
    if err != nil {
        return fmt.Errorf("error fetching JetStream info: %w", err)
    }

    log.Printf("\n=== JetStream Info ===")
    log.Printf("Total Streams: %d", jsInfo.Streams)
    log.Printf("Total Consumers: %d", jsInfo.Consumers)
    log.Printf("Total Messages: %d", jsInfo.Messages)
    log.Printf("Total Storage: %.2f MB", float64(jsInfo.Storage)/(1024*1024))

    // Log details for each stream
    if len(jsInfo.AccountDetails) > 0 {
        log.Printf("\n=== Stream Details ===")
        for _, account := range jsInfo.AccountDetails {
            for _, stream := range account.StreamDetail {
                log.Printf("Stream: %s", stream.Name)
                log.Printf("  Created: %s", stream.Created)
                log.Printf("  Messages: %d", stream.State.Messages)
                log.Printf("  Storage: %.2f MB", float64(stream.State.Bytes)/(1024*1024))
                log.Printf("  Consumers: %d", stream.State.ConsumerCount)
                if stream.State.Messages > 0 {
                    log.Printf("  Sequence: %d -> %d", stream.State.FirstSeq, stream.State.LastSeq)
                    log.Printf("  First Message: %s", stream.State.FirstTS.Format(time.RFC3339))
                    log.Printf("  Last Message: %s", stream.State.LastTS.Format(time.RFC3339))
                }
            }
        }
    }
    
    return nil
}

func (m *Monitor) getServerInfo() (*ServerInfo, error) {
    resp, err := m.client.Get(fmt.Sprintf("%s/varz", m.baseURL))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var info ServerInfo
    if err := json.Unmarshal(body, &info); err != nil {
        return nil, err
    }
    
    return &info, nil
}

func (m *Monitor) getJetStreamInfo() (*JetStreamResponse, error) {
    resp, err := m.client.Get(fmt.Sprintf("%s/jsz?streams=true", m.baseURL))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var response JetStreamResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("error unmarshaling JetStream info: %w", err)
    }
    
    return &response, nil
}
