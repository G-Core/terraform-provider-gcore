# Comprehensive Router Resource Testing

**JIRA:** GCLOUD2-21144
**Resource:** gcore_cloud_network_router
**Branch:** bugfix/routers-patch
**Date:** 2025-11-05

## Purpose

Comprehensive testing of router resource following the methodology in COMPREHENSIVE_MIGRATION_AND_TESTING_GUIDE.md, specifically validating the route deletion fix.

## Test Structure

```
test-router-comprehensive/
├── drift/          # Drift detection tests
├── update/         # Update operation tests
├── crud/          # CRUD verification tests
├── forcenew/      # ForceNew field tests
├── import/        # Import functionality tests
└── edge-cases/    # Edge case testing
```

## Key Areas to Test

### 1. Route Management (Critical - Bug Fix Area)
- Adding routes
- Removing routes (the fix we implemented)
- Updating routes
- Multiple routes

### 2. Drift Detection
- No routes configuration
- With routes configuration
- Routes removed from config (main bug scenario)

### 3. Update Operations
- Name changes (should PATCH)
- Route changes (should PATCH)
- External gateway changes (should PATCH)
- Interface changes (use attach/detach)

## Test Results

Will be documented after test execution.
