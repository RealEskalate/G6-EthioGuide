# Workspace Community Discussion Feature

A comprehensive Flutter feature module for managing community discussions in the EthioGuide workspace application.

## 🏗️ Architecture

This feature follows **Clean Architecture** principles with three distinct layers:

### **Domain Layer** (`/domain`)
- **Entities**: Core business objects (`Discussion`, `Comment`, `User`, `CommunityStats`)
- **Repository Contracts**: Abstract interfaces defining data operations
- **Use Cases**: Business logic for each feature operation

### **Data Layer** (`/data`)
- **Models**: Data models extending domain entities with JSON serialization
- **Data Sources**: Remote data source using Dio for API integration
- **Repository Implementation**: Concrete implementation of repository contracts

### **Presentation Layer** (`/presentation`)
- **BLoC**: State management for all discussion operations
- **Pages**: Main UI screens (list, detail, create)
- **Widgets**: Reusable UI components

## 🚀 Features

### **Core Functionality**
- ✅ Create new discussions with title, content, tags, and category
- ✅ View community statistics (members, discussions, active users)
- ✅ Browse discussions with filtering and search
- ✅ Like and report discussions
- ✅ Add comments to discussions
- ✅ Like and report comments

### **Filtering & Search**
- ✅ Filter by tags (trending topics)
- ✅ Filter by category
- ✅ Sort by most recent, trending, or popular
- ✅ Search discussions and procedures

### **User Experience**
- ✅ Modern, clean UI design matching the provided screenshots
- ✅ Fully responsive and scrollable layouts
- ✅ Real-time state management with BLoC
- ✅ Error handling and loading states
- ✅ Form validation for creating discussions

## 📱 UI Components

### **Main Screens**
1. **WorkspaceDiscussionPage**: Main discussion list with community stats
2. **CreateDiscussionPage**: Form for creating new discussions
3. **DiscussionDetailPage**: View discussion details and comments

### **Reusable Widgets**
- `CommunityStatsCard`: Overview of community metrics
- `TrendingTopics`: Interactive trending tags
- `FilterControls`: Category and sorting filters
- `DiscussionCard`: Individual discussion display
- `CommentItem`: Individual comment display

## 🔌 API Integration

### **Endpoints Used**
- `GET /api/discussions` - Fetch discussions with filtering
- `POST /api/discussions` - Create new discussion
- `POST /api/discussions/:id/like` - Like a discussion
- `POST /api/discussions/:id/report` - Report a discussion
- `GET /api/discussions/:id/comments` - Fetch comments
- `POST /api/discussions/:id/comments` - Add comment
- `POST /api/comments/:id/like` - Like a comment
- `POST /api/comments/:id/report` - Report a comment

### **Mock Data**
Currently uses mock data for development. Replace with real API calls by updating the `WorkspaceDiscussionRemoteDataSourceImpl`.

## 🧪 Testing

### **Test Coverage**
- ✅ Repository tests (success/failure scenarios)
- ✅ BLoC tests (state transitions)
- ✅ Mock Dio for API testing

### **Running Tests**
```bash
# Run all tests
flutter test

# Run specific test file
flutter test test/features/workspace_discussion/presentation/bloc/workspace_discussion_bloc_test.dart
```

## 📦 Dependencies

### **Core Dependencies**
- `flutter_bloc`: State management
- `dio`: HTTP client for API calls
- `equatable`: Value equality for entities
- `dartz`: Functional programming utilities

### **Dev Dependencies**
- `flutter_test`: Flutter testing framework
- `bloc_test`: BLoC testing utilities
- `mockito`: Mocking framework

## 🚀 Getting Started

### **1. Setup Dependencies**
```dart
// In your main.dart or dependency injection setup
final getIt = GetIt.instance;

// Register data source
getIt.registerLazySingleton<WorkspaceDiscussionRemoteDataSource>(
  () => WorkspaceDiscussionRemoteDataSourceImpl(dio: Dio()),
);

// Register repository
getIt.registerLazySingleton<WorkspaceDiscussionRepository>(
  () => WorkspaceDiscussionRepositoryImpl(getIt()),
);

// Register use cases
getIt.registerLazySingleton<GetCommunityStats>(
  () => GetCommunityStats(getIt()),
);
// ... register other use cases

// Register BLoC
getIt.registerFactory<WorkspaceDiscussionBloc>(
  () => WorkspaceDiscussionBloc(
    getCommunityStats: getIt(),
    getDiscussions: getIt(),
    createDiscussion: getIt(),
    likeDiscussion: getIt(),
    reportDiscussion: getIt(),
    getComments: getIt(),
    addComment: getIt(),
    likeComment: getIt(),
    reportComment: getIt(),
  ),
);
```

