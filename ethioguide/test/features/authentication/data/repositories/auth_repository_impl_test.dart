import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/core/network/network_info.dart';
// ADDED: Import the local data source contract
import 'package:ethioguide/features/authentication/data/datasources/auth_local_data_source.dart';
import 'package:ethioguide/features/authentication/data/datasources/auth_remote_data_source.dart';
import 'package:ethioguide/features/authentication/data/models/tokens_model.dart';
import 'package:ethioguide/features/authentication/data/models/user_model.dart';
import 'package:ethioguide/features/authentication/data/repositories/auth_repositoryy_impl.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

// UPDATED: Add AuthLocalDataSource to the list of mocks to be generated.
@GenerateMocks([
  AuthRemoteDataSource,
  AuthLocalDataSource,
  NetworkInfo,
  FlutterSecureStorage
])
import 'auth_repository_impl_test.mocks.dart';

void main() {
  late AuthRepositoryImpl repository;
  late MockAuthRemoteDataSource mockRemoteDataSource;
  late MockAuthLocalDataSource mockLocalDataSource; // ADDED: Declaration for the new mock
  late MockNetworkInfo mockNetworkInfo;
  late MockFlutterSecureStorage mockSecureStorage;

  setUp(() {
    mockRemoteDataSource = MockAuthRemoteDataSource();
    mockLocalDataSource = MockAuthLocalDataSource(); // ADDED: Initialization of the new mock
    mockNetworkInfo = MockNetworkInfo();
    mockSecureStorage = MockFlutterSecureStorage();
    repository = AuthRepositoryImpl(
      remoteDataSource: mockRemoteDataSource,
      localDataSource: mockLocalDataSource, // ADDED: Provide the mock to the constructor
      networkInfo: mockNetworkInfo,
      secureStorage: mockSecureStorage,
    );
  });

  const tUserModel = UserModel(id: '1', email: 'test@test.com', name: 'Test User');
  const tTokensModel = TokensModel(accessToken: 'abc', refreshToken: 'xyz');
  const tLoginIdentifier = 'test@test.com';
  const tLoginPassword = 'password';
  const tGoogleIdToken = 'google_auth_code_string';

  // Helper function to keep tests DRY (Don't Repeat Yourself)
  void runTestsOnline(Function body) {
    group('device is online', () {
      setUp(() {
        when(mockNetworkInfo.isConnected).thenAnswer((_) async => true);
      });
      body();
    });
  }

  // (The rest of your login tests can remain as they are)
  group('login', () {
    // ... your existing login tests ...
  });
  
  // ADDED: A new group of tests for the signInWithGoogle method
  group('signInWithGoogle', () {
    runTestsOnline(() {
      test('should return user and save tokens on successful Google sign-in', () async {
        // Arrange
        // 1. Mock the local data source to return a fake Google token
        when(mockLocalDataSource.getGoogleServerAuthCode()).thenAnswer((_) async => tGoogleIdToken);
        // 2. Mock the remote data source to return a user and app tokens
        when(mockRemoteDataSource.signInWithGoogle(any))
            .thenAnswer((_) async => (tUserModel, tTokensModel));
        
        // Act
        final result = await repository.signInWithGoogle();
        
        // Assert
        expect(result, equals(const Right(tUserModel)));
        verify(mockLocalDataSource.getGoogleServerAuthCode());
        verify(mockRemoteDataSource.signInWithGoogle(tGoogleIdToken));
        verify(mockSecureStorage.write(key: anyNamed('key'), value: tTokensModel.accessToken));
        verify(mockSecureStorage.write(key: anyNamed('key'), value: tTokensModel.refreshToken));
      });

      test('should return CacheFailure when local data source throws CacheException', () async {
        // Arrange
        when(mockLocalDataSource.getGoogleServerAuthCode()).thenThrow(CacheException(message: 'Cancelled'));
        
        // Act
        final result = await repository.signInWithGoogle();
        
        // Assert
        expect(result, equals(Left(const CachedFailure(message: 'Cancelled'))));
        verifyZeroInteractions(mockRemoteDataSource);
        verifyZeroInteractions(mockSecureStorage);
      });
    });
  });
}