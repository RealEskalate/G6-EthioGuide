import 'package:ethioguide/core/config/app_color.dart';
import 'package:flutter/material.dart';

// ignore: camel_case_types
class customTextField extends StatelessWidget {
  final String hintText;
  final bool obscureText;
  final TextEditingController controller;
  final IconData? prefixIcon;
  final TextInputType? keyboardType;
  final double? width; 
  final double? height; 
  final double borderRadius;
  final Widget? suffixIcon; 
  final String? errorText;

  const customTextField({super.key,
    required this.hintText,
    required this.controller,
    this.obscureText = false,
    this.prefixIcon,
    this.keyboardType,
    this.width,
    this.height,
    this.borderRadius = 12.0,
    this.suffixIcon,
     this.errorText,});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: width,
      height: height,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 25.0),
        child: TextField(
          controller: controller,
          obscureText: obscureText,
          keyboardType: keyboardType,
          decoration: InputDecoration(
            enabledBorder: OutlineInputBorder(
              //  Use the borderRadius property
              borderRadius: BorderRadius.circular(borderRadius),
              borderSide: BorderSide(color: AppColors.graycolor),
            ),
            focusedBorder: OutlineInputBorder(
              //  Use the borderRadius property
              borderRadius: BorderRadius.circular(borderRadius),
              borderSide: BorderSide(color: AppColors.darkGreenColor),
            ),
            fillColor: Colors.white,
            filled: true,
            hintText: hintText,
            hintStyle: TextStyle(color: AppColors.graycolor.withOpacity(0.8)),
            prefixIcon: prefixIcon != null
                ? Icon(
                    prefixIcon,
                    color: AppColors.graycolor,
                  )
                : null,
                suffixIcon: suffixIcon,
          ),
        ),
      ),
    );
  }
}