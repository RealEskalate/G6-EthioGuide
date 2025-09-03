import 'package:bloc_test/bloc_test.dart';
import 'package:dartz/dartz.dart';
import 'package:ethioguide/core/error/failures.dart';
import 'package:ethioguide/features/authentication/domain/entities/user.dart'; // IMPORTANT: Import User
import 'package:ethioguide/features/authentication/domain/usecases/forgot_password.dart';
import 'package:ethioguide/features/authentication/domain/usecases/login_user.dart';
import 'package:ethioguide/features/authentication/domain/usecases/register_user.dart';
import 'package:ethioguide/features/authentication/domain/usecases/reset_password.dart';
import 'package:ethioguide/features/authentication/domain/usecases/sign_in_with_google.dart';
import 'package:ethioguide/features/authentication/presentation/bloc/auth_bloc.dart';
import 'package:ethioguide/features/authentication/presentation/bloc/auth_event.dart';
import 'package:ethioguide/features/authentication/presentation/bloc/auth_state.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

@GenerateMocks([LoginUser, RegisterUser, ForgotPassword, ResetPassword, SignInWithGoogle])
import 'auth_bloc_test.mocks.dart';

void main() {
  late AuthBloc authBloc;
  late MockLoginUser mockLoginUser;
  late MockRegisterUser mockRegisterUser;
   late MockForgotPassword mockForgotPassword; 
  late MockResetPassword mockResetPassword;
   late MockSignInWithGoogle mockSignInWithGoogle;

  // CORRECTED: Create a mock user object to use in tests
  const tUser = User(id: '1', name: 'Test User', email: 'test@test.com');

  setUp(() {
    mockLoginUser = MockLoginUser();
    mockRegisterUser = MockRegisterUser();
    mockForgotPassword = MockForgotPassword(); 
    mockResetPassword = MockResetPassword();
     mockSignInWithGoogle = MockSignInWithGoogle();  
    authBloc = AuthBloc(loginUser: mockLoginUser, registerUser: mockRegisterUser, forgotPassword: mockForgotPassword, 
      resetPassword: mockResetPassword, signInWithGoogle: mockSignInWithGoogle, );
  });

  test('initial state is correct', () {
    expect(authBloc.state, const AuthState());
  });

  group('AuthViewSwitched', () {
    blocTest<AuthBloc, AuthState>(
      'emits state with signUp view when switched to signUp',
      build: () => authBloc,
      act: (bloc) => bloc.add(const AuthViewSwitched(AuthView.signUp)),
      expect: () => [const AuthState(authView: AuthView.signUp)],
    );
  });

  group('PasswordVisibilityToggled', () {
    blocTest<AuthBloc, AuthState>(
      'emits state with isPasswordVisible toggled to true',
      build: () => authBloc,
      act: (bloc) => bloc.add(PasswordVisibilityToggled()),
      expect: () => [const AuthState(isPasswordVisible: true)],
    );
  });

  group('LoginSubmitted', () {
    blocTest<AuthBloc, AuthState>(
      'emits [loading, success] when login is successful',
      build: () {
        // CORRECTED: Return Right(tUser) instead of Right(null)
        when(mockLoginUser(any)).thenAnswer((_) async => const Right(tUser));
        return authBloc;
      },
      act: (bloc) => bloc.add(const LoginSubmitted(identifier: 'test', password: '123')),
      expect: () => [
        const AuthState(status: AuthStatus.loading),
        const AuthState(status: AuthStatus.success),
      ],
    );

    blocTest<AuthBloc, AuthState>(
      'emits [loading, failure] when login fails',
      build: () {
        when(mockLoginUser(any)).thenAnswer((_) async => Left(ServerFailure(message: 'Invalid credentials')));
        return authBloc;
      },
      act: (bloc) => bloc.add(const LoginSubmitted(identifier: 'test', password: '123')),
      expect: () => [
        const AuthState(status: AuthStatus.loading),
        const AuthState(status: AuthStatus.failure, errorMessage: 'Invalid credentials'),
      ],
    );
  });

  // You can also add a similar group for SignUpSubmitted
  group('SignUpSubmitted', () {
    blocTest<AuthBloc, AuthState>(
      'emits [loading, success] when signup is successful',
      build: () {
        // The register use case correctly returns Right(null) so this is fine
        when(mockRegisterUser(any)).thenAnswer((_) async => const Right(null));
        return authBloc;
      },
      act: (bloc) => bloc.add(const SignUpSubmitted(name: 't', username: 't', email: 't', password: 't')),
      expect: () => [
        const AuthState(status: AuthStatus.loading),
        const AuthState(status: AuthStatus.success),
      ],
    );
  }
   
  );
   group('ForgotPasswordSubmitted', () {
    const tEmail = 'test@test.com';

    blocTest<AuthBloc, AuthState>(
      'emits [loading, resetLinkSent] when forgot password is successful',
      build: () {
        when(mockForgotPassword(any)).thenAnswer((_) async => const Right(null));
        return authBloc;
      },
      act: (bloc) => bloc.add(const ForgotPasswordSubmitted(email: tEmail)),
      expect: () => [
        const AuthState(status: AuthStatus.loading),
        const AuthState(status: AuthStatus.resetLinkSent),
      ],
      verify: (_) => verify(mockForgotPassword(tEmail)).called(1),
    );

    blocTest<AuthBloc, AuthState>(
      'emits [loading, failure] when forgot password fails',
      build: () {
        when(mockForgotPassword(any)).thenAnswer((_) async => Left(ServerFailure(message: 'Email not found')));
        return authBloc;
      },
      act: (bloc) => bloc.add(const ForgotPasswordSubmitted(email: tEmail)),
      expect: () => [
        const AuthState(status: AuthStatus.loading),
        const AuthState(status: AuthStatus.failure, errorMessage: 'Email not found'),
      ],
    );
  });

  group('GoogleSignInSubmitted', () {
    blocTest<AuthBloc, AuthState>(
      'emits [loading, success] when Google sign-in is successful',
      build: () {
        when(mockSignInWithGoogle()).thenAnswer((_) async => const Right(tUser));
        return authBloc;
      },
      act: (bloc) => bloc.add(GoogleSignInSubmitted()),
      expect: () => [
        const AuthState(status: AuthStatus.loading),
        const AuthState(status: AuthStatus.success),
      ],
      verify: (_) => verify(mockSignInWithGoogle()).called(1),
    );

    blocTest<AuthBloc, AuthState>(
      'emits [loading, failure] when Google sign-in fails',
      build: () {
        when(mockSignInWithGoogle()).thenAnswer((_) async => Left(CachedFailure(message: 'Cancelled by user')));
        return authBloc;
      },
      act: (bloc) => bloc.add(GoogleSignInSubmitted()),
      expect: () => [
        const AuthState(status: AuthStatus.loading),
        const AuthState(status: AuthStatus.failure, errorMessage: 'Cancelled by user'),
      ],
    );
  });
}
