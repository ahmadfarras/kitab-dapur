import '../domain/recipe.dart';
import 'ports.dart';

/// The app's use cases.
class KitabDapur {
  KitabDapur(this._api, this._store);

  final KitabApi _api;
  final KitabIdStore _store;

  Future<bool> hasKitab() async => await _store.read() != null;

  /// Creates the family's kitab and remembers it on this device.
  Future<void> startKitab(String name) async {
    final kitabId = await _api.createKitab(name);
    await _store.write(kitabId);
  }

  Future<Recipe> addRecipe(NewRecipe recipe) async =>
      _api.createRecipe(await _kitabId(), recipe);

  Future<List<Recipe>> listRecipes() async =>
      _api.listRecipes(await _kitabId());

  Future<String> _kitabId() async {
    final kitabId = await _store.read();
    if (kitabId == null) {
      throw StateError('no kitab yet: call startKitab first');
    }
    return kitabId;
  }
}
