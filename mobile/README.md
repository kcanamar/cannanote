# CannaNote Mobile Application

Cross-platform mobile application for CannaNote personal cannabis journaling platform, built with Flutter for iOS and Android deployment. Implements offline-first architecture with privacy-focused data synchronization and comprehensive local storage capabilities.

## Architecture Overview

The mobile application implements **privacy-first offline architecture** with selective cloud synchronization, providing users complete data ownership while enabling seamless cross-device experiences when desired.

### Architectural Principles

- **Offline-First Design** — Core functionality operates without network connectivity
- **Privacy by Design** — Local data storage with opt-in cloud synchronization
- **Cross-Platform Consistency** — Shared codebase for iOS and Android platforms
- **Modular Architecture** — Clean separation of concerns with dependency injection
- **Reactive State Management** — Efficient UI updates with minimal rebuilds
- **Security First** — Local encryption and secure authentication patterns

### Technology Stack

- **Framework**: Flutter 3.24+ with Dart 3.4+
- **State Management**: Provider pattern with ChangeNotifier
- **Local Database**: SQLite with sqflite for offline data persistence
- **Cloud Sync**: Supabase Flutter SDK for optional data synchronization
- **Authentication**: Supabase Auth with biometric authentication support
- **Local Storage**: Secure storage for sensitive data and preferences
- **Navigation**: Go Router for declarative navigation patterns
- **HTTP Client**: Dio with interceptors for API communication
- **Background Tasks**: WorkManager for data sync and notifications

## Project Structure

```
mobile/
├── android/                          # Android platform configuration
│   ├── app/
│   │   ├── src/main/
│   │   │   ├── kotlin/com/example/mobile/  # Android-specific code
│   │   │   ├── res/                   # Android resources and assets
│   │   │   └── AndroidManifest.xml    # Android app manifest
│   │   └── build.gradle.kts           # Android build configuration
│   └── gradle.properties              # Android build properties
├── ios/                              # iOS platform configuration
│   ├── Runner/
│   │   ├── Assets.xcassets/          # iOS app icons and launch images
│   │   ├── Info.plist               # iOS app configuration
│   │   └── AppDelegate.swift        # iOS application delegate
│   └── Runner.xcodeproj/            # Xcode project configuration
├── lib/                             # Dart application code
│   ├── main.dart                    # Application entry point
│   ├── core/                        # Core application infrastructure
│   │   ├── constants/               # Application constants and enums
│   │   ├── errors/                  # Error handling and exceptions
│   │   ├── network/                 # HTTP client and API configuration
│   │   ├── storage/                 # Local storage abstractions
│   │   └── utils/                   # Utility functions and helpers
│   ├── data/                        # Data layer implementation
│   │   ├── datasources/             # Local and remote data sources
│   │   │   ├── local/               # SQLite database implementation
│   │   │   └── remote/              # Supabase API integration
│   │   ├── models/                  # Data models and DTOs
│   │   └── repositories/            # Repository pattern implementations
│   ├── domain/                      # Business logic and entities
│   │   ├── entities/                # Core business entities
│   │   ├── repositories/            # Repository interfaces
│   │   └── usecases/                # Business use cases
│   ├── presentation/                # UI layer and state management
│   │   ├── providers/               # State management providers
│   │   ├── screens/                 # Application screens
│   │   │   ├── auth/                # Authentication screens
│   │   │   ├── home/                # Home dashboard
│   │   │   ├── journal/             # Cannabis journaling screens
│   │   │   ├── insights/            # Data insights and analytics
│   │   │   └── settings/            # User preferences and configuration
│   │   ├── widgets/                 # Reusable UI components
│   │   └── theme/                   # App theme and styling
│   └── services/                    # Application services
│       ├── auth_service.dart        # Authentication management
│       ├── sync_service.dart        # Data synchronization service
│       ├── notification_service.dart # Local notifications
│       └── analytics_service.dart   # Usage analytics (privacy-focused)
├── test/                            # Test suites
│   ├── unit/                        # Unit tests for business logic
│   ├── widget/                      # Widget tests for UI components
│   └── integration/                 # Integration tests for full workflows
├── assets/                          # Application assets
│   ├── images/                      # Image assets and icons
│   │   └── logos/                   # Brand logo assets
│   ├── fonts/                       # Custom font files
│   └── icons/                       # App icons and UI icons
├── pubspec.yaml                     # Flutter project configuration
├── analysis_options.yaml           # Dart static analysis configuration
└── README.md                        # This documentation file
```

## Current Implementation Status

### Implemented Features

