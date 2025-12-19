# Migration Plan: Express/MongoDB → Go/Supabase

## Executive Summary
This document outlines the migration strategy from the current Node.js/Express/MongoDB stack to Go/Supabase, designed to achieve HIPAA compliance while maintaining business continuity.

## Current State Analysis

### Existing Architecture
```
Node.js/Express → MongoDB → EJS Templates → Tailwind CSS
```

### Current Models (from CLAUDE.md)
- **Entry Model**: username, strain, type, amount, consumption, description, date, tags, meta (votes/favorites)
- **User Model**: Authentication and session management
- **Session Storage**: MongoDB with connect-mongo

### Current Routes
- `/` - Unauthenticated (login, signup, marketing)
- `/cannanote` - Main application (requires auth)
- `/htmx` - HTMX-specific routes

## Migration Strategy

### Phase 1: Foundation Setup (Week 1-2)

#### New Go Application Structure
```
cannanote-go/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   ├── entries/
│   ├── users/
│   └── audit/
├── pkg/
│   ├── database/
│   ├── middleware/
│   └── templates/
├── templates/
├── static/
├── config/
└── docs/
```

#### Database Schema Migration
```sql
-- Users table (Supabase auth handles most of this)
CREATE TABLE profiles (
    id UUID REFERENCES auth.users(id),
    username TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- HIPAA audit fields
    last_accessed TIMESTAMPTZ,
    access_count INTEGER DEFAULT 0
);

-- Entries table with enhanced HIPAA compliance
CREATE TABLE entries (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID REFERENCES auth.users(id) NOT NULL,
    
    -- Cannabis data
    strain TEXT NOT NULL,
    type TEXT NOT NULL, -- indica, sativa, hybrid
    amount DECIMAL(10,2),
    consumption_method TEXT,
    description TEXT,
    effects JSONB, -- structured effects data
    
    -- Metadata
    date TIMESTAMPTZ DEFAULT NOW(),
    tags TEXT[],
    votes INTEGER DEFAULT 0,
    favorites INTEGER DEFAULT 0,
    
    -- HIPAA compliance fields
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    encrypted_notes TEXT, -- PGP encrypted sensitive notes
    
    -- Audit trail
    created_by UUID NOT NULL,
    last_modified_by UUID,
    version INTEGER DEFAULT 1
);

-- Row Level Security Policies
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE entries ENABLE ROW LEVEL SECURITY;

-- Users can only see their own data
CREATE POLICY "Users can view own profile" ON profiles
    FOR SELECT USING (auth.uid() = id);

CREATE POLICY "Users can update own profile" ON profiles
    FOR UPDATE USING (auth.uid() = id);

CREATE POLICY "Users can view own entries" ON entries
    FOR SELECT USING (auth.uid() = user_id);

CREATE POLICY "Users can insert own entries" ON entries
    FOR INSERT WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update own entries" ON entries
    FOR UPDATE USING (auth.uid() = user_id);
```

### Phase 2: Core Service Implementation (Week 3-4)

#### Authentication Service
```go
package auth

import (
    "github.com/supabase-community/supabase-go"
    "github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
    client *supabase.Client
}

func (s *AuthService) LoginUser(email, password string) (*AuthResponse, error) {
    resp, err := s.client.Auth.SignInWithEmailPassword(email, password)
    if err != nil {
        // Log failed login attempt for HIPAA compliance
        auditLogger.LogEvent(AuditEvent{
            Action:    "LOGIN_FAILED",
            Email:     email,
            SourceIP:  getClientIP(),
            Timestamp: time.Now(),
            Outcome:   "FAILURE",
        })
        return nil, err
    }
    
    // Log successful login
    auditLogger.LogEvent(AuditEvent{
        Action:    "LOGIN_SUCCESS",
        UserID:    resp.User.ID,
        SourceIP:  getClientIP(),
        Timestamp: time.Now(),
        Outcome:   "SUCCESS",
    })
    
    return resp, nil
}
```

#### Entry Service with HIPAA Compliance
```go
package entries

type EntryService struct {
    db *sql.DB
    audit *audit.Service
}

func (s *EntryService) CreateEntry(entry *Entry) error {
    // Validate input
    if err := s.validateEntry(entry); err != nil {
        return err
    }
    
    // Encrypt sensitive data
    if entry.SensitiveNotes != "" {
        encrypted, err := s.encryptPHI(entry.SensitiveNotes)
        if err != nil {
            return err
        }
        entry.EncryptedNotes = encrypted
        entry.SensitiveNotes = "" // Clear plaintext
    }
    
    // Insert with audit trail
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Create entry
    _, err = tx.ExecContext(ctx, 
        `INSERT INTO entries (user_id, strain, type, amount, consumption_method, 
         description, encrypted_notes, created_by) 
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
        entry.UserID, entry.Strain, entry.Type, entry.Amount,
        entry.ConsumptionMethod, entry.Description, 
        entry.EncryptedNotes, entry.UserID)
    
    if err != nil {
        return err
    }
    
    // Log audit event
    s.audit.LogEvent(AuditEvent{
        Action:     "ENTRY_CREATED",
        UserID:     entry.UserID,
        ResourceID: entry.ID,
        SourceIP:   getClientIP(),
        Timestamp:  time.Now(),
        Outcome:    "SUCCESS",
    })
    
    return tx.Commit()
}
```

### Phase 3: HTMX Integration with HATEOAS (Week 5)

#### HATEOAS Response Structure
```go
type HTMXResponse struct {
    Data interface{} `json:"data"`
    Links Links      `json:"_links"`
    Meta Meta        `json:"meta"`
}

