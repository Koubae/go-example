Testing GO project definition
=============================

This is the canonical testing definition for this project. Follow it when writing or changing tests.

## Scope

- Unit tests live next to the code they cover (`internal/.../*_test.go`).
- Shared helpers live in `tests/testutil` (not named `*_test.go`).
- Default unit-test packages: `./internal/...` (`make test-unit`).
- Coverage target: 80% over `./internal/...` (`make test-coverage`).


## How to run

```bash
make test-unit
make test-coverage
```

## Package and file layout

- Use the **same package** as the code under test (white-box: `package foo`, not `foo_test`).
- Name files `<name>_test.go` beside the source file.
- Put reusable fixtures (DB, clock, HTTP) in `tests/testutil`. Persistence tests should use an isolated in-memory or per-test store from that helper, not a shared live database.

## Naming

- Top-level tests: `Test<Type>__<Method>` (double underscore), e.g. `TestStore__Create`, `TestStore__GetByID`.
- Pure functions may use `Test<Func>` or `Test<Func>OnSuccess` / `Test<Func>OnError`.
- Subtests: `t.Run` with short behavior names (`on success`, `on not found`, `on duplicate`). Table rows use `tt.name` or a distinctive field.

## Structure

- Call `t.Parallel()` on the top-level test.
- Tests must be autonomous: each test (and subtest) sets up its own data and must not depend on another test's order, state, or side effects.
- Prefer table-driven cases for many similar inputs; use nested `t.Run` for distinct behaviors.
- A unique case, or a case with unique setup, can be a standalone test instead of a table.
- Use `t.Context()` for repository and other context-aware APIs.
- Mark helpers with `t.Helper()`.
- Keep unique rows in DB tests (`"unit-test-" + uuid.NewString()[:8]`).

## Assertions

Use `github.com/stretchr/testify/require` and `assert`:

- `require` for conditions the rest of the test cannot proceed without (`NoError`, setup, IDs).
- `assert` for extra checks after the test can still continue.
- Compare sentinel errors with `ErrorIs` (`ErrNotFound`, `ErrDuplicate`, `ErrInvalidInput`, …).
- Prefer equality on domain values over field-by-field dumps unless a single field is the point of the case.

## What to test

- Aim for high coverage, but treat it as a **result** of asserting behavior and stability, not as the goal. Do not add tests only to raise a percentage.
- Keep tests unique. Do not re-prove the same behavior at every layer.
  - **Lower layer** (store, helper): a thin happy path, plus edge cases only that layer can see (constraints, mapping).
  - **Higher layer** (service, HTTP): the full contract (validation, orchestration, status codes, body).
  - Example: a service that uses a store should exercise the store for real and assert the *service* contract. Do not copy the store’s cases into the service test. The store still needs its own thin tests.
- Use SQLite (or another isolated in-memory / per-test engine) as the database. Do not mock the database except for true edge cases (driver failure, connection loss).
- Apply the **full** migration set in setup. Migrating a single table produces a schema the app never runs and leads to false confidence.
- HTTP tests boot the **entire** application (router, middleware, wiring). Do not unit-test a handler in isolation.

## Example (higher layer: full contract, real store)

```go
func TestService__Create(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t) // SQLite + all migrations
	svc := NewService(NewStore(db))

	t.Run("on success", func(t *testing.T) {
		created, err := svc.Create(t.Context(), CreateInput{
			Name: "unit-test-" + uuid.NewString()[:8],
		})

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created.ID)
		assert.Equal(t, StatusActive, created.Status)
	})

	t.Run("on duplicate name", func(t *testing.T) {
		name := "unit-test-" + uuid.NewString()[:8]
		_, err := svc.Create(t.Context(), CreateInput{Name: name})
		require.NoError(t, err)

		_, err = svc.Create(t.Context(), CreateInput{Name: name})
		require.ErrorIs(t, err, ErrDuplicate)
	})
}
```

## Example (HTTP: full application)

```go
func TestAPI__CreateItem(t *testing.T) {
	t.Parallel()

	app := testutil.StartApp(t) // router, middleware, store, full migrations

	body := `{"name":"unit-test-` + uuid.NewString()[:8] + `"}`
	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
}
```


## Example (pure function, table-driven)

```go
func TestFormat(t *testing.T) {
	t.Parallel()

	t.Run("on invalid input", func(t *testing.T) {
		result, err := Format("")

		require.Error(t, err)
		assert.Equal(t, "", result)
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "trims space", in: "  ok  ", want: "ok"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Format(tt.in)

			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}
```

## Example (store / repository)

```go
func TestStore__Create(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t)
	repository := NewStore(db)

	t.Run("on success", func(t *testing.T) {
		item := Item{Name: "unit-test-" + uuid.NewString()[:8]}
		created, err := repository.Create(t.Context(), item)

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created.ID)
		require.Equal(t, item.Name, created.Name)
	})
}
```

## What not to do

- Do not skip `t.Parallel()` without a reason written in the test.
- Do not use a shared live database for unit tests.
- Do not assert with `t.Fatal` / `t.Error` when testify already covers the check.
- Do not duplicate this file into `.cursor/rules/go-testing.mdc`; that rule only points here.
