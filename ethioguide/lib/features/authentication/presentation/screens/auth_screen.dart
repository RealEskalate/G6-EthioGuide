import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:get_it/get_it.dart';
import 'package:go_router/go_router.dart';
import 'package:ethioguide/core/config/app_color.dart';
import '../bloc/auth_bloc.dart';
import '../bloc/auth_state.dart';
import '../widgets/login_view.dart';
import '../widgets/signup_view.dart';
import '../widgets/forgot_password_view.dart';

class AuthScreen extends StatelessWidget {
  const AuthScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => GetIt.instance<AuthBloc>(),
      child: const AuthScreenView(),
    );
  }
}

class AuthScreenView extends StatelessWidget {
  const AuthScreenView({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[200],
      body: BlocConsumer<AuthBloc, AuthState>(
        listener: (context, state) {
          // Listen for one-time events, like successful login/signup
          if (state.status == AuthStatus.success) {
            // Navigate to the home screen and remove all previous routes
            context.go('/home');
          }
          if (state.status == AuthStatus.failure) {
            // Show a snackbar with the error message
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.errorMessage),
                backgroundColor: AppColors.redTagColor,
              ),
            );
          }
        },
        builder: (context, state) {
          return Center(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(24.0),
              child: AnimatedSwitcher(
                duration: const Duration(milliseconds: 300),
                child: _buildViewForState(state.authView),
              ),
            ),
          );
        },
      ),
    );
  }

  // Helper function to switch between views
  Widget _buildViewForState(AuthView view) {
    switch (view) {
      case AuthView.login:
        return const LoginView();
      case AuthView.signUp:
        return const SignUpView();
      case AuthView.forgotPassword:
        return const ForgotPasswordView();
      default:
        return const LoginView();
    }
  }
}