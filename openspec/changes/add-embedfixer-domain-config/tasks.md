## 1. Platform registry and persistence

- [x] 1.1 Define a hardcoded platform registry for current supported social media (Twitter/X, Reddit, Instagram) including source-host aliases and default replacement domains.
- [x] 1.2 Add a persistent model for custom replacement-domain overrides keyed by guild and platform.
- [x] 1.3 Register the new model in migration flow so config storage is created automatically.

## 2. Config command surface

- [x] 2.1 Add embedfixer slash command metadata/handlers for config management.
- [x] 2.2 Implement `show` command to list each supported platform with aliases, active domain, and default/custom status.
- [x] 2.3 Implement `set` command with strict input validation for supported platform and domain host value.
- [x] 2.4 Implement `reset` command to remove custom overrides and restore default behavior.

## 3. Embedfixer runtime integration

- [x] 3.1 Refactor embedfixer URL handling to resolve platform by alias from the hardcoded registry.
- [x] 3.2 Resolve replacement domain from custom override when present, otherwise fallback to default registry value.
- [x] 3.3 Preserve existing behavior for unsupported hosts and non-fixable message content.

## 4. Validation and rollout safety

- [x] 4.1 Add focused tests for alias mapping, custom override resolution, default fallback, and invalid command inputs.
- [ ] 4.2 Validate command behavior in a test guild, including `show` output for mixed default/custom states.
- [x] 4.3 Run `go build -o /dev/null ./...` and resolve compile issues.
- [x] 4.4 Document rollback operation to disable custom-domain lookup and force defaults.
