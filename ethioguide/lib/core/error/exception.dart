class ServerException implements Exception {
  final String message;
  final int? statusCode; // Optional: to store HTTP status code (e.g., 404, 500)

  ServerException({required this.message, this.statusCode});

  @override
  String toString() =>
      'ServerException: $message (Status code: ${statusCode ?? 'unknown'})';
}

class CacheException implements Exception {
  final String message;

  CacheException({required this.message});

  @override
  String toString() => 'CacheException: $message';
}