import 'package:flutter/material.dart';
import 'core/config/app_router.dart'; // <-- 1. Import your new router file.

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    // 2. Use the MaterialApp.router constructor.
    return MaterialApp.router(
      // 3. Pass your router configuration to the app.
      routerConfig: router,
      title: 'EthioGuide',
      debugShowCheckedModeBanner: false,
    );
  }
}