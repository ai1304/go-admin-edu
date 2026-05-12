# Frontend Workspace

This workspace separates the two frontend surfaces:

- `apps/admin`: administration console, currently linked to the cloned `go-admin-ui` Vue3 branch.
- `apps/portal`: public/user-facing portal for resources, courses, activities, experts, and student/teacher workflows.
- `packages/shared`: shared constants, request helpers, and cross-app types.

Run from `web/`:

```bash
pnpm dev:admin
pnpm dev:portal
```

The `apps/admin` path is now the physical Vue3 admin project directory.
