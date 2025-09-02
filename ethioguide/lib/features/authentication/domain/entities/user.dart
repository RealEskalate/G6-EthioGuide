import 'package:equatable/equatable.dart';

class User extends Equatable {
  final String id;
  final String email;
  final String name;
  final String? username;
  // Add any other user fields you need from the API response

  const User({
    required this.id,
    required this.email,
    required this.name,
    this.username,
  });

  @override
  List<Object?> get props => [id, email, name, username];
}