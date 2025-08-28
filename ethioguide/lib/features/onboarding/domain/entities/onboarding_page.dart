import 'package:equatable/equatable.dart';
import 'package:flutter/material.dart';

class OnboardingPage extends Equatable {
  final IconData icon;
  final String title;
  final String subtitle;
  final String description;

  const OnboardingPage({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.description,
  });

  @override
  List<Object?> get props => [icon, title, subtitle, description];
}