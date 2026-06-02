# OneTrack — Bid Workspace & Task Management API Documentation

## Base URL

```
http://localhost:8081/api/v1
```

## Authentication

All endpoints require a valid JWT Bearer token:

```
Authorization: Bearer <access_token>
```

---

## Response Envelope

### Success
```json
{
  "success": true,
  "message": "Operation description",
  "data": { }
}
```

### Error
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message",
    "details": null
  }
}
```

---

# PART 1 — BID WORKSPACE APIs

---

## Creation Mode Philosophy

When adding a tender, two modes are available:

| Mode           | Value          | Description                                                      |
| -------------- | -------------- | ---------------------------------------------------------------- |
| Manual         | `MANUAL`       | Operator fills all tender fields via structured form             |
| Intelligence   | `INTELLIGENCE` | AI processes uploaded document and auto-populates fields (future) |

Both modes share **identical lifecycle, task, and reporting APIs**. Only creation differs.

---

## POST /api/v1/bids

**Create a new Bid Workspace**

**Permission:** `bid.create`

### Request Body

| Field                | Type     | Required | Description                                          |
| -------------------- | -------- | -------- | ---------------------------------------------------- |
| `creation_mode`      | string   | Yes      | `MANUAL` or `INTELLIGENCE`                           |
| `title`              | string   | Yes      | Tender title                                         |
| `bid_no`             | string   | No       | Internal/portal tender number                        |
| `gem_bid_no`         | string   | No       | GeM Bid Number (GEM/2026/B/XXXXX)                   |
| `organization_name`  | string   | No       | Buying organization (e.g. NIC, DRDO)                 |
| `department_name`    | string   | No       | Ministry / department                                |
| `portal_source`      | string   | No       | `GeM` \| `CPPP` \| `eProcure` (default: `GeM`)       |
| `bid_type`           | string   | No       | `CUSTOM_BID` \| `REGULAR` \| `RA_BID`               |
| `gem_bid_type`       | string   | No       | GeM classification label                            |
| `category`           | string   | No       | Product/service category                             |
| `estimated_value`    | number   | No       | Estimated bid value (INR)                            |
| `emd_amount`         | number   | No       | EMD amount (INR)                                     |
| `emd_type`           | string   | No       | `ONLINE` \| `DD` \| `BG` \| `EXEMPTED`              |
| `emd_exempted`       | boolean  | No       | Whether EMD is exempted                              |
| `oem_required`       | boolean  | No       | Whether OEM authorization/MAF is required            |
| `has_tech_eval`      | boolean  | No       | Whether technical evaluation is required             |
| `opening_date`       | datetime | No       | Bid opening date (ISO 8601)                          |
| `closing_date`       | datetime | No       | Bid closing deadline (ISO 8601)                      |
| `bid_owner_id`       | UUID     | Yes      | Assigned bid owner user ID                           |
| `remarks`            | string   | No       | Internal remarks / notes                             |
| `metadata`           | object   | No       | Any additional key-value data                        |

### Example Request (Manual Mode)

```json
{
  "creation_mode": "MANUAL",
  "title": "Supply of Enterprise Firewall — NIC Delhi",
  "bid_no": "TENDER/2026/IT/001",
  "gem_bid_no": "GEM/2026/B/12345",
  "organization_name": "National Informatics Centre",
  "department_name": "MeitY",
  "portal_source": "GeM",
  "bid_type": "CUSTOM_BID",
  "gem_bid_type": "Custom Bid",
  "category": "Networking Equipment",
  "estimated_value": 2500000,
  "emd_amount": 50000,
  "emd_type": "ONLINE",
  "emd_exempted": false,
  "oem_required": true,
  "has_tech_eval": true,
  "opening_date": "2026-06-01T10:00:00Z",
  "closing_date": "2026-06-15T18:00:00Z",
  "bid_owner_id": "b3f1c2a4-1234-4abc-9def-000000000001",
  "remarks": "Cisco preferred OEM. Check MAF availability.",
  "metadata": {
    "priority_account": true
  }
}
```

### Response — 201 Created

```json
{
  "success": true,
  "message": "Bid workspace created successfully",
  "data": {
    "id": "c1d2e3f4-0000-4abc-9def-000000000099",
    "title": "Supply of Enterprise Firewall — NIC Delhi",
    "gem_bid_no": "GEM/2026/B/12345",
    "bid_no": "TENDER/2026/IT/001",
    "workflow_stage": "DISCOVERED",
    "bid_status": "ACTIVE",
    "creation_mode": "MANUAL",
    "bid_owner_id": "b3f1c2a4-1234-4abc-9def-000000000001",
    "created_at": "2026-06-02T12:00:00Z"
  }
}
```

### Error Codes

| Code                | HTTP | Description                         |
| ------------------- | ---- | ----------------------------------- |
| `VALIDATION_ERROR`  | 400  | Required fields missing/invalid     |
| `USER_NOT_FOUND`    | 404  | bid_owner_id does not exist         |
| `UNAUTHORIZED`      | 401  | No valid JWT token                  |
| `FORBIDDEN`         | 403  | Missing `bid.create` permission     |

---

## GET /api/v1/bids

**List / Search Bids (Paginated)**

**Permission:** `bid.view`

### Query Parameters

| Param            | Type    | Description                                             |
| ---------------- | ------- | ------------------------------------------------------- |
| `page`           | int     | Page number (default: 1)                                |
| `limit`          | int     | Results per page (default: 20, max: 100)                |
| `search`         | string  | Full-text search on title, bid_no, gem_bid_no           |
| `workflow_stage` | string  | Filter by stage (e.g. `DISCOVERED`, `SUBMITTED`)        |
| `bid_status`     | string  | `ACTIVE` \| `CANCELLED` \| `ARCHIVED` \| `WON` \| `LOST` |
| `bid_outcome`    | string  | `WON` \| `LOST` \| `PENDING`                           |
| `bid_owner_id`   | UUID    | Filter by owner                                         |
| `category`       | string  | Filter by category                                      |
| `creation_mode`  | string  | `MANUAL` \| `INTELLIGENCE`                             |
| `closing_before` | date    | ISO 8601 date — filter bids closing before this date    |
| `closing_after`  | date    | ISO 8601 date — filter bids closing after this date     |
| `oem_required`   | boolean | Filter by OEM requirement                               |

### Example Response

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "title": "Supply of Enterprise Firewall",
      "gem_bid_no": "GEM/2026/B/12345",
      "workflow_stage": "COMMERCIAL_PREPARATION",
      "bid_status": "ACTIVE",
      "estimated_value": 2500000,
      "closing_date": "2026-06-15T18:00:00Z",
      "bid_owner": {
        "id": "uuid",
        "full_name": "Rahul Sharma",
        "username": "rahul.sharma"
      },
      "creation_mode": "MANUAL",
      "created_at": "2026-06-02T12:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 42
  }
}
```

