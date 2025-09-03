# Go Project Test Writing Metaprompt

You are an AI agent tasked with analyzing a Go project and writing comprehensive tests for the `database.go` file. Follow this systematic approach:

## Phase 1: Project Documentation Analysis

### 1.1 Find and Analyze Info/Markdown Files
- Locate all `.md`, `.txt`, and documentation files in the project root and `/docs` directory
- Read and summarize:
    - `README.md` - project overview, setup instructions, key concepts
    - `PROJECT_PROGRESS.md` - coding standards and testing guidelines
    - 'telegram_organizer_schema.md' - database schema and relevant SQL queries
    - Any architecture or design documents
    - API documentation or specification files
- Extract key information about:
    - Project purpose and domain
    - Coding conventions and style guidelines
    - Testing frameworks and patterns used
    - Dependencies and external integrations
    - Error handling patterns

### 1.2 Understand Project Structure
- Map out the directory structure
- Identify main packages and their responsibilities
- Note any special build tags, configuration files, or environment setup

## Phase 2: Test Pattern Analysis

### 2.1 Examine Existing Test Files
Analyze at least 3-5 existing `*_test.go` files in the project, focusing on:

**Testing Framework and Patterns:**
- Which testing framework is used (standard `testing`, testify, ginkgo, etc.)
- Test file naming conventions
- Test function naming patterns
- Package declaration style (`package foo` vs `package foo_test`)

**Test Structure and Organization:**
- How test cases are organized (table-driven tests, subtests, etc.)
- Setup and teardown patterns
- Mock/stub usage and patterns
- Test data organization (fixtures, builders, etc.)

**Assertion and Error Handling:**
- How assertions are written
- Error testing patterns
- Expected vs actual value ordering
- Custom assertion helpers or utilities

**Coverage Patterns:**
- What types of scenarios are typically tested
- How edge cases are handled
- Integration vs unit test boundaries
- Performance/benchmark test patterns

### 2.2 Identify Common Utilities and Helpers
- Test helper functions
- Mock generators or factories
- Common test data builders
- Custom matchers or assertion helpers

## Phase 3: Target File Analysis

### 3.1 Analyze `database.go` Structure
Thoroughly examine the `database.go` file:

**Code Organization:**
- Package declaration and imports
- Exported vs unexported functions, types, and variables
- Constants and global variables
- Struct definitions and their fields

**Function Inventory:**
- List all functions with their signatures
- Identify pure functions vs functions with side effects
- Note functions that interact with external systems (files, network, databases)
- Identify error return patterns

**Dependencies and Interactions:**
- External package dependencies
- Internal package dependencies
- File system interactions
- Network calls or external API usage
- Database operations

**Business Logic Patterns:**
- Input validation logic
- Data transformation patterns
- Error handling and recovery
- Concurrency usage (goroutines, channels)

### 3.2 Identify Test Scenarios
For each function and method, identify:

**Happy Path Testing:**
- Normal input cases
- Boundary value cases
- Different input combinations

**Error Path Testing:**
- Invalid input handling
- External dependency failures
- Resource exhaustion scenarios
- Concurrent access issues

**Edge Cases:**
- Empty inputs, nil values
- Very large or very small inputs
- Malformed data handling
- State-dependent behavior

## Phase 4: Test Implementation

### 4.1 Follow Project Conventions
Based on your analysis, ensure your tests:
- Use the same testing framework and patterns as existing tests
- Follow the same naming conventions
- Use similar assertion patterns
- Match the project's error handling style
- Follow the same package organization (same package vs separate test package)

### 4.2 Test File Structure
Create `database_test.go` with:

```go
// Match existing test file headers and build tags
package [package_name] // or package [package_name]_test

import (
    // Include standard testing imports used in other test files
    // Include any testing utilities or assertion libraries used
    // Include mocking frameworks if used in other tests
)

// Include any test setup/teardown functions if common in the project
// Include any test helper functions needed
// Include any test data/fixtures following project patterns
```

### 4.3 Comprehensive Test Coverage
Write tests that cover:

**Unit Tests:**
- Each exported function with multiple test cases
- Each exported method on structs/interfaces
- Important unexported functions if they contain complex logic

**Integration Tests:**
- Functions that interact with external systems
- End-to-end workflows involving multiple functions
- Error propagation through the system

**Table-Driven Tests:**
- Use table-driven tests for functions with multiple input/output combinations
- Follow the project's table test structure and naming

**Benchmark Tests:**
- Add benchmark tests for performance-critical functions
- Follow existing benchmark patterns in the project

### 4.4 Test Quality Guidelines
Ensure your tests:
- Have clear, descriptive test names that explain what is being tested
- Include both positive and negative test cases
- Test edge cases and boundary conditions
- Are deterministic and can run in any order
- Clean up resources appropriately
- Use appropriate mocking for external dependencies
- Include helpful error messages in assertions
- Follow the AAA pattern (Arrange, Act, Assert) or similar

## Phase 5: Documentation and Finalization

### 5.1 Test Documentation
- Add comments explaining complex test scenarios
- Document any test data setup or special requirements
- Include examples of how to run the tests

### 5.2 Integration Verification
- Ensure tests follow the project's CI/CD patterns
- Verify tests work with existing build and test scripts
- Check that new tests integrate well with coverage tools

## Output Requirements

Provide:
1. **Analysis Summary**: Brief summary of project patterns and conventions discovered
2. **Complete Test File**: Full `database_test.go` implementation
3. **Test Coverage Report**: List of what is tested and any notable gaps
4. **Integration Notes**: Any special considerations for running or maintaining these tests

## Quality Checklist

Before finalizing, verify:
- [ ] Tests follow project conventions exactly
- [ ] All exported functions are tested
- [ ] Error cases are comprehensively covered
- [ ] Tests are readable and maintainable
- [ ] No hardcoded values that should be configurable
- [ ] Proper cleanup of resources
- [ ] Tests are deterministic
- [ ] Integration with existing test infrastructure
- [ ] Appropriate use of mocks and stubs
- [ ] Performance considerations addressed where relevant