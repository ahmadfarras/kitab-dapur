import 'package:flutter/material.dart';

import '../application/kitab_dapur.dart';
import '../domain/recipe.dart';
import 'recipe_detail_page.dart';
import 'recipe_form_page.dart';

class RecipeListPage extends StatefulWidget {
  const RecipeListPage({super.key, required this.kitabDapur});

  final KitabDapur kitabDapur;

  @override
  State<RecipeListPage> createState() => _RecipeListPageState();
}

class _RecipeListPageState extends State<RecipeListPage> {
  late Future<List<Recipe>> _recipes = widget.kitabDapur.listRecipes();

  void _reload() => setState(() {
    _recipes = widget.kitabDapur.listRecipes();
  });

  Future<void> _openForm() async {
    await Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => RecipeFormPage(kitabDapur: widget.kitabDapur),
      ),
    );
    _reload();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Kitab Dapur')),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: _openForm,
        icon: const Icon(Icons.add),
        label: const Text('Tulis resep baru'),
      ),
      body: FutureBuilder<List<Recipe>>(
        future: _recipes,
        builder: (context, snapshot) {
          if (snapshot.connectionState != ConnectionState.done) {
            return const Center(child: CircularProgressIndicator());
          }
          final recipes = snapshot.data;
          if (recipes == null) {
            return Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Text('Resep gagal dimuat.'),
                  TextButton(
                    onPressed: _reload,
                    child: const Text('Coba lagi'),
                  ),
                ],
              ),
            );
          }
          if (recipes.isEmpty) {
            return const Center(
              child: Text(
                'Belum ada resep.\nTulis resep pertama dari keluargamu.',
                textAlign: TextAlign.center,
              ),
            );
          }
          return ListView.separated(
            itemCount: recipes.length,
            separatorBuilder: (_, _) => const Divider(height: 1),
            itemBuilder: (context, index) {
              final recipe = recipes[index];
              return ListTile(
                key: ValueKey(recipe.id),
                title: Text(recipe.title),
                subtitle: Text('Resep ${recipe.author}'),
                onTap: () => Navigator.of(context).push(
                  MaterialPageRoute<void>(
                    builder: (_) => RecipeDetailPage(recipe: recipe),
                  ),
                ),
              );
            },
          );
        },
      ),
    );
  }
}
