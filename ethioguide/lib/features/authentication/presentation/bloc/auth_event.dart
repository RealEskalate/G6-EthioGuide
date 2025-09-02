import 'package:equatable/equatable.dart';
import 'auth_state.dart';

abstract class AuthEvent extends Equatable {
  const AuthEvent();

  @override
  List<Object> get props => [];
}

// Event when the user taps the Login/Sign Up toggle
class AuthViewSwitched extends AuthEvent {
  final AuthView newView;
  const AuthViewSwitched(this.newView);
}

// Event when the user taps the password visibility icon
class PasswordVisibilityToggled extends AuthEvent {}

// Event when the user submits the login form
class LoginSubmitted extends AuthEvent {
  final String identifier;
  final String password;
  const LoginSubmitted({required this.identifier, required this.password});
}

// Event when the user submits the sign up form
class SignUpSubmitted extends AuthEvent {
  final String name;
  final String username;
  final String email;
  final String password;
  const SignUpSubmitted({
    required this.name,
    required this.username,
    required this.email,
    required this.password,
  });
}

class ForgotPasswordSubmitted extends AuthEvent {
  final String email;
  const ForgotPasswordSubmitted({required this.email});
}

class ResetPasswordSubmitted extends AuthEvent {
  final String email;
  final String token;
  final String newPassword;
  const ResetPasswordSubmitted({
    required this.email,
    required this.token,
    required this.newPassword,
  });
  
}

class GoogleSignInSubmitted extends AuthEvent {}