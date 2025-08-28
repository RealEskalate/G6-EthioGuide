import 'package:flutter/material.dart';
import '../config/app_color.dart';

class DashboardCard extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final void Function()? onTap;

  const DashboardCard({
    super.key,
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    // 1. GestureDetector makes the entire card tappable.
    return GestureDetector(
      onTap: onTap,
      child: Container(
        // 2. This outer container creates the shadow effect.
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(16.0),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.08),
              blurRadius: 12,
              offset: const Offset(0, 4),
            ),
          ],
        ),
        // 3. ClipRRect is crucial! It forces its children (the gradient and white boxes)
        // to have the same rounded corners, creating the seamless look.
        child: ClipRRect(
          borderRadius: BorderRadius.circular(16.0),
          child: Column(
            children: [
              // --- TOP GRADIENT PART ---
              Container(
                height: 80, // Fixed height for the header
                width: double.infinity, // Fills the full width of the card
                padding: const EdgeInsets.all(12.0),
                decoration: const BoxDecoration(
                  gradient: AppColors.gradient, // Using the gradient from your colors file
                ),
                child: Align(
                  alignment: Alignment.topLeft,
                  child: Icon(icon, color: Colors.white, size: 28),
                ),
              ),

              // --- BOTTOM WHITE PART ---
              // 4. Expanded ensures this part fills all remaining vertical space in the card.
              Expanded(
                child: Container(
                  width: double.infinity,
                  color: Colors.white,
                  padding: const EdgeInsets.all(12.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start, // Aligns text to the left
                    mainAxisAlignment: MainAxisAlignment.center, // Centers the text block vertically
                    children: [
                      Text(
                        title,
                        style: const TextStyle(
                          fontWeight: FontWeight.bold,
                          fontSize: 16,
                          color: AppColors.darkGreenColor, // Assuming you add this to AppColors
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        subtitle,
                        style: const TextStyle(
                          fontSize: 14,
                          color: AppColors.graycolor, // Assuming you add this
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}