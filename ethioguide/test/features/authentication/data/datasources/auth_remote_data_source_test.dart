 import 'package:dio/dio.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/features/authentication/data/datasources/auth_remote_data_source.dart';
import 'package:ethioguide/features/authentication/data/models/tokens_model.dart';
import 'package:ethioguide/features/authentication/data/models/user_model.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

// Because we are mocking Dio, which comes from an external package,
// we often need to generate a separate mocks file for it.
@GenerateMocks([Dio])
import 'auth_remote_data_source_test.mocks.dart';

void main() {
  late AuthRemoteDataSourceImpl dataSource;
  late MockDio mockDio;

  setUp(() {
    mockDio = MockDio();
    dataSource = AuthRemoteDataSourceImpl(dio: mockDio);
  });

  const tUserModel = UserModel(id: '123', email: 'test@test.com', name: 'Lidiya Test', username: 'lidiyatest');
  const tTokensModel = TokensModel(accessToken: 'mock_access_token_12345', refreshToken: 'mock_refresh_token_67890');

  group('login', () {
    test('should return UserModel and TokensModel on successful mock login', () async {
      // Act
      final result = await dataSource.login('test@test.com', 'password');
      // Assert
      expect(result, equals((tUserModel, tTokensModel)));
    });

    test('should throw a ServerException on failed mock login', () async {
      // Act
      final call = dataSource.login;
      // Assert
      expect(() => call('wrong@test.com', 'wrongpassword'),
          throwsA(isA<ServerException>()));
    });
  });

  group('register', () {
    test('should complete successfully on mock registration', () async {
      // Act
      final call = dataSource.register(
          username: 'test', email: 'test@test.com', password: 'password', name: 'name');
      // Assert
      await expectLater(call, completes);
    });

    test('should throw a ServerException on mock registration with existing email', () async {
      // Act
      final call = dataSource.register;
      // Assert
      expect(() => call(username: 'test', email: 'exists@test.com', password: 'password', name: 'name'),
          throwsA(isA<ServerException>()));
    });
  });
}