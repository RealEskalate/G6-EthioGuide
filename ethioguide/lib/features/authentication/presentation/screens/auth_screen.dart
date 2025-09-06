import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:get_it/get_it.dart';
import 'package:go_router/go_router.dart';
import 'package:ethioguide/core/config/app_color.dart';
import '../bloc/auth_bloc.dart';
import '../bloc/auth_event.dart';
import '../bloc/auth_state.dart';
import '../widgets/login_view.dart';
import '../widgets/signup_view.dart';
import '../widgets/forgot_password_view.dart';
import '../widgets/reset_password_view.dart';

class AuthScreen extends StatelessWidget {
  final String? verificationToken;
  const AuthScreen({super.key, this.verificationToken});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) {
        // First, get an instance of the BLoC from GetIt.
        final bloc = GetIt.instance<AuthBloc>();
        
        // Then, check if a verification token was passed in.
        if (verificationToken != null && verificationToken!.isNotEmpty) {
          // If a token exists, immediately add the event to the BLoC.
          bloc.add(VerificationSubmitted(activationToken: verificationToken!));
        }
        
        // Finally, return the BLoC instance.
        return bloc;
      },
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
          
           // This existing logic is still correct.
    if (state.status == AuthStatus.success) {
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
          if (state.status == AuthStatus.registrationSuccess) {
            return const Center(
              child: Padding(
                padding: EdgeInsets.all(24.0),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(Icons.mark_email_read_outlined, size: 80, color: Colors.green),
                    SizedBox(height: 24),
                    Text("Registration Successful!", style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
                    SizedBox(height: 16),
                    Text("We've sent a verification link to your email. Please click the link to activate your account and log in.", textAlign: TextAlign.center, style: TextStyle(fontSize: 16)),
                  ],
                ),
              ),
            );
          }
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