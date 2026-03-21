---
description: "Generate CRUD handler scaffolding for a database table. Creates data layer, handler files, and paging handler following factory pattern conventions. Use after adding table DDL to dbutils.go."
arguments:
  - name: table
    description: "The table name (e.g., 'places', 'products')"
    required: true
  - name: singular
    description: "Singular name for the struct (defaults to title-cased table name)"
    required: false
---

You are scaffolding a new CRUD handler with pagination for the table `{{table}}`.

## Prerequisites
Verify the table exists in [data/dbutils.go](data/dbutils.go) initDB function.

## Step 1: Analyze Table Structure
Read the table DDL from dbutils.go and identify:
- Table name: `{{table}}`
- Column names and types
- Primary key
- Foreign keys (if any)

## Step 2: Generate Data Layer File

Create `data/{{table}}.go` following this pattern:

```go
package data

import "github.com/google/uuid"

type {{singular}} struct {
	// Map columns to Go fields:
	// TEXT -> string
	// INTEGER -> int or int64
	// REAL -> float64
}

func New{{singular}}() {{singular}} {
	return {{singular}}{
		ID: uuid.NewString(),
		// Initialize other fields with sensible defaults
	}
}
```

**Important:**
- Use PascalCase for struct name and fields
- Map SQL types correctly
- Always generate UUID for ID field in constructor

## Step 3: Generate Handler File

Create `handlers/{{table}}.go` with this structure:

```go
package handlers

import (
	"database/sql"
	"eneb/data"
	
	"github.com/gin-gonic/gin"
)

func Reg_{{table}}(r *gin.Engine, db *sql.DB) {

	// GET /{{table}} - List all
	getScanner := func(row data.RowScanner) (*data.{{singular}}, error) {
		item := data.New{{singular}}()
		err := row.Scan(&item.Field1, &item.Field2, /* all fields */)
		if err != nil {
			return nil, err
		}
		return &item, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany(db,
		`select field1, field2, ...
		from {{table}}
		order by /* appropriate column */`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	getHandler := MakeHandlerGetMany(cmdSelect)

	r.GET("/{{table}}",
		func(c *gin.Context) {
			getHandler(c, []any{})
		})

	// POST /{{table}} - Create/Update
	cmdSave, err := data.MakeDataCmdSaveOne(db,
		"insert or replace into {{table}}(field1, field2, ...) VALUES(?,?,?)",
		func(item *data.{{singular}}) []any {
			return []any{item.Field1, item.Field2, /* all fields */}
		})
	if err != nil {
		panic(err)
	}
	postHandler := MakeHandlerPostOne(cmdSave)

	r.POST("/{{table}}",
		func(c *gin.Context) {
			postHandler(c, []any{})
		})
}
```

**Important:**
- Use exact column order in SELECT and INSERT
- Scanner receives all columns from SELECT
- Slicer provides all values for INSERT in same order
- Use appropriate ORDER BY clause

## Step 4: Generate Paging Handler File

Add both registration calls to main.go:

1. Search for existing `handlers.Reg_` calls in main.go
2. Add both lines before the `// start` comment and `r.Run()` call:
   - `handlers.Reg_{{table}}(r, db)`
   - `handlers.Reg_{{table}}paging(r, db)`
3. Keep them grouped with other handler registrations

The lines should be inserted like:
```go
handlers.Reg_existinghandler(r, db)
handlers.Reg_{{table}}(r, db)          // <-- New line
handlers.Reg_{{table}}paging(r, db)    // <-- New line
// start
r.Run(...)
```

## Step 6: Summaryc(row data.RowScanner) (*data.{{singular}}, error) {
		item := data.New{{singular}}()
		err := row.Scan(&item.Field1, &item.Field2, /* all fields */)
		if err != nil {
			return nil, err
		}
		return &item, nil
	}

	// Previous page (before cursor)
	cmdSelectBefore, err := data.MakeDataCmdSelectMany(db,
		`select field1, field2, ...
		from {{table}}
		where (sortcol1, id) < (?, ?)
		order by sortcol1 desc, id desc limit ?`,
		true,  // reverse results
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany(cmdSelectBefore)

	// Next page (after cursor)
	cmdSelectAfter, err := data.MakeDataCmdSelectMany(db,
		`select field1, field2, ...
		from {{table}}
		where (sortcol1, id) > (?, ?)
		order by sortcol1, id limit ?`,
		false,  // don't reverse
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany(cmdSelectAfter)

	r.GET("/{{table}}/page/prev",
		func(c *gin.Context) {
			sortcol1 := ctxQParamInt(c, "sortcol1")  // or ctxQParamStr for TEXT
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*sortcol1, *id, *limit})
		})

	r.GET("/{{table}}/page/next",
		func(c *gin.Context) {
			sortcol1 := ctxQParamInt(c, "sortcol1")  // or ctxQParamStr for TEXT
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*sortcol1, *id, *limit})
		})
}
```

**Important:**
- Use the same ORDER BY columns from the regular handler
- For prev: `< (?, ?)` with DESC order and `reverse=true`
- For next: `> (?, ?)` with ASC order and `reverse=false`
- Query params match the cursor columns (first sort column + id)
- Use `ctxQParamInt` for INTEGER, `ctxQParamStr` for TEXT

## Step 5: Register Routes in main.go

Add the registration call to main.go:

1. Search for existing `handlers.Reg_` calls in main.go
2. Add `handlers.Reg_{{table}}(r, db)` before the `// start` comment and `r.Run()` call
3. Keep it grouped with other handler registrations

The line should be inserted like:
```go
handlers.Reg_existinghandler(r, db)
handlers.Reg_{{table}}(r, db)  // <-- New line
// start
r.Run(..6: Summary

After generating files, provide:
1. List of files created:
   - `data/{{table}}.go`
   - `handlers/{{table}}.go`
   - `handlers/{{table}}_page.go`
2. Confirmation that both registrations were added to main.go:
   - `handlers.Reg_{{table}}(r, db)`
   - `handlers.Reg_{{table}}paging(r, db)`
3. Example curl commands to test:
   - GET all: `curl http://localhost:8080/{{table}}`
   - POST: `curl -X POST http://localhost:8080/{{table}} -H "Content-Type: application/json" -d '{...}'`
   - GET next page: `curl "http://localhost:8080/{{table}}/page/next?col=value&id=xxx&limit=20"`
   - GET prev page: `curl "http://localhost:8080/{{table}}/page/prev?col=value&id=xxx&limit=20"`

## Design Principles Applied

- **Single Responsibility**: Data layer handles DB, handlers handle HTTP, pagination logic separated
- **Dependency Injection**: DB and router passed to Reg_ functions
- **Open/Closed**: Factory pattern allows extending without modifying base handlers
- **Interface Segregation**: Scanner and Slicer are minimal focused interfaces
- **DRY**: Cursor-based pagination reuses same scanner function

Execute all steps, create all three files,ws extending without modifying base handlers
- **Interface Segregation**: Scanner and Slicer are minimal focused interfaces

Execute all steps and create both files and update main.go.
