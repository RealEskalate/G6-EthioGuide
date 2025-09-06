import 'package:equatable/equatable.dart';

abstract class HomeEvent extends Equatable {
  const HomeEvent();
  @override
  List<Object> get props => [];
}

// The only event we need is one to load the data when the screen opens.
class LoadHomeData extends HomeEvent {}