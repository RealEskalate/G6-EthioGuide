import 'package:ethioguide/core/config/app_theme.dart';
import 'package:flutter/material.dart';
import 'core/config/app_router.dart';
import 'injection_container.dart' as di; // Import the container with a prefix

// The main function now needs to be `async`
Future<void> main() async {
  // This line is required to ensure that plugin services are initialized
  // before `runApp()` is called when `main` is async.
  WidgetsFlutterBinding.ensureInitialized();
  
  // Await the initialization of all our dependencies
  await di.init();
  
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  MyApp({super.key});

  // NOTE: This _themeMode variable should be managed by a state management solution
  // (like a ThemeCubit) in a real app, not as a local variable here.
  // For now, this is okay.
  final ThemeMode _themeMode = ThemeMode.light;

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      themeMode: _themeMode,
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      routerConfig: router,
      title: 'EthioGuide',
      debugShowCheckedModeBanner: false,
    );
  }
}