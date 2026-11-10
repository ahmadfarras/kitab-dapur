import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

import 'application/kitab_dapur.dart';
import 'infrastructure/http_kitab_api.dart';
import 'infrastructure/prefs_kitab_id_store.dart';
import 'presentation/app.dart';

// Android emulator reaches the host at 10.0.2.2: --dart-define=API_BASE_URL=http://10.0.2.2:8080
const _apiBaseUrl = String.fromEnvironment(
  'API_BASE_URL',
  defaultValue: 'http://localhost:8080',
);

/// Composition root: the only place that knows which implementations the app runs on.
void main() {
  final kitabDapur = KitabDapur(
    HttpKitabApi(http.Client(), Uri.parse(_apiBaseUrl)),
    PrefsKitabIdStore(),
  );
  runApp(KitabDapurApp(kitabDapur: kitabDapur));
}
