# Mizon Quick Start Guide

Welcome to Mizon! This guide will help you get started with development.

## 📖 Essential Reading

Before you begin, please read:

1. **[Project Constitution](../.specify/memory/constitution.md)** (15 min) - Non-negotiable principles
2. **[Architecture Overview](./architecture/README.md)** (10 min) - System design
3. **[Technology Decisions](./TECHNOLOGY_DECISIONS.md)** (5 min) - Tech stack

## 🛠️ Development Setup

### 1. Prerequisites

```bash
# Check you have required tools
docker --version
docker compose version
git --version
```

### 2. Clone and Run

```bash
git clone https://github.com/yesoreyeram/mizon.git
cd mizon
docker compose up -d
```

### 3. Verify Services

```bash
docker compose ps
docker compose logs -f
```

## 🔄 Development Workflow

### Creating a New Feature

1. **Understand Requirements**

   - Review functional/non-functional requirements
   - Check architecture alignment
   - Identify affected services

2. **Write Tests First (TDD)**

   ```bash
   # Write failing tests
   # Run tests - they should fail
   # Implement minimal code
   # Run tests - they should pass
   # Refactor
   ```

3. **Implement**

   ```

   ```

4. **Implement**

   - Follow SOLID principles
   - Keep functions small (<50 lines)
   - Use meaningful names
   - Add comments for complex logic

5. **Verify Quality Gates**

   ```bash
   # Run all tests
   npm test  # or equivalent for your language

   # Check coverage
   npm run coverage

   # Run linter
   npm run lint

   # Security scan
   npm audit
   ```

6. **Document**

   - Update service README
   - Update API docs
   - Update architecture docs (if needed)

7. **Submit PR**
   - Fill out [Constitution Checklist](./CONSTITUTION_CHECKLIST.md)
   - Request peer review
   - Address feedback

## 🎯 Key Principles

### Microservices

Each service must:

- Have ONE clear responsibility
- Own its data (no shared databases)
- Be independently deployable
- Communicate via APIs/messages

### Testing

- **Minimum 80% coverage** - Non-negotiable
- **TDD mandatory** - Tests before code
- **All types**: Unit, integration, E2E

### Security

- Always validate input
- Never hardcode secrets
- Encrypt data (transit & rest)
- Follow least privilege

### Performance

- API responses: p95 < 200ms
- Database queries: optimized and indexed
- Monitor: Rate, Errors, Duration (RED)

### Observability

- **Required**: Structured logging (JSON with correlation IDs)
- **Required**: OpenTelemetry instrumentation
- **Required**: Prometheus metrics for all endpoints
- **Required**: Distributed tracing for all requests
- **Stack**: Prometheus, Grafana, Loki, Tempo
- **Coverage**: Every endpoint, database query, external call

### Security & Auditability

- Always validate input
- Never hardcode secrets
- Encrypt data (transit & rest)
- Follow least privilege
- **Log all security events** (auth, access, changes)
- **Immutable audit trails** (90-day retention)

## 📁 Project Structure

```plaintext
mizon/
├── .specify/
│   └── memory/
│       └── constitution.md      # Core principles
├── docs/
│   ├── architecture/            # System design
│   ├── requirements/            # Requirements
│   ├── design/                  # Detailed designs
│   ├── TECHNOLOGY_DECISIONS.md  # Tech stack
│   └── CONSTITUTION_CHECKLIST.md
├── services/                    # Microservices
│   └── [service-name]/
│       ├── README.md
│       ├── src/
│       ├── tests/
│       └── docs/
└── docker-compose.yml
```

## 🔍 Common Tasks

### Adding a New Service

1. Create service directory structure
2. Add service README
3. Implement with TDD
4. Add to docker-compose.yml
5. Document API contracts
6. Update architecture docs

### Making API Changes

1. Update OpenAPI spec
2. Write contract tests
3. Implement changes
4. Update API documentation
5. Version appropriately

### Debugging

```bash
# View service logs
docker compose logs [service-name] -f

# Enter service container
docker compose exec [service-name] sh

# Restart service
docker compose restart [service-name]
```

## ✅ Pre-Commit Checklist

Before committing:

- [ ] All tests pass
- [ ] Coverage ≥80%
- [ ] Linter passes
- [ ] No security vulnerabilities
- [ ] Documentation updated
- [ ] Code reviewed by peer
- [ ] Observability instrumented (logs, metrics, traces)
- [ ] Security audit logging (if applicable)

## 🆘 Getting Help

- Read the [Constitution](../.specify/memory/constitution.md)
- Check [Architecture Docs](./architecture/)
- Review existing service implementations
- Ask team members

## 🚀 Next Steps

1. Read the constitution thoroughly
2. Set up your development environment
3. Explore existing services
4. Pick a task and start coding (TDD!)
5. Submit your first PR

Remember: **Quality over speed**. We follow enterprise-grade standards.

---

**Version**: 1.1.0 | **Last Updated**: 2025-11-19
