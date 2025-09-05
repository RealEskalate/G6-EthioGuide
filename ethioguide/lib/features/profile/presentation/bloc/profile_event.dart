import 'package:equatable/equatable.dart';

abstract class ProfileEvent extends Equatable {
  const ProfileEvent();
  @override
  List<Object> get props => [];
}

// Event to tell the BLoC to fetch the user's profile data.
class FetchProfileData extends ProfileEvent {}

// Event to tell the BLoC that the user has tapped the logout button.
class LogoutTapped extends ProfileEvent {}