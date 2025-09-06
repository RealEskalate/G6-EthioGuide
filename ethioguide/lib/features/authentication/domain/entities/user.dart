import 'package:equatable/equatable.dart';

// This is now the SINGLE SOURCE OF TRUTH for what a User is in our app.
class User extends Equatable {
  final String id;
  final String email;
  final String name;
  final String? username;
  // THE FIX: All new properties are nullable (using '?')
  final String? profilePicture;
  final String? role;
  final bool? isVerified;
  final DateTime? createdAt;

  const User({
    required this.id,
    required this.email,
    required this.name,
    this.username,
    this.profilePicture, // Optional
    this.role,           // Optional
    this.isVerified,     // Optional
    this.createdAt,      // Optional
  });
  
  String get initials {
    final names = name.split(' ');
    if (names.isNotEmpty && names.first.isNotEmpty) {
       if (names.length > 1 && names.last.isNotEmpty) {
         return '${names.first[0]}${names.last[0]}'.toUpperCase();
       }
       return names.first[0].toUpperCase();
    }
    return '?';
  }
  
  @override
  List<Object?> get props => [id, email, name, username, profilePicture, role, isVerified, createdAt];
}