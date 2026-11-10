import 'package:flutter_test/flutter_test.dart';
import 'package:kitab_dapur/application/kitab_dapur.dart';
import 'package:kitab_dapur/domain/recipe.dart';

import '../fakes.dart';

const _opor = NewRecipe(
  title: 'Opor Ayam',
  author: 'Ibu',
  ingredients: ['ayam'],
  steps: ['Masak'],
);

void main() {
  late FakeKitabApi api;
  late FakeKitabIdStore store;
  late KitabDapur kitabDapur;

  setUp(() {
    api = FakeKitabApi();
    store = FakeKitabIdStore();
    kitabDapur = KitabDapur(api, store);
  });

  test('hasKitab is false until a kitab is started', () async {
    expect(await kitabDapur.hasKitab(), isFalse);
    await kitabDapur.startKitab('Keluarga Syafrin');
    expect(await kitabDapur.hasKitab(), isTrue);
    expect(store.kitabId, 'k1');
  });

  test('startKitab remembers nothing when the name is invalid', () async {
    await expectLater(
      kitabDapur.startKitab(' '),
      throwsA(isA<ValidationFailure>()),
    );
    expect(store.kitabId, isNull);
  });

  test('recipes are added to and listed from the remembered kitab', () async {
    await kitabDapur.startKitab('Keluarga Syafrin');
    final added = await kitabDapur.addRecipe(_opor);

    expect(added.title, 'Opor Ayam');
    expect((await kitabDapur.listRecipes()).single.id, added.id);
    expect(api.recipesByKitab['k1'], hasLength(1));
  });

  test('recipe use cases fail clearly before a kitab exists', () async {
    await expectLater(kitabDapur.addRecipe(_opor), throwsStateError);
    await expectLater(kitabDapur.listRecipes(), throwsStateError);
  });

  test('ValidationFailure describes itself', () {
    expect(
      const ValidationFailure('title', 'wajib').toString(),
      'ValidationFailure(title: wajib)',
    );
  });
}
