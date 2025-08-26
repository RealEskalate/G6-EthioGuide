# EthioGuide Project

This document outlines the folder structure and organization of the codebase

## Folder Structure
```bash
ethioguide/
├── assets/                     # Static assets like images, fonts, and styles
│   ├── images/                 # Image files (e.g., PNG, JPG, SVG)
│   ├── fonts/                  # Custom font files (e.g., TTF, OTF)
│   └── misc/                   # Other static assets (e.g., icons, videos)
├── features/                   # Feature-specific modules
│   ├── test_feature/           # Example feature module
│   │   ├── domain/             # Business logic, models, and use cases
│   │   │   ├── entities/       # Data models or entities
│   │   │   ├── repositories/   # Interfaces for data access
│   │   │   └── usecases/       # Business logic and use cases
│   │   ├── presentation/       # UI-related code (views, widgets, screens)
│   │   │   ├── widgets/        # Reusable UI components
│   │   │   ├── screens/        # Feature-specific screens or pages
│   │   │   └── bloc/           # View models or presenters
│   │   └── data/               # Data layer (API calls, local storage)
│   │       ├── datasources/    # Remote and local data sources
│   │       ├── repositories/   # Implementation of domain repositories
│   │       └── models/         # Data models for API/storage
│   └── other_feature/          # Same structure for other features
├── core/                       # Shared utilities and core functionality
│   ├── error/                  # Error handling (custom exceptions, error models)
│   ├── network/                # Network utilities (e.g., API client, interceptors)
│   ├── success/                # Success handling (e.g., success models, response wrappers)
│   ├── utils/                  # General utilities (e.g., helpers, extensions, constants)
│   └── config/                 # Configuration files (e.g., environment, app settings)
├── tests/                      # Unit and integration tests
│   ├── features/               # Tests for feature-specific code
│   └── core/                   # Tests for core utilities
├── docs/                       # Documentation files
│   └── api/                    # API documentation or specifications
├── scripts/                    # Build, deployment, or automation scripts
└── README.md                   # Project documentation
```
