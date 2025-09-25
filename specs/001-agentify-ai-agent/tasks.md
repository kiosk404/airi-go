# Tasks: Agentify - AI Agent 开发调试平台

**Input**: Design documents from `/specs/001-agentify-ai-agent/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `src/`, `tests/` at repository root
- **Web app**: `backend/src/`, `frontend/src/`
- **Mobile**: `api/src/`, `ios/src/` or `android/src/`
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 3.1: Setup
- [X] T001 Create project structure per implementation plan in frontend/src/
- [X] T002 Initialize React + TypeScript project with rsbuild, Ant Design, Tailwind CSS dependencies
- [X] T003 [P] Configure linting and formatting tools (ESLint, Prettier)
- [X] T004 [P] Configure Tailwind CSS with Ant Design compatibility
- [X] T005 [P] Set up MSW for local JSON Mock data simulation
- [X] T006 [P] Configure Jest and React Testing Library for testing

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests
- [X] T007 [P] Contract test for Agent Management API in frontend/tests/contract/test_agent_management.py
- [X] T008 [P] Contract test for Workspace Data API in frontend/tests/contract/test_workspace_data.py
- [X] T009 [P] Contract test for Plugin Management API in frontend/tests/contract/test_plugin_management.py
- [X] T010 [P] Contract test for Workflow Management API in frontend/tests/contract/test_workflow_management.py

### Integration Tests (based on user stories)
- [X] T011 [P] Integration test for Agent creation in frontend/tests/integration/test_agent_creation.py
- [X] T012 [P] Integration test for Prompt editing and chatting in frontend/tests/integration/test_prompt_chat.py
- [X] T013 [P] Integration test for Plugin and Workflow creation in frontend/tests/integration/test_plugin_workflow_creation.py
- [X] T014 [P] Integration test for Analytics data display in frontend/tests/integration/test_analytics_display.py

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Data Models (based on data-model.md)
- [ ] T015 [P] AI Agent model in frontend/src/models/Agent.ts
- [ ] T016 [P] Plugin model in frontend/src/models/Plugin.ts
- [ ] T017 [P] Workflow model in frontend/src/models/Workflow.ts
- [ ] T018 [P] Prompt model in frontend/src/models/Prompt.ts
- [ ] T019 [P] AnalyticsData model in frontend/src/models/AnalyticsData.ts
- [ ] T020 [P] ChatMessage model in frontend/src/models/ChatMessage.ts

### Services
- [ ] T021 [P] AgentService for CRUD operations in frontend/src/services/AgentService.ts
- [ ] T022 [P] PluginService for plugin management in frontend/src/services/PluginService.ts
- [ ] T023 [P] WorkflowService for workflow management in frontend/src/services/WorkflowService.ts
- [ ] T024 [P] ChatService for chat interactions in frontend/src/services/ChatService.ts
- [ ] T025 [P] AnalyticsService for data retrieval in frontend/src/services/AnalyticsService.ts

### UI Components
- [ ] T026 [P] Develop workspace components in frontend/src/components/develop/
- [ ] T027 [P] Playground workspace components in frontend/src/components/playground/
- [ ] T028 [P] Explore workspace components in frontend/src/components/explore/
- [ ] T029 [P] Shared UI components in frontend/src/components/shared/
- [ ] T030 [P] Navigation and layout components in frontend/src/components/layout/

### Pages
- [ ] T031 Develop workspace page in frontend/src/pages/DevelopPage.tsx
- [ ] T032 Playground workspace page in frontend/src/pages/PlaygroundPage.tsx
- [ ] T033 Explore workspace page in frontend/src/pages/ExplorePage.tsx
- [ ] T034 Main layout and routing in frontend/src/pages/App.tsx

### API Endpoints Implementation
- [ ] T035 [P] Implement Agent Management API endpoints in frontend/src/api/agentApi.ts
- [ ] T036 [P] Implement Workspace Data API endpoints in frontend/src/api/workspaceApi.ts
- [ ] T037 [P] Implement Plugin Management API endpoints in frontend/src/api/pluginApi.ts
- [ ] T038 [P] Implement Workflow Management API endpoints in frontend/src/api/workflowApi.ts

### Utilities and Helpers
- [ ] T039 [P] Local storage utilities in frontend/src/utils/storage.ts
- [ ] T040 [P] Validation utilities in frontend/src/utils/validation.ts

## Phase 3.4: Integration
- [ ] T041 Connect services to MSW mock API
- [ ] T042 Implement state management (Zustand/Context) for global state
- [ ] T043 Add error handling and logging throughout the application
- [ ] T044 Implement responsive design for all components
- [ ] T045 Add internationalization support (if needed)

## Phase 3.5: Polish
- [ ] T046 [P] Unit tests for validation utilities in frontend/tests/unit/test_validation.py
- [ ] T047 [P] Unit tests for services in frontend/tests/unit/test_services.py
- [ ] T048 Performance optimization (<2s page load)
- [ ] T049 [P] Update documentation in frontend/docs/
- [ ] T050 Run manual testing based on quickstart.md scenarios
- [ ] T051 Add loading states and user feedback
- [ ] T052 Accessibility improvements
- [ ] T053 Final code review and cleanup

## Dependencies
- Tests (T007-T014) before implementation (T015-T040)
- Model tasks (T015-T020) before service tasks (T021-T025)
- Service tasks (T021-T025) before UI component tasks (T026-T030)
- UI components (T026-T030) before page tasks (T031-T034)
- API implementation (T035-T038) before integration (T041)
- Core implementation before polish (T046-T053)

## Parallel Example
```
# Launch model creation tasks together:
Task: "AI Agent model in frontend/src/models/Agent.ts"
Task: "Plugin model in frontend/src/models/Plugin.ts"
Task: "Workflow model in frontend/src/models/Workflow.ts"
Task: "Prompt model in frontend/src/models/Prompt.ts"
Task: "AnalyticsData model in frontend/src/models/AnalyticsData.ts"
Task: "ChatMessage model in frontend/src/models/ChatMessage.ts"

# Launch service creation tasks together:
Task: "AgentService for CRUD operations in frontend/src/services/AgentService.ts"
Task: "PluginService for plugin management in frontend/src/services/PluginService.ts"
Task: "WorkflowService for workflow management in frontend/src/services/WorkflowService.ts"
Task: "ChatService for chat interactions in frontend/src/services/ChatService.ts"
Task: "AnalyticsService for data retrieval in frontend/src/services/AnalyticsService.ts"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each task
- Avoid: vague tasks, same file conflicts

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts**:
   - Each contract file → contract test task [P]
   - Each endpoint → implementation task
   
2. **From Data Model**:
   - Each entity → model creation task [P]
   - Relationships → service layer tasks
   
3. **From User Stories**:
   - Each story → integration test [P]
   - Quickstart scenarios → validation tasks

4. **Ordering**:
   - Setup → Tests → Models → Services → Endpoints → Polish
   - Dependencies block parallel execution

## Validation Checklist
*GATE: Checked by main() before returning*

- [x] All contracts have corresponding tests
- [x] All entities have model tasks
- [x] All tests come before implementation
- [x] Parallel tasks truly independent
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task