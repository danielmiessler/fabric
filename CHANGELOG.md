# Changelog

## v1.4.246 (2025-07-13)

### Direct commits

- Add AI-powered changelog generation with high-performance Go tool and comprehensive caching
- Implement high-performance Go changelog generator with GraphQL and SQLite-based persistent caching for incremental updates
- Add one-pass git history walking algorithm with concurrent GitHub API processing and batching
- Create comprehensive CLI with cobra framework integration and tag-based caching for unreleased content detection
- Add content hashing for change detection optimization and AI summarization using Fabric CLI integration

## v1.4.245 (2025-07-11)

### PR [#1603](https://github.com/danielmiessler/Fabric/pull/1603) by [ksylvan](https://github.com/ksylvan): Together AI Support with OpenAI Fallback Mechanism Added

- Add direct model fetching support for non-standard providers
- Add `DirectlyGetModels` function to handle non-standard API responses
- Implement fallback to direct model fetching when standard method fails
- Enhance error messages in OpenAI compatible models endpoint with response body details
- Add context support to DirectlyGetModels method

### PR [#1599](https://github.com/danielmiessler/Fabric/pull/1599) by [ksylvan](https://github.com/ksylvan): Update file paths to reflect new data directory structure

- Update file paths to reflect new data directory structure
- Move fabric logo image path to docs directory
- Update patterns directory reference to data/patterns location
- Update strategies directory reference to data/strategies location
- Fix create_coding_feature README path reference

### Direct commits

- Broken image link

## v1.4.244 (2025-07-09)

### PR [#1598](https://github.com/danielmiessler/Fabric/pull/1598) by [jaredmontoya](https://github.com/jaredmontoya): flake: fixes and enhancements

- Updated Nix package to use self reference for better dependency management
- Renamed shell command for improved clarity
- Fixed generation path in update-mod functionality
- Corrected typo in shell configuration

## v1.4.243 (2025-07-09)

### PR [#1597](https://github.com/danielmiessler/Fabric/pull/1597) by [ksylvan](https://github.com/ksylvan): CLI Refactoring: Modular Command Processing and Pattern Loading Improvements

- Refactor CLI to modularize command handling with specialized handlers for setup, configuration, listing, management, and extensions
- Extract chat processing logic into separate function and improve patterns loader with migration support
- Add tool processing for YouTube and web scraping functionality with enhanced error handling
- Implement `handled` boolean return system across all command handlers for better control flow
- Improve error handling with proper wrapping, secure temporary directory creation, and context information

### Direct commits

- Nix:pkgs:fabric: use self reference
- Update-mod: fix generation path
- Shell: rename command

## v1.4.242 (2025-07-09)

### PR [#1596](https://github.com/danielmiessler/Fabric/pull/1596) by [ksylvan](https://github.com/ksylvan): Fix patterns zipping workflow

- Update workflow paths to reflect directory structure change
- Modify trigger path to `data/patterns/**`
- Update `git diff` command to new path
- Change zip command to include `data/patterns/` directory

## v1.4.241 (2025-07-09)

### PR [#1595](https://github.com/danielmiessler/Fabric/pull/1595) by [ksylvan](https://github.com/ksylvan): Restructure project to align with standard Go layout

- Restructured project to follow standard Go layout with `cmd` directory for binaries
- Moved all Go packages into `internal` directory and renamed `restapi` to `server`
- Consolidated patterns/strategies into `data` directory and scripts into `scripts` directory
- Updated all import paths and CI/CD workflows for new structure
- Added new patterns for content tagging and cognitive bias analysis

### PR [#1594](https://github.com/danielmiessler/Fabric/pull/1594) by [amancioandre](https://github.com/amancioandre): Adds check Dunning-Kruger Telos self-evaluation pattern

- Added pattern for Dunning-Kruger cognitive bias evaluation

## v1.4.240 (2025-07-07)

### PR [#1593](https://github.com/danielmiessler/Fabric/pull/1593) by [ksylvan](https://github.com/ksylvan): Refactor: Generalize OAuth flow for improved token handling

- Replace hardcoded "claude" with configurable `authTokenIdentifier` parameter for better flexibility
- Update `RunOAuthFlow`, `RefreshToken`, and `exchangeToken` functions to accept token identifier parameter
- Add token refresh attempt before full OAuth flow and improve existing token validation
- Create comprehensive OAuth testing suite with 434 lines coverage including mock token server
- Implement PKCE generation, token expiration logic, and performance benchmark tests

## v1.4.239 (2025-07-07)

### PR [#1592](https://github.com/danielmiessler/Fabric/pull/1592) by [ksylvan](https://github.com/ksylvan): Fix Streaming Error Handling in Chatter

- Improve error handling in streaming chat functionality with dedicated error channel
- Add proper goroutine synchronization using `done` channel for completion tracking
- Rename variables (`channel` to `responseChan`, `doneChan` to `done`) for better clarity
- Implement comprehensive testing with mock vendor for error propagation validation
- Streamline response aggregation and ensure proper resource cleanup in streaming operations

## v1.4.238 (2025-07-07)

### PR [#1591](https://github.com/danielmiessler/Fabric/pull/1591) by [ksylvan](https://github.com/ksylvan): Improved Anthropic Plugin Configuration Logic

- Add vendor configuration validation and OAuth auto-authentication
- Implement IsConfigured method for Anthropic client validation
- Add automatic OAuth flow when no valid token
- Add token expiration checking with 5-minute buffer
- Extract vendor token identifier constant and remove redundant configure call

## v1.4.237 (2025-07-07)

### PR [#1590](https://github.com/danielmiessler/Fabric/pull/1590) by [ksylvan](https://github.com/ksylvan): Do not pass non-default TopP values

- Add conditional check for TopP parameter in OpenAI client
- Add zero-value check before setting TopP parameter
- Prevent sending TopP when value is zero
- Apply fix to both chat completions method
- Apply fix to response parameters method

## v1.4.236 (2025-07-06)

### PR [#1587](https://github.com/danielmiessler/Fabric/pull/1587) by [ksylvan](https://github.com/ksylvan): Enhance bug report template

- Enhanced bug report template with detailed system info and installation method fields
- Added detailed instructions for bug reproduction steps
- Included operating system dropdown with specific architectures
- Added OS version textarea with command examples
- Created installation method dropdown with all options

## v1.4.235 (2025-07-06)

### PR [#1586](https://github.com/danielmiessler/Fabric/pull/1586) by [ksylvan](https://github.com/ksylvan): Fix to persist the CUSTOM_PATTERNS_DIRECTORY variable

- Make custom patterns persist correctly

## v1.4.234 (2025-07-06)

### PR [#1581](https://github.com/danielmiessler/Fabric/pull/1581) by [ksylvan](https://github.com/ksylvan): Fix Custom Patterns Directory Creation Logic

- Improve directory creation logic in `configure` method
- Add `fmt` package for logging errors
- Check directory existence before creating
- Log error without clearing directory value

## v1.4.233 (2025-07-06)

### PR [#1580](https://github.com/danielmiessler/Fabric/pull/1580) by [ksylvan](https://github.com/ksylvan): Alphabetical Pattern Sorting and Configuration Refactor

- Move custom patterns directory initialization to Configure method
- Add alphabetical sorting to pattern names retrieval
- Override ListNames method for PatternsEntity class
- Improve pattern listing with proper error handling
- Ensure custom patterns loaded after environment configuration

### PR [#1578](https://github.com/danielmiessler/Fabric/pull/1578) by [ksylvan](https://github.com/ksylvan): Document Custom Patterns Directory Support

- Add comprehensive custom patterns setup and usage guide
- Document priority system for custom vs built-in patterns
- Include step-by-step custom pattern creation workflow
- Explain update-safe custom pattern storage
- Document seamless integration with existing fabric commands

## v1.4.232 (2025-07-06)

### PR [#1577](https://github.com/danielmiessler/Fabric/pull/1577) by [ksylvan](https://github.com/ksylvan): Add Custom Patterns Directory Support

- Add custom patterns directory support with environment variable configuration
- Implement custom patterns plugin with registry integration
- Override main patterns with custom directory patterns
- Expand home directory paths in custom patterns config
- Add comprehensive test coverage for custom patterns functionality

## v1.4.231 (2025-07-05)

### PR [#1565](https://github.com/danielmiessler/Fabric/pull/1565) by [ksylvan](https://github.com/ksylvan): OAuth Authentication Support for Anthropic

- Add OAuth authentication support for Anthropic Claude with PKCE flow and browser integration
- Implement automatic OAuth token refresh and persistent storage for seamless authentication
- Support both API key and OAuth authentication methods with fallback re-authentication
- Extract OAuth functionality to separate module for cleaner code organization
- Standardize all API calls to use v2 endpoint and simplify base URL configuration

## v1.4.230 (2025-07-05)

### PR [#1575](https://github.com/danielmiessler/Fabric/pull/1575) by [ksylvan](https://github.com/ksylvan): Advanced image generation parameters for OpenAI models

- Add four new image generation CLI flags for enhanced control
- Implement validation for image parameter combinations
- Support size, quality, compression, and background controls
- Add comprehensive test coverage for new parameters
- Update shell completions and README with detailed examples

## v1.4.229 (2025-07-05)

### PR [#1574](https://github.com/danielmiessler/Fabric/pull/1574) by [ksylvan](https://github.com/ksylvan): Add Model Validation for Image Generation and Fix CLI Flag Mapping

- Add model validation for image generation support with `supportsImageGeneration` function
- Add model field to `BuildChatOptions` method for proper CLI flag mapping
- Extract supported models list to shared constant `ImageGenerationSupportedModels` for reusability
- Implement validation in `sendResponses` to ensure model supports image generation
- Add comprehensive tests for model validation logic in `TestModelValidationLogic`

## v1.4.228 (2025-07-05)

### PR [#1573](https://github.com/danielmiessler/Fabric/pull/1573) by [ksylvan](https://github.com/ksylvan): Add Image File Validation and Dynamic Format Support

- Add image file validation and format detection for image generation
- Implement dynamic output format detection from file extensions
- Add comprehensive test coverage for image file validation
- Upgrade YAML library from v2 to v3
- Support PNG, JPEG, JPG, and WEBP image formats

### Direct commits

- Added tutorial as a tag

## v1.4.227 (2025-07-04)

### PR [#1572](https://github.com/danielmiessler/Fabric/pull/1572) by [ksylvan](https://github.com/ksylvan): Add Image Generation Support to Fabric

- Add image generation support with OpenAI image generation model
- Add `--image-file` flag for saving generated images
- Implement image generation tool integration with OpenAI
- Add web search tool for Anthropic and OpenAI models
- Support PNG, JPG, JPEG, GIF, BMP image formats

### Direct commits

- Fixed ul tag applier
- Updated ul tag prompt
- Added the UL tags pattern

## v1.4.226 (2025-07-04)

### PR [#1569](https://github.com/danielmiessler/Fabric/pull/1569) by [ksylvan](https://github.com/ksylvan): OpenAI Plugin Now Supports Web Search Functionality

- Add web search tool support for OpenAI models with citation formatting
- Enable web search tool for OpenAI models with location parameter support
- Extract and format citations from search responses with deduplication
- Implement comprehensive test coverage for search functionality
- Update CLI flag description and README with new web search feature details

## v1.4.225 (2025-07-04)

### PR [#1568](https://github.com/danielmiessler/Fabric/pull/1568) by [ksylvan](https://github.com/ksylvan): Runtime Web Search Control via Command-Line Flag

- Add web search tool support for Anthropic models with --search flag
- Add --search-location for timezone-based results through ChatOptions struct
- Implement web search tool in Anthropic client with formatted citations
- Add comprehensive tests and remove plugin-level web search configuration
- Extract web search tool constants and optimize string building with sources header

### Direct commits

- Merge branch 'main' of <https://github.com/amancioandre/Fabric>
- Sections as heading 1, typos
- Merge branch 'danielmiessler:main' into main
- Adds pattern telos check dunning kruger

