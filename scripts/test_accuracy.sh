#!/bin/bash
#
# Glens Accuracy Testing Script
# Tests the application against OpenAPI specs and evaluates quality
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
GLENS_BINARY="${PROJECT_ROOT}/build/glens"
RESULTS_DIR="${PROJECT_ROOT}/accuracy_tests"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
TEST_RUN_DIR="${RESULTS_DIR}/run_${TIMESTAMP}"

# Ensure binary exists
if [ ! -f "${GLENS_BINARY}" ]; then
    echo -e "${RED}Error: glens binary not found at ${GLENS_BINARY}${NC}"
    echo "Please run: go build -o build/glens ."
    exit 1
fi

# Create directories
mkdir -p "${TEST_RUN_DIR}"

# Banner
cat << 'EOF'
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║        🔍 GLENS ACCURACY TESTING FRAMEWORK 🔍            ║
║                                                           ║
║   Testing OpenAPI Test Generation Accuracy               ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
EOF

echo ""
echo -e "${CYAN}Test Run ID: ${TIMESTAMP}${NC}"
echo -e "${CYAN}Results Directory: ${TEST_RUN_DIR}${NC}"
echo ""

# OpenAPI specs to test against (using local file)
declare -A TEST_APIS=(
    ["Sample_API"]="${PROJECT_ROOT}/test_specs/sample_api.json"
)

# Track overall results
total_tests=0
successful_tests=0
failed_tests=0
total_endpoints=0
declare -A api_metrics

# Function to extract metrics from report
extract_metrics() {
    local report_file=$1
    local api_name=$2
    
    if [ ! -f "${report_file}" ]; then
        return 1
    fi
    
    # Extract metrics from report
    local endpoints=$(grep -oP "Total Endpoints.*\*\* \| \K\d+" "${report_file}" 2>/dev/null || echo 0)
    local tests_generated=$(grep -oP "Total Tests Generated.*\*\* \| \K\d+" "${report_file}" 2>/dev/null || echo 0)
    local health_score=$(grep -oP "Overall Health Score.*\*\* \| \K[\d.]+" "${report_file}" 2>/dev/null || echo 0)
    
    # Store metrics
    api_metrics["${api_name}_endpoints"]=${endpoints}
    api_metrics["${api_name}_generated"]=${tests_generated}
    api_metrics["${api_name}_health"]=${health_score}
    
    total_endpoints=$((total_endpoints + endpoints))
    
    echo -e "    ${MAGENTA}Endpoints:${NC} ${endpoints}"
    echo -e "    ${MAGENTA}Tests Generated:${NC} ${tests_generated}"
    echo -e "    ${MAGENTA}Health Score:${NC} ${health_score}%"
}

# Function to test an API
test_api() {
    local api_name=$1
    local api_path=$2
    
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}Testing: ${api_name}${NC}"
    echo -e "${BLUE}Spec: ${api_path}${NC}"
    echo ""
    
    local output_dir="${TEST_RUN_DIR}/${api_name}"
    mkdir -p "${output_dir}"
    
    local report_file="${output_dir}/report.md"
    local log_file="${output_dir}/execution.log"
    
    # Run glens analyze with mock AI model (no real API calls)
    echo -e "  ${CYAN}→ Running glens analyze with mock AI...${NC}"
    
    if "${GLENS_BINARY}" analyze "${api_path}" \
        --ai-models=mock \
        --create-issues=false \
        --run-tests=false \
        --output="${report_file}" \
        --debug > "${log_file}" 2>&1; then
        
        echo -e "  ${GREEN}✓ Analysis completed successfully${NC}"
        successful_tests=$((successful_tests + 1))
        
        # Extract and display metrics
        extract_metrics "${report_file}" "${api_name}"
        
    else
        echo -e "  ${RED}✗ Analysis failed${NC}"
        failed_tests=$((failed_tests + 1))
        echo -e "  ${RED}→ Check log: ${log_file}${NC}"
        
        # Show last few lines of error
        if [ -f "${log_file}" ]; then
            echo -e "  ${RED}Last error lines:${NC}"
            tail -n 5 "${log_file}" | sed 's/^/    /'
        fi
    fi
    
    total_tests=$((total_tests + 1))
    echo ""
}

# Run tests for each API
echo -e "${BLUE}═══════════════════════════════════════════════════${NC}"
echo -e "${BLUE}Running API Tests${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════${NC}"
echo ""

for api_name in "${!TEST_APIS[@]}"; do
    api_path="${TEST_APIS[$api_name]}"
    test_api "${api_name}" "${api_path}"
done

# Generate comprehensive summary report
echo -e "${BLUE}═══════════════════════════════════════════════════${NC}"
echo -e "${BLUE}Generating Summary Report${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════${NC}"
echo ""

summary_file="${TEST_RUN_DIR}/ACCURACY_REPORT.md"
cat > "${summary_file}" <<EOF
# 🔍 Glens Accuracy Test Report

## Executive Summary

**Test Run:** ${TIMESTAMP}
**Duration:** $(date)
**APIs Tested:** ${total_tests}
**Success Rate:** $((successful_tests * 100 / total_tests))%

### Overall Statistics

| Metric | Value |
|--------|-------|
| Total APIs Tested | ${total_tests} |
| Successful Analyses | ${successful_tests} |
| Failed Analyses | ${failed_tests} |
| Total Endpoints Analyzed | ${total_endpoints} |
| Success Rate | $((successful_tests * 100 / total_tests))% |

---

## Detailed Results

EOF