---

## GET /api/v1/bids/{id}

**Get Bid Details (Full)**

**Permission:** `bid.view`

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "title": "Supply of Enterprise Firewall — NIC Delhi",
    "bid_no": "TENDER/2026/IT/001",
    "gem_bid_no": "GEM/2026/B/12345",
    "organization_name": "National Informatics Centre",
    "department_name": "MeitY",
    "portal_source": "GeM",
    "bid_type": "CUSTOM_BID",
    "gem_bid_type": "Custom Bid",
    "category": "Networking Equipment",
    "creation_mode": "MANUAL",
    "workflow_stage": "QUALIFICATION_REVIEW",
    "bid_status": "ACTIVE",
    "estimated_value": 2500000,
    "emd_amount": 50000,
    "emd_type": "ONLINE",
    "emd_exempted": false,
    "oem_required": true,
    "has_tech_eval": true,
    "opening_date": "2026-06-01T10:00:00Z",
    "closing_date": "2026-06-15T18:00:00Z",
    "submission_date": null,
    "result_date": null,
    "ra_date": null,
    "final_bid_value": null,
    "l1_price": null,
    "quoted_price": null,
    "bid_outcome": null,
    "outcome_reason": null,
    "qualification_status": "PENDING",
    "tech_compliance_status": null,
    "remarks": "Cisco preferred OEM. Check MAF availability.",
    "competitor_info": [],
    "metadata": { "priority_account": true },
    "bid_owner": {
      "id": "uuid",
      "full_name": "Rahul Sharma",
      "username": "rahul.sharma"
    },
    "members": [
      {
        "user_id": "uuid",
        "full_name": "Anita Singh",
        "role": "MEMBER",
        "added_at": "2026-06-02T12:30:00Z"
      }
    ],
    "created_by": "uuid",
    "created_at": "2026-06-02T12:00:00Z",
    "updated_at": "2026-06-02T12:00:00Z"
  }
}
```

---

## PATCH /api/v1/bids/{id}

**Update Bid Fields**

**Permission:** `bid.edit`

All fields are optional. Send only fields you want to update.

```json
{
  "title": "Updated Tender Title",
  "emd_amount": 75000,
  "closing_date": "2026-06-20T18:00:00Z",
  "remarks": "Updated remarks",
  "oem_required": true
}
```

---

## POST /api/v1/bids/{id}/transition

**Move Bid to Next Workflow Stage**

**Permission:** `bid.edit`

### Request

```json
{
  "target_stage": "COMMERCIAL_PREPARATION",
  "reason": "Qualification verified by bid manager"
}
```

### Allowed Transitions (Manual Mode)

| From Stage               | To Stage(s)                                     |
| ------------------------ | ----------------------------------------------- |
| `DISCOVERED`             | `QUALIFICATION_REVIEW`, `CANCELLED`             |
| `QUALIFICATION_REVIEW`   | `DOCUMENT_COMPILATION`, `CANCELLED`             |
| `DOCUMENT_COMPILATION`   | `OEM_COORDINATION`, `COMMERCIAL_PREPARATION`    |
| `OEM_COORDINATION`       | `COMMERCIAL_PREPARATION`, `CANCELLED`           |
| `COMMERCIAL_PREPARATION` | `INTERNAL_REVIEW`                               |
| `INTERNAL_REVIEW`        | `FINAL_APPROVAL`, `COMMERCIAL_PREPARATION`      |
| `FINAL_APPROVAL`         | `READY_FOR_SUBMISSION`, `INTERNAL_REVIEW`       |
| `READY_FOR_SUBMISSION`   | `SUBMITTED`                                     |
| `SUBMITTED`              | `RA_ACTIVE`, `AWAITING_RESULT`                  |
| `RA_ACTIVE`              | `AWAITING_RESULT`                               |
| `AWAITING_RESULT`        | `WON`, `LOST`, `CANCELLED`                      |

### Response

```json
{
  "success": true,
  "message": "Bid transitioned to COMMERCIAL_PREPARATION",
  "data": {
    "bid_id": "uuid",
    "previous_stage": "QUALIFICATION_REVIEW",
    "current_stage": "COMMERCIAL_PREPARATION",
    "transitioned_at": "2026-06-05T10:00:00Z"
  }
}
```

### Error Codes

| Code                         | HTTP | Description                      |
| ---------------------------- | ---- | -------------------------------- |
| `INVALID_STAGE_TRANSITION`   | 409  | Target stage not allowed         |
| `BID_NOT_FOUND`              | 404  | Bid ID does not exist            |

---

## GET /api/v1/bids/{id}/stage-history

**Get Stage Transition History (Immutable Audit)**

**Permission:** `bid.view`

### Response

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "from_stage": null,
      "to_stage": "DISCOVERED",
      "transition_reason": "Bid created",
      "transitioned_by": {
        "id": "uuid",
        "full_name": "Admin User"
      },
      "created_at": "2026-06-02T12:00:00Z"
    },
    {
      "id": "uuid",
      "from_stage": "DISCOVERED",
      "to_stage": "QUALIFICATION_REVIEW",
      "transition_reason": "Manually qualified",
      "transitioned_by": {
        "id": "uuid",
        "full_name": "Rahul Sharma"
      },
      "created_at": "2026-06-03T09:00:00Z"
    }
  ]
}
```

