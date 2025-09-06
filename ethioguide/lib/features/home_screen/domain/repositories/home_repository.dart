import '../entities/home_data.dart';

// This is a simplified contract. We are fetching all data in one go.
abstract class HomeRepository {
  List<QuickAction> getQuickActions();
  List<ContentCard> getContentCards();
  List<PopularService> getPopularServices();
}