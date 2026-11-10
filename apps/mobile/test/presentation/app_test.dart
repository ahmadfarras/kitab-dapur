import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:kitab_dapur/application/kitab_dapur.dart';
import 'package:kitab_dapur/domain/recipe.dart';
import 'package:kitab_dapur/presentation/app.dart';

import '../fakes.dart';

void main() {
  late FakeKitabApi api;
  late FakeKitabIdStore store;

  setUp(() {
    api = FakeKitabApi();
    store = FakeKitabIdStore();
  });

  Future<void> pumpApp(WidgetTester tester) async {
    await tester.pumpWidget(KitabDapurApp(kitabDapur: KitabDapur(api, store)));
    await tester.pumpAndSettle();
  }

  Future<void> startWithKitab(WidgetTester tester) async {
    store.kitabId = await api.createKitab('Keluarga');
    await pumpApp(tester);
  }

  Finder field(String label) => find.widgetWithText(TextField, label);

  testWidgets('first launch asks for a kitab name, then shows the empty list', (
    tester,
  ) async {
    await pumpApp(tester);
    await tester.tap(find.text('Buat kitab'));
    await tester.pumpAndSettle();
    expect(find.text('Nama kitab wajib diisi'), findsOneWidget);

    await tester.enterText(field('Nama kitab'), 'Keluarga Syafrin');
    await tester.tap(find.text('Buat kitab'));
    await tester.pumpAndSettle();

    expect(find.textContaining('Belum ada resep.'), findsOneWidget);
    expect(store.kitabId, 'k1');
  });

  testWidgets('a later launch goes straight to the recipes', (tester) async {
    await startWithKitab(tester);
    expect(find.text('Tulis resep baru'), findsOneWidget);
    expect(find.text('Buat kitab'), findsNothing);
  });

  testWidgets('writes a recipe, sees it in the list and opens it', (
    tester,
  ) async {
    await startWithKitab(tester);
    await tester.tap(find.text('Tulis resep baru'));
    await tester.pumpAndSettle();

    await tester.enterText(field('Nama masakan'), 'Soto Ayam');
    await tester.enterText(field('Resep dari siapa?'), 'Bude');
    await tester.enterText(field('Bahan (satu per baris)'), 'ayam\nkunyit');
    await tester.enterText(
      field('Cara memasak (satu langkah per baris)'),
      'Rebus ayam\nSajikan',
    );
    await tester.tap(find.text('Simpan resep'));
    await tester.pumpAndSettle();

    expect(find.text('Soto Ayam'), findsOneWidget);
    expect(find.text('Resep Bude'), findsOneWidget);
    expect(api.recipesByKitab['k1']!.single.ingredients, ['ayam', 'kunyit']);

    await tester.tap(find.text('Soto Ayam'));
    await tester.pumpAndSettle();
    expect(find.text('• kunyit'), findsOneWidget);
    expect(find.text('2. Sajikan'), findsOneWidget);
  });

  testWidgets('a rejected recipe shows the error and keeps what was typed', (
    tester,
  ) async {
    await startWithKitab(tester);
    await tester.tap(find.text('Tulis resep baru'));
    await tester.pumpAndSettle();

    await tester.enterText(field('Resep dari siapa?'), 'Bude');
    await tester.tap(find.text('Simpan resep'));
    await tester.pumpAndSettle();

    expect(find.text('Judul resep wajib diisi'), findsOneWidget);
    expect(find.text('Bude'), findsOneWidget);
    expect(api.recipesByKitab['k1'], isEmpty);
  });

  testWidgets('a failed load offers a retry that works', (tester) async {
    api.listError = Exception('offline');
    await startWithKitab(tester);
    expect(find.text('Resep gagal dimuat.'), findsOneWidget);

    api.listError = null;
    api.recipesByKitab['k1']!.add(
      const Recipe(
        id: 'r1',
        title: 'Rawon',
        author: 'Mbah',
        ingredients: ['kluwek'],
        steps: ['Rebus'],
      ),
    );
    await tester.tap(find.text('Coba lagi'));
    await tester.pumpAndSettle();

    expect(find.text('Rawon'), findsOneWidget);
  });
}