---

## POST /api/v1/bids/{id}/members

**Add Workspace Member**

**Permission:** `bid.edit`

```json
{
  "user_id": "uuid",
  "role": "MEMBER"
}
```

Roles: `OWNER | MANAGER | MEMBER | REVIEWER | OBSERVER`

---

## DELETE /api/v1/bids/{id}/members/{user_id}

**Remove Workspace Member**

**Permission:** `bid.edit`

---

## PATCH /api/v1/bids/{id}/outcome

**Record Bid Outcome (Win/Loss)**

**Permission:** `bid.edit`

```json
{
  "bid_outcome": "WON",
  "final_bid_value": 2350000,
  "l1_price": 2350000,
  "quoted_price": 2350000,
  "outcome_reason": "L1 bidder. Technically compliant. OEM authorized.",
  "result_date": "2026-07-01T00:00:00Z",
  "competitor_info": [
    { "name": "ABC Technologies", "quoted_price": 2500000, "rank": "L2" },
    { "name": "XYZ Infra", "quoted_price": 2700000, "rank": "L3" }
  ]
}
```

Values for `bid_outcome`: `WON | LOST | CANCELLED`

---

## DELETE /api/v1/bids/{id}

**Archive / Soft-Delete Bid**

**Permission:** `bid.delete`

