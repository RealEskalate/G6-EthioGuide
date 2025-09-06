import '../entities/home_data.dart';
import '../repositories/home_repository.dart';

class GetHomeData {
  final HomeRepository repository;

  GetHomeData(this.repository);

  // The use case returns a tuple with all the data lists.
  (List<QuickAction>, List<ContentCard>, List<PopularService>) call() {
    final quickActions = repository.getQuickActions();
    final contentCards = repository.getContentCards();
    final popularServices = repository.getPopularServices();
    return (quickActions, contentCards, popularServices);
  }
}