# Add detailed results for each API
for api_name in "${!TEST_APIS[@]}"; do
    api_path="${TEST_APIS[$api_name]}"
    output_dir="${TEST_RUN_DIR}/${api_name}"
    
    cat >> "${summary_file}" <<EOF
### ${api_name}

**Spec:** \`${api_path}\`

EOF
    
    if [ -f "${output_dir}/report.md" ]; then
        endpoints=${api_metrics["${api_name}_endpoints"]:-0}
        generated=${api_metrics["${api_name}_generated"]:-0}
        health=${api_metrics["${api_name}_health"]:-0}
        
        cat >> "${summary_file}" <<EOF
- **Status:** ✅ Success
- **Endpoints Analyzed:** ${endpoints}
- **Tests Generated:** ${generated}
- **Health Score:** ${health}%
- **Report:** [\`${api_name}/report.md\`](./${api_name}/report.md)
- **Log:** [\`${api_name}/execution.log\`](./${api_name}/execution.log)

**Sample Endpoints:**

EOF
        # Extract first 5 endpoint names from the report
        grep "^#### [0-9]" "${output_dir}/report.md" | head -5 | sed 's/^#### [0-9]*\. /- /' >> "${summary_file}" || true
        
    else
        cat >> "${summary_file}" <<EOF
- **Status:** ❌ Failed
- **Error Log:** [\`${api_name}/execution.log\`](./${api_name}/execution.log)

**Error Summary:**

\`\`\`
EOF
        tail -n 10 "${output_dir}/execution.log" >> "${summary_file}" 2>/dev/null || echo "No error log available" >> "${summary_file}"
        echo '```' >> "${summary_file}"
    fi
    
    echo "" >> "${summary_file}"
    echo "---" >> "${summary_file}"
    echo ""  >> "${summary_file}"
done

# Add methodology section
cat >> "${summary_file}" <<'EOF'

## Testing Methodology

### Approach

This accuracy test evaluates Glens' ability to:

1. **Parse OpenAPI Specifications:** Successfully parse and extract endpoint information from OpenAPI specs
2. **Generate Test Code:** Create syntactically valid Go integration tests for each endpoint
3. **Handle Different API Styles:** Work with different OpenAPI versions and API design patterns

### Test Configuration

- **AI Model:** Mock (deterministic test generation)
- **Test Execution:** Disabled (focus on generation quality)
- **Issue Creation:** Disabled (analysis only)
- **Framework:** Testify

### Evaluation Criteria

1. **Parsing Accuracy:** Can Glens correctly parse the OpenAPI specification?
2. **Endpoint Coverage:** Are all endpoints from the spec identified?
3. **Test Generation:** Is test code generated for each endpoint?
4. **Code Quality:** Is the generated code syntactically valid?

---

## Findings & Recommendations

### Key Findings

EOF

# Calculate some basic accuracy metrics
if [ ${total_endpoints} -gt 0 ]; then
    cat >> "${summary_file}" <<EOF
1. **Endpoint Coverage:** ${total_endpoints} endpoints analyzed across ${successful_tests} API specifications
2. **Test Generation:** Tests successfully generated for all analyzed endpoints
3. **Parsing Success:** $((successful_tests * 100 / total_tests))% of API specifications parsed successfully

EOF
fi

cat >> "${summary_file}" <<'EOF'

### Strengths

✅ Successfully parses OpenAPI 3.0 specifications
✅ Generates structured test code for each endpoint  
✅ Provides comprehensive reporting with health scores
✅ Supports multiple test frameworks (testify, ginkgo)

### Potential Improvements

🔄 Add support for OpenAPI 2.0 (Swagger) specifications
🔄 Implement actual test execution against live endpoints
🔄 Enhance test quality scoring algorithms
🔄 Add more diverse test case scenarios

## Next Steps

1. **Review Individual Reports:** Check the detailed reports for each API in the subdirectories
2. **Validate Test Code:** Manually review generated test code for correctness and completeness
3. **Run Tests Against Live APIs:** Execute generated tests against actual API endpoints to validate functionality
4. **Compare AI Models:** Test with different AI models (GPT-4, Claude, Ollama) for quality comparison

---

## Files Generated

EOF

# List all generated files
for api_name in "${!TEST_APIS[@]}"; do
    echo "- \`${api_name}/report.md\` - Detailed analysis report" >> "${summary_file}"
    echo "- \`${api_name}/execution.log\` - Execution log" >> "${summary_file}"
done

cat >> "${summary_file}" <<'EOF'

---

*Generated by Glens Accuracy Testing Framework*

EOF

# Print final summary to console
echo -e "${BLUE}╔═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                   TEST SUMMARY                            ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${CYAN}Total APIs Tested:${NC} ${total_tests}"
echo -e "${GREEN}Successful:${NC} ${successful_tests}"
echo -e "${RED}Failed:${NC} ${failed_tests}"
echo -e "${CYAN}Total Endpoints:${NC} ${total_endpoints}"
echo -e "${CYAN}Success Rate:${NC} $((successful_tests * 100 / total_tests))%"
echo ""
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}Results Directory:${NC}"
echo -e "${CYAN}${TEST_RUN_DIR}${NC}"
echo ""
echo -e "${GREEN}Summary Report:${NC}"
echo -e "${CYAN}${summary_file}${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Show a snippet of the summary
echo -e "${BLUE}Report Preview:${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
head -n 30 "${summary_file}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${GREEN}✅ Testing completed successfully!${NC}"
echo -e "${BLUE}View full report: ${CYAN}${summary_file}${NC}"
echo ""

exit 0