#### Core Infrastructure
- **Flutter Project Setup** — Basic Flutter project structure with platform configurations
- **Dependency Configuration** — Essential packages configured in pubspec.yaml
- **Platform Integration** — Android and iOS platform-specific configurations
- **Development Environment** — Hot reload and debugging capabilities
- **Asset Management** — Image and font asset pipeline configuration

#### Authentication System (Planned)
- **Supabase Auth Integration** — JWT-based authentication with refresh tokens
- **Biometric Authentication** — Face ID, Touch ID, and fingerprint support
- **Secure Storage** — Encrypted local storage for authentication tokens
- **Session Management** — Automatic session refresh and logout handling

#### Offline-First Architecture (Planned)
- **Local Database Schema** — SQLite database for cannabis journal entries
- **Data Synchronization** — Bi-directional sync with backend when online
- **Conflict Resolution** — Last-write-wins with manual resolution options
- **Offline Indicators** — UI feedback for network connectivity status

### Current Development State

The mobile application is currently in initial development phase with basic Flutter project structure established. Core dependencies and platform configurations are in place, with implementation focus on offline-first cannabis journaling functionality.

#### Dependencies Configuration

Current `pubspec.yaml` includes:

```yaml
dependencies:
  flutter:
    sdk: flutter
  supabase_flutter: ^2.5.6      # Supabase integration
  provider: ^6.1.2              # State management
  sqflite: ^2.3.3              # Local SQLite database
  shared_preferences: ^2.2.3    # Simple local storage
  http: ^1.2.1                  # HTTP client
  
dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^3.0.0        # Dart linting rules
  mockito: ^5.4.4               # Testing framework
```

## Development Environment

### Prerequisites

- **Flutter 3.24+** — Latest stable Flutter SDK
- **Dart 3.4+** — Dart language SDK (included with Flutter)
- **Android Studio** — Android development environment
- **Xcode 15+** — iOS development environment (macOS only)
- **VS Code** — Recommended IDE with Flutter extensions
- **Git** — Version control system

### Platform-Specific Requirements

#### Android Development
- **Android SDK 34+** — Target API level for Android deployment
- **Gradle 8.0+** — Build system for Android applications
- **Java 17+** — Required for Android build tools
- **Android Device/Emulator** — Testing environment

#### iOS Development (macOS only)
- **Xcode 15+** — iOS development tools and simulator
- **iOS 12.0+** — Minimum deployment target
- **macOS 12+** — Required for iOS development
- **Apple Developer Account** — Required for device testing and App Store deployment

### Environment Configuration

Required environment variables and configurations:

```bash
# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your_anonymous_key

# Development Configuration
FLUTTER_ENV=development
DEBUG_MODE=true

# Platform Configurations
ANDROID_SDK_ROOT=/path/to/android/sdk
JAVA_HOME=/path/to/java/17
```

### Development Commands

```bash
# Project Setup
flutter create .                    # Initialize Flutter project (if needed)
flutter pub get                     # Install dependencies
flutter pub upgrade                 # Update dependencies

# Development Workflow
flutter run                         # Start development with hot reload
flutter run --debug                 # Debug mode with additional logging
flutter run --release               # Release mode for performance testing
flutter run -d android              # Run on Android device/emulator
flutter run -d ios                  # Run on iOS device/simulator

# Code Quality
flutter analyze                     # Static analysis and linting
flutter test                        # Run unit and widget tests
flutter test --coverage            # Generate test coverage reports
dart format .                      # Format code according to Dart style

# Build and Deployment
flutter build apk                   # Build Android APK
flutter build appbundle            # Build Android App Bundle
flutter build ios                  # Build iOS application
flutter build ipa                  # Build iOS App Store package
```

### Local Development Setup

1. **Flutter Environment Setup:**
   ```bash
   # Verify Flutter installation
   flutter doctor
   
   # Clone project and navigate to mobile directory
   git clone <repository-url>
   cd cannanote/mobile
   ```

2. **Dependency Installation:**
   ```bash
   # Install Flutter dependencies
   flutter pub get
   
   # Generate platform-specific files
   flutter packages pub run build_runner build
   ```

3. **Device/Emulator Setup:**
   ```bash
   # List available devices
   flutter devices
   
   # Start Android emulator (if available)
   flutter emulators --launch <emulator_id>
   
   # Start iOS simulator (macOS only)
   open -a Simulator
   ```

4. **Development Server:**
   ```bash
   # Start development server with hot reload
   flutter run --debug
   
   # Enable hot reload in development
   # Press 'r' for hot reload, 'R' for hot restart
   ```

## Planned Feature Development

### Phase 1: Core Cannabis Journaling

