import 'package:equatable/equatable.dart';
import 'package:ethioguide/features/home_screen/domain/entities/home_data.dart';

// We can use a single state for this simple, synchronous feature.
class HomeState extends Equatable {
  final List<QuickAction> quickActions;
  final List<ContentCard> contentCards;
  final List<PopularService> popularServices;

  const HomeState({
    this.quickActions = const [],
    this.contentCards = const [],
    this.popularServices = const [],
  });

  @override
  List<Object?> get props => [quickActions, contentCards, popularServices];
}