# DRPS: Detailed Microservices Architecture & Technical Specification

## 1. Architectural Vision & Motivation
The DRPS (Digital Retail Platform System) project is undergoing a fundamental architectural shift from a tightly coupled Go monolith (using the Gin framework, global state, and shared packages) into a decentralized, domain-driven microservices ecosystem. 

**Motivations for Refactoring:**
*   **Eliminate Tight Coupling:** The monolithic architecture relied on a global `shared` module, intertwining the Auth, Sellers, Customers, and Admin domains. This prevented independent scaling and deployments, and caused cascading test failures.
*   **Domain Isolation & Polyglot Persistence:** Migrating to isolated domains allows precise scaling and optimized database choices (e.g., PostgreSQL for ACID compliance in Orders; MongoDB for schema flexibility in Products).
*   **Standardization:** Moving away from heavy frameworks (like Gin) toward Go's standard `net/http` (with routers like `chi`) for lower overhead, native `context` handling, and reduced framework lock-in.

## 2. Target Architecture & MVP Scope

The final architectural vision encompasses 7 distinct microservices (Auth, Shop, Product, Order, Customer, Analytics, Theme). However, to prevent scope creep, the **MVP (Minimum Viable Product)** focuses strictly on building the foundational infrastructure with 3 core domains and an API Gateway.

### 🟢 MVP Scope (In Scope)
1.  **API Gateway (BFF - Backend for Frontend):**
    *   **Tech:** Go `net/http` + `chi` router.
    *   **Role:** Single entry point for external REST/JSON HTTP requests. Handles rate limiting, authentication/JWT validation, and routes requests to internal domain services via gRPC.
2.  **Auth Service:**
    *   **Tech:** Go, gRPC, PostgreSQL.
    *   **Role:** Centralized identity management. Handles user credentials securely, session management, and issues/validates JWTs.
3.  **Product Service:**
    *   **Tech:** Go, gRPC, MongoDB.
    *   **Role:** High-read catalog and inventory management. Manages product variants and strict stock deductions.
4.  **Order Service:**
    *   **Tech:** Go, gRPC, PostgreSQL.
    *   **Role:** Transactional engine managing cart lifecycles and state machine logic (Pending -> Paid -> Shipped).

### 🔴 Deferred (Out of Scope for MVP)
*   Shop, Customer, Analytics, and Theme Services.
*   Message Brokers (NATS/Kafka/RabbitMQ) for Eventual Consistency (MVP uses synchronous gRPC).
*   Payment Gateway Integration (Mock successful payments only).
*   OAuth2 (Google/Facebook Login) and complex notification systems.
*   Kubernetes deployments (stick to `docker-compose` for local orchestrations).

---

## 3. Technical Implementation Standards

### 3.1 Clean Architecture & Code Organization
Every microservice must strictly adhere to Clean Architecture principles to decouple business rules from infrastructure implementations:

```text
service-name/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point, Dependency Injection
├── internal/
│   ├── domain/               # Core business entities, interface definitions
│   ├── usecase/              # Application business logic (orchestration)
│   ├── delivery/             # Presentation layer (HTTP handlers, gRPC servers)
│   └── infrastructure/       # Concrete DB implementations (Postgres/Mongo)
├── pkg/                      # Shared internal utilities (e.g., custom errors)
└── proto/                    # Protobuf definitions and generated code
```

**Strict Coding Guidelines:**
*   **No Global State:** Global variables (like `var DB *mongo.Client`) are strictly prohibited. All dependencies (configs, DB connections, loggers) must be injected via constructors.
*   **Dependency Inversion:** Repositories and external service clients must be defined as interfaces in the `domain` package and injected into the `usecase` layer.
*   **Separation of Concerns:** Controllers/Delivery handlers must NOT contain business logic. They are only responsible for request parsing, input validation, delegating to the `usecase`, and formatting the response.
*   **Context Propagation:** `context.Context` must be passed down through all layers to support timeouts, cancellation, and distributed tracing.

### 3.2 Inter-Service Communication Strategy
The system adopts a highly optimized hybrid communication model:

1.  **External Communication (Frontend ↔ Gateway):** RESTful JSON over HTTP.
2.  **Internal Communication (Gateway ↔ Service / Service ↔ Service):** Synchronous gRPC over HTTP/2 using Protocol Buffers (v3).

**Example Internal gRPC Contract (Product Service):**
```protobuf
syntax = "proto3";
package product.v1;

service ProductService {
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  rpc GetProductsByIDs(GetProductsByIDsRequest) returns (GetProductsByIDsResponse);
  rpc ReduceStock(ReduceStockRequest) returns (ReduceStockResponse);
}
```
*   **Backward Compatibility Strategy:** Ensure schema evolution by adding new fields and deprecating old ones; never reuse or modify existing protobuf field numbers.

### 3.3 Data Management & Consistency
*   **Database per Service:** A strict isolation policy. Services are completely forbidden from connecting to another service's database. If the Order Service needs product data, it must make a gRPC call to the Product Service.
*   **Synchronous Consistency (MVP Focus):** For the MVP, distributed transactions are managed via direct gRPC calls. E.g., the Checkout UseCase in the Order Service creates an order record, synchronously calls the Product Service `ReduceStock` RPC, and updates the order status based on the result.
*   **Future Saga Migration:** Once the MVP is stable, complex operations will migrate to the **Saga Pattern** or use **Event Sourcing** via an event bus to handle compensations seamlessly.

---

## 4. Phased Engineering Roadmap

**Phase 1: Contracts & Foundation**
1.  Draft and finalize `.proto` files for the Auth, Product, and Order services.
2.  Generate Go gRPC client/server stubs.
3.  Establish the `docker-compose.yml` infrastructure, configuring PostgreSQL, MongoDB, and networking.

**Phase 2: Authentication & Gateway Routing**
1.  Scaffold the Auth Service with Clean Architecture. Implement user registration and token generation using PostgreSQL.
2.  Scaffold the API Gateway using `net/http` and `chi`.
3.  Establish routing logic in the Gateway that maps incoming REST endpoints (e.g., `POST /api/v1/auth/login`) to the Auth Service's gRPC methods.

**Phase 3: Core E-Commerce Domain**
1.  Build the Product Service (MongoDB) focusing on catalog CRUD and the critical `ReduceStock` gRPC endpoint.
2.  Build the Order Service (PostgreSQL) focusing on the Cart and Checkout flows.
3.  Implement the synchronous Checkout orchestration: The Order Service coordinates the transaction by executing local state changes and calling the Product Service for inventory reduction.

**Phase 4: Polish & Production Readiness**
1.  Implement Custom Error Handling with HTTP status code mapping and standardized JSON error structures.
2.  Add gRPC Interceptors for structured logging (e.g., `slog` or `zap`), request tracing, and panic recovery.
3.  Implement graceful shutdown mechanisms across all services to prevent dropped connections during redeployments.
4.  Configure database connection pooling, query timeouts, and API rate limiting.
