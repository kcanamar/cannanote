# CannaNote Supabase Configuration

Database schema, migrations, and Row Level Security policies for HIPAA-compliant data handling.

## Setup

```bash
# Initialize Supabase (run from project root)
supabase init

# Link to remote project
supabase link --project-ref your-project-ref

# Start local development
supabase start

# Push schema changes
supabase db push
```

## Database Structure

### Tables
- `humans`: User profiles with privacy settings
- `entries`: Cannabis experience tracking with RLS
- `strains`: Cannabis strain database
- `wellness_logs`: Calendar correlation data

### Security
- Row Level Security (RLS) enabled on all tables
- `USING (auth.uid() = human_id)` policies for data privacy
- Microsoft OAuth with calendar scopes configured
- Audit logging for compliance tracking

## Planned Features

- HIPAA-compliant PHI storage (Phase 3)
- Real-time subscriptions for mobile sync
- Microsoft Graph OAuth integration
- Automated data retention policies