## v1.4.224 (2025-07-01)

### PR [#1564](https://github.com/danielmiessler/Fabric/pull/1564) by [ksylvan](https://github.com/ksylvan): Add code_review pattern and updates in Pattern_Descriptions

- Add comprehensive code review pattern with systematic analysis framework and principal engineer reviewer role
- Add new patterns: `review_code`, `extract_alpha`, and `extract_mcp_servers` for enhanced functionality
- Refactor pattern extraction script with improved error handling and docstrings for better clarity
- Add JSONDecodeError handling in `load_existing_file` with graceful fallback to empty list
- Fix typo in `analyze_bill_short` pattern description and improve formatting in pattern management README

## v1.4.223 (2025-07-01)

### PR [#1563](https://github.com/danielmiessler/Fabric/pull/1563) by [ksylvan](https://github.com/ksylvan): Fix Cross-Platform Compatibility in Release Workflow

- Update GitHub Actions to use bash shell in release job
- Adjust repository_dispatch type spacing for consistency
- Use bash shell for creating release if absent

## v1.4.222 (2025-07-01)

### PR [#1559](https://github.com/danielmiessler/Fabric/pull/1559) by [ksylvan](https://github.com/ksylvan): OpenAI Plugin Migrates to New Responses API

- Migrate OpenAI plugin to use new responses API instead of chat completions
- Add chat completions API fallback for non-Responses API providers
- Implement `sendChatCompletions` and `sendStreamChatCompletions` methods
- Add `ImplementsResponses` flag to track provider API capabilities
- Extract common message conversion logic to reduce duplication

### Direct commits

- Updated alpha post
- Updated extract alpha
- Added extract_alpha as kind of an experiment

## v1.4.221 (2025-06-28)

### PR [#1556](https://github.com/danielmiessler/Fabric/pull/1556) by [ksylvan](https://github.com/ksylvan): feat: Migrate to official openai-go SDK

- Abstract chat message structs and migrate to official openai-go SDK
- Introduce local `chat` package for message abstraction
- Replace sashabaranov/go-openai with official openai-go SDK
- Update OpenAI, Azure, and Exolab plugins for new client
- Refactor all AI providers to use internal chat types

## v1.4.220 (2025-06-28)

### PR [#1555](https://github.com/danielmiessler/Fabric/pull/1555) by [ksylvan](https://github.com/ksylvan): fix: Race condition in GitHub actions release flow

- Improve release creation to gracefully handle pre-existing tags
- Check if a release exists before attempting creation
- Suppress error output from `gh release view` command
- Add an informative log when release already exists

## v1.4.219 (2025-06-28)

### PR [#1553](https://github.com/danielmiessler/Fabric/pull/1553) by [ksylvan](https://github.com/ksylvan): docs: add DeepWiki badge and fix minor typos in README

- Add DeepWiki badge to README header
- Fix typo "chatbots" to "chat-bots" and "Perlexity" to "Perplexity"
- Correct "distro" to "Linux distribution"
- Add alt text to contributor images
- Update dependency versions in go.mod and remove unused soup dependency

### PR [#1552](https://github.com/danielmiessler/Fabric/pull/1552) by [nawarajshahi](https://github.com/nawarajshahi): Fix typos in README.md

- Fix typos on README.md

## v1.4.218 (2025-06-27)

### PR [#1550](https://github.com/danielmiessler/Fabric/pull/1550) by [ksylvan](https://github.com/ksylvan): Add Support for OpenAI Search and Research Model Variants

- Add support for new OpenAI search and research model variants
- Add slices import for array operations
- Define new search preview model names and mini search preview variants
- Include deep research model support with June 2025 dated model versions
- Replace hardcoded check with slices.Contains for both prefix and exact model matching

## v1.4.217 (2025-06-26)

### PR [#1546](https://github.com/danielmiessler/Fabric/pull/1546) by [ksylvan](https://github.com/ksylvan): New YouTube Transcript Endpoint Added to REST API

- Add dedicated YouTube transcript API endpoint
- Create `/youtube/transcript` POST endpoint route
- Add request/response types for YouTube API
- Support language and timestamp options
- Update frontend to use new endpoint

### Direct commits