### **2. Navigate to Discussion Page**
```dart
Navigator.push(
  context,
  MaterialPageRoute(
    builder: (context) => BlocProvider(
      create: (context) => getIt<WorkspaceDiscussionBloc>(),
      child: const WorkspaceDiscussionPage(),
    ),
  ),
);
```

### **3. Use in Your App**
```dart
class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: BlocProvider(
        create: (context) => getIt<WorkspaceDiscussionBloc>(),
        child: const WorkspaceDiscussionPage(),
      ),
    );
  }
}
```

## 🎨 Customization

### **Theming**
- Update colors in the widget files
- Modify text styles using Theme.of(context)
- Customize card elevations and borders

### **API Integration**
- Replace mock data in `WorkspaceDiscussionRemoteDataSourceImpl`
- Update base URL and endpoint paths
- Add authentication headers if required

### **UI Modifications**
- Modify widget layouts in the presentation layer
- Add new UI components as needed
- Update form validation rules

## 📋 File Structure

```
features/workspace_discussion/
├── data/
│   ├── models/
│   │   ├── user_model.dart
│   │   ├── comment_model.dart
│   │   ├── discussion_model.dart
│   │   └── community_stats_model.dart
│   ├── datasources/
│   │   └── workspace_discussion_remote_data_source.dart
│   └── repositories/
│       └── workspace_discussion_repository_impl.dart
├── domain/
│   ├── entities/
│   │   ├── user.dart
│   │   ├── comment.dart
│   │   ├── discussion.dart
│   │   └── community_stats.dart
│   ├── repositories/
│   │   └── workspace_discussion_repository.dart
│   └── usecases/
│       ├── get_community_stats.dart
│       ├── get_discussions.dart
│       ├── create_discussion.dart
│       ├── like_discussion.dart
│       ├── report_discussion.dart
│       ├── get_comments.dart
│       ├── add_comment.dart
│       ├── like_comment.dart
│       └── report_comment.dart
├── presentation/
│   ├── bloc/
│   │   └── workspace_discussion_bloc.dart
│   ├── pages/
│   │   ├── workspace_discussion_page.dart
│   │   ├── create_discussion_page.dart
│   │   └── discussion_detail_page.dart
│   └── widgets/
│       ├── community_stats_card.dart
│       ├── trending_topics.dart
│       ├── filter_controls.dart
│       └── discussion_card.dart
└── test/
    ├── data/repositories/
    │   └── workspace_discussion_repository_impl_test.dart
    └── presentation/bloc/
        └── workspace_discussion_bloc_test.dart
```

## 🔧 Troubleshooting

### **Common Issues**
1. **BLoC not found**: Ensure proper dependency injection setup
2. **API errors**: Check network connectivity and API endpoints
3. **State not updating**: Verify BLoC event dispatching
4. **UI not rendering**: Check if BLoC state is properly handled

### **Debug Tips**
- Use `BlocListener` to monitor state changes
- Add print statements in BLoC methods
- Check console for error messages
- Verify mock data is being returned correctly

## 📈 Future Enhancements

### **Planned Features**
- [ ] Real-time notifications for new comments
- [ ] Rich text editor for discussions
- [ ] File/image attachments
- [ ] User profiles and reputation system
- [ ] Moderation tools for admins
- [ ] Discussion categories management
- [ ] Search with advanced filters
- [ ] Offline support with local caching

### **Performance Optimizations**
- [ ] Pagination for large discussion lists
- [ ] Image lazy loading
- [ ] Debounced search input
- [ ] Optimistic updates for likes/comments

## 🤝 Contributing

1. Follow the existing code structure and patterns
2. Add tests for new functionality
3. Update documentation as needed
4. Ensure all tests pass before submitting

## 📄 License

This feature module is part of the EthioGuide application and follows the same licensing terms.
