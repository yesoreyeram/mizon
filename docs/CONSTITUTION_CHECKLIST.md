# Constitution Compliance Checklist

Use this checklist to verify that your pull request adheres to the [Project Constitution](../../.specify/memory/constitution.md).

## I. Microservices Architecture

- [ ] Service has a single, well-defined responsibility
- [ ] Service is independently deployable
- [ ] Service owns its own data/database (no shared databases)
- [ ] Service communication via well-defined APIs
- [ ] Service has proper health checks (liveness/readiness)

## II. Test-Driven Development

- [ ] Tests written before implementation
- [ ] All tests passing
- [ ] Code coverage ≥80%
- [ ] Unit tests cover business logic
- [ ] Integration tests for API contracts
- [ ] E2E tests for critical user journeys (if applicable)

## III. Security-First Design

- [ ] Authentication/Authorization implemented where required
- [ ] All inputs validated and sanitized
- [ ] No secrets in code (using environment variables)
- [ ] Dependencies scanned for vulnerabilities
- [ ] TLS/HTTPS for all external communication
- [ ] Least privilege principle applied

## IV. Clean Code & Documentation

- [ ] Code follows SOLID principles
- [ ] Functions are small and focused (<50 lines)
- [ ] Meaningful, descriptive naming used
- [ ] Complex logic has explanatory comments
- [ ] No code duplication (DRY principle)
- [ ] Service README.md updated
- [ ] API documentation updated (OpenAPI/Swagger)
- [ ] Architecture/design docs updated (if applicable)

## V. Enterprise-Grade Quality

### Performance

- [ ] API endpoints meet response time requirements (p95 < 200ms)
- [ ] Database queries are optimized and indexed
- [ ] Performance benchmarks passing (no regressions)

### Observability

- [ ] Structured logging implemented (JSON format)
- [ ] Correlation IDs used for tracing
- [ ] RED metrics exposed (Rate, Errors, Duration)
- [ ] Error handling with proper logging

### Observability

- [ ] Structured logging implemented (JSON format)
- [ ] Correlation IDs used for tracing
- [ ] RED metrics exposed (Rate, Errors, Duration)
- [ ] Error handling with proper logging
- [ ] OpenTelemetry instrumentation added
- [ ] All endpoints instrumented with metrics
- [ ] Database queries logged with execution time
- [ ] External API calls traced

### Auditability

- [ ] Security-relevant events logged (auth, authz, data access)
- [ ] Audit logs include required fields (timestamp, user, action, resource, outcome, IP)
- [ ] Audit logs are immutable/append-only
- [ ] Centralized audit log aggregation configured

### Infrastructure Monitoring

- [ ] Database monitoring configured (if applicable)
- [ ] Cache monitoring configured (if applicable)
- [ ] Message queue monitoring configured (if applicable)
- [ ] Grafana dashboards created/updated
- [ ] Prometheus metrics exported
- [ ] Alerts configured with runbooks

### Reliability

- [ ] Error rate acceptable (<0.1% for critical ops)
- [ ] Graceful degradation implemented
- [ ] Retry/circuit breaker patterns (where applicable)

## VI. User Experience Consistency

- [ ] Consistent API response formats
- [ ] Clear, actionable error messages
- [ ] Accessibility considerations (if UI changes)
- [ ] Responsive design (if UI changes)

## VII. Open Source & Best Practices

- [ ] Technology choices documented in TECHNOLOGY_DECISIONS.md
- [ ] Design patterns properly applied
- [ ] 12-factor app principles followed
- [ ] Infrastructure as Code (if infrastructure changes)

## Quality Gates

- [ ] All automated tests pass
- [ ] Linting and static analysis pass
- [ ] Security scans pass (no critical/high vulnerabilities)
- [ ] Code coverage ≥80%
- [ ] Peer code review completed
- [ ] CI/CD pipeline passes
- [ ] Observability instrumentation verified
- [ ] Security audit logging verified (if applicable)

## Documentation

- [ ] README.md updated (if applicable)
- [ ] API documentation updated (if applicable)
- [ ] Architecture docs updated (if applicable)
- [ ] Design docs created/updated (if applicable)

## Additional Notes

_Add any additional context, deviations, or explanations here:_

---

**Reviewer**: Please verify this checklist has been completed before approving.

**Constitution Version**: 1.1.0
