# C4 Model Documentation

## Overview

This directory contains C4 model documentation for the Go Backend API with Clean Architecture. The C4 model provides different levels of abstraction to understand the system architecture:

## C4 Model Levels

### 1. [Context Diagram](01-context.md) 
- System context and external interactions
- Shows users, external systems, and high-level purpose

### 2. [Container Diagram](02-container.md)
- High-level technology choices  
- Shows applications, data stores, and major technology boundaries

### 3. [Component Diagram](03-component.md) 
- Components within the Go API container
- Shows Clean Architecture layers and their interactions
- Internal structure และ component responsibilities

### 4. [Code Diagram](04-code.md)
- Key classes and their relationships
- Focus on domain entities and critical interfaces
- Implementation details และ code structure

## Architecture Summary

This Go backend follows **Clean Architecture** principles with clear separation of concerns:

- **Domain Layer**: Core business logic and entities
- **Use Case Layer**: Application workflows and business rules  
- **Interface Layer**: HTTP handlers, DTOs, and external communication
- **Infrastructure Layer**: Database, authentication, and external services
- **Application Layer**: Dependency injection and configuration

## Technologies Used

- **Runtime**: Go 1.24.3
- **Web Framework**: Chi Router
- **Database**: SQLite with SQL driver
- **Authentication**: JWT tokens
- **API Documentation**: Swagger/OpenAPI
- **Testing**: Go built-in testing with comprehensive coverage

## Quick Navigation

- 🌍 [System Context](01-context.md) - Who uses the system?
- 📦 [Containers](02-container.md) - What does the system contain?  
- 🔧 [Components](03-component.md) - How is the API structured?
- 💻 [Code](04-code.md) - What are the key classes?

## Additional Documentation

- 🗄️ [Database Schema](../database.md) - ER diagram และ database structure
- 🏗️ [Architecture Overview](../architecture.md) - Complete architecture documentation

## Diagrams Notation

These diagrams use C4 model notation:
- **Person**: Users of the system
- **System**: The software system boundary
- **Container**: Applications, databases, file systems, etc.
- **Component**: Logical components within containers
- **Relationship**: Dependencies and data flow
