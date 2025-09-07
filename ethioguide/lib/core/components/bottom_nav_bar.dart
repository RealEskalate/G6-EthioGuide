import 'package:ethioguide/core/config/route_names.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

Widget bottomNav({required BuildContext context, required selectedIndex}) {
  return Container(
    decoration: BoxDecoration(
      color: Colors.white,
      borderRadius: const BorderRadius.vertical(top: Radius.circular(20)),
      boxShadow: [
        BoxShadow(
          color: Colors.grey.withOpacity(0.2),
          spreadRadius: 1,
          blurRadius: 4,
          offset: const Offset(0, -2),
        ),
      ],
    ),
    child: BottomNavigationBar(
      type: BottomNavigationBarType.fixed,
      currentIndex: selectedIndex, // AI page
      selectedItemColor: Colors.teal,
      unselectedItemColor: Colors.grey[600],
      selectedIconTheme: const IconThemeData(color: Colors.teal, size: 28),
      unselectedIconTheme: const IconThemeData(color: Colors.grey, size: 24),
      backgroundColor: Colors.transparent,
      elevation: 0,
      selectedLabelStyle: const TextStyle(fontWeight: FontWeight.bold),
      unselectedLabelStyle: const TextStyle(fontWeight: FontWeight.normal),
      items: [
        BottomNavigationBarItem(
          icon: AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            transform: Matrix4.identity()..scale(2 == 0 ? 1.2 : 1.0),
            child: const Icon(Icons.home),
          ),
          label: 'Home',
        ),
        BottomNavigationBarItem(
          icon: AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            transform: Matrix4.identity()..scale(2 == 1 ? 1.2 : 1.0),
            child: const Icon(Icons.work),
          ),
          label: 'Workspace',
        ),
        BottomNavigationBarItem(
          icon: AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            transform: Matrix4.identity()..scale(2 == 2 ? 1.2 : 1.0),
            child: const Icon(Icons.chat),
          ),
          label: 'AI',
        ),
        BottomNavigationBarItem(
          icon: AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            transform: Matrix4.identity()..scale(2 == 3 ? 1.2 : 1.0),
            child: const Icon(Icons.person),
          ),
          label: 'Profile',
        ),
      ],
      onTap: (index) {
        switch (index) {
          case 0:
            if (index != selectedIndex) {
              context.goNamed(RouteNames.home);
            }
            break;
          case 1:
            if (index != selectedIndex) {
              context.goNamed(RouteNames.workspace);
            }
            break;
          case 2:
            if (index != selectedIndex) {
              context.goNamed(RouteNames.aiChat);
            }
            break;
          case 3:
            if (index != selectedIndex) {
              context.goNamed(RouteNames.profile);
            }
            break;
        }
      },
    ),
  );
}
