import 'package:ethioguide/core/data/repositories/auth_repository_impl.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

@GenerateMocks([FlutterSecureStorage])
import 'auth_repository_impl_test.mocks.dart';

void main() {
  late MockFlutterSecureStorage mockSecureStorage;
  late CoreAuthRepositoryImpl authRepository;

  setUp(() {
    mockSecureStorage = MockFlutterSecureStorage();
    authRepository = CoreAuthRepositoryImpl(secureStorage: mockSecureStorage);
  });

  group('AuthRepositoryImpl', () {
    test(
      'should save both access and refresh tokens when saveTokens is called',
      () async {
        // Arrange
        const accessToken = 'acess123';
        const refreshToken = 'refresh456';

        // Act
        await authRepository.saveTokens(
          accessToken: accessToken,
          refreshToken: refreshToken,
        );

        // Assert
        verify(
          mockSecureStorage.write(key: 'accessToken', value: accessToken),
        ).called(1);

        verify(
          mockSecureStorage.write(key: 'refreshToken', value: refreshToken),
        ).called(1);
      },
    );

    test(
      'should return stored access token when getAccessToken is called',
      () async {
        // Arrange
        const accessToken = 'access123';
        when(
          mockSecureStorage.read(key: 'accessToken'),
        ).thenAnswer((_) async => accessToken);

        // Act
        final result = await authRepository.getAccessToken();

        // Assert
        expect(result, accessToken);
        verify(mockSecureStorage.read(key: 'accessToken')).called(1);
      },
    );

    test(
      'should return stored refresh token when getRefreshToken is called',
      () async {
        // Arrange
        const refreshToken = 'refresh456';
        when(
          mockSecureStorage.read(key: 'refreshToken'),
        ).thenAnswer((_) async => refreshToken);

        // Act
        final result = await authRepository.getRefreshToken();

        // Assert
        expect(result, refreshToken);
        verify(mockSecureStorage.read(key: 'refreshToken')).called(1);
      },
    );

    test('should delete both tokens when clearTokens is called', () async {
      // Act
      await authRepository.clearTokens();

      // Assert
      verify(mockSecureStorage.delete(key: 'accessToken')).called(1);
      verify(mockSecureStorage.delete(key: 'refreshToken')).called(1);
    });

    test(
      'should return true if access token exists when isAuthenticated is called',
      () async {
        // Arrange
        when(
          mockSecureStorage.read(key: 'accessToken'),
        ).thenAnswer((_) async => 'access123');

        // Act
        final result = await authRepository.isAuthenticated();

        // Assert
        expect(result, true);
      },
    );

    test(
      'should return false if access token doesn\'t exists when isAuthenticated is called',
      () async {
        // Arrange
        when(
          mockSecureStorage.read(key: 'accessToken'),
        ).thenAnswer((_) async => null);

        // Act
        final result = await authRepository.isAuthenticated();

        // Assert
        expect(result, false);
      },
    );

    test(
      'should overWrite access token when updateAccessToken is called',
      () async {
        // Arrange
        const newAccessToken = 'newAccess123';

        // Act
        await authRepository.updateAccessToken(newAccessToken);

        // Assert
        verify(
          mockSecureStorage.write(key: 'accessToken', value: newAccessToken),
        ).called(1);
      },
    );
  });
}
