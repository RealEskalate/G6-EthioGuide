import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:get_it/get_it.dart';
import 'package:dio/dio.dart';
import '../../data/datasources/workspace_discussion_remote_data_source.dart';
import '../../data/repositories/workspace_discussion_repository_impl.dart';
import '../../domain/repositories/workspace_discussion_repository.dart';
import '../../domain/usecases/get_comments.dart';
import '../../domain/usecases/get_community_stats.dart';
import '../../domain/usecases/get_discussions.dart';
import '../bloc/workspace_discussion_bloc.dart';
import 'workspace_discussion_page.dart';

/// Example of how to set up dependency injection and integrate the feature
class ExampleIntegration {
  static void setupDependencies() {
    final getIt = GetIt.instance;

    // Register Dio instance
    getIt.registerLazySingleton<Dio>(() => Dio());

    // Register data source
    getIt.registerLazySingleton<WorkspaceDiscussionRemoteDataSource>(
      () => WorkspaceDiscussionRemoteDataSourceImpl(dio: getIt<Dio>()),
    );

    // Register repository
    getIt.registerLazySingleton<WorkspaceDiscussionRepository>(
      () => WorkspaceDiscussionRepositoryImpl(getIt<WorkspaceDiscussionRemoteDataSource>()),
    );

    // Register use cases
    getIt.registerLazySingleton<GetCommunityStats>(
      () => GetCommunityStats(getIt<WorkspaceDiscussionRepository>()),
    );
    getIt.registerLazySingleton<GetDiscussions>(
      () => GetDiscussions(getIt<WorkspaceDiscussionRepository>()),
    );
    getIt.registerLazySingleton<CreateDiscussion>(
      () => CreateDiscussion(getIt<WorkspaceDiscussionRepository>(), title: '', content: '', tags: [], category: ''),
    );
    getIt.registerLazySingleton<LikeDiscussion>(
      () => LikeDiscussion(getIt<WorkspaceDiscussionRepository>() as String),
    );
    getIt.registerLazySingleton<ReportDiscussion>(
      () => ReportDiscussion(getIt<WorkspaceDiscussionRepository>() as String),
    );
    getIt.registerLazySingleton<GetComments>(
      () => GetComments(getIt<WorkspaceDiscussionRepository>()),
    );
    getIt.registerLazySingleton<AddComment>(
      () => AddComment(getIt<WorkspaceDiscussionRepository>(), discussionId: '', content: ''),
    );
    getIt.registerLazySingleton<LikeComment>(
      () => LikeComment(getIt<WorkspaceDiscussionRepository>() as String),
    );
    getIt.registerLazySingleton<ReportComment>(
      () => ReportComment(getIt<WorkspaceDiscussionRepository>() as String),
    );

    // Register bloc
    getIt.registerFactory<WorkspaceDiscussionBloc>(
      () => WorkspaceDiscussionBloc(
        getCommunityStats: getIt<GetCommunityStats>(),
        getDiscussions: getIt<GetDiscussions>(),
        createDiscussion: getIt<CreateDiscussion>(),
        likeDiscussion: getIt<LikeDiscussion>(),
        reportDiscussion: getIt<ReportDiscussion>(),
        getComments: getIt<GetComments>(),
        addComment: getIt<AddComment>(),
        likeComment: getIt<LikeComment>(),
        reportComment: getIt<ReportComment>(),
      ),
    );
  }

  /// Example of how to navigate to the workspace discussion page
  static void navigateToDiscussionPage(BuildContext context) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => BlocProvider(
          create: (context) => GetIt.instance<WorkspaceDiscussionBloc>(),
          child: const WorkspaceDiscussionPage(),
        ),
      ),
    );
  }

  /// Example of how to use the feature in a list item
  static Widget buildDiscussionListItem(BuildContext context, String title, String subtitle) {
    return ListTile(
      leading: Icon(
        Icons.chat_bubble,
        color: Colors.teal[600],
      ),
      title: Text(title),
      subtitle: Text(subtitle),
      trailing: const Icon(Icons.arrow_forward_ios),
      onTap: () => navigateToDiscussionPage(context),
    );
  }
}

/// Example usage in a main app or other page
class ExampleUsagePage extends StatelessWidget {
  const ExampleUsagePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Example Usage'),
      ),
      body: ListView(
        children: [
          ExampleIntegration.buildDiscussionListItem(
            context,
            'Community Discussions',
            'Share knowledge and get help',
          ),
          const Divider(),
          ListTile(
            leading: const Icon(Icons.info),
            title: const Text('Feature Information'),
            subtitle: const Text('This is the Workspace Community Discussion feature'),
            onTap: () {
              showDialog(
                context: context,
                builder: (context) => AlertDialog(
                  title: const Text('Workspace Discussion Feature'),
                  content: const Text(
                    'This feature provides a complete community discussion system with:\n\n'
                    '• Create and manage discussions\n'
                    '• Community statistics and trending topics\n'
                    '• Comment system with likes and reports\n'
                    '• Advanced filtering and search\n'
                    '• Clean architecture with BLoC state management',
                  ),
                  actions: [
                    TextButton(
                      onPressed: () => Navigator.pop(context),
                      child: const Text('OK'),
                    ),
                  ],
                ),
              );
            },
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => ExampleIntegration.navigateToDiscussionPage(context),
        child: const Icon(Icons.chat),
      ),
    );
  }
}

/// Example of how to initialize the app with dependencies
class ExampleApp extends StatelessWidget {
  const ExampleApp({super.key});

  @override
  Widget build(BuildContext context) {
    // Setup dependencies when app starts
    ExampleIntegration.setupDependencies();

    return MaterialApp(
      title: 'EthioGuide - Workspace Discussion Example',
      theme: ThemeData(
        primarySwatch: Colors.teal,
        useMaterial3: true,
      ),
      home: const ExampleUsagePage(),
    );
  }
}

/// Example of how to integrate with existing navigation
class NavigationExample {
  static Route<dynamic> generateRoute(RouteSettings settings) {
    switch (settings.name) {
      case '/discussions':
        return MaterialPageRoute(
          builder: (context) => BlocProvider(
            create: (context) => GetIt.instance<WorkspaceDiscussionBloc>(),
            child: const WorkspaceDiscussionPage(),
          ),
        );
      default:
        return MaterialPageRoute(
          builder: (context) => const Scaffold(
            body: Center(
              child: Text('Route not found'),
            ),
          ),
        );
    }
  }
}

/// Example of how to use in a bottom navigation app
class BottomNavigationExample extends StatefulWidget {
  const BottomNavigationExample({super.key});

  @override
  State<BottomNavigationExample> createState() => _BottomNavigationExampleState();
}

class _BottomNavigationExampleState extends State<BottomNavigationExample> {
  int _currentIndex = 0;

  final List<Widget> _pages = [
    const Center(child: Text('Home Page')),
    BlocProvider(
      create: (context) => GetIt.instance<WorkspaceDiscussionBloc>(),
      child: const WorkspaceDiscussionPage(),
    ),
    const Center(child: Text('Profile Page')),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _pages[_currentIndex],
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _currentIndex,
        onTap: (index) {
          setState(() {
            _currentIndex = index;
          });
        },
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.home),
            label: 'Home',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.chat_bubble),
            label: 'Discussions',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.person),
            label: 'Profile',
          ),
        ],
      ),
    );
  }
}
