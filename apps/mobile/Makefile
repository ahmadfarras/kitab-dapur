.PHONY: run test lint

## run: run on a connected device or simulator (Android emulator: API_BASE_URL=http://10.0.2.2:8080)
run:
	flutter run --dart-define=API_BASE_URL=$(or $(API_BASE_URL),http://localhost:8080)

test:
	flutter test

lint:
	dart format --output=none --set-exit-if-changed lib test
	flutter analyze
