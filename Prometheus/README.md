# Prometheus Monitoring Workshop

## Overview

A **1-hour** hands-on workshop on Prometheus monitoring, starting from monitoring fundamentals and progressively covering Prometheus architecture, deployment, alerting with Alertmanager, and building a complete monitoring stack with Grafana. By the end, you'll be able to build a monitoring system from scratch.

## Prerequisites

- [ ] **Docker** installed — `docker --version`
- [ ] **Docker Compose** installed — `docker compose version`
- [ ] **Git** installed — `git --version`
- [ ] **Discord** account — for receiving alert notifications

## Curriculum

| Chapter | Topic | Duration |
|---------|-------|----------|
| 01 | Monitoring Concepts | 15 min |
| 02 | Prometheus Fundamentals | 20 min |
| 03 | Hands-on: Deploy Prometheus | 25 min |
| 04 | Alertmanager & Alert Rules | 25 min |
| 05 | Full Monitoring Stack | 30 min |
| -- | Wrap-up & Q&A | 5 min |

> Total: 1 hour

## Project Structure

```
Prometheus/
├── README.md                          # This file
├── README.zh-TW.md                    # Chinese version
├── 01-monitoring-intro.md             # Monitoring concepts
├── 02-prometheus-fundamentals.md      # Prometheus architecture & PromQL
├── 03-first-prometheus.md             # Hands-on: deploy Prometheus
├── 04-alerting.md                     # Alertmanager & alert rules
├── 05-full-monitoring-stack.md        # Complete monitoring stack
├── prometheus-guide.md                # Full reference guide (from PDF)
├── prometheus_schedule.md             # Original schedule (Chinese)
├── Prometheus_Guide_*.pdf             # Original PDF guide
└── examples/                         # Example configurations
    └── monitoring-stack/              # Complete monitoring stack
        ├── docker-compose.yml         # All services definition
        └── config/
            ├── prometheus.yml         # Prometheus configuration
            ├── alertmanager.yml       # Alertmanager configuration
            ├── blackbox.yml           # Blackbox Exporter configuration
            ├── rules/
            │   └── alerts.yml         # Alert rules
            └── grafana/
                └── provisioning/
                    └── datasources/
                        └── prometheus.yml  # Grafana auto-provisioning
```

## Final Outcome

After completing all chapters, you will have a fully functional monitoring stack:

```
Your Monitoring Stack
├── Prometheus       — Collects metrics, evaluates alert rules
├── Node Exporter    — Reports host hardware metrics (CPU, memory, disk)
├── Blackbox Exporter — Probes HTTP endpoints from outside
├── Alertmanager     — Receives alerts, routes notifications to Discord
└── Grafana          — Visualizes metrics in real-time dashboards
```
