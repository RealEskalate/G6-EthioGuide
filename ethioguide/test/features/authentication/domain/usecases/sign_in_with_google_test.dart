import 'package:dartz/dartz.dart';
import 'package:ethioguide/features/authentication/domain/entities/user.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';
import 'package:ethioguide/features/authentication/domain/usecases/sign_in_with_google.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';

import 'login_user_test.mocks.dart'; // We can reuse the same mock

void main() {
  late SignInWithGoogle usecase;
  late MockAuthRepository mockAuthRepository;

  setUp(() {
    mockAuthRepository = MockAuthRepository();
    usecase = SignInWithGoogle(mockAuthRepository);
  });

  const tUser = User(id: 'google_user_123', email: 'google@test.com', name: 'Google User');

  test('should get user from the repository on successful Google sign-in', () async {
    // Arrange
    when(mockAuthRepository.signInWithGoogle())
        .thenAnswer((_) async => const Right(tUser));
    // Act
    final result = await usecase();
    // Assert
    expect(result, const Right(tUser));
    verify(mockAuthRepository.signInWithGoogle());
    verifyNoMoreInteractions(mockAuthRepository);
  });
}