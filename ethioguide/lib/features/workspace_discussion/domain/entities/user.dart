import 'package:equatable/equatable.dart';

/// Domain entity representing a user in the workspace discussion system
class User extends Equatable {
  final String id;
  final String name;
  final String? avatar;
  final String? role; // e.g., "Moderator", "Member"

  const User({
    required this.id,
    required this.name,
    this.avatar,
    this.role,
  });

  @override
  List<Object?> get props => [id, name, avatar, role];
}