- Add extract_mcp_servers pattern
New pattern to extract mentions of MCP (Model Context Protocol) servers from content. Identifies server names, features, capabilities, and usage examples.
🤖 Generated with [Claude Code](<https://claude.ai/code)>
Co-Authored-By: Claude <noreply@anthropic.com>

## v1.4.216 (2025-06-26)

### PR [#1545](https://github.com/danielmiessler/Fabric/pull/1545) by [ksylvan](https://github.com/ksylvan): Update Message Handling for Attachments and Multi-Modal content

- Allow combining user messages and attachments with patterns
- Refactor chat request builder for improved clarity and enhanced dryrun client to display multi-content user messages
- Handle multi-content messages for user role and display image URLs from user messages in output
- Fix duplicate user message issue when applying patterns and ensure multi-part content is always included in session
- Extract message and option formatting logic into reusable methods to reduce code duplication and improve maintainability

## v1.4.215 (2025-06-25)

### PR [#1543](https://github.com/danielmiessler/Fabric/pull/1543) by [ksylvan](https://github.com/ksylvan): fix: Revert multiline tags in generated json files

- Reformat `pattern_descriptions.json` to improve readability
- Reformat JSON `tags` array to display on new lines
- Update `write_essay` pattern description for clarity
- Apply consistent formatting to both data files

## v1.4.214 (2025-06-25)

### PR [#1542](https://github.com/danielmiessler/Fabric/pull/1542) by [ksylvan](https://github.com/ksylvan): Add `write_essay_by_author` and update Pattern metadata

- Refactor ProviderMap for dynamic URL template handling with environment variables
- Add new patterns: `analyze_terraform_plan`, `write_essay_by_author`, `summarize_board_meeting`, `create_mnemonic_phrases`
- Rename `write_essay` to `write_essay_pg` for Paul Graham style specificity
- Update pattern metadata files with tags and descriptions for new analytical patterns
- Sort pattern explanations alphabetically and clean up duplicate entries

## v1.4.213 (2025-06-23)

### PR [#1538](https://github.com/danielmiessler/Fabric/pull/1538) by [andrewsjg](https://github.com/andrewsjg): Bug/bedrock region handling

- Updated hasAWSCredentials to check for AWS_DEFAULT_REGION when access keys are configured
- Fixed bedrock region handling with correct pointer reference and region value setting
- Refactored Bedrock client with improved error handling and ai.Vendor interface compliance
- Added AWS region validation logic and enhanced resource cleanup in SendStream
- Improved code documentation and added user agent constants with proper context usage

## v1.4.212 (2025-06-23)

### PR [#1540](https://github.com/danielmiessler/Fabric/pull/1540) by [ksylvan](https://github.com/ksylvan): Add Langdock AI and enhance generic OpenAI compatible support

- Refactor ProviderMap for dynamic URL template handling with environment variables
- Add `os` and `strings` packages to imports for template processing
- Implement dynamic URL handling using environment variables or default values
- Reorder providers for consistent key order in ProviderMap
- Extract and parse template variables from BaseURL

### Direct commits

- Refactor Bedrock client with improved error handling and interface compliance
- Add AWS region validation logic and fix resource cleanup in SendStream
- Enhanced code documentation and user agent constants
- Fixed Bedrock region handling with proper pointer reference resolution
- Updated paper analyzer functionality

## v1.4.211 (2025-06-19)

### PR [#1533](https://github.com/danielmiessler/Fabric/pull/1533) by [ksylvan](https://github.com/ksylvan): REST API and Web UI Now Support Dynamic Pattern Variables

- Add pattern variables support to REST API chat endpoint with Variables field in PromptRequest struct
- Add pattern variables UI in web interface with JSON textarea for variable input
- Add `ApplyPattern` route for applying patterns with variables via POST /patterns/:name/apply
- Refactor ChatService to clean up message stream and pattern output methods
- Remove unnecessary raycast scripts directory from patterns/ folder

### Direct commits

- Updated paper analyzer format and sanitization instructions
- Updated markdown cleaner functionality

## v1.4.210 (2025-06-18)

### PR [#1530](https://github.com/danielmiessler/Fabric/pull/1530) by [ksylvan](https://github.com/ksylvan): Add Citation Support to Perplexity Response

- Add citation support to perplexity AI responses
- Add citation extraction from API responses
- Append citations section to response content
- Format citations as numbered markdown list
- Handle citations in streaming responses

### Direct commits

- Update README.md
- Updated readme and intro text describing Fabric's utility

## v1.4.208 (2025-06-17)

### PR [#1527](https://github.com/danielmiessler/Fabric/pull/1527) by [ksylvan](https://github.com/ksylvan): Add Perplexity AI Provider with Token Limits Support

- Add Perplexity AI provider support with token limits and streaming
- Add `MaxTokens` field to `ChatOptions` struct for response control
- Integrate Perplexity client into core plugin registry initialization
- Implement stream handling in Perplexity client using sync.WaitGroup
- Update README with Perplexity AI support instructions

### PR [#1526](https://github.com/danielmiessler/Fabric/pull/1526) by [ConnorKirk](https://github.com/ConnorKirk): Check for AWS_PROFILE or AWS_ROLE_SESSION_NAME environment variables

- Check for AWS_PROFILE or AWS_ROLE_SESSION_NAME environment variables

## v1.4.207 (2025-06-17)

### PR [#1525](https://github.com/danielmiessler/Fabric/pull/1525) by [ksylvan](https://github.com/ksylvan): Refactor yt-dlp Transcript Logic and Fix Language Bug

- Extract common yt-dlp logic to reduce code duplication in YouTube plugin
- Add processVTTFileFunc parameter for flexible VTT processing
- Implement language matching for 2-char language codes
- Refactor transcript methods to use new helper function
- Maintain existing functionality with cleaner structure

### Direct commits

- Updated extract insights

## v1.4.206 (2025-06-16)

### PR [#1523](https://github.com/danielmiessler/Fabric/pull/1523) by [ksylvan](https://github.com/ksylvan): Conditional AWS Bedrock Plugin Initialization

- Add AWS credential detection for Bedrock client initialization
- Add hasAWSCredentials helper function to check for AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
- Look for AWS shared credentials file with support for custom AWS_SHARED_CREDENTIALS_FILE path
- Default to ~/.aws/credentials location for credential detection
- Only initialize Bedrock client if credentials exist to prevent AWS SDK credential search failures

### Direct commits

- Updated prompt.

## v1.4.205 (2025-06-16)

### PR [#1519](https://github.com/danielmiessler/Fabric/pull/1519) by [ConnorKirk](https://github.com/ConnorKirk): feat: Dynamically list AWS Bedrock models

- Dynamically fetch and list available foundation models and inference profiles

### PR [#1518](https://github.com/danielmiessler/Fabric/pull/1518) by [ksylvan](https://github.com/ksylvan): chore: remove duplicate/outdated patterns

- Remove duplicate/outdated patterns

### Direct commits

- Updated markdown sanitizer
- Updated markdown cleaner

## v1.4.204 (2025-06-15)

### PR [#1517](https://github.com/danielmiessler/Fabric/pull/1517) by [ksylvan](https://github.com/ksylvan): Fix: Prevent race conditions in versioning workflow

- Improve version update workflow to prevent race conditions
- Add concurrency control to prevent simultaneous runs
- Pull latest main branch changes before tagging
- Fetch all remote tags before calculating version

## v1.4.203 (2025-06-14)

### PR [#1512](https://github.com/danielmiessler/Fabric/pull/1512) by [ConnorKirk](https://github.com/ConnorKirk): feat:Add support for Amazon Bedrock

- Add Bedrock plugin for Amazon Bedrock integration within fabric

### PR [#1513](https://github.com/danielmiessler/Fabric/pull/1513) by [marcas756](https://github.com/marcas756): feat: create mnemonic phrase pattern

- Create mnemonic phrase pattern for generating phrases from diceware words
- Add markdown files with user guide and system implementation details

### PR [#1516](https://github.com/danielmiessler/Fabric/pull/1516) by [ksylvan](https://github.com/ksylvan): Fix REST API pattern creation

- Add Save method to PatternsEntity for persisting patterns to filesystem
- Create pattern directory with proper permissions and write content to system files
- Add comprehensive tests for Save functionality with error handling

## v1.4.202 (2025-06-12)

### PR [#1510](https://github.com/danielmiessler/Fabric/pull/1510) by [ksylvan](https://github.com/ksylvan): Cross-Platform fix for Youtube Transcript extraction

- Replace hardcoded `/tmp` with `os.TempDir()` for cross-platform compatibility
- Use `filepath.Join()` instead of string concatenation for proper path handling
- Remove Unix `find` command dependency completely
- Add new `findVTTFiles()` method using `filepath.Walk()` for Windows support
- Improve error handling for file operations while maintaining backward compatibility

## v1.4.201 (2025-06-12)

### PR [#1503](https://github.com/danielmiessler/Fabric/pull/1503) by [dependabot[bot]](https://github.com/apps/dependabot): chore(deps): bump brace-expansion from 1.1.11 to 1.1.12 in /web in the npm_and_yarn group across 1 directory

- Updated brace-expansion dependency from 1.1.11 to 1.1.12 in /web directory
- Indirect dependency update in npm_and_yarn group

### PR [#1508](https://github.com/danielmiessler/Fabric/pull/1508) by [ksylvan](https://github.com/ksylvan): feat: cleanup after `yt-dlp` addition

- Updated README with yt-dlp requirement for transcripts
- Improved error messages for better clarity and actionability
- General cleanup following yt-dlp integration

## v1.4.200 (2025-06-11)

### PR [#1507](https://github.com/danielmiessler/Fabric/pull/1507) by [ksylvan](https://github.com/ksylvan): Refactor: No more web scraping, just use yt-dlp

- Replace web scraping with yt-dlp for YouTube transcript extraction
- Remove unreliable YouTube API scraping methods
- Add yt-dlp integration for transcript extraction
- Implement VTT subtitle parsing functionality
- Add timestamp preservation for transcripts

## v1.4.199 (2025-06-11)

### PR [#1506](https://github.com/danielmiessler/Fabric/pull/1506) by [eugeis](https://github.com/eugeis): fix: fix web search tool location

- Fix web search tool location

## v1.4.198 (2025-06-11)

### PR [#1504](https://github.com/danielmiessler/Fabric/pull/1504) by [marcas756](https://github.com/marcas756): fix: Add configurable HTTP timeout for Ollama client

- Added configurable HTTP timeout for Ollama client
- Introduced new setup question to configure timeout duration for Ollama requests
- Set default timeout value to 20 minutes
- Improved request handling reliability for Ollama integration

## v1.4.197 (2025-06-11)

### PR [#1502](https://github.com/danielmiessler/Fabric/pull/1502) by [eugeis](https://github.com/eugeis): Feat/antropic tool

- Search tool working
- Search tool result collection

### PR [#1499](https://github.com/danielmiessler/Fabric/pull/1499) by [noamsiegel](https://github.com/noamsiegel): feat: Enhance the PRD Generator's identity and purpose

- Enhanced PRD Generator identity and purpose for better clarity
- Added structured output format with Markdown, sections, and bullet points
- Defined key PRD sections: Overview, Objectives, Features, User Stories, Requirements
- Improved instructions for highlighting priorities and MVP features

### PR [#1497](https://github.com/danielmiessler/Fabric/pull/1497) by [ksylvan](https://github.com/ksylvan): feat: add Terraform plan analyzer pattern for infrastructure changes

- Added Terraform plan analyzer pattern for infrastructure change assessment
- Included security, cost, and compliance focus areas
- Created structured output with summaries, critical changes, and key takeaways
- Required numbered lists and specific word limits for sections

### Direct commits

- Dependency update: bumped brace-expansion from 1.1.11 to 1.1.12
- Added configurable HTTP timeout for Ollama client with 20-minute default

## v1.4.196 (2025-06-07)

### PR [#1495](https://github.com/danielmiessler/Fabric/pull/1495) by [ksylvan](https://github.com/ksylvan): Add AIML provider configuration

- Add AIML provider to OpenAI compatible providers configuration
- Set AIML base URL to api.aimlapi.com/v1
- Expand supported OpenAI compatible providers list
- Enable AIML API integration support
- Added simpler paper analyzer with updated output

## v1.4.195 (2025-05-24)

### PR [#1487](https://github.com/danielmiessler/Fabric/pull/1487) by [ksylvan](https://github.com/ksylvan): Dependency Updates and PDF Worker Refactoring

- Upgrade PDF.js to v4.2 and refactor worker initialization
- Add `.browserslistrc` to define target browser versions
- Upgrade `pdfjs-dist` dependency from v2.16 to v4.2.67
- Upgrade `nanoid` dependency from v4.0.2 to v5.0.9
- Introduce `pdf-config.ts` for centralized PDF.js worker setup

## v1.4.194 (2025-05-24)

### PR [#1485](https://github.com/danielmiessler/Fabric/pull/1485) by [ksylvan](https://github.com/ksylvan): Web UI: Centralize Environment Configuration and Make Fabric Base URL Configurable

- Add centralized environment configuration for Fabric base URL
- Create environment config module for URL handling
- Add getFabricBaseUrl() function with server/client support
- Add getFabricApiUrl() helper for API endpoints
- Configure Vite to inject FABRIC_BASE_URL client-side

## v1.4.193 (2025-05-24)

### PR [#1484](https://github.com/danielmiessler/Fabric/pull/1484) by [ksylvan](https://github.com/ksylvan): Web UI update all packages, reorganize docs, add install scripts

- Reorganize web documentation and add installation scripts
- Update all package dependencies to latest versions
- Add PDF-to-Markdown installation steps to README
- Move legacy documentation files to web/legacy/
- Add convenience scripts for npm and pnpm installation

### PR [#1481](https://github.com/danielmiessler/Fabric/pull/1481) by [skibum1869](https://github.com/skibum1869): Add board meeting summary pattern template

- Add board meeting summary pattern template
- Update meeting summary template with word count requirement
- Add minimum word count for context section in board summary

### Direct commits

- Add centralized environment configuration for Fabric base URL
- Create getFabricBaseUrl() function with server/client support
- Configure Vite to inject FABRIC_BASE_URL client-side
- Support FABRIC_BASE_URL env var with fallback
- Add TypeScript definitions for window config

## v1.4.192 (2025-05-23)

### PR [#1480](https://github.com/danielmiessler/Fabric/pull/1480) by [ksylvan](https://github.com/ksylvan): Automatic setting of "raw mode" for some models

- Add NeedsRawMode method to AI vendor interface for automatic raw mode detection
- Implement NeedsRawMode in all AI clients with model-specific detection logic
- Enable automatic raw mode for Ollama llama2/llama3 models and OpenAI o1/o3/o4 models
- Auto-enable raw mode when vendor requires it based on model configuration
- Import strings package for prefix matching functionality

## v1.4.191 (2025-05-22)

### PR [#1478](https://github.com/danielmiessler/Fabric/pull/1478) by [ksylvan](https://github.com/ksylvan): Claude 4 Integration and README Updates

- Add support for Anthropic Claude 4 models and update SDK to v1.2.0
- Upgrade `anthropic-sdk-go` dependency to version `v1.2.0`
- Integrate new Anthropic Claude 4 Opus and Sonnet models
- Remove deprecated Claude 2.0 and 2.1 models from list
- Adjust model type casting for `anthropic-sdk-go v1.2.0` compatibility

## v1.4.190 (2025-05-20)

### PR [#1475](https://github.com/danielmiessler/Fabric/pull/1475) by [ksylvan](https://github.com/ksylvan): refactor: improve raw mode handling in BuildSession

- Improved raw mode handling in BuildSession with better system message processing
- Fixed duplicate inputs when using patterns in raw mode
- Added conditional logic to distinguish between pattern and non-pattern scenarios
- Simplified message construction with clearer variable names
- Enhanced code comments for better readability

## v1.4.189 (2025-05-19)

### PR [#1473](https://github.com/danielmiessler/Fabric/pull/1473) by [roumy](https://github.com/roumy): add authentification for ollama instance

- Add authentification for ollama instance

## v1.4.188 (2025-05-19)

### PR [#1474](https://github.com/danielmiessler/Fabric/pull/1474) by [ksylvan](https://github.com/ksylvan): feat: update `BuildSession` to handle message appending logic

- Improve message handling for raw mode and Anthropic client
- Fix pattern-based message handling in non-raw mode
- Add proper handling for empty message arrays
- Implement user/assistant message alternation for Anthropic
- Preserve system messages in Anthropic conversations

### PR [#1467](https://github.com/danielmiessler/Fabric/pull/1467) by [joshuafuller](https://github.com/joshuafuller): Typos, spelling, grammar and other minor updates

- Fix spelling in pattern management guide
- Correct Anthropic spelling in notes
- Fix typos in web README
- Fix grammar in nuclei template instructions
- Fix spelling in PR 1284 update notes

### PR [#1468](https://github.com/danielmiessler/Fabric/pull/1468) by [NavNab](https://github.com/NavNab): Refactor content structure in create_hormozi_offer system.md for clarity and readability

- Improved formatting of introduction and content summary sections
- Consolidated repetitive sentences and enhanced text coherence
- Adjusted bullet points and numbering for consistency
- Enhanced visual distinction of key concepts
- Ensured clear articulation of important information

### Direct commits

- Add authentification for ollama instance

## v1.4.187 (2025-05-10)

### PR [#1463](https://github.com/danielmiessler/Fabric/pull/1463) by [CodeCorrupt](https://github.com/CodeCorrupt): Add completion to the build output for Nix

- Add completion files to the build output for Nix

## v1.4.186 (2025-05-06)

### PR [#1459](https://github.com/danielmiessler/Fabric/pull/1459) by [ksylvan](https://github.com/ksylvan): chore: Repository cleanup and .gitignore Update

- Added `coverage.out` to `.gitignore` for ignoring coverage output
- Removed `Alma.md` documentation file from the repository
- Deleted `rate_ai_result.txt` stitch script from `stitches` folder
- Removed `readme.md` for `rate_ai_result` stitch documentation
- Updated `.gitignore` and removed obsolete files

## v1.4.185 (2025-04-28)

### PR [#1453](https://github.com/danielmiessler/Fabric/pull/1453) by [ksylvan](https://github.com/ksylvan): Fix for default model setting

- Introduce `getSortedGroupsItems` for consistent sorting logic
- Add centralized sorting method for groups and items (alphabetical, case-insensitive)
- Replace inline sorting in `Print` with new method
- Update `GetGroupAndItemByItemNumber` to use sorted data
- Ensure original `GroupsItems` remains unmodified

## v1.4.184 (2025-04-25)

### PR [#1447](https://github.com/danielmiessler/Fabric/pull/1447) by [ksylvan](https://github.com/ksylvan): More shell completion scripts: Zsh, Bash, and Fish

- Add shell completion scripts for Zsh, Bash, and Fish
- Create standardized completion scripts in completions/ directory
- Add --shell-complete-list flag for machine-readable output
- Update Print() methods to support plain output format
- Document installation steps for each shell in README

## v1.4.183 (2025-04-23)

### PR [#1431](https://github.com/danielmiessler/Fabric/pull/1431) by [KenMacD](https://github.com/KenMacD): Add a completion script for fish

- Add a completion script for fish

## v1.4.182 (2025-04-23)

### PR [#1441](https://github.com/danielmiessler/Fabric/pull/1441) by [ksylvan](https://github.com/ksylvan): Update go toolchain and go module packages to latest versions

- Updated Go to version 1.24.2 across Dockerfile and Nix configurations
- Refreshed Go module dependencies and updated go.mod/go.sum files
- Updated Nix flake lock file inputs and configured packages for Go 1.24
- Centralized Go version definition with `getGoVersion` function in flake.nix
- Fixed "nix flake check" errors and removed redundant Go version definitions

## v1.4.181 (2025-04-22)

### PR [#1433](https://github.com/danielmiessler/Fabric/pull/1433) by [ksylvan](https://github.com/ksylvan): chore: update Anthropic SDK to v0.2.0-beta.3 and migrate to V2 API

- Upgrade Anthropic SDK from alpha.11 to beta.3
- Update API endpoint from v1 to v2
- Replace anthropic.F() with direct assignment and anthropic.Opt() for optional params
- Simplify event delta handling in streaming
- Change client type from pointer to value type

## v1.4.180 (2025-04-22)

### PR [#1435](https://github.com/danielmiessler/Fabric/pull/1435) by [ksylvan](https://github.com/ksylvan): chore: Fix user input handling when using raw mode and `--strategy` flag

- Unify raw mode message handling and preserve env vars in extension executor
- Refactor BuildSession raw mode to prepend system to user content
- Ensure raw mode messages always have User role
- Append systemMessage separately in non-raw mode sessions
- Store original cmd.Env before context-based exec command creation

### Direct commits

- Update Anthropic SDK to v0.2.0-beta.3 and migrate to V2 API
- Replace anthropic.F() with direct assignment and anthropic.Opt() for optional params
- Change client type from pointer to value type
- Update API endpoint from v1 to v2
- Simplify event delta handling in streaming

## v1.4.179 (2025-04-21)

### PR [#1432](https://github.com/danielmiessler/Fabric/pull/1432) by [ksylvan](https://github.com/ksylvan): chore: fix fabric setup mess-up introduced by sorting lists (tools and models)

- Alphabetize the order of plugin tools
- Sort AI models alphabetically for consistent listing
- Import `sort` and `strings` packages for sorting functionality
- Sort retrieved AI model names alphabetically, ignoring case
- Ensure consistent ordering of AI models in lists

### Direct commits

- Add a completion script for fish

## v1.4.178 (2025-04-21)

### PR [#1427](https://github.com/danielmiessler/Fabric/pull/1427) by [ksylvan](https://github.com/ksylvan): Refactor OpenAI-compatible AI providers and add `--listvendors` flag

- Add `--listvendors` command to list AI vendors
- Introduce `--listvendors` flag to display all AI vendors
- Refactor OpenAI-compatible providers into a unified configuration
- Remove individual vendor packages for streamlined management
- Add sorting for consistent vendor listing output

## v1.4.177 (2025-04-21)

### PR [#1428](https://github.com/danielmiessler/Fabric/pull/1428) by [ksylvan](https://github.com/ksylvan): feat: Alphabetical case-insensitive sorting for groups and items

- Added alphabetical sorting to groups and items in Print method
- Imported `sort` and `strings` packages for sorting functionality
- Implemented case-insensitive sorting for both groups and items
- Created stable copies of groups and items before sorting
- Enhanced display iteration to use sorted collections

## v1.4.176 (2025-04-21)

### PR [#1429](https://github.com/danielmiessler/Fabric/pull/1429) by [ksylvan](https://github.com/ksylvan): feat: enhance StrategyMeta with Prompt field and dynamic naming

- Add `Prompt` field to `StrategyMeta` struct for storing JSON prompt data
- Implement dynamic strategy naming using filename with `strings.TrimSuffix`
- Add alphabetical sorting to groups and items in Print method with case-insensitive ordering
- Introduce `--listvendors` command to display all AI vendors with consistent output
- Refactor OpenAI-compatible providers into unified configuration, removing individual vendor packages

## v1.4.175 (2025-04-19)

### PR [#1418](https://github.com/danielmiessler/Fabric/pull/1418) by [dependabot[bot]](https://github.com/apps/dependabot): chore(deps): bump golang.org/x/net from 0.36.0 to 0.38.0 in the go_modules group across 1 directory

- Updated golang.org/x/net dependency from version 0.36.0 to 0.38.0
- Dependency update applied to go_modules group in root directory
- Indirect dependency type update managed by dependabot
- Security and performance improvements included in newer version
- Automated dependency maintenance to keep project current

## v1.4.174 (2025-04-19)

### PR [#1425](https://github.com/danielmiessler/Fabric/pull/1425) by [ksylvan](https://github.com/ksylvan): feat: add Cerebras AI plugin to plugin registry

- Add Cerebras AI plugin to plugin registry
- Introduce Cerebras AI plugin import in plugin registry
- Register Cerebras client in the NewPluginRegistry function

## v1.4.173 (2025-04-18)

### PR [#1420](https://github.com/danielmiessler/Fabric/pull/1420) by [sherif-fanous](https://github.com/sherif-fanous): Fix error in deleting patterns due to non empty directory

- Fix error in deleting patterns due to non empty directory

### PR [#1421](https://github.com/danielmiessler/Fabric/pull/1421) by [ksylvan](https://github.com/ksylvan): feat: add Atom-of-Thought (AoT) strategy and prompt definition

- Add Atom-of-Thought (AoT) strategy and prompt definition
- Add new aot.json for Atom-of-Thought (AoT) prompting
- Define AoT strategy description and detailed prompt instructions
- Update strategies.json to include AoT in available strategies list
- Ensure AoT strategy appears alongside CoD, CoT, and LTM options

### Direct commits

- Chore(deps): bump golang.org/x/net from 0.36.0 to 0.38.0

## v1.4.172 (2025-04-16)

### PR [#1415](https://github.com/danielmiessler/Fabric/pull/1415) by [ksylvan](https://github.com/ksylvan): feat: add Grok AI provider support

- Add Grok AI provider support for AI model interactions
- Integrate Grok AI client into plugin registry
- Include Grok AI API key in REST API configuration endpoints
- Update README with Grok documentation

### PR [#1411](https://github.com/danielmiessler/Fabric/pull/1411) by [ksylvan](https://github.com/ksylvan): docs: add contributors section to README with contrib.rocks image

- Add contributors section with visual representation to README
- Include link to project contributors page
- Add attribution to contrib.rocks tool

## v1.4.171 (2025-04-15)

### PR [#1407](https://github.com/danielmiessler/Fabric/pull/1407) by [sherif-fanous](https://github.com/sherif-fanous): Update Dockerfile so that Go image version matches go.mod version

- Bump golang version to match go.mod

### Direct commits

- Multiple README.md updates (12 commits)

## v1.4.170 (2025-04-13)

### PR [#1406](https://github.com/danielmiessler/Fabric/pull/1406) by [jmd1010](https://github.com/jmd1010): Fix chat history LLM response sequence in ChatInput.svelte

- Fix chat history LLM response sequence in ChatInput.svelte
- Finalize WEB UI V2 loose ends fixes
- Update pattern_descriptions.json
- Bump golang version to match go.mod

## v1.4.169 (2025-04-11)

### PR [#1403](https://github.com/danielmiessler/Fabric/pull/1403) by [jmd1010](https://github.com/jmd1010): Strategy flag enhancement - Web UI implementation

- Integrated strategy flag enhancement from fabric CLI into web UI
- Updated strategies.json configuration
- Added new excalidraw pattern for visual documentation
- Implemented bill analyzer functionality with shorter version
- Enhanced analyze bill pattern for improved processing

## v1.4.168 (2025-04-02)

### PR [#1399](https://github.com/danielmiessler/Fabric/pull/1399) by [HaroldFinchIFT](https://github.com/HaroldFinchIFT): feat: add simple optional api key management for protect routes in --serve mode

- Add simple optional API key management for protected routes in --serve mode
- Refactor API key middleware based on code review feedback
- Fix formatting issues

## v1.4.167 (2025-03-31)

### PR [#1397](https://github.com/danielmiessler/Fabric/pull/1397) by [HaroldFinchIFT](https://github.com/HaroldFinchIFT): feat: add it lang to the chat drop down menu lang in web gui

- Add it lang to the chat drop down menu lang in web gui

## v1.4.166 (2025-03-29)

### PR [#1392](https://github.com/danielmiessler/Fabric/pull/1392) by [ksylvan](https://github.com/ksylvan): chore: enhance argument validation in `code_helper` tool

- Streamline code_helper CLI interface and require explicit instructions
- Require exactly two arguments: directory and instructions
- Remove dedicated help flag, use flag.Usage instead
- Improve directory validation to check if it's a directory
- Inline pattern parsing, removing separate function

### PR [#1390](https://github.com/danielmiessler/Fabric/pull/1390) by [PatrickCLee](https://github.com/PatrickCLee): docs: improve README link

- Fix broken what-and-why link reference

## v1.4.165 (2025-03-26)

### PR [#1389](https://github.com/danielmiessler/Fabric/pull/1389) by [ksylvan](https://github.com/ksylvan): Create Coding Feature

- Add `code_helper` tool (renamed from `fabric_code`) for AI-driven codebase modifications
- Implement `create_coding_feature` pattern with file management API for code changes
- Add secure file parsing and validation system with JSON escape sequence handling
- Update README with installation instructions and usage examples
- Replace deprecated io/ioutil with modern alternatives and improve error handling

### Direct commits

- Improve README link
- Fix broken what-and-why link reference

## v1.4.164 (2025-03-22)

### PR [#1380](https://github.com/danielmiessler/Fabric/pull/1380) by [jmd1010](https://github.com/jmd1010): Add flex windows sizing to web interface + raw text input fix

- Add flex windows sizing to web interface
- Fixed processing message not stopping after pattern output completion

### PR [#1379](https://github.com/danielmiessler/Fabric/pull/1379) by [guilhermechapiewski](https://github.com/guilhermechapiewski): Fix typo on fallacies instruction

- Fix typo on fallacies instruction

### PR [#1382](https://github.com/danielmiessler/Fabric/pull/1382) by [ksylvan](https://github.com/ksylvan): docs: improve README formatting and fix some broken links

- Improve README formatting and add clipboard support section
- Fix broken installation link reference and environment variables link
- Replace code tags with backticks and improve code block formatting

### PR [#1376](https://github.com/danielmiessler/Fabric/pull/1376) by [vaygr](https://github.com/vaygr): Add installation instructions for OS package managers

- Add installation instructions for OS package managers

### Direct commits

- Added find_female_life_partner pattern
- Updated find prompt multiple times

## v1.4.163 (2025-03-19)

### PR [#1362](https://github.com/danielmiessler/Fabric/pull/1362) by [dependabot[bot]](https://github.com/apps/dependabot): Bump golang.org/x/net from 0.35.0 to 0.36.0 in the go_modules group across 1 directory

- Updated golang.org/x/net dependency from version 0.35.0 to 0.36.0

### PR [#1372](https://github.com/danielmiessler/Fabric/pull/1372) by [rube-de](https://github.com/rube-de): fix: set percentEncoded to false

- Fixed YouTube link encoding issue by setting percentEncoded to false
- Prevents URL encoding errors when using YouTube links like youtu.be/sHIlFKKaq0A

### PR [#1373](https://github.com/danielmiessler/Fabric/pull/1373) by [ksylvan](https://github.com/ksylvan): Remove unnecessary `system.md` file at top level

- Removed redundant system.md file from top-level directory
- File was an RPG session summarization prompt superseded by create_rpg_summary and summarize_rpg_session patterns

## v1.4.162 (2025-03-19)

### PR [#1374](https://github.com/danielmiessler/Fabric/pull/1374) by [ksylvan](https://github.com/ksylvan): Fix Default Model Change Functionality

- Improve error handling in ChangeDefaultModel flow and save environment file
- Add early return on setup error
- Save environment file after successful setup
- Maintain proper error propagation

### Direct commits

- Remove redundant file system.md at top level
- Set percentEncoded to false for YouTube links to prevent encoding errors
- Removed RPG session summarization prompt (system.md) - replaced by create_rpg_summary and summarize_rpg_session patterns
- Fix YouTube link processing to avoid URL encoding issues that caused fabric errors

## v1.4.161 (2025-03-17)

### PR [#1363](https://github.com/danielmiessler/Fabric/pull/1363) by [garkpit](https://github.com/garkpit): clipboard operations now work on Mac and PC

- Clipboard operations now work on Mac and PC

## v1.4.160 (2025-03-17)

### PR [#1368](https://github.com/danielmiessler/Fabric/pull/1368) by [vaygr](https://github.com/vaygr): Standardize sections for no repeat guidelines

- Standardize sections for no repeat guidelines

### Direct commits

- Moved system file to proper directory
- Added activity extractor
- Merge branch 'main' of github.com:danielmiessler/fabric

## v1.4.159 (2025-03-16)

### Direct commits

- Added flashcard generator.

## v1.4.158 (2025-03-16)

### PR [#1367](https://github.com/danielmiessler/Fabric/pull/1367) by [ksylvan](https://github.com/ksylvan): Remove Generic Type Parameters from StorageHandler Initialization

- Remove generic type parameters from NewStorageHandler calls
- Remove explicit type parameters from StorageHandler initialization
- Update contexts handler constructor implementation
- Update patterns handler constructor implementation
- Update sessions handler constructor implementation

## v1.4.157 (2025-03-16)

### PR [#1365](https://github.com/danielmiessler/Fabric/pull/1365) by [ksylvan](https://github.com/ksylvan): Implement Prompt Strategies in Fabric

- Add prompt strategies like Chain of Thought (CoT) with `--strategy` flag
- Implement `--liststrategies` command to view available strategies
- Support applying strategies to system prompts
- Improve README with platform-specific installation instructions
- Refactor git operations with new githelper package

### Direct commits

- Clipboard operations now work on Mac and PC
- Bump golang.org/x/net from 0.35.0 to 0.36.0 in go_modules group

## v1.4.156 (2025-03-11)

### PR [#1356](https://github.com/danielmiessler/Fabric/pull/1356) by [ksylvan](https://github.com/ksylvan): chore: add .vscode to `.gitignore` and fix typos and markdown linting  in `Alma.md`

- Add .vscode to `.gitignore` and fix typos and markdown linting in `Alma.md`

### PR [#1352](https://github.com/danielmiessler/Fabric/pull/1352) by [matmilbury](https://github.com/matmilbury): pattern_explanations.md: fix typo

- Pattern_explanations.md: fix typo

### PR [#1354](https://github.com/danielmiessler/Fabric/pull/1354) by [jmd1010](https://github.com/jmd1010): Fix Chat history window scrolling behavior

- Fix Chat history window sizing
- Update Web V2 Install Guide with improved instructions

## v1.4.155 (2025-03-09)

### PR [#1350](https://github.com/danielmiessler/Fabric/pull/1350) by [jmd1010](https://github.com/jmd1010): Implement Pattern Tile search functionality

- Implement Pattern Tile search functionality
- Implement  column resize functionnality

## v1.4.154 (2025-03-09)

### PR [#1349](https://github.com/danielmiessler/Fabric/pull/1349) by [ksylvan](https://github.com/ksylvan): Fix: v1.4.153 does not compile because of extra version declaration

- Remove unnecessary `version` variable from `main.go`
- Update Azure client API version access path in tests
- Implement column resize functionality
- Implement Pattern Tile search functionality

## v1.4.153 (2025-03-08)

### PR [#1348](https://github.com/danielmiessler/Fabric/pull/1348) by [liyuankui](https://github.com/liyuankui): feat: Add LiteLLM AI plugin support with local endpoint configuration

- Add LiteLLM AI plugin support with local endpoint configuration

## v1.4.152 (2025-03-07)

### Direct commits

- Fix pipe handling

## v1.4.151 (2025-03-07)

### PR [#1339](https://github.com/danielmiessler/Fabric/pull/1339) by [Eckii24](https://github.com/Eckii24): Feature/add azure api version

- Update azure.go
- Update azure_test.go
- Update openai.go

## v1.4.150 (2025-03-07)

### PR [#1343](https://github.com/danielmiessler/Fabric/pull/1343) by [jmd1010](https://github.com/jmd1010): Rename input.svelte to Input.svelte for proper component naming convention

- Rename input.svelte to Input.svelte for proper component naming convention

## v1.4.149 (2025-03-05)

### PR [#1340](https://github.com/danielmiessler/Fabric/pull/1340) by [ksylvan](https://github.com/ksylvan): Fix for youtube live links plus new youtube_summary pattern

- Update YouTube regex to support live URLs and timestamped transcripts
- Add argument validation and -t flag for transcript with timestamps
- Refactor PowerShell yt function with parameter switch
- Introduce youtube_summary pattern with documentation
- Update README to dynamically select transcript option

### PR [#1338](https://github.com/danielmiessler/Fabric/pull/1338) by [jmd1010](https://github.com/jmd1010): Update Web V2 Install Guide layout

- Update Web V2 Install Guide layout improvements

### PR [#1330](https://github.com/danielmiessler/Fabric/pull/1330) by [jmd1010](https://github.com/jmd1010): Fixed ALL CAP DIR as requested and processed minor updates to documentation

- Reorganize documentation with consistent directory naming

### PR [#1333](https://github.com/danielmiessler/Fabric/pull/1333) by [asasidh](https://github.com/asasidh): Update QUOTES section to include speaker names for clarity

- Update QUOTES section to include speaker names for clarity

### Direct commits

- Update azure.go, openai.go, and azure_test.go files

## v1.4.148 (2025-03-03)

## Summary of Changes

### Direct commits

- Rework LM Studio plugin
- Update QUOTES section to include speaker names for clarity
- Update Web V2 Install Guide with improved instructions V2
- Update Web V2 Install Guide with improved instructions
- Reorganize documentation with consistent directory naming and updated guides

## v1.4.147 (2025-02-28)

### PR [#1326](https://github.com/danielmiessler/Fabric/pull/1326) by [pavdmyt](https://github.com/pavdmyt): fix: continue fetching models even if some vendors fail

- Continue fetching models even if some vendors fail
- Remove cancellation of remaining goroutines when vendor collection fails
- Ensure other vendor collections continue even if one fails
- Fix model listing via `fabric -L` when localhost models are down
- Fix using non-default models via `fabric -m custom_model`

### PR [#1329](https://github.com/danielmiessler/Fabric/pull/1329) by [jmd1010](https://github.com/jmd1010): Svelte Web V2 Installation Guide

- Add Web V2 Installation Guide
- Update install guide with Plain Text instructions

## v1.4.146 (2025-02-27)

### PR [#1319](https://github.com/danielmiessler/Fabric/pull/1319) by [jmd1010](https://github.com/jmd1010): Enhancement: PDF to Markdown Conversion Functionality to the Web Svelte Chat Interface

- Add PDF to Markdown conversion functionality to the web svelte chat interface
- Add PDF to Markdown integration documentation
- Add Svelte implementation files for PDF integration
- Update README files directory structure and naming convention
- Add required UI image assets for feature implementation

## v1.4.145 (2025-02-26)

### PR [#1324](https://github.com/danielmiessler/Fabric/pull/1324) by [jaredmontoya](https://github.com/jaredmontoya): flake: fix/update and enhance

- Flake: fix/update

## v1.4.144 (2025-02-26)

### Direct commits

- Upgrade upload artifacts to v4

## v1.4.143 (2025-02-26)

### PR [#1264](https://github.com/danielmiessler/Fabric/pull/1264) by [eugeis](https://github.com/eugeis): feat: implement support for exolab

- Implement support for <https://github.com/exo-explore/exo>
- Merge branch 'main' into feat/exolab

## v1.4.142 (2025-02-25)

### Direct commits

- Build problems

## v1.4.141 (2025-02-25)

### PR [#1260](https://github.com/danielmiessler/Fabric/pull/1260) by [bluPhy](https://github.com/bluPhy): Fixing typo

- Fixed typos in codebase
- Updated version to v1.4.80
- Reverted previous v1.4.79 version update
- Merged changes from main branch
- Applied version control corrections

## v1.4.140 (2025-02-25)

### PR [#1313](https://github.com/danielmiessler/Fabric/pull/1313) by [cx-ken-swain](https://github.com/cx-ken-swain): Updated ollama.go to fix a couple of potential DoS issues

- Fixed security vulnerabilities in ollama.go to prevent potential DoS attacks
- Resolved multiple medium-severity vulnerabilities
- Updated application version to v..1
- Removed version.nix and version.go files
- Merged changes from main branch

## v1.4.139 (2025-02-25)

### PR [#1321](https://github.com/danielmiessler/Fabric/pull/1321) by [jmd1010](https://github.com/jmd1010): Update demo video link in PR-1309 documentation

- Updated demo video link in PR-1284 documentation
- Added complete PDF to Markdown conversion functionality
- Implemented Svelte integration files for PDF processing
- Added comprehensive PDF to Markdown documentation
- Enhanced web Svelte chat interface with PDF conversion capabilities

## v1.4.138 (2025-02-24)

### PR [#1317](https://github.com/danielmiessler/Fabric/pull/1317) by [ksylvan](https://github.com/ksylvan): chore: update Anthropic SDK and add Claude 3.7 Sonnet model support

- Updated anthropic-sdk-go from v0.2.0-alpha.4 to v0.2.0-alpha.11
- Added Claude 3.7 Sonnet models to available model list
- Added ModelClaude3_7SonnetLatest to model options
- Added ModelClaude3_7Sonnet20250219 to model options
- Removed ModelClaude_Instant_1_2 from available models

## v1.4.80 (2025-02-24)

### Direct commits

- Impl. multi-model / attachments, images

## v1.4.79 (2025-02-24)

### PR [#1257](https://github.com/danielmiessler/Fabric/pull/1257) by [jessefmoore](https://github.com/jessefmoore): Create analyze_threat_report_cmds

- Create system.md
Create pattern to extract commands from videos and threat reports to obtain commands so pentesters or red teams or Threat hunters can use to either threat hunt or simulate the threat actor.

### PR [#1256](https://github.com/danielmiessler/Fabric/pull/1256) by [JOduMonT](https://github.com/JOduMonT): Update README.md

- Update README.md

1. Windows Command: Because actually curl does not exist natively on Windows
2. Syntax: Because like this; it makes the "click, cut and paste" easier

### PR [#1247](https://github.com/danielmiessler/Fabric/pull/1247) by [kevnk](https://github.com/kevnk): Update suggest_pattern: refine summaries and add recently added patterns

- Update summaries and add recently added patterns

### PR [#1252](https://github.com/danielmiessler/Fabric/pull/1252) by [jeffmcjunkin](https://github.com/jeffmcjunkin): Update README.md: Add PowerShell aliases

- Update README.md: Add PowerShell aliases

### PR [#1253](https://github.com/danielmiessler/Fabric/pull/1253) by [abassel](https://github.com/abassel): Fixed few typos that I could find

- Fixed few typos that I could find

**Key Changes:**

- Added threat analysis pattern for extracting commands from security reports
- Improved Windows compatibility with PowerShell aliases and curl alternatives
- Enhanced pattern suggestions with updated summaries
- Added markdown callout pattern and prediction generator
- Implemented multi-model support with image attachments

## v1.4.137 (2025-02-24)

### PR [#1296](https://github.com/danielmiessler/Fabric/pull/1296) by [dependabot[bot]](https://github.com/apps/dependabot): Bump github.com/go-git/go-git/v5 from 5.12.0 to 5.13.0 in the go_modules group across 1 directory

- Updated github.com/go-git/go-git/v5 dependency from version 5.12.0 to 5.13.0
- Automated dependency update in go_modules group
- Direct production dependency upgrade
- Includes release notes and commit history links for transparency
- Signed-off by dependabot bot for automated maintenance

## v1.4.136 (2025-02-24)

### Direct commits

- Update to upload-artifact@v4 because upload-artifact@v3 is deprecated
- Merge branch 'danielmiessler:main' into main
- Update Anthropic SDK and add Claude 3.7 Sonnet model support
- Updated anthropic-sdk-go from v0.2.0-alpha.4 to v0.2.0-alpha.11
- Added Claude 3.7 Sonnet models to available model list

## v1.4.135 (2025-02-24)

### PR [#1309](https://github.com/danielmiessler/Fabric/pull/1309) by [jmd1010](https://github.com/jmd1010): Feature/Web Svelte GUI Enhancements: Pattern Descriptions, Tags, Favorites, Search Bar, Language Integration, PDF file conversion, etc

- Enhanced Web UI with pattern descriptions, tags, favorites, and search functionality
- Improved chat interface and pattern handling capabilities
- Updated dependencies and backup configuration
- Cleaned up sensitive files and improved .gitignore structure

### PR [#1312](https://github.com/danielmiessler/Fabric/pull/1312) by [junaid18183](https://github.com/junaid18183): Added Create LOE Document Prompt

- Added create_loe_document prompt for Level of Effort documentation

### PR [#1302](https://github.com/danielmiessler/Fabric/pull/1302) by [verebes1](https://github.com/verebes1): feat: Add LM Studio compatibility

- Added LM Studio as new plugin with base URL configuration
- Updated plugin registry to include LM Studio integration

### PR [#1297](https://github.com/danielmiessler/Fabric/pull/1297) by [Perchycs](https://github.com/Perchycs): Create pattern_explanations.md

- Created comprehensive pattern explanations with one-line summaries

### Direct commits

- Fixed security vulnerabilities in ollama.go and updated version to v1.1
- Added and updated extract_domains functionality

## v1.4.134 (2025-02-11)

### PR [#1289](https://github.com/danielmiessler/Fabric/pull/1289) by [thevops](https://github.com/thevops): Add the ability to grab YouTube video transcript with timestamps

- Added new `--transcript-with-timestamps` flag for YouTube video processing
- Timestamps formatted as HH:MM:SS and prepended to each transcript line
- Enables quick navigation to specific video segments in summaries
- Similar functionality to existing `--transcript` flag but with time markers
- Useful for creating timestamped video summaries and references

## v1.4.133 (2025-02-11)

### PR [#1294](https://github.com/danielmiessler/Fabric/pull/1294) by [TvisharajiK](https://github.com/TvisharajiK): Improved unit-test coverage from 0 to 100 (AI module) using Keploy's agent

- Increased unit test coverage from 0 to 100% in the AI module using Keploy's Agent
- Added YouTube video transcript with timestamps feature via `--transcript-with-timestamps` flag
- Bumped github.com/go-git/go-git/v5 from 5.12.0 to 5.13.0
- Added multiple new TELOS patterns and challenge handling pattern
- Added panel topic extractor and intro sentences pattern

## v1.4.132 (2025-02-02)

### PR [#1278](https://github.com/danielmiessler/Fabric/pull/1278) by [aicharles](https://github.com/aicharles): feat(anthropic): enable custom API base URL support

- Enable custom API base URL configuration for Anthropic integration
- Add proper handling of v1 endpoint for UUID-containing URLs
- Implement URL formatting logic for consistent endpoint structure
- Clean up commented code and improve configuration flow
- Enhance API flexibility for different deployment environments

## v1.4.131 (2025-01-30)

### PR [#1270](https://github.com/danielmiessler/Fabric/pull/1270) by [wmahfoudh](https://github.com/wmahfoudh): Added output filename support for to_pdf

- Added output filename support for to_pdf

### PR [#1271](https://github.com/danielmiessler/Fabric/pull/1271) by [wmahfoudh](https://github.com/wmahfoudh): Adding deepseek support

- Added Deepseek AI integration

### PR [#1258](https://github.com/danielmiessler/Fabric/pull/1258) by [tuergeist](https://github.com/tuergeist): Minor README fix and additional Example

- Doc: Custom patterns also work with Claude models
- Doc: Add scrape URL example. Fix Example 4

### Direct commits

- Implement support for <https://github.com/exo-explore/exo>

## v1.4.130 (2025-01-03)

### PR [#1240](https://github.com/danielmiessler/Fabric/pull/1240) by [johnconnor-sec](https://github.com/johnconnor-sec): Updates: ./web

- Moved pattern loader to ModelConfig and added page fly transitions
- Updated UI components and improved responsive layout for chat interface
- Added NotesDrawer component that saves notes to lib/content/inbox
- Centered chat and NotesDrawer in viewport for better user experience
- Restructured project organization: moved types to lib/interfaces and lib/api

## v1.4.129 (2025-01-03)

### PR [#1242](https://github.com/danielmiessler/Fabric/pull/1242) by [CuriouslyCory](https://github.com/CuriouslyCory): Adding youtube --metadata flag

- Added metadata lookup to youtube helper
- Better metadata

### PR [#1230](https://github.com/danielmiessler/Fabric/pull/1230) by [iqbalabd](https://github.com/iqbalabd): Update translate pattern to use curly braces

- Update translate pattern to use curly braces

### Direct commits

- Chat and NotesDrawer now centered in viewport
- Enhanced enrich pattern and added enrich_blog_post
- Major file reorganization: moved types to lib/interfaces, components restructured
- Updated Post page styling and NotesDrawer saves to lib/content/inbox
- Version updates to v..1 with corresponding .nix and .go files

## v1.4.128 (2024-12-26)

### PR [#1227](https://github.com/danielmiessler/Fabric/pull/1227) by [mattjoyce](https://github.com/mattjoyce): Feature/template extensions

- Implemented stdout template extensions with path-based registry storage and hash verification
- Added file-based output handling with proper cleanup of temporary files for local and remote operations
- Fixed pattern file usage without stdin by initializing empty message for template processing
- Added extension manager tests, registration, execution validation, and example files with tutorial
- Enhanced extension listing with better error messages when hash verification fails

### Direct commits

- Updated story format to shorter bullets and improved notes drawer with rocket theme
- Updated POSTS for main 24-12-08 release and fixed import statements

## v1.4.127 (2024-12-23)

### PR [#1218](https://github.com/danielmiessler/Fabric/pull/1218) by [sosacrazy126](https://github.com/sosacrazy126): streamlit ui

- Added comprehensive Streamlit application for managing and executing patterns
- Implemented pattern creation, execution, and analysis with advanced UI components
- Enhanced logging configuration with color-coded console and detailed file handlers
- Added pattern chain execution functionality for sequential pattern processing
- Integrated output management with starring/favoriting and persistent storage features

### PR [#1225](https://github.com/danielmiessler/Fabric/pull/1225) by [wmahfoudh](https://github.com/wmahfoudh): Added Humanize Pattern

- Added Humanize Pattern

## v1.4.126 (2024-12-22)

### PR [#1212](https://github.com/danielmiessler/Fabric/pull/1212) by [wrochow](https://github.com/wrochow): Significant updates to Duke and Socrates

- Significant thematic rewrite incorporating classical philosophical texts
- Ingested 8 key documents including Plato's Apology, Phaedrus, Symposium, and Republic
- Added Xenophon's works: The Economist, Memorabilia, Memorable Thoughts, and Symposium
- Enhanced with specific steps for research, analysis, and code reviews
- Version updates and branch merging activities

## v1.4.125 (2024-12-22)

### PR [#1222](https://github.com/danielmiessler/Fabric/pull/1222) by [wmahfoudh](https://github.com/wmahfoudh): Fix cross-filesystem file move in to_pdf plugin (issue 1221)

- Fix cross-filesystem file move in to_pdf plugin (issue 1221)

### Direct commits

- Update version to v..1 and commit
- Don't quite know how I screwed this up, I wasn't even working there.
- Update version to v..1 and commit
- Merge branch 'main' into main
- Merge branch 'main' into main

## v1.4.124 (2024-12-21)

### PR [#1215](https://github.com/danielmiessler/Fabric/pull/1215) by [infosecwatchman](https://github.com/infosecwatchman): Add Endpoints to facilitate Ollama based chats

- Add endpoints to facilitate Ollama based chats
- Built to use with Open WebUI

### PR [#1214](https://github.com/danielmiessler/Fabric/pull/1214) by [iliaross](https://github.com/iliaross): Fix the typo in the sentence

- Fix typo in sentence

### PR [#1213](https://github.com/danielmiessler/Fabric/pull/1213) by [AnirudhG07](https://github.com/AnirudhG07): Spelling Fixes

- Spelling fixes in patterns and README
- Spelling fixes in create_quiz pattern

### Direct commits

- Enhanced pattern management with improved creation, editing, and deletion
- Improved logging configuration and error handling for better debugging

## v1.4.123 (2024-12-20)

### PR [#1208](https://github.com/danielmiessler/Fabric/pull/1208) by [mattjoyce](https://github.com/mattjoyce): Fix: Issue with the custom message and added example config file

- Fixed custom message issue and added example config file

### Direct commits

- Added Streamlit application for pattern management and execution with logging, session state, and UI components
- Added Ollama chat endpoints for Open WebUI integration
- Fixed multiple spelling errors across patterns and documentation
- Updated version to v1.1 with significant research and analysis workflow improvements
- Major Socrates pattern rewrite incorporating Plato and Xenophon source materials from Project Gutenberg

## v1.4.122 (2024-12-14)

### PR [#1201](https://github.com/danielmiessler/Fabric/pull/1201) by [mattjoyce](https://github.com/mattjoyce): feat: Add YAML configuration support

- Add YAML configuration support for persistent settings
- Add --config flag for specifying YAML config file path
- Support standard option precedence (CLI > YAML > defaults)
- Add type-safe YAML parsing with reflection
- Add tests for YAML config functionality

## v1.4.121 (2024-12-13)

### PR [#1200](https://github.com/danielmiessler/Fabric/pull/1200) by [mattjoyce](https://github.com/mattjoyce): Fix: Mask input token to prevent var substitution in patterns

- Mask input token to prevent var substitution in patterns

### Direct commits

- Added new instruction trick.

## v1.4.120 (2024-12-10)

### PR [#1189](https://github.com/danielmiessler/Fabric/pull/1189) by [mattjoyce](https://github.com/mattjoyce): Add --input-has-vars flag to control variable substitution in input

- Add --input-has-vars flag to control variable substitution in input
- Add InputHasVars field to ChatRequest struct
- Only process template variables in user input when flag is set
- Fixes issue with Ansible/Jekyll templates that use {{var}} syntax
- Makes template variable substitution opt-in, preserving literal curly braces by default

### PR [#1182](https://github.com/danielmiessler/Fabric/pull/1182) by [jessefmoore](https://github.com/jessefmoore): analyze_risk pattern

- Created analyze_risk pattern for 3rd party vendor risk analysis

## v1.4.119 (2024-12-07)

### PR [#1181](https://github.com/danielmiessler/Fabric/pull/1181) by [mattjoyce](https://github.com/mattjoyce): Bugfix/1169 symlinks

- Fix #1169: Add robust handling for paths and symlinks in GetAbsolutePath
- Update version to v..1 and commit
- Revert "Update version to v..1 and commit"

### Direct commits

- Added tutorial and example files
- Add cards component and update packages, main page, styles
- Check extension names don't have spaces
- Added test pattern

## v1.4.118 (2024-12-05)

### PR [#1174](https://github.com/danielmiessler/Fabric/pull/1174) by [mattjoyce](https://github.com/mattjoyce): Curly brace templates

- Fixed pattern file usage without stdin to prevent segfault by initializing empty message
- Removed redundant template processing of message content
- Simplified template processing flow for both stdin and non-stdin use cases
- Added proper handling for patterns with variables but no stdin input

### PR [#1179](https://github.com/danielmiessler/Fabric/pull/1179) by [sluosapher](https://github.com/sluosapher): added a new pattern create_newsletter_entry

- Added new pattern create_newsletter_entry

## v1.4.117 (2024-11-30)

### Direct commits

- Close #1173

## v1.4.116 (2024-11-28)

### Direct commits

- Cleanup style

## v1.4.115 (2024-11-28)

### PR [#1168](https://github.com/danielmiessler/Fabric/pull/1168) by [johnconnor-sec](https://github.com/johnconnor-sec): Update README.md

- Updated README.md documentation
- Cleaned up code styling
- Enhanced message handling to use custom messages with piped input
- Improved overall project documentation and formatting
- Streamlined user experience with better message processing

## v1.4.114 (2024-11-26)

### PR [#1164](https://github.com/danielmiessler/Fabric/pull/1164) by [MegaGrindStone](https://github.com/MegaGrindStone): fix: provide default message content to avoid nil pointer dereference

- Provide default message content to avoid nil pointer dereference

## v1.4.113 (2024-11-26)

### PR [#1166](https://github.com/danielmiessler/Fabric/pull/1166) by [dependabot[bot]](https://github.com/apps/dependabot): build(deps-dev): bump @sveltejs/kit from 2.6.1 to 2.8.4 in /web in the npm_and_yarn group across 1 directory

- Updated @sveltejs/kit from version 2.6.1 to 2.8.4 in the /web directory
- Dependency update in the npm_and_yarn group for development dependencies
- Automated security and maintenance update by dependabot
- Direct development dependency type update
- Includes release notes and changelog references for the SvelteKit framework

## v1.4.112 (2024-11-26)

### PR [#1165](https://github.com/danielmiessler/Fabric/pull/1165) by [johnconnor-sec](https://github.com/johnconnor-sec): feat: Fabric Web UI

- Update version to v..1 and commit
- Update Obsidian.md documentation
- Update README.md file
- Multiple commits by John on 2024-11-26

### Direct commits

- Provide default message content to avoid nil pointer dereference

## v1.4.111 (2024-11-26)

### Direct commits

- Integrate code formating

## v1.4.110 (2024-11-26)

### PR [#1135](https://github.com/danielmiessler/Fabric/pull/1135) by [mrtnrdl](https://github.com/mrtnrdl): Add `extract_recipe`

- Update version to v..1 and commit
- Add extract_recipe to easily extract the necessary information from cooking-videos
- Merge branch 'main' into main

## v1.4.109 (2024-11-24)

### PR [#1157](https://github.com/danielmiessler/Fabric/pull/1157) by [mattjoyce](https://github.com/mattjoyce): fix: process template variables in raw input

- Process template variables ({{var}}) consistently in both pattern files and raw input messages
- Add template variable processing for raw input in BuildSession with explicit messageContent initialization
- Remove errantly committed build artifact (fabric binary)
- Fix template.go to handle missing variables in stdin input with proper error messages
- Fix raw mode doubling user input by streamlining context staging

### Direct commits

- Added analyze_mistakes

## v1.4.108 (2024-11-21)

### PR [#1155](https://github.com/danielmiessler/Fabric/pull/1155) by [mattjoyce](https://github.com/mattjoyce): Curly brace templates and plugins

- Introduced template package for variable substitution with {{variable}} syntax
- Moved substitution logic from patterns to centralized template system
- Updated patterns.go and chatter.go to use new template package
- Added support for special {{input}} handling and nested variables
- Implemented core plugin system with utility plugins (datetime, fetch, file, sys, text)

## v1.4.107 (2024-11-19)

### PR [#1149](https://github.com/danielmiessler/Fabric/pull/1149) by [mathisto](https://github.com/mathisto): Fix typo in md_callout

- Fixed typo in md_callout pattern
- Updated patterns zip workflow
- Removed patterns zip workflow

## v1.4.106 (2024-11-19)

### Direct commits

- Migrate to official anthropics Go SDK

## v1.4.105 (2024-11-19)

### PR [#1147](https://github.com/danielmiessler/Fabric/pull/1147) by [mattjoyce](https://github.com/mattjoyce): refactor: unify pattern loading and variable handling

- Unify pattern loading and variable handling
- Stronger separation of concerns between chatter.go and patterns.go
- Consolidate pattern loading logic into GetPattern method
- Support both file and database patterns through single interface
- Handle variable substitution in one place

### PR [#1146](https://github.com/danielmiessler/Fabric/pull/1146) by [mrwadams](https://github.com/mrwadams): Add summarize_meeting

- Add summarize_meeting pattern for creating meeting summaries from audio transcripts
- Outputs Key Points, Tasks, Decisions, and Next Steps sections
- Provides structured format for meeting documentation
- Supports audio transcript processing
- Enables organized meeting follow-up workflow

### Direct commits

- Introduce template package for variable substitution with {{variable}} syntax
- Move substitution logic from patterns to centralized template system
- Support special {{input}} handling for pattern content
- Enable multiple passes to handle nested variables
- Report errors for missing required variables

## v1.4.104 (2024-11-18)

### PR [#1142](https://github.com/danielmiessler/Fabric/pull/1142) by [mattjoyce](https://github.com/mattjoyce): feat: add file-based pattern support

- Add file-based pattern support for loading patterns directly from files using path prefixes
- Supports relative paths (./pattern.txt), home directory expansion (~/patterns/test.txt), and absolute paths
- Maintains backwards compatibility with named patterns
- Requires explicit path markers to distinguish from pattern names
- Example usage: `fabric --pattern ./draft-pattern.txt` or `fabric --pattern ~/patterns/my-pattern.txt`

### Direct commits

- Add summarize_meeting pattern to create meeting summaries from audio transcripts with sections for Key Points, Tasks, Decisions, and Next Steps

## v1.4.103 (2024-11-18)

### PR [#1133](https://github.com/danielmiessler/Fabric/pull/1133) by [igophper](https://github.com/igophper): fix: fix default gin

- Fix default gin configuration

### PR [#1129](https://github.com/danielmiessler/Fabric/pull/1129) by [xyb](https://github.com/xyb): add a screenshot of fabric

- Add a screenshot of fabric to documentation

## v1.4.102 (2024-11-18)

### PR [#1143](https://github.com/danielmiessler/Fabric/pull/1143) by [mariozig](https://github.com/mariozig): Update docker image

- Update docker image

### Direct commits

- Add file-based pattern support for loading patterns directly from files using explicit path prefixes
- Support relative paths (./pattern.txt, ../pattern.txt)
- Support home directory expansion (~/patterns/test.txt)
- Support absolute paths with backwards compatibility for named patterns
- Require explicit path markers to distinguish from pattern names

## v1.4.101 (2024-11-15)

### Direct commits

- Improve logging for missing setup steps
- Add extract_recipe to easily extract the necessary information from cooking-videos
- Fix default gin
- Update version to v..1 and commit
- Add a screenshot of fabric

## v1.4.100 (2024-11-13)

### Direct commits

- Added our first formal stitch
- Upgraded AI result rater (multiple iterations)

## v1.4.99 (2024-11-10)

### PR [#1126](https://github.com/danielmiessler/Fabric/pull/1126) by [jaredmontoya](https://github.com/jaredmontoya): flake: add gomod2nix auto-update

- Added gomod2nix auto-update functionality to flake

### Direct commits

- Upgraded AI result rater (multiple iterations)

## v1.4.98 (2024-11-09)

### Direct commits

- Zip patterns

## v1.4.97 (2024-11-09)

### Direct commits

- Update dependencies; improve vendors setup/default model

## v1.4.96 (2024-11-09)

### PR [#1060](https://github.com/danielmiessler/Fabric/pull/1060) by [noamsiegel](https://github.com/noamsiegel): Analyze Candidates Pattern

- Added system and user prompts

### Direct commits

- Add claude-3-5-haiku-latest model

## v1.4.95 (2024-11-09)

### PR [#1123](https://github.com/danielmiessler/Fabric/pull/1123) by [polyglotdev](https://github.com/polyglotdev): :sparkles: Added unaliasing to pattern setup

- Added unaliasing step to pattern setup process
- Prevents conflicts between dynamically defined functions and existing aliases

### PR [#1119](https://github.com/danielmiessler/Fabric/pull/1119) by [verebes1](https://github.com/verebes1): Add auto save functionality

- Added auto save functionality to aliases
- Updated README with autogenerating aliases information for Obsidian integration
- Updated table of contents

### Direct commits

- Merged main branch updates
- Updated README documentation

## v1.4.94 (2024-11-06)

### PR [#1108](https://github.com/danielmiessler/Fabric/pull/1108) by [butterflyx](https://github.com/butterflyx): [add] RegEx for YT shorts

- Added VideoID support for YouTube shorts
- Merged branch 'main' into fix/yt-shorts

### PR [#1117](https://github.com/danielmiessler/Fabric/pull/1117) by [verebes1](https://github.com/verebes1): Add alias generation information

- Added alias generation documentation to README
- Updated table of contents
- Included information about generating aliases for prompts including YouTube transcripts
- Merged branch 'main' into add-aliases-for-patterns

### PR [#1115](https://github.com/danielmiessler/Fabric/pull/1115) by [ignacio-arce](https://github.com/ignacio-arce): Added create_diy

- Added create_diy functionality

## v1.4.93 (2024-11-06)

### Direct commits

- Short YouTube url pattern
- Add alias generation information
- Updated the readme with information about generating aliases for each prompt including one for youtube transcripts
- Updated the table of contents
- Added create_diy
- [add] VideoID for YT shorts

## v1.4.92 (2024-11-05)

### PR [#1109](https://github.com/danielmiessler/Fabric/pull/1109) by [leonsgithub](https://github.com/leonsgithub): Add docker

- Add docker

## v1.4.91 (2024-11-05)

### Direct commits

- Bufio.Scanner message too long
- Add docker

## v1.4.90 (2024-11-04)

### Direct commits

- Impl. Youtube PlayList support
- Close #1103, Update Readme hpt to install to_pdf

## v1.4.89 (2024-11-04)

### PR [#1102](https://github.com/danielmiessler/Fabric/pull/1102) by [jholsgrove](https://github.com/jholsgrove): Create user story pattern

- Create user story pattern

### Direct commits

- Close #1106, fix pipe reading
- YouTube PlayList support

## v1.4.88 (2024-10-30)

### PR [#1098](https://github.com/danielmiessler/Fabric/pull/1098) by [jaredmontoya](https://github.com/jaredmontoya): Fix nix package update workflow

- Fix nix package version auto update workflow

## v1.4.87 (2024-10-30)

### PR [#1096](https://github.com/danielmiessler/Fabric/pull/1096) by [jaredmontoya](https://github.com/jaredmontoya): Implement automated ci nix package version update

- Modularize nix flake
- Automate nix package version update

## v1.4.86 (2024-10-30)

### PR [#1088](https://github.com/danielmiessler/Fabric/pull/1088) by [jaredmontoya](https://github.com/jaredmontoya): feat: add DEFAULT_CONTEXT_LENGTH setting

- Add model context length setting

## v1.4.85 (2024-10-30)

### Direct commits

- Write tools output also to output file if defined; fix XouTube transcript &#39; character

## v1.4.84 (2024-10-30)

### Direct commits

- Deactivate build triggering at changes of patterns or docu

## v1.4.83 (2024-10-30)

### PR [#1089](https://github.com/danielmiessler/Fabric/pull/1089) by [jaredmontoya](https://github.com/jaredmontoya): Introduce Nix to the project

- Add trailing newline
- Add Nix Flake

## v1.4.82 (2024-10-30)

### PR [#1094](https://github.com/danielmiessler/Fabric/pull/1094) by [joshmedeski](https://github.com/joshmedeski): feat: add md_callout pattern

- Add md_callout pattern
Add a pattern that can convert text into an appropriate markdown callout

## v1.4.81 (2024-10-29)

### Direct commits

- Split tools messages from use message

## v1.4.78 (2024-10-28)

### PR [#1059](https://github.com/danielmiessler/Fabric/pull/1059) by [noamsiegel](https://github.com/noamsiegel): Analyze Proposition Pattern

- Added system and user prompts

## v1.4.77 (2024-10-28)

### PR [#1073](https://github.com/danielmiessler/Fabric/pull/1073) by [mattjoyce](https://github.com/mattjoyce): Five patterns to explore a project, opportunity or brief

- Added five new DSRP patterns for project exploration and strategic analysis
- Enhanced prompts to increase divergent thinking capabilities
- Implemented identify_job_stories pattern for user story identification
- Created S7 Strategy profiling with headwinds and tailwinds analysis
- Added comprehensive metadata and style guide documentation

### Direct commits

- Add Nix Flake

## v1.4.76 (2024-10-28)

### Direct commits

- Simplify isChatRequest

## v1.4.75 (2024-10-28)

### PR [#1090](https://github.com/danielmiessler/Fabric/pull/1090) by [wrochow](https://github.com/wrochow): A couple of patterns

- Added "Dialog with Socrates" pattern for deep philosophical conversations with a modern-day philosopher
- Added "Ask uncle Duke" pattern for Java development expertise, specializing in Spring Framework and Maven
- Fixed formatting by adding trailing newline to files

### Direct commits

- Add trailing newline

## v1.4.74 (2024-10-27)

### PR [#1077](https://github.com/danielmiessler/Fabric/pull/1077) by [xvnpw](https://github.com/xvnpw): feat: add pattern refine_design_document

- Add pattern refine_design_document

## v1.4.73 (2024-10-27)

### PR [#1086](https://github.com/danielmiessler/Fabric/pull/1086) by [NuCl34R](https://github.com/NuCl34R): Create a basic translator pattern, edit file to add desired language

- Create system.md

### Direct commits

- Added metadata and styleguide
- Added structure to prompt
- Added headwinds and tailwinds
- Initial draft of s7 Strategy profiling
- Merge operations from main branch

## v1.4.72 (2024-10-25)

### PR [#1070](https://github.com/danielmiessler/Fabric/pull/1070) by [xvnpw](https://github.com/xvnpw): feat: create create_design_document pattern

- Create create_design_document pattern

## v1.4.71 (2024-10-25)

### PR [#1072](https://github.com/danielmiessler/Fabric/pull/1072) by [xvnpw](https://github.com/xvnpw): feat: add review_design pattern

- Add review_design pattern

## v1.4.70 (2024-10-25)

### PR [#1064](https://github.com/danielmiessler/Fabric/pull/1064) by [rprouse](https://github.com/rprouse): Update README.md with pbpaste section

- Update README.md with pbpaste section

### Direct commits

- Add pattern refine_design_document
- Added identify_job_stories
- Add review_design pattern
- Create create_design_document pattern
- Added system and user prompts

## v1.4.69 (2024-10-21)

### Direct commits

- Updated the Alma.md file.

## v1.4.68 (2024-10-21)

### Direct commits

- Setup does not overwrites old values

## v1.4.67 (2024-10-19)

### Direct commits

- Merge remote-tracking branch 'origin/main'
- Plugins arch., new setup procedure

## v1.4.66 (2024-10-19)

### Direct commits

- Plugins arch., new setup procedure

## v1.4.65 (2024-10-16)

### PR [#1045](https://github.com/danielmiessler/Fabric/pull/1045) by [Fenicio](https://github.com/Fenicio): Update patterns/analyze_answers/system.md - Fixed a bunch of typos

- Update patterns/analyze_answers/system.md - Fixed a bunch of typos

## v1.4.64 (2024-10-14)

### Direct commits

- Updated readme

## v1.4.63 (2024-10-13)

### PR [#862](https://github.com/danielmiessler/Fabric/pull/862) by [Thepathakarpit](https://github.com/Thepathakarpit): Create setup_fabric.bat, a batch script to automate setup and running…

- Created setup_fabric.bat batch script to automate Fabric setup and execution on Windows
- Merged branch 'main' into patch-1

## v1.4.62 (2024-10-13)

### PR [#1044](https://github.com/danielmiessler/Fabric/pull/1044) by [eugeis](https://github.com/eugeis): Feat/rest api

- Work on Rest API
- Restructure for better reuse
- Merge branch 'main' into feat/rest-api

## v1.4.61 (2024-10-13)

### Direct commits

- Updated extract sponsors.
- Merge branch 'main' into feat/rest-api
- Restructure for better reuse
- Restructure for better reuse
- Restructure for better reuse

## v1.4.60 (2024-10-12)

### Direct commits

- IsChatRequest rule; Close #1042 is

## v1.4.59 (2024-10-11)

### Direct commits

- Added ctw to Raycast.

## v1.4.58 (2024-10-11)

### Direct commits

- We don't need tp configure DryRun vendor
- Close #1040. Configure vendors separately that were not configured yet

## v1.4.57 (2024-10-11)

### Direct commits

- Close #1035, provide better example for pattern variables

## v1.4.56 (2024-10-11)

### PR [#1039](https://github.com/danielmiessler/Fabric/pull/1039) by [hallelujah-shih](https://github.com/hallelujah-shih): Feature/set default lang

- Support set default output language
- Updated cli/cli.go
- Modified core/fabric.go with formatting changes

### Direct commits

- Updated all dsrp prompts to increase divergent thinking
- Fixed mix up with system
- Initial dsrp prompts

## v1.4.55 (2024-10-09)

### Direct commits

- Close #1036

## v1.4.54 (2024-10-07)

### PR [#1021](https://github.com/danielmiessler/Fabric/pull/1021) by [joshuafuller](https://github.com/joshuafuller): Corrected spelling and grammatical errors for consistency and clarity for transcribe_minutes

- Fixed grammatical accuracy: "agreed within" → "agreed upon within"
- Added missing periods for consistency across list items
- Corrected spelling: "highliting" → "highlighting"
- Fixed spelling: "exxactly" → "exactly"
- Updated phrasing: "Write NEXT STEPS a 2-3 sentences" → "Write NEXT STEPS as 2-3 sentences"

## v1.4.53 (2024-10-07)

### Direct commits

- Fix NP if response is empty, close #1026, #1027

## v1.4.52 (2024-10-06)

### Direct commits

- Merge branch 'main' of github.com:danielmiessler/fabric
- Added extract_core_message functionality
- Extensive work on Rest API development and implementation
- Corrected spelling and grammatical errors for consistency and clarity:
- Fixed "agreed within" to "agreed upon within"
- Added missing periods for consistency
- Corrected "highliting" to "highlighting"
- Fixed "exxactly" to "exactly"
- Updated "Write NEXT STEPS a 2-3 sentences" to "Write NEXT STEPS as 2-3 sentences"

These changes improve document readability and API functionality.

## v1.4.51 (2024-10-05)

### Direct commits

- Tests

## v1.4.50 (2024-10-05)

### Direct commits

- Windows release

## v1.4.49 (2024-10-05)

### Direct commits

- Windows release

## v1.4.48 (2024-10-05)

### Direct commits

- Add 'meta' role to store meta info to session, like source of input content.

## v1.4.47 (2024-10-05)

### Direct commits

- Add 'meta' role to store meta info to session, like source of input content.
- Add 'meta' role to store meta info to session, like source of input content.

## v1.4.46 (2024-10-04)

### Direct commits

- Close #1018
- Implement print session and context
- Implement print session and context

## v1.4.45 (2024-10-04)

### Direct commits

- Setup for specific vendor, e.g. --setup-vendor=OpenAI

## v1.4.44 (2024-10-03)

### Direct commits

- Use the latest tag by date

## v1.4.43 (2024-10-03)

### Direct commits

- Use the latest tag by date

## v1.4.42 (2024-10-03)

### Direct commits

- Use the latest tag by date
- Use the latest tag by date

## v1.4.41 (2024-10-03)

### Direct commits

- Trigger release workflow ony tag_created

## v1.4.40 (2024-10-03)

### Direct commits

- Create repo dispatch

## v1.4.39 (2024-10-03)

### Direct commits

- Test tag creation

## v1.4.38 (2024-10-03)

### Direct commits

- Test tag creation
- Commit version changes only if it changed
- Use TAG_PAT instead of secrets.GITHUB_TOKEN for tag push
- Merge branch 'main' of github.com:danielmiessler/fabric
- Updated predictions pattern

## v1.4.36 (2024-10-03)

### Direct commits

- Merge branch 'main' of github.com:danielmiessler/fabric
- Added redeeming thing.

## v1.4.35 (2024-10-02)

### Direct commits

- Clean up html readability; add autm. tag creation

## v1.4.34 (2024-10-02)

### Direct commits

- Clean up html readability; add autm. tag creation

## v1.4.33 (2024-10-02)

### Direct commits

- Clean up html readability; add autm. tag creation
- Clean up html readability; add autm. tag creation
- Clean up html readability; add autm. tag creation

## v1.5.0 (2024-10-02)

### Direct commits

- Clean up html readability; add autm. tag creation

## v1.4.32 (2024-10-02)

### PR [#1007](https://github.com/danielmiessler/Fabric/pull/1007) by [hallelujah-shih](https://github.com/hallelujah-shih): support turn any web page into clean view content

- Support turn any web page into clean view content

### PR [#1005](https://github.com/danielmiessler/Fabric/pull/1005) by [fn5](https://github.com/fn5): Update patterns/solve_with_cot/system.md typos

- Fixed multiple typos in solve_with_cot pattern
- Corrected opening/closing brackets format

### PR [#962](https://github.com/danielmiessler/Fabric/pull/962) by [alucarded](https://github.com/alucarded): Update prompt in agility_story

- Updated agility_story pattern to make topic more logical

### PR [#994](https://github.com/danielmiessler/Fabric/pull/994) by [OddDuck11](https://github.com/OddDuck11): Add pattern analyze_military_strategy

- Added new pattern for analyzing historic or fictional battle strategies

### PR [#1008](https://github.com/danielmiessler/Fabric/pull/1008) by [MattBash17](https://github.com/MattBash17): Update system.md in transcribe_minutes

- Updated transcribe_minutes pattern system file

## v1.4.31 (2024-10-01)

### PR [#987](https://github.com/danielmiessler/Fabric/pull/987) by [joshmedeski](https://github.com/joshmedeski): feat: remove cli list label and indentation

- Remove cli list label and indentation

### PR [#1011](https://github.com/danielmiessler/Fabric/pull/1011) by [fooman[org]](https://github.com/fooman): Grab transcript from youtube matching the user's language

- Grab transcript from youtube matching the user's language instead of the first one

### Direct commits

- Added epp pattern
- Added create_story_explanation pattern
- Add version updater bot
- Update system.md in transcribe_minutes
- Support turn any web page into clean view content

## v1.4.30 (2024-09-29)

### Direct commits

- Add version updater bot

## v1.4.29 (2024-09-29)

### PR [#996](https://github.com/danielmiessler/Fabric/pull/996) by [hallelujah-shih](https://github.com/hallelujah-shih): add wipe flag for ctx and session

- Add wipe flag for ctx and session

### PR [#984](https://github.com/danielmiessler/Fabric/pull/984) by [riccardo1980](https://github.com/riccardo1980): adding flag for pinning seed in openai and compatible APIs

- Adding flag for pinning seed in openai and compatible APIs
- Updating README with the new flag

### PR [#991](https://github.com/danielmiessler/Fabric/pull/991) by [aculich](https://github.com/aculich): Fix GOROOT path for Apple Silicon Macs

- Fix GOROOT path for Apple Silicon Macs in setup instructions
- Updated instructions to dynamically determine correct GOROOT path using Homebrew

### PR [#861](https://github.com/danielmiessler/Fabric/pull/861) by [noamsiegel](https://github.com/noamsiegel): Scrape url

- Add ScrapeURL flag for CLI to scrape website URL to markdown using Jina AI
- Add Jina AI integration for web scraping and question search
- Made jina api key optional

### PR [#970](https://github.com/danielmiessler/Fabric/pull/970) by [mark-kazakov](https://github.com/mark-kazakov): add mistral vendor

- Add mistral vendor
