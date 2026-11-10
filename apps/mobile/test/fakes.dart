import 'package:kitab_dapur/application/ports.dart';
import 'package:kitab_dapur/domain/recipe.dart';

/// Behaves like the real API for the rules the app's screens depend on.
class FakeKitabApi implements KitabApi {
  final recipesByKitab = <String, List<Recipe>>{};
  Object? listError;

  @override
  Future<String> createKitab(String name) async {
    if (name.trim().isEmpty) {
      throw const ValidationFailure('name', 'Nama kitab wajib diisi');
    }
    final id = 'k${recipesByKitab.length + 1}';
    recipesByKitab[id] = [];
    return id;
  }

  @override
  Future<Recipe> createRecipe(String kitabId, NewRecipe recipe) async {
    if (recipe.title.trim().isEmpty) {
      throw const ValidationFailure('title', 'Judul resep wajib diisi');
    }
    final recipes = recipesByKitab[kitabId]!;
    final created = Recipe(
      id: 'r${recipes.length + 1}',
      title: recipe.title,
      author: recipe.author,
      ingredients: recipe.ingredients,
      steps: recipe.steps,
    );
    recipes.insert(0, created);
    return created;
  }

  @override
  Future<List<Recipe>> listRecipes(String kitabId) async {
    if (listError case final error?) throw error;
    return List.of(recipesByKitab[kitabId]!);
  }
}

class FakeKitabIdStore implements KitabIdStore {
  String? kitabId;

  @override
  Future<String?> read() async => kitabId;

  @override
  Future<void> write(String kitabId) async => this.kitabId = kitabId;
}
