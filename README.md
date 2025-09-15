# EthioGuide - Government Service Navigator

EthioGuide is a modern, AI-powered mobile application designed to simplify Ethiopia's complex government procedures for its citizens. Built with Flutter, it provides a clean, user-friendly interface for users to get step-by-step guidance on administrative tasks, making bureaucracy accessible and empowering Ethiopians to handle official matters with confidence.

![App Screenshot](./ethioguide/assets/images/light_logo.jpg) 
*(Note: Replace the placeholder with a real screenshot of your app's home or auth screen)*

---

## ✨ Features

- **User Authentication:** Secure user registration and login with email/password, including a complete password reset flow and Google Sign-In integration.
- **Onboarding:** A smooth, welcoming onboarding experience for first-time users.
- **Home Dashboard:** A central hub to access all the app's main features like starting a new procedure, accessing your workspace, and getting AI assistance.
- **Profile Management:** A dedicated profile screen where users can view their information and manage their account, including updating their password and an inline "edit profile" feature.
- **AI Legal Assistant:** (In Progress) A feature to provide instant, AI-powered guidance on government processes.
- **Procedure & Workspace Management:** (In Progress) The core functionality for browsing government procedures and managing personal applications.

---

## 🏗 Architecture

This project is built using a feature-driven **Clean Architecture** approach, ensuring a clear separation of concerns between layers. This makes the codebase scalable, maintainable, and highly testable.

- **Domain Layer:** Contains the core business logic, entities, and use case contracts. It has no dependencies on any other layer.
- **Data Layer:** Implements the repositories defined in the Domain Layer. It is responsible for orchestrating data from various sources (remote API, local cache, device SDKs).
- **Presentation Layer:** The UI of the application, built with Flutter. State management is handled reactively using the **BLoC** pattern, separating UI from business logic.

---

## 🛠 Technology Stack & Key Packages

- **Framework:** [Flutter](https://flutter.dev/)
- **State Management:** [flutter_bloc](https://pub.dev/packages/flutter_bloc)
- **Architecture:** Clean Architecture with TDD (Test-Driven Development)
- **Dependency Injection:** [get_it](https://pub.dev/packages/get_it)
- **Routing:** [go_router](https://pub.dev/packages/go_router)
- **Networking:** [Dio](https://pub.dev/packages/dio) (with a custom `AuthInterceptor` for token management)
- **Authentication:**
  - [google_sign_in](https://pub.dev/packages/google_sign_in) (for Social Login)
  - [flutter_secure_storage](https://pub.dev/packages/flutter_secure_storage) (for secure token caching)
- **Testing:**
  - [flutter_test](https://api.flutter.dev/flutter/flutter_test/flutter_test-library.html) (for Unit and Widget Testing)
  - [bloc_test](https://pub.dev/packages/bloc_test)
  - [mockito](https://pub.dev/packages/mockito) (for mock generation)

---