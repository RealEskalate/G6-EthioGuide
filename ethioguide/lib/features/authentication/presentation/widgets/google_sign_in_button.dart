import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/auth_bloc.dart';
import '../bloc/auth_event.dart';

class GoogleSignInButton extends StatelessWidget {
  const GoogleSignInButton({super.key});

  @override
  Widget build(BuildContext context) {
    return OutlinedButton.icon(
      icon: Image.asset('assets/images/google_logo.png', height: 24.0), 
      label: const Text(
        'Continue with Google',
        style: TextStyle(color: Colors.black54, fontWeight: FontWeight.bold),
      ),
      style: OutlinedButton.styleFrom(
        foregroundColor: Colors.grey,
        side: BorderSide(color: Colors.grey.shade300),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        padding: const EdgeInsets.symmetric(vertical: 16),
      ),
      onPressed: () {
        // Dispatch the event to the AuthBloc
        context.read<AuthBloc>().add(GoogleSignInSubmitted());
      },
    );
  }
}