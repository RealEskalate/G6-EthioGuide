import 'package:dartz/dartz.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';
import 'package:ethioguide/features/authentication/domain/entities/user.dart';
import 'package:ethioguide/features/authentication/domain/usecases/login_user.dart';

@GenerateMocks([AuthRepository])
import 'login_user_test.mocks.dart';

void main() {
  late LoginUser usecase;
  late MockAuthRepository mockAuthRepository;

  setUp(() {
    mockAuthRepository = MockAuthRepository();
    usecase = LoginUser(mockAuthRepository);
  });

  const tIdentifier = 'test@test.com';
  const tPassword = 'password';
  const tUser = User(id: '1', email: 'test@test.com', name: 'Test User');

  test('should get user from the repository on successful login', () async {
    // Arrange
    when(mockAuthRepository.login(any, any))
        .thenAnswer((_) async => const Right(tUser));

    // Act
    final result = await usecase(const LoginParams(identifier: tIdentifier, password: tPassword));

    // Assert
    expect(result, const Right(tUser));
    verify(mockAuthRepository.login(tIdentifier, tPassword));
    verifyNoMoreInteractions(mockAuthRepository);
  });
}