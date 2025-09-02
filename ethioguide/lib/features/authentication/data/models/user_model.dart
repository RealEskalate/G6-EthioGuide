import 'package:ethioguide/features/authentication/domain/entities/user.dart';


// UserModel extends User to inherit its properties but adds data-specific functionality.
class UserModel extends User {
  const UserModel({
    required super.id,
    required super.email,
    required super.name,
    super.username,
  });

  // Factory constructor to create a UserModel from a JSON map.
  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'],
      email: json['email'],
      name: json['name'],
      username: json['username'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'name': name,
      'username': username,
    };
  }
}