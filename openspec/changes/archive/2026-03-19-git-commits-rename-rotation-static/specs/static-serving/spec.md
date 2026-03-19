## ADDED Requirements

### Requirement: Static files configuration

The system SHALL provide configuration for static files directory path in config.toml.

#### Scenario: Static files path not configured

- **WHEN** `static_files_path` is not specified in configuration
- **THEN** static file serving is disabled
- **THEN** server only serves API endpoints

#### Scenario: Custom static files path configured

- **WHEN** `static_files_path` is set in config.toml
- **THEN** server serves static files from the configured directory
- **THEN** server enables SPA routing fallback

### Requirement: Static file serving

The system SHALL serve static files from the configured directory at the root URL path when enabled.

#### Scenario: Serve index.html

- **WHEN** user requests `/`
- **THEN** server returns `index.html` from static files directory

#### Scenario: Serve static assets

- **WHEN** user requests `/assets/logo.png`
- **THEN** server returns the file from `./assets/logo.png` in static directory
- **THEN** server sets appropriate Content-Type header

#### Scenario: File not found

- **WHEN** user requests a non-existent static file
- **THEN** server returns 404 Not Found response

### Requirement: SPA routing support

The system SHALL support single-page application routing by falling back to index.html for unknown routes.

#### Scenario: Unknown route fallback

- **WHEN** user requests a non-existent route like `/projects/123`
- **THEN** server returns `index.html` (for client-side routing)
- **THEN** browser handles the route via JavaScript router

#### Scenario: API routes excluded

- **WHEN** user requests `/api/*` endpoints
- **THEN** server handles as API request, not static file
- **THEN** no fallback to index.html occurs

### Requirement: Cache headers for static assets

The system SHALL set appropriate cache headers for static assets.

#### Scenario: Cache immutable assets

- **WHEN** serving files from `/assets/` or `/static/` directories
- **THEN** server sets `Cache-Control: public, max-age=31536000, immutable`

#### Scenario: No cache for index.html

- **WHEN** serving `index.html`
- **THEN** server sets `Cache-Control: no-cache, no-store, must-revalidate`
