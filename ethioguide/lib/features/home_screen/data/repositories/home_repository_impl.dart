import 'package:ethioguide/features/home_screen/domain/entities/home_data.dart';
import 'package:ethioguide/features/home_screen/domain/repositories/home_repository.dart';
import 'package:flutter/material.dart';

class HomeRepositoryImpl implements HomeRepository {
  @override
  List<QuickAction> getQuickActions() {
    
    return const [
      QuickAction(
        icon: Icons.add,
        title: 'Start New Process',
        subtitle: 'Begin a new government service',
        routeName: '/Procedure', 
      ),
      QuickAction(
        icon: Icons.work_outline,
        title: 'My Workspace',
        subtitle: 'Manage your ongoing procedures',
        routeName: '/workspace',
      ),
    ];
  }

  @override
  List<ContentCard> getContentCards() {
    // Mocked data for the larger cards
    return const [
      ContentCard(
        sectionTitle: 'Organizations',
        icon: Icons.business_outlined,
        title: 'Government Organizations',
        subtitle: 'Browse and connect with government offices, agencies, and public service organizations',
        details: ['15+ organizations', 'Multiple locations'],
        routeName: '/placeholder',
      ),
      ContentCard(
        sectionTitle: 'AI Assistant',
        icon: Icons.chat_bubble_outline,
        title: 'Legal Assistant Chat',
        subtitle: 'Get instant AI-powered guidance for Ethiopian government processes and legal requirements',
        details: ['Instant responses', 'Step-by-step guides'],
        routeName: '/aiChat', // This one goes to a specific route
      ),
      ContentCard(
        sectionTitle: 'Community & Support',
        icon: Icons.people_outline,
        title: 'Community Discussions',
        subtitle: 'Ask questions, share experiences, and help others navigate government processes',
        details: ['1.2k members', '42 active discussions'],
        routeName: '/discussion',
      ),
    ];
  }

  @override
  List<PopularService> getPopularServices() {
    // Mocked data for the popular services list
    return const [
      PopularService(
        icon: Icons.article_outlined,
        title: 'Passport Application',
        category: 'Travel',
        timeEstimate: '2-6 weeks',
        routeName: '/placeholder',
      ),
      PopularService(
        icon: Icons.directions_car_outlined,
        title: "Driver's License Renewal",
        category: 'Transportation',
        timeEstimate: '1-2 hours',
        routeName: '/placeholder',
      ),
      // Add more services as needed
    ];
  }
}