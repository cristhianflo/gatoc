## Context

The finance feature currently exposes `/dollar all` and `/dollar status`, with `/dollar all` depending on a single upstream USD-focused endpoint. The change introduces EUR data into the same command output while preserving the existing command contract and embed-driven Discord response style.

Upstream contract context:
- `https://ve.dolarapi.com/v1/dolares` returns an array of USD rate objects (for example, `oficial` and `paralelo` sources).
- `https://ve.dolarapi.com/v1/euros/oficial` returns a single official EUR rate object.
- Both payloads include `moneda`, `fuente`, `nombre`, `promedio`, and `fechaActualizacion`; `compra` and `venta` may be null.

Current implementation constraints:
- Slash command responses are deferred and later edited.
- External data is fetched inline in the command handler.
- A single upstream failure currently fails the command response.

Stakeholders:
- Discord users who consume exchange-rate information.
- Bot maintainers responsible for operational reliability of external API integrations.

## Goals / Non-Goals

**Goals:**
- Include EUR rate information in `/dollar all` output together with current USD rates.
- Keep response understandable by labeling currency, source, and update timestamp.
- Reduce total latency impact by performing upstream requests concurrently.
- Avoid full command failure when only one upstream source is unavailable.

**Non-Goals:**
- Renaming `/dollar` command to a broader currency command.
- Adding persistent storage or caching for exchange rates.
- Refactoring `/dollar status` behavior.

## Decisions

1. Keep the existing `/dollar all` command surface and extend its payload.
   - Rationale: Lowest user disruption and no new command discovery burden.
   - Alternative considered: introduce new `/fx all` command; rejected to avoid duplicate functionality and migration overhead.

2. Fetch USD and EUR data concurrently and normalize into a shared in-memory representation before embed rendering.
   - Rationale: Concurrent IO keeps interaction latency closer to the slower endpoint rather than sum of both endpoints.
   - Normalization details:
     - Flatten mixed upstream shapes (USD array + EUR object) into a single list of normalized rate entries.
     - Treat `compra` and `venta` as optional/nullable fields.
     - Parse `fechaActualizacion` as RFC3339-compatible timestamps with both `Z` and offset formats.
     - Use `moneda` and `fuente` fields for stable labeling/grouping in embed output.
   - Alternative considered: sequential requests; rejected because it increases timeout risk and worsens UX under normal network conditions.

3. Allow partial success with explicit degradation messaging in the embed.
   - Rationale: Users should receive available rates even when one endpoint is down.
   - Alternative considered: fail-fast if any endpoint fails; rejected because it creates avoidable outages for users.

4. Maintain rollback simplicity by making EUR retrieval additive and isolated.
   - Rationale: If operational instability appears, maintainers can remove/disable EUR retrieval with minimal impact to existing USD flow.
   - Alternative considered: deep abstraction rewrite now; rejected as unnecessary scope for this change.

### Sequence Diagram

```text
User            Discord           Bot /dollar all handler        USD API        EUR API
 |                 |                        |                      |              |
 | /dollar all     |                        |                      |              |
 |---------------->| InteractionCreate      |                      |              |
 |                 |----------------------->| defer reply          |              |
 |                 |                        |--------------------->| GET USD      |
 |                 |                        |------------------------------->| GET EUR |
 |                 |                        |<---------------------| USD array    |
 |                 |                        |<-------------------------------| EUR object |
 |                 |                        | normalize (array+obj merge)    |
 |                 |                        | compose embed                  |
 |                 |<-----------------------| edit deferred response         |
 | view embed      |                        |                      |              |
```

## Risks / Trade-offs

- [Higher dependency count] -> Mitigation: implement partial-failure behavior and clear labels for unavailable currency data.
- [Slower worst-case response due to upstream latency] -> Mitigation: concurrent requests and request timeouts.
- [Inconsistent upstream schemas between USD/EUR endpoints] -> Mitigation: normalize into an internal type and validate parse errors per source.
- [Command naming mismatch (`/dollar` including EUR)] -> Mitigation: keep explicit currency labels in output and revisit command taxonomy in a separate change.

## Migration Plan

1. Implement EUR retrieval as an additive path in `/dollar all`.
2. Deploy with existing USD flow intact.
3. Monitor command reliability and response quality in production.
4. Rollback (if needed): disable/remove EUR fetch path and return to USD-only embed composition.

## Open Questions

- Desired presentation order in embed fields when one currency is missing.
