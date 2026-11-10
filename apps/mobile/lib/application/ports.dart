import '../domain/recipe.dart';

/// The Kitab Dapur backend, as the app needs it. Follows api/openapi.yaml.
abstract interface class KitabApi {
  /// Returns the new kitab's id. Throws [ValidationFailure] on invalid input.
  Future<String> createKitab(String name);

  /// Throws [ValidationFailure] on invalid input.
  Future<Recipe> createRecipe(String kitabId, NewRecipe recipe);

  /// Newest first.
  Future<List<Recipe>> listRecipes(String kitabId);
}

/// Remembers on this device which kitab the family uses.
abstract interface class KitabIdStore {
  Future<String?> read();
  Future<void> write(String kitabId);
}
