# CannaNote Mobile App

Flutter application for iOS and Android with offline capabilities and Supabase synchronization.

## Setup

```bash
# Create Flutter project (run from this directory)
flutter create .

# Add dependencies
flutter pub add supabase_flutter http provider

# Run the app
flutter run
```

## Architecture

- **Screens**: Auth, journaling, calendar integration, analytics
- **State Management**: Provider pattern for reactive UI
- **Offline Support**: Local storage with Supabase sync
- **Microsoft Graph**: Calendar API integration for wellness correlation

## Planned Features

- Cannabis experience journaling with offline support
- Real-time sync with Supabase backend
- Microsoft Graph calendar integration
- Wellness correlation insights
- Social sharing and community features