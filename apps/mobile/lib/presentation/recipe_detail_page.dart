import 'package:flutter/material.dart';

import '../domain/recipe.dart';

class RecipeDetailPage extends StatelessWidget {
  const RecipeDetailPage({super.key, required this.recipe});

  final Recipe recipe;

  @override
  Widget build(BuildContext context) {
    final text = Theme.of(context).textTheme;
    return Scaffold(
      appBar: AppBar(title: Text(recipe.title)),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Text('Resep ${recipe.author}', style: text.bodyMedium),
          const SizedBox(height: 16),
          Semantics(header: true, child: Text('Bahan', style: text.titleLarge)),
          for (final ingredient in recipe.ingredients) Text('• $ingredient'),
          const SizedBox(height: 16),
          Semantics(
            header: true,
            child: Text('Cara memasak', style: text.titleLarge),
          ),
          for (final (index, step) in recipe.steps.indexed)
            Text('${index + 1}. $step'),
        ],
      ),
    );
  }
}
