import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';
import 'package:ethioguide/features/authentication/domain/usecases/reset_password.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';

// We can reuse the mock from the forgot_password_test, but it's good practice to generate it here too.
@GenerateMocks([AuthRepository])
import 'reset_password_test.mocks.dart';

void main() {
  late ResetPassword usecase;
  late MockAuthRepository mockAuthRepository;

  setUp(() {
    mockAuthRepository = MockAuthRepository();
    usecase = ResetPassword(mockAuthRepository);
  });

  // Define the test parameters
  const tParams = ResetPasswordParams(
    email: 'test@test.com',
    token: '123456',
    newPassword: 'newPassword123',
  );

  test('should call resetPassword on the repository with correct parameters', () async {
    // Arrange
    // THE FIX: Use `thenAnswer` with an explicit return type.
    when(mockAuthRepository.resetPassword(
      email: anyNamed('email'),
      token: anyNamed('token'),
      newPassword: anyNamed('newPassword'),
    )).thenAnswer((_) async => const Right<Failure, void>(null));

    // Act
    final result = await usecase(tParams);

    // Assert
    expect(result, const Right(null));
    
    // Verify that the repository's method was called with the exact parameters.
    verify(mockAuthRepository.resetPassword(
      email: tParams.email,
      token: tParams.token,
      newPassword: tParams.newPassword,
    ));
    
    verifyNoMoreInteractions(mockAuthRepository);
  });
}