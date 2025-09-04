import 'package:ethioguide/core/components/button.dart';
import 'package:ethioguide/core/components/textfield.dart';
import 'package:ethioguide/features/authentication/presentation/bloc/auth_bloc.dart';
import 'package:ethioguide/features/authentication/presentation/bloc/auth_event.dart';
import 'package:ethioguide/features/authentication/presentation/bloc/auth_state.dart';
import 'package:ethioguide/features/authentication/presentation/widgets/google_sign_in_button.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'auth_toggle_buttons.dart';
// import 'package:ethioguide/core/utils/validators.dart';

class LoginView extends StatelessWidget {
  const LoginView({super.key});

  @override
  Widget build(BuildContext context) {
    // Create controllers here. In a real app with more complex state,
    // you might manage these within the BLoC, but for a simple form, this is fine.
    final emailController = TextEditingController();
    final passwordController = TextEditingController();

    return Column(
      key: const ValueKey('loginView'),
      children: [
        const Text("Welcome to EthioGuide", style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        const Text("Navigate government processes with ease", style: TextStyle(fontSize: 16, color: Colors.grey)),
        const SizedBox(height: 32),
        const AuthToggleButtons(activeView: AuthView.login),
        const SizedBox(height: 24),
        Container(
          padding: const EdgeInsets.all(24.0),
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(16),
            boxShadow: [BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 10)],
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text("Email", style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              customTextField(hintText: "Enter your email", controller: emailController, prefixIcon: Icons.email_outlined),
              const SizedBox(height: 16),
              const Text("Password", style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              BlocBuilder<AuthBloc, AuthState>(
                builder: (context, state) {
                  return customTextField(
                    hintText: "Enter your password",
                    controller: passwordController,
                    prefixIcon: Icons.lock_outline,
                    obscureText: !state.isPasswordVisible, 
                    suffixIcon: IconButton(
                      icon: Icon(
                        state.isPasswordVisible ? Icons.visibility_off : Icons.visibility,
                        color: Colors.grey,
                      ),
                      onPressed: () {
                        // Dispatch the event to toggle the state in the BLoC
                        context.read<AuthBloc>().add(PasswordVisibilityToggled());
                      },
                    ),
                  );
                },
              ),
              const SizedBox(height: 8),
              Align(
                alignment: Alignment.centerRight,
                child: TextButton(
                  onPressed: () => context.read<AuthBloc>().add(const AuthViewSwitched(AuthView.forgotPassword)),
                  child: const Text("Forgot your password?"),
                ),
              ),
              const SizedBox(height: 24),
              BlocBuilder<AuthBloc, AuthState>(
                builder: (context, state) {
                  return state.status == AuthStatus.loading
                      ? const Center(child: CircularProgressIndicator())
                      : CustomButton(
                          text: "Sign In",
                          onTap: () {
                            // TODO: Add form validation
                            context.read<AuthBloc>().add(LoginSubmitted(
                                identifier: emailController.text,
                                password: passwordController.text,
                              ));
                          },
                        );

                },
              ),
              const SizedBox(height: 24),
              const Row(
                children: [
                  Expanded(child: Divider()),
                  Padding(
                    padding: EdgeInsets.symmetric(horizontal: 8.0),
                    child: Text("OR CONTINUE WITH", style: TextStyle(color: Colors.grey)),
                  ),
                  Expanded(child: Divider()),
                ],
              ),
              const SizedBox(height: 24),
              const SizedBox(
                width: double.infinity,
                child: GoogleSignInButton(),
              ),
            ],
          ),
        ),
      ],
    );
  }
}