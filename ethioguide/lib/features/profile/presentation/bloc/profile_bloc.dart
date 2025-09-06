import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/usecases/get_user_profile.dart';
import '../../domain/usecases/logout_user.dart';
import 'profile_event.dart';
import 'profile_state.dart';

class ProfileBloc extends Bloc<ProfileEvent, ProfileState> {
  final GetUserProfile getUserProfile;
  final LogoutUser logoutUser;

  ProfileBloc({
    required this.getUserProfile,
    required this.logoutUser,
  }) : super(const ProfileState()) {
    on<FetchProfileData>(_onFetchProfileData);
    on<LogoutTapped>(_onLogoutTapped);
  }

  Future<void> _onFetchProfileData(FetchProfileData event, Emitter<ProfileState> emit) async {
    emit(state.copyWith(status: ProfileStatus.loading));
    final result = await getUserProfile();
    result.fold(
      (failure) => emit(state.copyWith(status: ProfileStatus.failure, errorMessage: failure.message)),
      (user) => emit(state.copyWith(status: ProfileStatus.success, user: user)),
    );
  }

  Future<void> _onLogoutTapped(LogoutTapped event, Emitter<ProfileState> emit) async {
    // We don't need a loading state for logout as it's very fast.
    final result = await logoutUser();
    result.fold(
      (failure) => emit(state.copyWith(status: ProfileStatus.failure, errorMessage: failure.message)),
      (_) => emit(state.copyWith(status: ProfileStatus.loggedOut)),
    );
  }
}