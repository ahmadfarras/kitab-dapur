import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:kitab_dapur/domain/recipe.dart';
import 'package:kitab_dapur/infrastructure/http_kitab_api.dart';

const _recipeJson = {
  'id': 'r1',
  'kitabId': 'k1',
  'title': 'Rendang',
  'author': 'Ibu',
  'ingredients': ['daging', 'santan'],
  'steps': ['Masak lama'],
  'createdAt': '2026-11-10T00:00:00Z',
};

const _newRecipe = NewRecipe(
  title: 'Rendang',
  author: 'Ibu',
  ingredients: ['daging', 'santan'],
  steps: ['Masak lama'],
);

HttpKitabApi _api(MockClientHandler handler) =>
    HttpKitabApi(MockClient(handler), Uri.parse('http://api.test:8080'));

http.Response _json(Object? body, int status) => http.Response(
  jsonEncode(body),
  status,
  headers: {'content-type': 'application/json; charset=utf-8'},
);

void main() {
  test('createKitab posts the name and returns the id', () async {
    late http.Request sent;
    final api = _api((request) async {
      sent = request;
      return _json({'id': 'k1', 'name': 'Keluarga'}, 201);
    });

    expect(await api.createKitab('Keluarga'), 'k1');
    expect(sent.method, 'POST');
    expect(sent.url.toString(), 'http://api.test:8080/kitabs');
    expect(sent.headers['Content-Type'], startsWith('application/json'));
    expect(jsonDecode(sent.body), {'name': 'Keluarga'});
  });

  test('createRecipe posts the recipe into its kitab', () async {
    late http.Request sent;
    final api = _api((request) async {
      sent = request;
      return _json(_recipeJson, 201);
    });

    final recipe = await api.createRecipe('k1', _newRecipe);

    expect(sent.url.path, '/kitabs/k1/recipes');
    expect(jsonDecode(sent.body), {
      'title': 'Rendang',
      'author': 'Ibu',
      'ingredients': ['daging', 'santan'],
      'steps': ['Masak lama'],
    });
    expect(recipe.id, 'r1');
    expect(recipe.ingredients, ['daging', 'santan']);
  });

  test('listRecipes returns the items in server order', () async {
    final api = _api((request) async {
      expect(request.method, 'GET');
      expect(request.url.path, '/kitabs/k1/recipes');
      return _json({
        'items': [
          _recipeJson,
          {..._recipeJson, 'id': 'r2', 'title': 'Sayur Asem'},
        ],
      }, 200);
    });

    final recipes = await api.listRecipes('k1');
    expect(recipes.map((recipe) => recipe.title), ['Rendang', 'Sayur Asem']);
  });

  test('a 400 with a field becomes a ValidationFailure', () async {
    final api = _api(
      (_) async =>
          _json({'field': 'title', 'message': 'Judul resep wajib diisi'}, 400),
    );

    await expectLater(
      api.createRecipe('k1', _newRecipe),
      throwsA(
        isA<ValidationFailure>()
            .having((failure) => failure.field, 'field', 'title')
            .having(
              (failure) => failure.message,
              'message',
              'Judul resep wajib diisi',
            ),
      ),
    );
  });

  final unusable = <String, http.Response>{
    'a 400 without a field': _json({'message': 'Permintaan tidak valid'}, 400),
    'a 404': _json({'message': 'Kitab tidak ditemukan'}, 404),
    'a 500': _json({'message': 'Terjadi kesalahan di server'}, 500),
    'a non-JSON body': http.Response('<html>bad gateway</html>', 502),
    'a JSON array': _json([1], 200),
    'items that is not a list': _json({'items': 'x'}, 200),
    'a recipe that is not an object': _json({
      'items': ['x'],
    }, 200),
    'a recipe without a title': _json({
      'items': [
        {..._recipeJson, 'title': null},
      ],
    }, 200),
    'ingredients that are not strings': _json({
      'items': [
        {
          ..._recipeJson,
          'ingredients': [1],
        },
      ],
    }, 200),
  };
  for (final MapEntry(key: name, value: response) in unusable.entries) {
    test('$name becomes an ApiException', () async {
      final api = _api((_) async => response);
      await expectLater(api.listRecipes('k1'), throwsA(isA<ApiException>()));
    });
  }

  test('ApiException describes itself', () {
    expect(
      const ApiException('status 500').toString(),
      'ApiException(status 500)',
    );
  });
}
