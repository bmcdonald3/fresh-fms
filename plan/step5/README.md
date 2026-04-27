## Implement Reconciler Logic

### Step outcome

Go code implementing the background execution logic.

### What to give

The concrete operations from Step 1 and the Go structs from Step 3.

### Prompt

Generate the Go code for the reconcilers. The code should be formatted to fit within the `reconcile[Resource]` methods located in the generated `pkg/reconcilers/[resource]_reconciler.go` files. The logic must execute the concrete operations identified in Step 1. Ensure the code includes idempotency checks, progressive status updates, calls to `r.Client.Update` to save status changes, and appropriate error returns to trigger requeuing on failure.

### Context

The custom logic must be implemented in the safe-to-edit user stub files (`pkg/reconcilers/<resource>_reconciler.go`). The generated orchestration wrapper handles event ingestion and requeuing.

1. Idempotency: The reconciler may be called multiple times for the same event. It must check the `Status.Phase` first and return immediately if the work is already done.
2. Progressive Updates: The reconciler should update `Status.Phase` to intermediate states (e.g., "Provisioning") before starting long-running tasks.
3. State Storage: Any changes to the resource status must be explicitly saved via the storage client before the function returns.

```go
// Example Reconciler Implementation
func (r *RackReconciler) reconcileRack(ctx context.Context, res *rack.Rack) error {
    // 1. Idempotency Check
    if res.Status.Phase == "Ready" {
        return nil
    }

    // 2. Perform domain logic (e.g. create child resources)
    template := r.loadTemplate(ctx, res.Spec.TemplateUID)
    
    // 3. Update Status
    res.Status.Phase = "Ready"
    res.Status.TotalChassis = template.Spec.ChassisCount
    
    // 4. Save state
    if err := r.Client.Update(ctx, res); err != nil {
        return fmt.Errorf("failed to update status: %w", err)
    }

    return nil
}
```