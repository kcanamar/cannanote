import 'package:flutter_test/flutter_test.dart';

import 'package:cannanote_mobile/main.dart';
import 'package:cannanote_mobile/di/injection.dart';

void main() {
  setUpAll(() async {
    // Initialize dependencies before tests
    await configureDependencies();
  });

  tearDownAll(() async {
    await resetDependencies();
  });

  testWidgets('App renders home page', (WidgetTester tester) async {
    // Build our app and trigger a frame.
    await tester.pumpWidget(const CannaNote());

    // Verify that the app title is displayed
    expect(find.text('CannaNote'), findsOneWidget);

    // Verify that the Log Session section exists
    expect(find.text('Log Session'), findsOneWidget);

    // Verify that the Sessions section exists
    expect(find.text('Sessions'), findsOneWidget);
  });
}
