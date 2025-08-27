import 'package:ethioguide/core/config/app_theme.dart';
import 'package:flutter/material.dart';
import 'core/config/app_router.dart'; // <-- 1. Import your new router file.

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  
  MyApp({super.key});

    // default light

    ThemeMode _themeMode = ThemeMode.light; 


  @override
  Widget build(BuildContext context) {

    // 2. Use the MaterialApp.router constructor.
    return MaterialApp.router(
      themeMode: _themeMode,
      theme:AppTheme.lightTheme,
      darkTheme:AppTheme.darkTheme,
      routerConfig: router,
      title: 'EthioGuide',
      debugShowCheckedModeBanner: false,
    );
  }
}