Sets `bid_status = ARCHIVED` and `archived_at`. Data is preserved.

---

# PART 2 — TASK MANAGEMENT APIs

Tasks are always scoped to a Bid Workspace. They support a full hierarchy:

```
Task
└── Subtask
    └── Checklist Items
```

Activities (comments, call logs, status changes) are tracked at the Task level.

---

## POST /api/v1/bids/{id}/tasks

**Create a Task for a Bid**

**Permission:** `task.create`

### Request Body

| Field         | Type     | Required | Description                    |
| ------------- | -------- | -------- | ------------------------------ |
| `task_type`   | string   | Yes      | Classification of the task     |
| `title`       | string   | Yes      | Task title                     |
| `description` | string   | No       | Detailed instructions          |
| `priority`    | string   | No       | `LOW \| MEDIUM \| HIGH \| CRITICAL` (default: `MEDIUM`) |
| `assigned_to` | UUID     | No       | User to assign task to         |
| `due_date`    | datetime | No       | Task deadline (ISO 8601)       |
| `metadata`    | object   | No       | Task-specific metadata         |

**Task Types:**
`GENERAL | DOCUMENT_COLLECTION | OEM_COORDINATION | QUALIFICATION | COMMERCIAL | APPROVAL | REVIEW | SUBMISSION | CHECKLIST | CUSTOM`

### Example Request

```json
{
  "task_type": "DOCUMENT_COLLECTION",
  "title": "Upload EMD Certificate",
  "description": "Upload scanned EMD receipt from HDFC Bank. Amount: ₹50,000.",
  "priority": "HIGH",
  "assigned_to": "b3f1c2a4-1234-4abc-9def-000000000001",
  "due_date": "2026-06-10T10:00:00Z",
  "metadata": {
    "document_required": "EMD Receipt",
    "bank": "HDFC"
  }
}
```

### Response — 201 Created

```json
{
  "success": true,
  "message": "Task created successfully",
  "data": {
    "id": "task-uuid-0001",
    "bid_id": "bid-uuid",
    "parent_task_id": null,
    "task_type": "DOCUMENT_COLLECTION",
    "title": "Upload EMD Certificate",
    "status": "ASSIGNED",
    "priority": "HIGH",
    "assigned_to": "uuid",
    "due_date": "2026-06-10T10:00:00Z",
    "completion_percentage": 0,
    "created_at": "2026-06-02T12:00:00Z"
  }
}
```

