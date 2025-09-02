import 'package:ethioguide/core/components/button.dart';
import 'package:ethioguide/core/components/textfield.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/auth_bloc.dart';
import '../bloc/auth_event.dart';
import '../bloc/auth_state.dart';

class ForgotPasswordView extends StatelessWidget {
  const ForgotPasswordView({super.key});

  @override
  Widget build(BuildContext context) {
    final emailController = TextEditingController();

    return Column(
      key: const ValueKey('forgotPasswordView'),
      children: [
        const Text("Forgot Password", style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        const Text("No more confusion. EthioGuide shows the way.", textAlign: TextAlign.center, style: TextStyle(fontSize: 16, color: Colors.grey)),
        const SizedBox(height: 32),
        Container(
          padding: const EdgeInsets.all(24.0),
          decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(16)),
          child: BlocBuilder<AuthBloc, AuthState>(
            builder: (context, state) {
              // If the reset link has been sent, show a success message
              if (state.status == AuthStatus.resetLinkSent) {
                return const Column(
                  children: [
                    Icon(Icons.check_circle_outline, color: Colors.green, size: 60),
                    SizedBox(height: 16),
                    Text("Reset Link Sent!", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
                    SizedBox(height: 8),
                    Text("Please check your email to reset your password.", textAlign: TextAlign.center),
                  ],
                );
              }

              // Otherwise, show the form
              return Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text("Email", style: TextStyle(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 8),
                  customTextField(hintText: "Enter your email", controller: emailController, prefixIcon: Icons.email_outlined),
                  const SizedBox(height: 24),
                  state.status == AuthStatus.loading
                      ? const Center(child: CircularProgressIndicator())
                      : CustomButton(
                          text: "Send Reset Link",
                          onTap: () {
                            context.read<AuthBloc>().add(ForgotPasswordSubmitted(email: emailController.text));
                          },
                        ),
                  const SizedBox(height: 16),
                  Center(
                    child: TextButton(
                      onPressed: () => context.read<AuthBloc>().add(const AuthViewSwitched(AuthView.login)),
                      child: const Text("Back to login"),
                    ),
                  )
                ],
              );
            },
          ),
        ),
      ],
    );
  }
}