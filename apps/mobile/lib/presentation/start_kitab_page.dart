import 'package:flutter/material.dart';

import '../application/kitab_dapur.dart';
import '../domain/recipe.dart';

class StartKitabPage extends StatefulWidget {
  const StartKitabPage({
    super.key,
    required this.kitabDapur,
    required this.onStarted,
  });

  final KitabDapur kitabDapur;
  final VoidCallback onStarted;

  @override
  State<StartKitabPage> createState() => _StartKitabPageState();
}

class _StartKitabPageState extends State<StartKitabPage> {
  final _name = TextEditingController();
  String? _error;
  bool _isSaving = false;

  @override
  void dispose() {
    _name.dispose();
    super.dispose();
  }

  Future<void> _start() async {
    setState(() {
      _isSaving = true;
      _error = null;
    });
    try {
      await widget.kitabDapur.startKitab(_name.text);
      widget.onStarted();
    } on ValidationFailure catch (failure) {
      setState(() => _error = failure.message);
    } catch (_) {
      setState(
        () => _error = 'Kitab gagal dibuat. Periksa koneksi lalu coba lagi.',
      );
    } finally {
      if (mounted) setState(() => _isSaving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Kitab Dapur')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Text('Beri nama kitab resep keluargamu.'),
            const SizedBox(height: 16),
            TextField(
              controller: _name,
              decoration: InputDecoration(
                labelText: 'Nama kitab',
                hintText: 'mis. Keluarga Syafrin',
                errorText: _error,
                border: const OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 16),
            FilledButton(
              onPressed: _isSaving ? null : _start,
              child: Text(_isSaving ? 'Membuat…' : 'Buat kitab'),
            ),
          ],
        ),
      ),
    );
  }
}