#### Offline-First Journal System
- **Local Database Implementation** — SQLite schema for cannabis experiences
- **Entry Creation Interface** — Comprehensive form for logging consumption data
- **Data Validation** — Client-side validation for data integrity
- **Offline Storage** — Full functionality without network connectivity
- **Background Sync** — Automatic synchronization when network available

#### User Interface Development
- **Material Design 3** — Modern UI components with platform adaptations
- **Custom Theme System** — Brand-aligned color schemes and typography
- **Responsive Design** — Optimized layouts for various screen sizes
- **Accessibility Support** — Screen reader and navigation accessibility
- **Dark Mode Support** — User preference-based theme switching

#### Data Privacy Implementation
- **Local Encryption** — SQLCipher integration for database encryption
- **Secure Storage** — Encrypted storage for authentication tokens
- **Privacy Controls** — Granular data sharing and sync preferences
- **Data Export** — User-owned data export in standard formats
- **Consent Management** — Clear privacy preferences and controls

### Phase 2: Advanced Features and Insights

#### Pattern Recognition and Analytics
- **Local Analytics Engine** — Client-side pattern recognition algorithms
- **Visualization Components** — Charts and graphs for consumption patterns
- **Trend Analysis** — Historical data analysis and insights generation
- **Correlation Detection** — Environmental and lifestyle factor analysis
- **Recommendation System** — Personalized suggestions based on patterns

#### Enhanced User Experience
- **Biometric Authentication** — Touch ID, Face ID, fingerprint support
- **Push Notifications** — Harm reduction reminders and tracking prompts
- **Widget Support** — Home screen widgets for quick entry logging
- **Apple Watch Integration** — Companion app for quick logging and insights
- **Voice Input** — Speech-to-text for hands-free entry creation

#### Social Features (Privacy-Optional)
- **Community Sharing** — Opt-in sharing with privacy controls
- **Anonymous Insights** — Aggregated anonymous pattern sharing
- **Educational Content** — Integrated cannabis education and research
- **Expert Content** — Curated content from cannabis professionals
- **Discussion Forums** — Moderated community discussions

### Phase 3: Platform Integration and Advanced Capabilities

#### External Integrations
- **Health App Integration** — Apple Health and Google Fit data sharing
- **Calendar Integration** — Lifestyle correlation with calendar events
- **Wearable Device Support** — Heart rate and activity data integration
- **Third-Party APIs** — Integration with other health tracking platforms
- **Dispensary Integrations** — Product information and availability

#### Enterprise and Scaling Features
- **Multi-Device Sync** — Seamless experience across multiple devices
- **Family Sharing** — Secure sharing between trusted family members
- **Caregiver Access** — Medical caregiver access with permissions
- **Healthcare Provider Integration** — Secure data sharing with medical professionals
- **Research Participation** — Opt-in anonymous data contribution for research

#### Advanced Technical Features
- **Machine Learning Models** — On-device ML for personalized insights
- **Advanced Analytics** — Sophisticated pattern recognition and prediction
- **Real-Time Sync** — Live synchronization with WebSocket connections
- **Offline ML** — Machine learning models that work without internet
- **Performance Optimization** — Advanced caching and rendering optimizations

## Testing Strategy

### Testing Architecture

#### Unit Testing
- **Business Logic Testing** — Use case and repository testing
- **Model Validation** — Data model serialization and validation testing
- **Utility Function Testing** — Helper function and calculation testing
- **Service Testing** — Authentication and sync service testing

#### Widget Testing
- **UI Component Testing** — Individual widget behavior and rendering
- **Screen Testing** — Full screen layouts and interactions
- **Navigation Testing** — Route navigation and deep linking
- **Theme Testing** — UI component appearance across themes

#### Integration Testing
- **End-to-End Workflows** — Complete user journey testing
- **Database Integration** — Local database operations and migrations
- **API Integration** — Backend service integration testing
- **Platform Integration** — Native platform feature testing

#### Testing Commands

```bash
# Test Execution
flutter test                        # Run all tests
flutter test test/unit/             # Run unit tests only
flutter test test/widget/           # Run widget tests only
flutter test integration_test/      # Run integration tests

# Test Coverage
flutter test --coverage            # Generate coverage reports
genhtml coverage/lcov.info -o coverage/html  # Generate HTML coverage report

# Performance Testing
flutter drive --target=test_driver/app.dart  # Performance profiling
flutter test --platform chrome     # Web platform testing
```

## Security and Privacy

### Data Protection Implementation

#### Local Data Security
- **Database Encryption** — SQLCipher for encrypted local database
- **Secure Storage** — Platform-native secure storage for sensitive data
- **Biometric Authentication** — Hardware-backed biometric security
- **App Lock** — Additional app-level security with PIN/biometric
- **Screen Recording Protection** — Prevent sensitive data capture