type Links struct {
    Self   Link `json:"self"`
    Edit   Link `json:"edit,omitempty"`
    Delete Link `json:"delete,omitempty"`
    Share  Link `json:"share,omitempty"`
}

type Link struct {
    Href   string `json:"href"`
    Method string `json:"method"`
    Type   string `json:"type"`
}

func (s *EntryService) GetEntryHTML(entryID string, userID string) (string, error) {
    entry, err := s.GetEntry(entryID, userID)
    if err != nil {
        return "", err
    }
    
    // Add HATEOAS links
    entry.Links = Links{
        Self: Link{
            Href:   fmt.Sprintf("/api/entries/%s", entry.ID),
            Method: "GET",
            Type:   "application/json",
        },
    }
    
    // Only add edit/delete if user owns entry
    if entry.UserID == userID {
        entry.Links.Edit = Link{
            Href:   fmt.Sprintf("/entries/%s/edit", entry.ID),
            Method: "GET",
            Type:   "text/html",
        }
        entry.Links.Delete = Link{
            Href:   fmt.Sprintf("/api/entries/%s", entry.ID),
            Method: "DELETE",
            Type:   "application/json",
        }
    }
    
    return s.templates.ExecuteTemplate("entry-card.html", entry)
}
```

### Phase 4: Data Migration (Week 6)

#### Migration Script
```go
package main

import (
    "encoding/json"
    "log"
    
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/supabase-community/supabase-go"
)

func migrateData() error {
    // Connect to MongoDB
    mongoClient := connectMongo()
    defer mongoClient.Disconnect(ctx)
    
    // Connect to Supabase
    supabaseClient := supabase.CreateClient(supabaseURL, supabaseKey)
    
    // Migrate users
    if err := migrateUsers(mongoClient, supabaseClient); err != nil {
        return err
    }
    
    // Migrate entries
    if err := migrateEntries(mongoClient, supabaseClient); err != nil {
        return err
    }
    
    return nil
}

func migrateEntries(mongo *mongo.Client, supabase *supabase.Client) error {
    collection := mongo.Database("cannanote").Collection("entries")
    
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    for cursor.Next(ctx) {
        var mongoEntry MongoEntry
        if err := cursor.Decode(&mongoEntry); err != nil {
            log.Printf("Error decoding entry: %v", err)
            continue
        }
        
        // Transform to new schema
        newEntry := Entry{
            UserID:            mongoEntry.Username, // Map to user ID
            Strain:           mongoEntry.Strain,
            Type:             mongoEntry.Type,
            Amount:           mongoEntry.Amount,
            ConsumptionMethod: mongoEntry.Consumption,
            Description:      mongoEntry.Description,
            Date:             mongoEntry.Date,
            Tags:             mongoEntry.Tags,
            Votes:            mongoEntry.Meta.Votes,
            Favorites:        mongoEntry.Meta.Favorites,
        }
        
        // Insert into Supabase
        _, err := supabase.From("entries").Insert(newEntry).Execute()
        if err != nil {
            log.Printf("Error inserting entry: %v", err)
            continue
        }
        
        log.Printf("Migrated entry %s", mongoEntry.ID)
    }
    
    return nil
}
```

## Risk Mitigation

### Business Continuity
- **Zero-downtime deployment**: Blue-green deployment strategy
- **Rollback plan**: Maintain MongoDB for 30 days post-migration
- **Data validation**: Automated testing of migrated data
- **User communication**: Advance notice of system improvements

### Compliance Risks
- **PHI handling**: All patient data encrypted during migration
- **Audit trail**: Complete migration logs for compliance review
- **Access controls**: User permissions verified post-migration
- **Backup strategy**: Complete system backup before migration

## Success Criteria

### Technical Metrics
- [ ] 100% data migration accuracy
- [ ] <2 second page load times
- [ ] 99.9% uptime during migration
- [ ] All HIPAA compliance checks pass
- [ ] Zero security vulnerabilities

### Business Metrics
- [ ] Zero data loss
- [ ] <1% user churn during migration
- [ ] Feature parity with existing system
- [ ] Improved performance metrics
- [ ] Positive user feedback

## Timeline Summary

| Week | Phase | Key Deliverables | Success Criteria |
|------|-------|------------------|------------------|
| 1-2 | Foundation | Go app structure, Docker setup | Working development environment |
| 3-4 | Core Services | Auth, Entries, Audit services | All APIs functional |
| 5 | HTMX Integration | HATEOAS implementation | Interactive frontend working |
| 6 | Data Migration | Complete data transfer | 100% data migrated successfully |
| 7 | Testing & Launch | Production deployment | System live with monitoring |

## Post-Migration Activities

### Week 8-9: Monitoring & Optimization
- [ ] Performance monitoring setup
- [ ] User feedback collection
- [ ] System optimization
- [ ] Documentation updates

### Week 10-12: HIPAA Compliance Validation
- [ ] Security audit
- [ ] Compliance assessment
- [ ] Penetration testing
- [ ] Documentation review