---

## GET /api/v1/bids/{id}/tasks

**List All Tasks for a Bid**

**Permission:** `task.view`

### Query Parameters

| Param         | Description                              |
| ------------- | ---------------------------------------- |
| `status`      | Filter by task status                    |
| `priority`    | Filter by priority                       |
| `assigned_to` | Filter by assigned user UUID             |
| `task_type`   | Filter by task type                      |
| `parent_only` | `true` — only top-level tasks (no subtasks) |

---

## GET /api/v1/tasks/{id}

**Get Task Details**

**Permission:** `task.view`

Returns full task object including subtask count, checklist items count, and last activity.

---

## PATCH /api/v1/tasks/{id}

**Update Task Fields**

**Permission:** `task.edit`

```json
{
  "title": "Updated task title",
  "description": "Updated description",
  "priority": "CRITICAL",
  "due_date": "2026-06-12T10:00:00Z"
}
```

---

## PATCH /api/v1/tasks/{id}/status

**Update Task Status**

**Permission:** `task.edit`

```json
{
  "status": "IN_PROGRESS"
}
```

**Valid Status Values:**
`OPEN | ASSIGNED | IN_PROGRESS | WAITING_EXTERNAL | BLOCKED | UNDER_REVIEW | COMPLETED | CANCELLED | ESCALATED | REOPENED`

### Auto-Activity Logged

Status changes automatically create a `STATUS_CHANGED` activity log entry.

---

## PATCH /api/v1/tasks/{id}/assign

**Assign or Reassign Task**

**Permission:** `task.assign`

```json
{
  "assigned_to": "new-user-uuid"
}
```

---

## DELETE /api/v1/tasks/{id}

**Cancel / Delete Task**

**Permission:** `task.delete`

Sets status to `CANCELLED`. Tasks with completed subtasks cannot be deleted.

---

## POST /api/v1/tasks/{id}/subtasks

**Create a Subtask**

**Permission:** `task.create`

```json
{
  "title": "Confirm EMD amount with finance team",
  "priority": "MEDIUM",
  "assigned_to": "uuid",
  "due_date": "2026-06-08T10:00:00Z"
}
```

Subtasks inherit `bid_id` from parent task. `parent_task_id` is set automatically.

---

## GET /api/v1/tasks/{id}/subtasks

**List Subtasks**

**Permission:** `task.view`

---

## POST /api/v1/tasks/{id}/activities

**Add Activity (Comment / Call Log / Status Note)**

**Permission:** `task.edit`

### Request

| Field           | Type   | Required | Description                   |
| --------------- | ------ | -------- | ----------------------------- |
| `activity_type` | string | Yes      | Type of activity              |
| `activity_data` | object | Yes      | Payload based on activity type |

**Activity Types:**
`COMMENT | STATUS_CHANGED | ASSIGNED | ATTACHMENT_ADDED | CALL_LOGGED | EMAIL_LOGGED | ESCALATED | DEADLINE_CHANGED`

### Example — Comment

```json
{
  "activity_type": "COMMENT",
  "activity_data": {
    "message": "EMD cheque deposited at HDFC Janakpuri branch. Expecting receipt by EOD."
  }
}
```

### Example — Call Logged

```json
{
  "activity_type": "CALL_LOGGED",
  "activity_data": {
    "contact": "Rajesh (OEM Manager)",
    "summary": "Discussed MAF availability. Will send in 2 days.",
    "duration_minutes": 15
  }
}
```

### Response — 201 Created

```json
{
  "success": true,
  "message": "Activity logged",
  "data": {
    "id": "activity-uuid",
    "task_id": "task-uuid",
    "activity_type": "COMMENT",
    "activity_data": { "message": "EMD cheque deposited..." },
    "performed_by": {
      "id": "uuid",
      "full_name": "Rahul Sharma"
    },
    "created_at": "2026-06-05T14:30:00Z"
  }
}
```

