import 'package:flutter/material.dart';

import '../application/kitab_dapur.dart';
import '../domain/recipe.dart';

class RecipeFormPage extends StatefulWidget {
  const RecipeFormPage({super.key, required this.kitabDapur});

  final KitabDapur kitabDapur;

  @override
  State<RecipeFormPage> createState() => _RecipeFormPageState();
}

class _RecipeFormPageState extends State<RecipeFormPage> {
  final _title = TextEditingController();
  final _author = TextEditingController();
  final _ingredients = TextEditingController();
  final _steps = TextEditingController();
  ValidationFailure? _failure;
  String? _error;
  bool _isSaving = false;

  @override
  void dispose() {
    _title.dispose();
    _author.dispose();
    _ingredients.dispose();
    _steps.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    setState(() {
      _isSaving = true;
      _failure = null;
      _error = null;
    });
    try {
      await widget.kitabDapur.addRecipe(
        NewRecipe(
          title: _title.text,
          author: _author.text,
          ingredients: _ingredients.text.split('\n'),
          steps: _steps.text.split('\n'),
        ),
      );
      if (mounted) Navigator.of(context).pop();
    } on ValidationFailure catch (failure) {
      setState(() => _failure = failure);
    } catch (_) {
      // The typed text stays in the fields, so nothing is lost on failure.
      setState(
        () => _error = 'Resep gagal disimpan. Periksa koneksi lalu coba lagi.',
      );
    } finally {
      if (mounted) setState(() => _isSaving = false);
    }
  }

  InputDecoration _decoration(String field, String label, {String? hint}) {
    return InputDecoration(
      labelText: label,
      hintText: hint,
      errorText: _failure?.field == field ? _failure?.message : null,
      border: const OutlineInputBorder(),
      alignLabelWithHint: true,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Resep baru')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          TextField(
            controller: _title,
            decoration: _decoration('title', 'Nama masakan'),
          ),
          const SizedBox(height: 16),
          TextField(
            controller: _author,
            decoration: _decoration(
              'author',
              'Resep dari siapa?',
              hint: 'mis. Ibu, Om Budi',
            ),
          ),
          const SizedBox(height: 16),
          TextField(
            controller: _ingredients,
            decoration: _decoration('ingredients', 'Bahan (satu per baris)'),
            minLines: 4,
            maxLines: null,
          ),
          const SizedBox(height: 16),
          TextField(
            controller: _steps,
            decoration: _decoration(
              'steps',
              'Cara memasak (satu langkah per baris)',
            ),
            minLines: 4,
            maxLines: null,
          ),
          if (_error case final error?) ...[
            const SizedBox(height: 16),
            Semantics(
              liveRegion: true,
              child: Text(
                error,
                style: TextStyle(color: Theme.of(context).colorScheme.error),
              ),
            ),
          ],
          const SizedBox(height: 16),
          FilledButton(
            onPressed: _isSaving ? null : _save,
            child: Text(_isSaving ? 'Menyimpan…' : 'Simpan resep'),
          ),
        ],
      ),
    );
  }
}
