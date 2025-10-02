#!/bin/bash

# Test runner script that reports passing and failing test counts
# while maintaining proper exit codes for CI

echo "Running Go tests with summary reporting..."

# Run tests and capture both stdout and stderr
TEST_OUTPUT=$(go test ./... -v 2>&1)
TEST_EXIT_CODE=$?

# Count passing and failing tests
PASS_COUNT=$(echo "$TEST_OUTPUT" | grep -c "^--- PASS" || true)
FAIL_COUNT=$(echo "$TEST_OUTPUT" | grep -c "^--- FAIL" || true)
SKIP_COUNT=$(echo "$TEST_OUTPUT" | grep -c "^--- SKIP" || true)

# Count package-level results
PACKAGE_PASS=$(echo "$TEST_OUTPUT" | grep -c "^PASS$" || true)
PACKAGE_FAIL=$(echo "$TEST_OUTPUT" | grep -c "^FAIL$" || true)

# Display the full test output
echo "$TEST_OUTPUT"

echo ""
echo "=========================================="
echo "TEST SUMMARY"
echo "=========================================="
echo "Individual test results:"
echo "  ✅ Passing tests: $PASS_COUNT"
echo "  ❌ Failing tests: $FAIL_COUNT"
echo "  ⚠️  Skipped tests: $SKIP_COUNT"
echo ""
echo "Package results:"
echo "  📦 Packages passed: $PACKAGE_PASS"
echo "  📦 Packages failed: $PACKAGE_FAIL"
echo ""

# Calculate total tests run
TOTAL_TESTS=$((PASS_COUNT + FAIL_COUNT + SKIP_COUNT))
echo "Total tests run: $TOTAL_TESTS"

# Show success rate if there are tests
if [ $TOTAL_TESTS -gt 0 ]; then
    SUCCESS_RATE=$((PASS_COUNT * 100 / TOTAL_TESTS))
    echo "Success rate: $SUCCESS_RATE%"
fi

echo "=========================================="

# Exit with the original test exit code to maintain CI failure behavior
if [ $TEST_EXIT_CODE -ne 0 ]; then
    echo "❌ CI Status: FAILURE (exit code: $TEST_EXIT_CODE)"
    echo "Note: Build will be marked as failed due to failing tests"
else
    echo "✅ CI Status: SUCCESS"
fi

exit $TEST_EXIT_CODE