import 'package:ethioguide/core/components/button.dart';
import 'package:ethioguide/core/components/textfield.dart';
import 'package:ethioguide/features/authentication/presentation/bloc/auth_bloc.dart';
import 'package:ethioguide/features/authentication/presentation/bloc/auth_event.dart';
import 'package:ethioguide/features/authentication/presentation/bloc/auth_state.dart';
import 'package:ethioguide/features/authentication/presentation/widgets/auth_toggle_buttons.dart';
import 'package:ethioguide/features/authentication/presentation/widgets/google_sign_in_button.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';


class SignUpView extends StatelessWidget {
  const SignUpView({super.key});

  @override
  Widget build(BuildContext context) {
    final nameController = TextEditingController();
    final usernameController = TextEditingController();
    final emailController = TextEditingController();
    final passwordController = TextEditingController();
    final confirmPasswordController = TextEditingController();

    return Column(
      key: const ValueKey('signUpView'),
      children: [
        // ... (Header Text)
        const Text("Create an Account", style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
        const SizedBox(height: 32),
        const AuthToggleButtons(activeView: AuthView.signUp),
        const SizedBox(height: 24),
        Container(
          padding: const EdgeInsets.all(24.0),
          decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(16)),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text("Name", style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              customTextField(hintText: "Enter your  name", controller: nameController, prefixIcon: Icons.person_outline),
              const SizedBox(height: 16),

               const Text("Username", style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              customTextField(hintText: "Enter your username", controller: usernameController, prefixIcon: Icons.alternate_email),
              
              const SizedBox(height: 16),

              const Text("Email", style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              customTextField(hintText: "Enter your email", controller: emailController, prefixIcon: Icons.email_outlined),
               const SizedBox(height: 16),
              const Text("Password", style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),

               BlocBuilder<AuthBloc, AuthState>(
                builder: (context, state) {
                  return customTextField(
                    hintText: "Create a password",
                    controller: passwordController,
                    prefixIcon: Icons.lock_outline,
                    obscureText: !state.isPasswordVisible, 
                    suffixIcon: IconButton(
                      icon: Icon(
                        state.isPasswordVisible ? Icons.visibility_off : Icons.visibility,
                        color: Colors.grey,
                      ),
                      onPressed: () {
                        context.read<AuthBloc>().add(PasswordVisibilityToggled());
                      },
                    ),
                  );
                },
              ),

               const SizedBox(height: 16),

              // ADDED: Confirm Password field
              const Text("Confirm Password", style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              BlocBuilder<AuthBloc, AuthState>(
                builder: (context, state) {
                   return customTextField(
                    hintText: "Confirm your password",
                    controller: confirmPasswordController,
                    prefixIcon: Icons.lock_outline,
                    obscureText: !state.isPasswordVisible,
                  );
                },
              ),

              const SizedBox(height: 24),
              BlocBuilder<AuthBloc, AuthState>(
                builder: (context, state) {
                   return state.status == AuthStatus.loading
                      ? const Center(child: CircularProgressIndicator())
                      : CustomButton(
                          text: "Create Account",
                          onTap: () {
                            context.read<AuthBloc>().add(SignUpSubmitted(
                                name: nameController.text,
                                username: usernameController.text, // Add a username field if needed
                                email: emailController.text,
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