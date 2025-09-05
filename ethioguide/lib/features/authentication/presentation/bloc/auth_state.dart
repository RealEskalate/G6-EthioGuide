import 'package:equatable/equatable.dart';

// Enum to represent the two different views on the screen
enum AuthView { login, signUp, forgotPassword, resetPassword }

// Enum to represent the status of form submission
enum AuthStatus { initial, loading, success, failure, resetLinkSent }

class AuthState extends Equatable {
  final AuthView authView;
  final AuthStatus status;
  final bool isPasswordVisible;
  final String errorMessage;

  const AuthState({
    this.authView = AuthView.login, // Default to the login view
    this.status = AuthStatus.initial,
    this.isPasswordVisible = false,
    this.errorMessage = '',
  });

  AuthState copyWith({
    AuthView? authView,
    AuthStatus? status,
    bool? isPasswordVisible,
    String? errorMessage,
  }) {
    return AuthState(
      authView: authView ?? this.authView,
      status: status ?? this.status,
      isPasswordVisible: isPasswordVisible ?? this.isPasswordVisible,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }

  @override
  List<Object?> get props => [authView, status, isPasswordVisible, errorMessage];
}