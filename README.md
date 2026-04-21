# SRE Workshop

A hands-on workshop series covering essential SRE (Site Reliability Engineering) skills — containerization with Docker, CI/CD automation with GitHub Actions, and monitoring with Prometheus.

## Workshops

| Workshop | Duration | Description |
|----------|----------|-------------|
| [Docker](Docker/) | 2 hours | Docker fundamentals, Dockerfile, and Docker Compose |
| [CI/CD](CI-CD/) | 2 hours | GitHub Actions, CI pipeline, release and deployment automation |
| [Prometheus](Prometheus/) | 1 hour | Metrics collection, alerting, and Grafana dashboards |

> **Recommended order:** Docker → CI/CD → Prometheus — later workshops assume knowledge from earlier ones.

## Prerequisites

- [ ] **Git** — `git --version`
- [ ] **Go 1.22+** — `go version`
- [ ] **Docker** — `docker --version`
- [ ] **Docker Compose** — `docker compose version`
- [ ] **GitHub account** — [github.com](https://github.com)
- [ ] **Discord account** — for receiving Prometheus alert notifications (optional)
- [ ] **Editor** — [VS Code](https://code.visualstudio.com/) recommended (with YAML extension)

## Repository Structure

```
sre-workshop/
├── README.md              # This file
├── Docker/                # Docker workshop
│   ├── README.md          # Workshop overview
│   └── docker-workshop.md # Full teaching material
├── CI-CD/                 # CI/CD workshop
│   ├── README.md          # Workshop overview
│   ├── 01-cicd-intro.md   # CI/CD concepts
│   ├── 02-github-actions-basics.md
│   ├── 03-go-ci-pipeline.md
│   ├── 04-deployment.md
│   ├── examples/          # Sample Go application
│   └── exercises/         # Hands-on exercises
└── Prometheus/            # Prometheus workshop
    ├── README.md          # Workshop overview
    ├── 01-monitoring-intro.md
    ├── 02-prometheus-fundamentals.md
    ├── 03-first-prometheus.md
    ├── 04-alerting.md
    ├── 05-full-monitoring-stack.md
    └── examples/          # Monitoring stack configurations
```

## License

This material is intended for educational use within SDC workshops.
