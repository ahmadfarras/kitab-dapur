import 'dart:convert';

import 'package:http/http.dart' as http;

import '../application/ports.dart';
import '../domain/recipe.dart';

/// The server answered with something the app cannot use.
class ApiException implements Exception {
  const ApiException(this.message);

  final String message;

  @override
  String toString() => 'ApiException($message)';
}

/// [KitabApi] over HTTP.
// shortcut: hand-written against api/openapi.yaml, generate it once the contract has more than a few operations
class HttpKitabApi implements KitabApi {
  HttpKitabApi(this._client, this._baseUrl);

  final http.Client _client;
  final Uri _baseUrl;

  static const _timeout = Duration(seconds: 15);

  @override
  Future<String> createKitab(String name) async {
    final body = await _send('POST', '/kitabs', body: {'name': name});
    return _string(body, 'id');
  }

  @override
  Future<Recipe> createRecipe(String kitabId, NewRecipe recipe) async {
    final body = await _send(
      'POST',
      '/kitabs/$kitabId/recipes',
      body: {
        'title': recipe.title,
        'author': recipe.author,
        'ingredients': recipe.ingredients,
        'steps': recipe.steps,
      },
    );
    return _recipe(body);
  }

  @override
  Future<List<Recipe>> listRecipes(String kitabId) async {
    final body = await _send('GET', '/kitabs/$kitabId/recipes');
    final items = body['items'];
    if (items is! List) throw const ApiException('items is not a list');
    return [for (final item in items) _recipe(item)];
  }

  Future<Map<String, Object?>> _send(
    String method,
    String path, {
    Object? body,
  }) async {
    final request = http.Request(method, _baseUrl.replace(path: path));
    if (body != null) {
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(body);
    }
    final response = await http.Response.fromStream(await _client.send(request))
        .timeout(_timeout);

    final Object? decoded;
    try {
      decoded = jsonDecode(utf8.decode(response.bodyBytes));
    } on FormatException {
      throw ApiException('status ${response.statusCode}, body is not JSON');
    }
    if (decoded is! Map<String, Object?>) {
      throw ApiException(
        'status ${response.statusCode}, body is not an object',
      );
    }
    if (response.statusCode >= 200 && response.statusCode < 300) return decoded;

    final field = decoded['field'];
    final message = decoded['message'];
    if (response.statusCode == 400 && field is String && message is String) {
      throw ValidationFailure(field, message);
    }
    throw ApiException('status ${response.statusCode}');
  }
}

Recipe _recipe(Object? json) {
  if (json is! Map<String, Object?>) {
    throw const ApiException('recipe is not an object');
  }
  return Recipe(
    id: _string(json, 'id'),
    title: _string(json, 'title'),
    author: _string(json, 'author'),
    ingredients: _strings(json, 'ingredients'),
    steps: _strings(json, 'steps'),
  );
}

String _string(Map<String, Object?> json, String key) {
  final value = json[key];
  if (value is! String) throw ApiException('$key is not a string');
  return value;
}

List<String> _strings(Map<String, Object?> json, String key) {
  final value = json[key];
  if (value is! List || value.any((item) => item is! String)) {
    throw ApiException('$key is not a list of strings');
  }
  return value.cast<String>();
}
