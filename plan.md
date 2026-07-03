# Improvement Plan - Todo TUI App

## Current Status
A functional Go-based CLI application with a Bubble Tea interface to manage tasks. It supports categories, priority levels, sorting, and basic filtering of completed items.

## Identified Issues & Opportunities
1. **Dual codebase**: Incomplete refactoring to `pkg/` + `cmd/` with import path mismatch — dead code that cannot build.
2. **Monolithic Update()**: 120-line method with deeply nested switch handling 5 modes and 15+ keybindings.
3. **No separation of concerns**: UI rendering, state management, and persistence all in `main` package.
4. **Global mutable state**: `var filename` in storage.go break parallel testability.
5. **No abstraction over persistence**: `loadTodos()` / `saveTodos()` are concrete functions, not interfaces.
6. **Missing empty-state UX**: No "No tasks found" messages for empty categories.

---

## Phase 1: Architecture & Refactoring (High Priority) ✅ 2026-07-03

- [x] **1-1: Resolve dual codebase** — Delete dead `pkg/` and `cmd/` directories. All active code is in the root package.
- [x] **1-2: Interface-ize storage** — Create `Storage` interface, inject via constructor. Eliminate package-level `var filename`.
- [x] **1-3: Split Update()** — Extract mode-specific handler methods (`handleInputKey`, `handleCategoryDeleteKey`, `handleViewKey`).
- [x] **1-4: Separate View()** — Move UI rendering to `view.go`.

**Result**: Build + 29 tests passing. File structure: `main.go` (entry + Update), `model.go` (types + sort), `view.go` (UI), `storage.go` (interface + FileStorage).

## Phase 2: Bug Fixes & Stability (Medium Priority) ✅ 2026-07-03

- [x] **2-1: Fix sortTodos()** — Extract active-category items, sort separately, then place back. Preserves non-active item positions.
- [x] **2-2: Empty-state messages** — Show "タスクがありません" when category empty, "表示できるタスクがありません" when all filtered out.
- [x] **2-3: Configurable default category** — Replace hardcoded "Home" with `defaultCategory` constant in model.go.
- [x] **2-4: Debounce saveTodos()** — Save only on quit (`q` / `ctrl+c`). Removed per-mutation saves.

**Result**: Build + 29 tests passing. `sortTodo()` now uses proper strict weak ordering via extraction. Empty categories show helpful messages. Save-on-quit reduces disk I/O.

## Phase 3: Features (Medium-Low Priority)

- [ ] **3-1: Custom due date input** — Allow users to set arbitrary dates instead of always "tomorrow".
- [ ] **3-2: Text search** — Filter by title substring match in addition to completed/uncompleted.
- [ ] **3-3: Task description field** — Add optional notes/body to `Todo` struct.
- [ ] **3-4: Undo** — Basic undo for delete and edit operations.

## Phase 4: Polish (Low Priority)

- [ ] **4-1: Config file support** — Load settings (default category, colors, keybindings) from `~/.config`.
- [ ] **4-2: Color theme** — User-customizable color schemes.
- [ ] **4-3: i18n groundwork** — Externalize Japanese strings for future multi-language support.
