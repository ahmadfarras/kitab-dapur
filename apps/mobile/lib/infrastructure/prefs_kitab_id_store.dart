import 'package:shared_preferences/shared_preferences.dart';

import '../application/ports.dart';

/// [KitabIdStore] in the device's preferences.
class PrefsKitabIdStore implements KitabIdStore {
  static const _key = 'kitab_id';

  @override
  Future<String?> read() async =>
      (await SharedPreferences.getInstance()).getString(_key);

  @override
  Future<void> write(String kitabId) async {
    final saved = await (await SharedPreferences.getInstance()).setString(
      _key,
      kitabId,
    );
    // Losing the id would orphan the family's kitab, so a failed write must not pass silently.
    if (!saved) throw StateError('could not save the kitab id');
  }
}
