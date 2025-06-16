# Lemonade DDD

## Setup Instructions

1. **Install dependencies**:
   ```bash
   make setup
   ```

2. **Run the application**:
   ```bash
   make run
   ```

3. **Run with hot reload (development)**:
   ```bash
   make dev
   ```

4. **Run tests**:
   ```bash
   make test
   ```

5. **Test gRPC endpoints**:
   ```bash
   # In another terminal
   make grpc-test
   ```

## Additional Setup for Proto Development

1. **Install buf CLI**:
   ```bash
   # macOS
   brew install bufbuild/buf/buf

   # or using Go
   go install github.com/bufbuild/buf/cmd/buf@latest
   ```

2. **Install grpcurl for testing**:
   ```bash
   # macOS
   brew install grpcurl

   # or using Go
   go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
   ```
