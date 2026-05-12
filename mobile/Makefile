.PHONY: dev build-android build-ios test clean gen deps lint

# Get dependencies
deps:
	flutter pub get

# Run in development mode
dev:
	flutter run

# Build Android APK
build-android:
	flutter build apk --release

# Build iOS (requires macOS)
build-ios:
	flutter build ios --release

# Run tests
test:
	flutter test

# Generate code (Drift tables, etc)
gen:
	flutter pub run build_runner build --delete-conflicting-outputs

# Watch and regenerate code
gen-watch:
	flutter pub run build_runner watch --delete-conflicting-outputs

# Clean build artifacts
clean:
	flutter clean

# Analyze code
lint:
	flutter analyze

# Format code
fmt:
	dart format lib test
