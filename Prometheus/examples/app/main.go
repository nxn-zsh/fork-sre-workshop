package main

import (
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	//CounterVec is a collection of counters for the same metric, use labels to differentiate
	httpRequestsTotal *prometheus.CounterVec

	//httpRequestsFailed prometheus.Counter

	//Gauge measures a val that can go up or down
	activeUsers prometheus.Gauge
	// A Histogram observes distributions of values in configurable buckets
	httpRequestDuration *prometheus.HistogramVec
	// A Summary also observes distributions but calculates quantiles on the client.
	httpRequestSummary *prometheus.SummaryVec

	//cpuTemp       prometheus.Gauge
	//hdFailures    *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "home_requests_total",
				Help: "Total number of requests and status to the /home endpoint.",
			},
			[]string{"method", "status"},
		),
		activeUsers: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of active users.",
		}),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds.",
				Buckets: []float64{0.1, 0.5, 1, 5},
			},
			[]string{"endpoint"},
		),
		httpRequestSummary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "http_request_summary_seconds",
				Help:       "Summary of HTTP request duration in seconds.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"endpoint"},
		),
	}
	reg.MustRegister(m.httpRequestsTotal)
	reg.MustRegister(m.activeUsers)
	reg.MustRegister(m.httpRequestDuration)
	reg.MustRegister(m.httpRequestSummary)
	return m
}

// Alert
type Alert struct {
	Name        string
	Description string
	Severity    string
	FiredAt     time.Time
	Resolved    bool
}

// AlertManager: Handles rules and notifications
type AlertManager struct {
	mu     sync.RWMutex
	alerts map[string]*Alert
	rules  []AlertRule
}

type AlertRule struct {
	Name        string
	Description string
	Severity    string
	Condition   func() bool
	Cooldown    time.Duration
	lastFired   time.Time
}

func NewAlertManager() *AlertManager {
	return &AlertManager{
		alerts: make(map[string]*Alert),
		rules:  []AlertRule{},
	}
}

// AddRule adds an alert rule to the manager
func (am *AlertManager) AddRule(rule AlertRule) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.rules = append(am.rules, rule)
}

// CheckRules evaluates all alert rules
func (am *AlertManager) CheckRules() {
	am.mu.Lock()
	defer am.mu.Unlock()

	for i := range am.rules {
		rule := &am.rules[i]

		// Check cooldown period
		if time.Since(rule.lastFired) < rule.Cooldown {
			continue
		}

		// Evaluate condition
		if rule.Condition() {
			// Fire alert
			alert := &Alert{
				Name:        rule.Name,
				Description: rule.Description,
				Severity:    rule.Severity,
				FiredAt:     time.Now(),
				Resolved:    false,
			}
			am.alerts[rule.Name] = alert
			rule.lastFired = time.Now()
			log.Printf("ALERT FIRED: [%s] %s - %s", rule.Severity, rule.Name, rule.Description)
		} else {
			// Resolve alert if it was previously fired
			if existing, ok := am.alerts[rule.Name]; ok && !existing.Resolved {
				existing.Resolved = true
				log.Printf("ALERT RESOLVED: %s", rule.Name)
			}
		}
	}
}

// GetActiveAlerts returns all unresolved alerts
func (am *AlertManager) GetActiveAlerts() []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var active []*Alert
	for _, alert := range am.alerts {
		if !alert.Resolved {
			active = append(active, alert)
		}
	}
	return active
}

// StartMonitoring starts the alert checking loop
func (am *AlertManager) StartMonitoring(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			am.CheckRules()
		}
	}()
}

func main() {
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	// Initialize Alert Manager
	alertMgr := NewAlertManager()

	// Counter to track requests for alert rules
	var requestCount int64
	var mu sync.Mutex

	// Define alert rules
	alertMgr.AddRule(AlertRule{
		Name:        "HighRequestRate",
		Description: "Request rate exceeded 10 requests in last check",
		Severity:    "critical",
		Cooldown:    1 * time.Minute,
		Condition: func() bool {
			mu.Lock()
			defer mu.Unlock()
			rate := requestCount
			requestCount = 0 // Reset counter
			return rate > 10
		},
	})

	alertMgr.AddRule(AlertRule{
		Name:        "SlowRequests",
		Description: "Average request duration exceeds 2 seconds",
		Severity:    "warning",
		Cooldown:    45 * time.Second,
		Condition: func() bool {
			// Simplified check - in production, calculate from histogram
			return false
		},
	})

	// Start alert monitoring every 10 seconds
	alertMgr.StartMonitoring(10 * time.Second)
	log.Println("Alert Manager started - checking rules every 10 seconds")

	// Home handler
	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.activeUsers.Inc()
		defer m.activeUsers.Dec()

		// Increment request counter for alerts
		mu.Lock()
		requestCount++
		mu.Unlock()

		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

		duration := time.Since(start).Seconds()
		m.httpRequestDuration.WithLabelValues("/home").Observe(duration)
		m.httpRequestSummary.WithLabelValues("/home").Observe(duration)
		m.httpRequestsTotal.With(prometheus.Labels{"method": r.Method, "status": "200"}).Inc()
		w.Write([]byte("Welcome!"))
	}

	// Alerts endpoint - view active alerts
	alertsHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		activeAlerts := alertMgr.GetActiveAlerts()

		html := "<h1>Active Alerts</h1>"
		if len(activeAlerts) == 0 {
			html += "<p>No active alerts</p>"
		} else {
			html += "<ul>"
			for _, alert := range activeAlerts {
				html += "<li><strong>" + alert.Severity + ":</strong> " +
					alert.Name + " - " + alert.Description +
					" (Fired at: " + alert.FiredAt.Format(time.RFC3339) + ")</li>"
			}
			html += "</ul>"
		}
		w.Write([]byte(html))
	}

	// Set up endpoints
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/alerts", alertsHandler)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	log.Println("Server started on http://localhost:8088")
	log.Println("  - Main page: http://localhost:8088/home")
	log.Println("  - Metrics:   http://localhost:8088/metrics")
	log.Println("  - Alerts:    http://localhost:8088/alerts")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
