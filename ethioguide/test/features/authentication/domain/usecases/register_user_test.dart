import 'package:dartz/dartz.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:ethioguide/features/authentication/domain/repositories/auth_repositoryy.dart';
import 'package:ethioguide/features/authentication/domain/usecases/register_user.dart';

// We can reuse the mock generated from the login test.
// If you keep tests for the same feature in the same folder,
// you often only need one @GenerateMocks annotation.
// However, it's good practice to have it in each file in case they are moved.
@GenerateMocks([AuthRepository])
import 'register_user_test.mocks.dart';

void main() {
  late RegisterUser usecase;
  late MockAuthRepository mockAuthRepository;

  setUp(() {
    mockAuthRepository = MockAuthRepository();
    usecase = RegisterUser(mockAuthRepository);
  });

  // Define the test parameters
  const tUsername = 'newuser';
  const tEmail = 'new@user.com';
  const tPassword = 'newpassword123';
  const tName = 'New User';

  final tRegisterParams = RegisterParams(
    username: tUsername,
    email: tEmail,
    password: tPassword,
    name: tName,
  );

  test('should call the register method on the repository with correct parameters', () async {
    // Arrange: Program the mock to return a successful result (Right(null)).
    // `anyNamed` is used for named parameters.
    when(mockAuthRepository.register(
      username: anyNamed('username'),
      email: anyNamed('email'),
      password: anyNamed('password'),
      name: anyNamed('name'),
      phone: anyNamed('phone'),
    )).thenAnswer((_) async => const Right(null));

    // Act: Execute the use case with our test parameters.
    final result = await usecase(tRegisterParams);

    // Assert: Check if the result is a success.
    expect(result, const Right(null));

    // Verify that the repository's register method was called exactly once
    // with the exact parameters we passed to the use case.
    verify(mockAuthRepository.register(
      username: tUsername,
      email: tEmail,
      password: tPassword,
      name: tName,
      phone: null, // We expect phone to be null since we didn't provide it
    ));

    // Verify that no other methods were called on the repository.
    verifyNoMoreInteractions(mockAuthRepository);
  });
}