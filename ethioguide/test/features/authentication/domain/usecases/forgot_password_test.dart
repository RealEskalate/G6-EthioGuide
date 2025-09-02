import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';
import 'package:ethioguide/features/authentication/domain/usecases/forgot_password.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';

// We must generate mocks in the file that uses them.
@GenerateMocks([AuthRepository])
import 'forgot_password_test.mocks.dart';

void main() {
  late ForgotPassword usecase;
  late MockAuthRepository mockAuthRepository;

  setUp(() {
    mockAuthRepository = MockAuthRepository();
    usecase = ForgotPassword(mockAuthRepository);
  });

  const tEmail = 'test@test.com';

  test('should call forgotPassword on the repository with the correct email', () async {
    // Arrange
    // THE FIX: Use `thenAnswer` with an explicit return type.
    when(mockAuthRepository.forgotPassword(any))
        .thenAnswer((_) async => const Right<Failure, void>(null));
        
    // Act
    final result = await usecase(tEmail);
    
    // Assert
    expect(result, const Right(null));
    verify(mockAuthRepository.forgotPassword(tEmail));
    verifyNoMoreInteractions(mockAuthRepository);
  });
}