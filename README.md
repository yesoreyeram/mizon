# Mizon - Minimal E-commerce Platform

A microservices-based e-commerce platform with minimal functionality for demonstration and learning purposes.

## 🏛️ Project Constitution

This project follows strict architectural and quality principles defined in the [Project Constitution](./.specify/memory/constitution.md).

Key principles:

- ✅ **Microservices Architecture** - Single responsibility per service
- ✅ **Test-Driven Development** - TDD mandatory, 80%+ coverage
- ✅ **Security-First** - OWASP compliance, encryption, least privilege, audit logging
- ✅ **Enterprise Quality** - 99.9% uptime SLA, p95 < 200ms response times
- ✅ **Full Observability** - Prometheus, Grafana, Loki, Tempo monitoring stack
- ✅ **Clean Code** - SOLID principles, comprehensive documentation
- ✅ **Open Source First** - Leverage proven OSS technologies

All contributors must read and adhere to the constitution.

## 📚 Documentation

- [**Constitution**](./.specify/memory/constitution.md) - Core principles and governance
- [**Architecture**](./docs/architecture/) - System design and ADRs
- [**Requirements**](./docs/requirements/) - Functional and non-functional requirements
- [**Design**](./docs/design/) - Detailed component designs
- [**Technology Decisions**](./docs/TECHNOLOGY_DECISIONS.md) - Technology stack choices

## 🚀 Getting Started

### Prerequisites

- Docker and Docker Compose
- Git

### Running Locally

```bash
# Clone the repository
git clone https://github.com/yesoreyeram/mizon.git
cd mizon

# Start services with Docker Compose
docker compose up -d

# View logs
docker compose logs -f
```

## 🏗️ Architecture

Mizon follows a microservices architecture where each service:

- Has a single, well-defined responsibility
- Owns its own data and database
- Communicates via REST APIs or message queues
- Is independently deployable and testable
- Implements comprehensive observability

See [Architecture Documentation](./docs/architecture/) for details.

## 🧪 Development Workflow

1. **Review Constitution** - Understand core principles
2. **Create Tests First** - TDD is non-negotiable
3. **Implement** - Follow clean code principles
4. **Code Review** - Peer review required
5. **Automated Checks** - All quality gates must pass
6. **Deploy** - Automated CI/CD pipeline

See [Project Constitution](./.specify/memory/constitution.md) for complete development workflow.

## 📋 Quality Gates

All pull requests must pass:

- ✅ All automated tests (unit, integration, e2e)
- ✅ Code coverage ≥80%
- ✅ Linting and static analysis
- ✅ Security vulnerability scans
- ✅ Performance benchmarks (no regressions)
- ✅ Peer code review approval
- ✅ Observability instrumentation
- ✅ Security audit logging (if applicable)

## 🤝 Contributing

1. Read the [Project Constitution](./.specify/memory/constitution.md)
2. Review [Architecture Documentation](./docs/architecture/)
3. Follow the development workflow
4. Ensure all quality gates pass
5. Submit pull request with constitution compliance checklist

## 📄 License

[Add license information]

## 🔗 Links

- [GitHub Repository](https://github.com/yesoreyeram/mizon)
- [Project Constitution](./.specify/memory/constitution.md)
- [Documentation](./docs/)
