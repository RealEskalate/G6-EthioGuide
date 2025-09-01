import 'package:equatable/equatable.dart';

class Procedure extends Equatable {
  final String id;
  final String name;

  const Procedure({required this.id, required this.name});

  @override
  List<Object?> get props => [id, name];
}

class Conversation extends Equatable {
  final String id;
  final String request;
  final String response;
  final String source; // "official" or "ai-generated"
  final List<Procedure> procedures;

  const Conversation({
    required this.id,
    required this.request,
    required this.response,
    required this.source,
    required this.procedures,
  });

  @override
  List<Object?> get props => [request, response, source, procedures];
}
