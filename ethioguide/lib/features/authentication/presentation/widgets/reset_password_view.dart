import 'package:ethioguide/core/components/button.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
// Make sure this is your correct textfield
import 'package:ethioguide/core/utils/validators.dart';
import '../bloc/auth_bloc.dart';
import '../bloc/auth_event.dart';
import '../bloc/auth_state.dart';
import 'package:ethioguide/core/components/textfield.dart';
class ResetPasswordView extends StatefulWidget {
  final String resetToken;
  
  const ResetPasswordView({super.key, required this.resetToken});

  @override
  State<ResetPasswordView> createState() => _ResetPasswordViewState();
}

class _ResetPasswordViewState extends State<ResetPasswordView> {
  final _formKey = GlobalKey<FormState>();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();

  @override
  void dispose() {
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  void _submitForm() {
    if (_formKey.currentState?.validate() ?? false) {
      if (_passwordController.text != _confirmPasswordController.text) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text("Passwords do not match!")),
        );
        return;
      }
      
      // Use the resetToken passed from the router
      context.read<AuthBloc>().add(ResetPasswordSubmitted(
            resetToken: widget.resetToken,
            newPassword: _passwordController.text,
          ));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: Column(
        key: const ValueKey('resetPasswordView'),
        children: [
          const Text("Reset Password", style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          const Text("Create a new password for your account.", textAlign: TextAlign.center, style: TextStyle(fontSize: 16, color: Colors.grey)),
          const SizedBox(height: 32),
          Container(
            padding: const EdgeInsets.all(24.0),
            decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(16)),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                
                const Text("New Password", style: TextStyle(fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),
                TextFormField(
                  controller: _passwordController,
                  validator: Validators.validatePassword,
                  obscureText: true,
                  decoration: const InputDecoration(hintText: 'Enter your new password'),
                ),
                
                const SizedBox(height: 16),

                const Text("Confirm New Password", style: TextStyle(fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),
                TextFormField(
                  controller: _confirmPasswordController,
                  validator: Validators.validatePassword,
                  obscureText: true,
                  decoration: const InputDecoration(hintText: 'Confirm your new password'),
                ),

                const SizedBox(height: 24),
                
                BlocBuilder<AuthBloc, AuthState>(
                  builder: (context, state) {
                    return state.status == AuthStatus.loading
                        ? const Center(child: CircularProgressIndicator())
                        : CustomButton(
                            text: "Update Password",
                            onTap: _submitForm,
                          );
                  },
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}