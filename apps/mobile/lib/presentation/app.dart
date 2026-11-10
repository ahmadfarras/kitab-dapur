import 'package:flutter/material.dart';

import '../application/kitab_dapur.dart';
import 'recipe_list_page.dart';
import 'start_kitab_page.dart';

class KitabDapurApp extends StatefulWidget {
  const KitabDapurApp({super.key, required this.kitabDapur});

  final KitabDapur kitabDapur;

  @override
  State<KitabDapurApp> createState() => _KitabDapurAppState();
}

class _KitabDapurAppState extends State<KitabDapurApp> {
  late Future<bool> _hasKitab = widget.kitabDapur.hasKitab();

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Kitab Dapur',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFFB5462B)),
      ),
      home: FutureBuilder<bool>(
        future: _hasKitab,
        builder: (context, snapshot) {
          if (snapshot.connectionState != ConnectionState.done) {
            return const Scaffold(
              body: Center(child: CircularProgressIndicator()),
            );
          }
          if (snapshot.data ?? false) {
            return RecipeListPage(kitabDapur: widget.kitabDapur);
          }
          return StartKitabPage(
            kitabDapur: widget.kitabDapur,
            onStarted: () => setState(() {
              _hasKitab = widget.kitabDapur.hasKitab();
            }),
          );
        },
      ),
    );
  }
}