#### Network Security
- **Certificate Pinning** — TLS certificate validation for API calls
- **Request Signing** — Cryptographic signing of API requests
- **Token Management** — Secure JWT token storage and refresh
- **Network Monitoring** — Detection of network security threats
- **VPN Detection** — Optional VPN requirement for enhanced privacy

#### Privacy Controls
- **Granular Permissions** — Fine-grained data sharing controls
- **Consent Management** — Clear consent flows for data collection
- **Data Minimization** — Collect only necessary data for functionality
- **Right to Deletion** — Complete data removal capabilities
- **Privacy Dashboard** — Transparent data usage and control interface

## Performance Optimization

### Application Performance

#### Rendering Optimization
- **Widget Rebuilding** — Efficient state management to minimize rebuilds
- **Image Optimization** — Lazy loading and caching for images
- **List Performance** — Virtual scrolling for large data sets
- **Animation Performance** — Optimized animations with proper curves
- **Memory Management** — Efficient memory usage and garbage collection

#### Data Performance
- **Database Optimization** — Indexed queries and efficient schema design
- **Caching Strategy** — Multi-level caching for frequently accessed data
- **Background Processing** — Heavy operations moved to background threads
- **Pagination** — Efficient data loading with pagination
- **Sync Optimization** — Intelligent sync with delta updates

### Monitoring and Analytics

#### Performance Monitoring
- **Frame Rate Monitoring** — Real-time FPS tracking and optimization
- **Memory Usage Tracking** — Memory leak detection and optimization
- **Network Performance** — API call performance and error tracking
- **Battery Usage Optimization** — Efficient background processing
- **Crash Reporting** — Comprehensive crash detection and reporting

#### Privacy-Focused Analytics
- **Local Analytics** — On-device usage analytics without data transmission
- **Feature Usage Tracking** — Understanding feature adoption patterns
- **Performance Metrics** — App performance metrics for optimization
- **Error Tracking** — Error pattern analysis for stability improvements
- **User Flow Analysis** — Understanding user journey patterns

## Deployment and Distribution

### Build Configuration

#### Android Deployment
- **App Bundle Generation** — Optimized Android App Bundle for Play Store
- **APK Generation** — Direct APK installation for testing and distribution
- **Signing Configuration** — Release signing with secure key management
- **ProGuard Optimization** — Code obfuscation and size optimization
- **Multiple Architecture Support** — ARM64, ARMv7, and x86_64 support

#### iOS Deployment
- **App Store Build** — Optimized build for App Store submission
- **Ad Hoc Distribution** — Internal testing and beta distribution
- **TestFlight Integration** — Beta testing through Apple TestFlight
- **Code Signing** — Automatic provisioning and certificate management
- **App Store Connect Integration** — Automated metadata and screenshot uploads

### CI/CD Pipeline

#### Automated Testing
- **Test Automation** — Comprehensive test suite execution on code changes
- **Code Quality Gates** — Linting, formatting, and security scanning
- **Platform Testing** — Automated testing on multiple devices and OS versions
- **Performance Testing** — Automated performance regression detection

#### Deployment Automation
- **Build Automation** — Automated builds for multiple platforms and environments
- **Beta Distribution** — Automated beta releases to testing groups
- **Store Deployment** — Automated submission to App Store and Play Store
- **Rollback Capabilities** — Quick rollback mechanisms for problematic releases

## Contributing Guidelines

### Development Standards

#### Code Quality
- **Dart Style Guide** — Follow official Dart style conventions
- **Flutter Best Practices** — Implement Flutter-recommended patterns
- **Architecture Compliance** — Maintain clean architecture boundaries
- **Test Coverage** — Minimum 85% test coverage for new features
- **Documentation** — Comprehensive code documentation and examples

#### Code Review Process
- **Pull Request Reviews** — Mandatory code review for all changes
- **Design Review** — UI/UX review for interface changes
- **Performance Review** — Performance impact assessment for features
- **Security Review** — Security assessment for sensitive functionality
- **Platform Review** — Platform-specific implementation review

### Development Workflow

#### Feature Development
- **Feature Branches** — Isolated development with clear naming conventions
- **Commit Standards** — Conventional commit messages with clear descriptions
- **Progressive Enhancement** — Incremental feature development with regular integration
- **Documentation Updates** — Maintain documentation alongside code changes
- **Testing Integration** — Test-driven development with comprehensive coverage

This mobile application provides the foundation for CannaNote's cross-platform cannabis journaling experience, emphasizing privacy protection, offline functionality, and user data ownership while delivering a superior mobile experience across iOS and Android platforms.