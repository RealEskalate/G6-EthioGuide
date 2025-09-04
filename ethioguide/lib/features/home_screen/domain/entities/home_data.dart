import 'package:equatable/equatable.dart';
import 'package:flutter/material.dart';

// Entity for the "Quick Actions" grid
class QuickAction extends Equatable {
  final IconData icon;
  final String title;
  final String subtitle;
  final String routeName; // To handle navigation

  const QuickAction({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.routeName,
  });

  @override
  List<Object?> get props => [icon, title, subtitle, routeName];
}

// Entity for the main content cards (Organizations, AI, Community)
class ContentCard extends Equatable {
  final String sectionTitle;
  final IconData icon;
  final String title;
  final String subtitle;
  final List<String> details;
  final String routeName;

  const ContentCard({
    required this.sectionTitle,
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.details,
    required this.routeName,
  });
  
  @override
  List<Object?> get props => [sectionTitle, icon, title, subtitle, details, routeName];
}

// Entity for a single service in the "Popular Services" list
class PopularService extends Equatable {
  final IconData icon;
  final String title;
  final String category;
  final String timeEstimate;
  final String routeName;

  const PopularService({
    required this.icon,
    required this.title,
    required this.category,
    required this.timeEstimate,
    required this.routeName,
  });
  
  @override
  List<Object?> get props => [icon, title, category, timeEstimate, routeName];
}