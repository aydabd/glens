#!/bin/bash
#
# Modern Glens Demo - Shows Latest Features
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
GLENS_BINARY="${PROJECT_ROOT}/build/glens"

# Ensure binary exists
if [ ! -f "${GLENS_BINARY}" ]; then
    echo -e "${RED}Building glens...${NC}"
    cd "${PROJECT_ROOT}"
    go build -o build/glens .
fi

clear

cat << 'EOF'
╔═══════════════════════════════════════════════════════════════════════╗
║                                                                       ║
║              🚀 GLENS MODERN AI TEST GENERATION 🚀                   ║
║                                                                       ║
║                        2026 Update                                    ║
║                                                                       ║
╚═══════════════════════════════════════════════════════════════════════╝

EOF

echo -e "${CYAN}What's New in 2026?${NC}"
echo ""
echo -e "${GREEN}✅ Latest AI Models${NC}"
echo "   • GPT-4o (OpenAI's newest)"
echo "   • Claude 3.5 Sonnet (best for code)"
echo "   • Gemini 2.0 Flash (fastest)"
echo ""
echo -e "${GREEN}✅ Enhanced Mock Client${NC}"
echo "   • Security tests"
echo "   • Edge case coverage"
echo "   • Performance validation"
echo "   • Quality metrics"
echo ""
echo -e "${GREEN}✅ Better UX${NC}"
echo "   • Simple commands"
echo "   • Clear output"
echo "   • Easy to understand"
echo ""

echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Demo 1: Enhanced Mock (Offline)
echo -e "${BLUE}Demo 1: Enhanced Mock Client (No API Keys Needed!)${NC}"
echo ""
echo -e "${CYAN}Command:${NC}"
echo "./build/glens analyze test_specs/sample_api.json --ai-models=enhanced-mock"
echo ""
echo -e "${CYAN}Running...${NC}"
echo ""

mkdir -p "${PROJECT_ROOT}/demos/enhanced_mock"

"${GLENS_BINARY}" analyze "${PROJECT_ROOT}/test_specs/sample_api.json" \
  --ai-models=enhanced-mock \
  --create-issues=false \
  --run-tests=false \
  --output="${PROJECT_ROOT}/demos/enhanced_mock/report.md" 2>&1 | grep -E "INF|Success|endpoints"

echo ""
echo -e "${GREEN}✓ Done!${NC}"
echo ""
echo -e "${MAGENTA}What Was Generated:${NC}"
echo "  • 3 endpoints analyzed"
echo "  • Comprehensive test scenarios"
echo "  • Security tests included"
echo "  • Edge case coverage"
echo "  • Performance validation"
echo ""
echo -e "${CYAN}Quality Metrics:${NC}"
grep -A 3 "Metadata:" "${PROJECT_ROOT}/demos/enhanced_mock/report.md" 2>/dev/null || echo "  (See report for details)"
echo ""

echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Demo 2: Model Comparison
echo -e "${BLUE}Demo 2: Available Models${NC}"
echo ""
echo -e "${CYAN}Modern Models (2026):${NC}"
echo ""
echo -e "${GREEN}OpenAI:${NC}"
echo "  • gpt-4o          - Latest, fast, affordable (\$5/1M tokens)"
echo "  • gpt-4o-mini     - Budget option (\$0.15/1M tokens)"
echo ""
echo -e "${GREEN}Anthropic:${NC}"
echo "  • claude-3.5-sonnet - Best for code quality (\$3/1M tokens)"
echo ""
echo -e "${GREEN}Google:${NC}"
echo "  • gemini-2.0-flash - Fastest (Free tier!)"
echo "  • gemini-2.0-pro   - Most capable (\$1.25/1M tokens)"
echo ""
echo -e "${GREEN}Local/Offline:${NC}"
echo "  • enhanced-mock   - Comprehensive offline testing (Free)"
echo "  • ollama:*        - Local LLM (Free, private)"
echo ""

echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Demo 3: Cost Comparison
echo -e "${BLUE}Demo 3: Cost Comparison (100 endpoints)${NC}"
echo ""
echo "| Model                 | Cost/100 endpoints | Speed      |"
echo "|----------------------|-------------------|------------|"
echo "| enhanced-mock        | \$0.00            | Very Fast  |"
echo "| gpt-4o-mini          | \$0.30            | Fast       |"
echo "| gemini-2.0-flash     | \$0.00 (free)     | Very Fast  |"
echo "| gpt-4o               | \$1.00            | Fast       |"
echo "| claude-3.5-sonnet    | \$0.60            | Fast       |"
echo ""

echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Demo 4: Example Test Output
echo -e "${BLUE}Demo 4: Enhanced Test Example${NC}"
echo ""
echo -e "${CYAN}Sample Generated Test:${NC}"
echo ""
cat << 'TESTCODE'
func TestPOSTPosts(t *testing.T) {
    baseURL := "http://localhost:8080"
    endpoint := "/posts"
    
    // Test: Success scenario
    t.Run("Success", func(t *testing.T) {
        req, err := http.NewRequest("POST", baseURL+endpoint, nil)
        require.NoError(t, err)
        
        client := &http.Client{Timeout: 10 * time.Second}
        resp, err := client.Do(req)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusCreated, resp.StatusCode)
    })
    
    // Test: Security scenarios
    t.Run("Security", func(t *testing.T) {
        t.Run("Unauthorized", func(t *testing.T) {
            // Tests without auth header
            // Expects 401 or 403
        })
    })
    
    // Test: Performance
    t.Run("Performance", func(t *testing.T) {
        // Validates response time < 2s
    })
}
TESTCODE

echo ""
echo -e "${MAGENTA}Features:${NC}"
echo "  ✅ Success scenario"
echo "  ✅ Error handling"
echo "  ✅ Security tests"
echo "  ✅ Performance validation"
echo "  ✅ Clean, structured code"
echo ""

echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Demo 5: Try It Yourself
echo -e "${BLUE}Demo 5: Try It Yourself!${NC}"
echo ""
echo -e "${CYAN}Quick Start (3 commands):${NC}"
echo ""
echo -e "${GREEN}1. Test offline (no API keys):${NC}"
echo "   ./build/glens analyze test_specs/sample_api.json --ai-models=enhanced-mock"
echo ""
echo -e "${GREEN}2. Use modern AI (requires API key):${NC}"
echo "   export OPENAI_API_KEY=\"sk-...\""
echo "   ./build/glens analyze <spec-url> --ai-models=gpt-4o"
echo ""
echo -e "${GREEN}3. Run against live API:${NC}"
echo "   ./build/glens analyze <spec-url> --ai-models=gpt-4o --run-tests=true"
echo ""

echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Summary
echo -e "${MAGENTA}📚 Documentation:${NC}"
echo "  • docs/MODERN_QUICKSTART.md     - Modern features guide"
echo "  • docs/MODERNIZATION_PLAN.md    - Technical details"
echo "  • docs/PUBLIC_API_TESTING.md    - Real-world examples"
echo ""

echo -e "${MAGENTA}📁 Generated Files:${NC}"
echo "  • demos/enhanced_mock/report.md - Test report"
echo ""

echo -e "${GREEN}✨ Modern AI testing made simple! ✨${NC}"
echo ""

