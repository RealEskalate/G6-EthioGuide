import 'package:ethioguide/features/authentication/domain/usecases/forgot_password.dart';
import 'package:ethioguide/features/authentication/domain/usecases/reset_password.dart';
import 'package:ethioguide/features/authentication/domain/usecases/sign_in_with_google.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:ethioguide/features/authentication/domain/usecases/login_user.dart';
import 'package:ethioguide/features/authentication/domain/usecases/register_user.dart';
import 'auth_event.dart';
import 'auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final LoginUser loginUser;
  final RegisterUser registerUser;
   final ForgotPassword forgotPassword; 
  final ResetPassword resetPassword;
  final SignInWithGoogle signInWithGoogle;

  AuthBloc({required this.loginUser, required this.registerUser,  required this.forgotPassword, // ADDED
    required this.resetPassword,
    required this.signInWithGoogle,}) : super(const AuthState()) {
    on<AuthViewSwitched>(_onAuthViewSwitched);
    on<PasswordVisibilityToggled>(_onPasswordVisibilityToggled);
    on<LoginSubmitted>(_onLoginSubmitted);
    on<SignUpSubmitted>(_onSignUpSubmitted);
    on<GoogleSignInSubmitted>(_onGoogleSignInSubmitted);
    on<ForgotPasswordSubmitted>(_onForgotPasswordSubmitted);
    on<ResetPasswordSubmitted>(_onResetPasswordSubmitted);
  }

  void _onAuthViewSwitched(AuthViewSwitched event, Emitter<AuthState> emit) {
    emit(state.copyWith(authView: event.newView, status: AuthStatus.initial));
  }

  void _onPasswordVisibilityToggled(PasswordVisibilityToggled event, Emitter<AuthState> emit) {
    emit(state.copyWith(isPasswordVisible: !state.isPasswordVisible));
  }

  Future<void> _onLoginSubmitted(LoginSubmitted event, Emitter<AuthState> emit) async {
    emit(state.copyWith(status: AuthStatus.loading));
    final result = await loginUser(LoginParams(identifier: event.identifier, password: event.password));
    result.fold(
      (failure) => emit(state.copyWith(status: AuthStatus.failure, errorMessage: failure.message)),
      (user) => emit(state.copyWith(status: AuthStatus.success)),
    );
  }

  Future<void> _onSignUpSubmitted(SignUpSubmitted event, Emitter<AuthState> emit) async {
    emit(state.copyWith(status: AuthStatus.loading));
    final result = await registerUser(RegisterParams(
      username: event.username,
      email: event.email,
      password: event.password,
      name: event.name,
    ));
    result.fold(
      (failure) => emit(state.copyWith(status: AuthStatus.failure, errorMessage: failure.message)),
      (_) => emit(state.copyWith(status: AuthStatus.success)), // Or a specific 'needs verification' state
    );
  }
  Future<void> _onGoogleSignInSubmitted(GoogleSignInSubmitted event, Emitter<AuthState> emit) async {
    emit(state.copyWith(status: AuthStatus.loading));
    final result = await signInWithGoogle();
    result.fold(
      (failure) => emit(state.copyWith(status: AuthStatus.failure, errorMessage: failure.message)),
      (user) => emit(state.copyWith(status: AuthStatus.success)),
    );
}
  // Handler for the forgot password event
  Future<void> _onForgotPasswordSubmitted(ForgotPasswordSubmitted event, Emitter<AuthState> emit) async {
    emit(state.copyWith(status: AuthStatus.loading));
    final result = await forgotPassword(event.email);
    result.fold(
      (failure) => emit(state.copyWith(status: AuthStatus.failure, errorMessage: failure.message)),
      (_) => emit(state.copyWith(status: AuthStatus.resetLinkSent)),
    );
  }

  // Handler for the reset password event
  Future<void> _onResetPasswordSubmitted(ResetPasswordSubmitted event, Emitter<AuthState> emit) async {
    emit(state.copyWith(status: AuthStatus.loading));
    final result = await resetPassword(ResetPasswordParams(
      resetToken: event.resetToken,
      newPassword: event.newPassword,
    ));
    result.fold(
      (failure) => emit(state.copyWith(status: AuthStatus.failure, errorMessage: failure.message)),
      // On success, switch the view back to the login page.
      (_) => emit(state.copyWith(status: AuthStatus.initial, authView: AuthView.login)),
    );
  }
}