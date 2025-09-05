import 'package:ethioguide/features/authentication/domain/entities/user.dart';

class UserModel extends User {
  const UserModel({
    required super.id,
    required super.email,
    required super.name,
    super.username,
    super.profilePicture,
    super.role,
    super.isVerified,
    super.createdAt,
  });

  // This factory can now safely parse JSON from both login and /auth/me.
  // If a field is missing from the JSON (e.g., profile_picture during login),
  // it will safely be assigned null.
  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'],
      email: json['email'],
      name: json['name'],
      username: json['username'],
      profilePicture: json['profile_picture'],
      role: json['role'],
      isVerified: json['is_verified'],
      createdAt: json['created_at'] != null ? DateTime.tryParse(json['created_at']) : null,
    );
  }
}