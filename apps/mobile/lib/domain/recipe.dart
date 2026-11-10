class Recipe {
  const Recipe({
    required this.id,
    required this.title,
    required this.author,
    required this.ingredients,
    required this.steps,
  });

  final String id;
  final String title;

  /// Who the recipe comes from, e.g. "Ibu".
  final String author;
  final List<String> ingredients;
  final List<String> steps;
}

class NewRecipe {
  const NewRecipe({
    required this.title,
    required this.author,
    required this.ingredients,
    required this.steps,
  });

  final String title;
  final String author;
  final List<String> ingredients;
  final List<String> steps;
}

/// Input that breaks a rule of the kitab or the recipe. [field] is the API's field name.
class ValidationFailure implements Exception {
  const ValidationFailure(this.field, this.message);

  final String field;
  final String message;

  @override
  String toString() => 'ValidationFailure($field: $message)';
}
