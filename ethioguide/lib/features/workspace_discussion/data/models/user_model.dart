import '../../domain/entities/user.dart';

/// Data model for users
class UserModel extends User {
  const UserModel({
    required super.id,
    required super.name,
    super.avatar,
    super.role,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'] as String,
      name: json['name'] as String,
      avatar: json['avatar'] as String?,
      role: json['role'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'avatar': avatar,
      'role': role,
    };
  }
}
