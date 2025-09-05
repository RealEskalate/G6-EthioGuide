import 'package:ethioguide/core/config/app_color.dart';
import 'package:flutter/material.dart';


class IconContainer extends StatelessWidget {
  final IconData icon;

  const IconContainer({
    super.key,
    required this.icon,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 120,
      height: 120,
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [
            AppColors.patina,
            AppColors.darkGreenColor,
          ],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(30), 
      ),
      child: Center(
        child: Icon(
          icon,
          color: Colors.white,
          size: 60,
        ),
      ),
    );
  }
}