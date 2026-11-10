import 'package:flutter_test/flutter_test.dart';
import 'package:kitab_dapur/infrastructure/prefs_kitab_id_store.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  test('reads null before anything is written, then the written id', () async {
    SharedPreferences.setMockInitialValues({});
    final store = PrefsKitabIdStore();

    expect(await store.read(), isNull);
    await store.write('k1');
    expect(await store.read(), 'k1');
  });

  test('reads an id saved by an earlier app launch', () async {
    SharedPreferences.setMockInitialValues({'kitab_id': 'k9'});
    expect(await PrefsKitabIdStore().read(), 'k9');
  });
}
