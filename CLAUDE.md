# CannaNote Development Roadmap

## Core Philosophy: Kina'ole Development

The primary metric for success is reducing cannabis session logging to under 30 seconds on mobile devices. All features and infrastructure decisions must align with our privacy promises—kina'ole development where our actions match our values.

## Architecture Principles

### Minimal Dependencies
- Prefer custom implementations over third-party libraries
- Each dependency must justify its inclusion
- Lifetime project mentality: code that stands the test of time

### Hexagonal Architecture (Client + Server)
Both frontend and backend follow hexagonal architecture:
- **Core Domain** - Pure business logic, no external dependencies
- **Ports** - Interfaces defining external needs
- **Adapters** - Concrete implementations for external services

This enables swapping implementations without changing business logic.

### Local-First Data
- Cannabis data never leaves the device unless user explicitly opts into paid sync
- Browser is the distribution vehicle; app works offline indefinitely
- Server only handles authentication and optional sync

---

## 🌱 GERMINATION CYCLE: Offline-First Foundation

**Status**: Current focus. Building the local-first PWA infrastructure.

### Client-Side Architecture

#### Storage Layer (Ports & Adapters)
```
cmd/web/assets/js/
  core/
    domain/
      session.js           # Domain entity definitions
      strain.js
    ports/
      storage-port.js      # Storage interface contract
      sync-port.js         # Sync interface contract (future)
  adapters/
    storage/
      indexed-db-adapter.js  # Implements storage-port
    sync/
      rest-sync-adapter.js   # Implements sync-port (future)
  application/
    session-service.js     # Business logic, depends on ports
```

#### Storage Port Contract
```javascript
// What any storage adapter must implement
{
  init: async () => {},
  addSession: async (session) => {},
  getSession: async (id) => {},
  updateSession: async (id, updates) => {},
  deleteSession: async (id) => {},
  getAllSessions: async () => {},
  getSessionsByStatus: async (status) => {},
  getSessionsInRange: async (startTime, endTime) => {},
  getPendingSessions: async () => {},
  markSessionsSynced: async (ids) => {}
}
```

#### IndexedDB Adapter
- Custom implementation (~150 lines, zero dependencies)
- Promise-based API
- Schema versioning and migrations
- Index-based querying

#### Data Schema
```javascript
sessions: {
  id: uuid,
  timestamp: number,
  strain: string,
  method: string,  // vape, smoke, edible, tincture, topical
  amount: string,
  effects: string[],
  intensity: number,  // 1-10
  duration: number,   // minutes
  notes: string,
  rating: number,     // 1-5
  syncStatus: string  // local, pending, synced
}
```

### Service Worker
- Cache app shell for offline loading
- Static asset caching (CSS, JS, fonts)
- Network-first for API calls
- Offline fallback page

### PWA Manifest Enhancement
- `start_url: "/app"`
- Proper icons for all sizes
- `display: standalone`
- Offline capability declaration

---

## 🌿 SEEDLING CYCLE: Core Cannabis Experience

### 30-Second Session Logging
- **Quick entry form** - Touch-optimized, minimal taps
- **Smart defaults** - Remember last strain, method
- **Offline-first** - Works without network
- **Immediate feedback** - Session saved confirmation

### Session List & History
- **Chronological view** - Recent sessions first
- **Filter by strain** - Find patterns
- **Filter by date range** - Weekly/monthly views
- **Search notes** - Full-text search in IndexedDB

### Data Export
- **JSON export** - Complete data portability
- **CSV export** - Spreadsheet compatible
- **User-initiated** - Manual export button
- **No server involvement** - Client-side generation

---

## 🍃 VEGETATIVE CYCLE: Sync Infrastructure

**Prerequisite**: Germination and Seedling cycles complete.

### Paid Sync Feature
- **Opt-in only** - Explicit user consent required
- **NeonDB storage** - PostgreSQL with Row Level Security
- **Simple REST sync** - POST pending sessions, receive merged state
- **Conflict resolution** - Last-write-wins with timestamp comparison

### Sync Protocol
```
Client                           Server
  |                                |
  |-- POST /api/sessions/sync ---->|
  |   { sessions: [...],           |
  |     lastSyncTimestamp }        |
  |                                |
  |<-- { merged: [...],         ---|
  |      serverTimestamp }         |
  |                                |
  |-- Mark synced locally          |
```

### Server-Side (Go)
- **sqlc** for type-safe database access
- **Sync endpoint** - `/api/sessions/sync`
- **Conflict detection** - Compare timestamps
- **User isolation** - RLS policies on all tables

---

## 🌸 PRE-FLOWER CYCLE: Pattern Recognition

### Client-Side Analytics
- **Strain effectiveness** - Which strains work best for you
- **Time-of-day patterns** - When do you consume
- **Method preferences** - Consumption method trends
- **Effect correlations** - What predicts good experiences

### Visualization
- **Simple charts** - Consumption frequency over time
- **Effect heatmaps** - Common effects by strain
- **Personal insights** - "Evening indicas improve your sleep"

---

## 💐 FLOWERING CYCLE: Mobile Application

**Prerequisite**: PWA beta validates core experience and API design.

### Flutter + Drift
- **iOS and Android** - Single codebase
- **Drift database** - SQLite-based local storage
- **Same sync protocol** - Reuse server endpoints
- **Offline-first** - Same architecture as PWA

### Cross-Platform Sync
- **PWA ↔ Mobile** - Seamless data sync
- **Conflict resolution** - Same algorithm
- **Subscription sharing** - One payment, all platforms

---

## 🔥 HARVEST CYCLE: Advanced Features

### Harm Reduction Tools
- **Dosage guidance** - Evidence-based recommendations
- **Tolerance tracking** - Break suggestions
- **Session spacing** - Mindful consumption reminders
- **Safety information** - Contextual education

### Privacy Dashboard
- **Data location transparency** - Where is your data
- **Sync status** - What's local vs synced
- **Export history** - When you exported
- **Delete options** - Granular data removal

---

## Implementation Principles

### Development Philosophy
- **Evidence-based features** - Backed by research or user data
- **Privacy by design** - Privacy at architecture level
- **Harm reduction first** - Safety over engagement
- **Cannabis culture alignment** - Language reflects community values

### Code Standards
- **Pure JavaScript** - No TypeScript, no transpilation
- **Zero framework frontend** - Vanilla JS only
- **Hexagonal architecture** - Ports and adapters pattern
- **Minimal dependencies** - Justify every import

### Quality Standards
- **Offline functionality** - 100% features work offline
- **30-second logging** - Maximum time for session entry
- **Sub-2-second startup** - App loads instantly
- **Zero tracking** - No analytics, no telemetry

---

## Risk Mitigation

### Technical
- **Database abstraction** - Can swap NeonDB if needed
- **Storage abstraction** - Can swap IndexedDB adapter
- **No vendor lock-in** - All custom implementations

### Business
- **Open source core** - Premium features fund development
- **Beta grandfathering** - Early supporters get lifetime access
- **Data portability** - Users can always export and leave

---

## Metrics

### Primary (User Experience)
- **Logging time** - Target: <30 seconds
- **Offline availability** - Target: 100%
- **App startup** - Target: <2 seconds

### Secondary (Business)
- **PWA installs** - Home screen additions
- **Beta retention** - Monthly active users
- **Sync conversion** - Free → paid upgrades
