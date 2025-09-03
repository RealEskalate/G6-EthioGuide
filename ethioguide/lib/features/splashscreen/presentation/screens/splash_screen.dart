import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../data/repositories/splash_repository_impl.dart';
import '../../domain/usecases/complete_splash_task.dart';
import '../bloc/splash_bloc.dart';
import '../bloc/splash_event.dart';
import '../bloc/splash_state.dart';

class SplashScreen extends StatelessWidget {
  const SplashScreen({super.key});

  @override
  Widget build(BuildContext context) {
    // BlocProvider creates the BLoC and provides it to the widgets below.
    return BlocProvider(
      // 1. We create the entire dependency chain here.
      create: (context) => SplashBloc(
        completeSplashTask: CompleteSplashTask(
          SplashRepositoryImpl(),
        ),
      )..add(StartSplashTimer()), // 2. Add the event to start the timer right away.
      
      // BlocListener listens for state changes to perform actions like navigation.
      child: BlocListener<SplashBloc, SplashState>(
        listener: (context, state) {
          // 3. When the BLoC emits SplashCompleted...
          if (state is SplashCompleted) {
            context.go('/onboarding');
          }
        },
        child: const Scaffold(
          backgroundColor: Colors.black,
          body: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                
                Image(
                  image: AssetImage('assets/images/dark_logo.jpg'),
                  width: 300.0, 
                ),
                SizedBox(height: 24),
                // 7. A loading indicator to show work is being done.
                CircularProgressIndicator(
                  valueColor: AlwaysStoppedAnimation<Color>(Color(0xFFa7b3b9)), // Hit Gray
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}