---

## GET /api/v1/tasks/{id}/activities

**Get Task Activity Timeline**

**Permission:** `task.view`

### Response

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "activity_type": "STATUS_CHANGED",
      "activity_data": {
        "from": "OPEN",
        "to": "ASSIGNED"
      },
      "performed_by": { "id": "uuid", "full_name": "Admin" },
      "created_at": "2026-06-02T12:00:00Z"
    },
    {
      "id": "uuid",
      "activity_type": "COMMENT",
      "activity_data": {
        "message": "EMD cheque deposited at HDFC."
      },
      "performed_by": { "id": "uuid", "full_name": "Rahul Sharma" },
      "created_at": "2026-06-05T14:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 50,
    "total": 2
  }
}
```

---

## POST /api/v1/tasks/{id}/checklists

**Add Checklist Item to Task**

**Permission:** `task.edit`

```json
{
  "title": "EMD receipt scanned and attached",
  "sort_order": 1
}
```

---

## PATCH /api/v1/tasks/{id}/checklists/{checklist_id}

**Toggle Checklist Item Done/Undone**

**Permission:** `task.edit`

```json
{
  "is_done": true
}
```

---

## DELETE /api/v1/tasks/{id}/checklists/{checklist_id}

**Remove Checklist Item**

**Permission:** `task.edit`

---

# PART 3 — PERMISSIONS REFERENCE

## Required Permissions

| Permission    | Description                              |
| ------------- | ---------------------------------------- |
| `bid.create`  | Create a new bid workspace               |
| `bid.view`    | View bid list and details                |
| `bid.edit`    | Update bid, manage members, transition   |
| `bid.delete`  | Archive / soft-delete a bid              |
| `task.create` | Create tasks and subtasks                |
| `task.view`   | View tasks, subtasks, activities         |
| `task.edit`   | Update tasks, add activities, checklists |
| `task.assign` | Assign tasks to users                    |
| `task.delete` | Cancel / delete tasks                    |

## Default Role Access

| Role          | bid.create | bid.view | bid.edit | task.create | task.edit | task.assign |
| ------------- | ---------- | -------- | -------- | ----------- | --------- | ----------- |
| SUPER_ADMIN   | ✅         | ✅       | ✅       | ✅          | ✅        | ✅          |
| ADMIN         | ✅         | ✅       | ✅       | ✅          | ✅        | ✅          |
| BID_MANAGER   | ✅         | ✅       | ✅       | ✅          | ✅        | ✅          |
| BID_OWNER     | ✅         | ✅       | ✅       | ✅          | ✅        | ✅          |
| REVIEWER      | ❌         | ✅       | ❌       | ❌          | ✅        | ❌          |
| FINANCE       | ❌         | ✅       | ❌       | ❌          | ✅        | ❌          |
| MANAGEMENT    | ❌         | ✅       | ❌       | ❌          | ❌        | ❌          |
| OPERATOR      | ❌         | ✅       | ✅       | ✅          | ✅        | ✅          |

---

# PART 4 — ERROR CODES REFERENCE

| Code                       | HTTP | Description                                |
| -------------------------- | ---- | ------------------------------------------ |
| `BID_NOT_FOUND`            | 404  | Bid with given ID does not exist           |
| `TASK_NOT_FOUND`           | 404  | Task with given ID does not exist          |
| `INVALID_STAGE_TRANSITION` | 409  | Target stage transition is not allowed     |
| `BID_ARCHIVED`             | 409  | Cannot modify an archived bid              |
| `TASK_COMPLETED`           | 409  | Cannot modify a completed task             |
| `MEMBER_ALREADY_EXISTS`    | 409  | User is already a workspace member         |
| `VALIDATION_ERROR`         | 400  | Request payload validation failed          |
| `UNAUTHORIZED`             | 401  | Missing or invalid JWT token               |
| `FORBIDDEN`                | 403  | Missing required permission                |
| `INTERNAL_ERROR`           | 500  | Unexpected server error                    |
