# Authorization

This document describes how the system authenticates requests and manages session lifecycle using access and refresh tokens.

## Overview

Authorization relies on two token types working together:

| Token         | Purpose                              | Lifetime                | Format                   |
| ------------- | ------------------------------------ | ----------------------- | ------------------------ |
| Access token  | Grants access to protected resources | Short (15 min – 1 hour) | JWT                      |
| Refresh token | Issues a new access token            | Long (1 day – 1 month)  | Opaque unique identifier |

## Access tokens

- Required on every request to a non-public resource.
- Issued as a **JWT**, signed by the server, so the client and server can trust its claims without a database lookup.
- Carries the minimum claims needed to authorize a request, typically:
  - `user_id`
  - `roles`
  - `permissions`
- Once expired, the token is rejected. The client must use a valid refresh token to obtain a new one.

## Refresh tokens

- Opaque, unique identifiers — they carry no claims and cannot be decoded like a JWT.
- Used exclusively against the refresh endpoint to obtain a new access token (and a new refresh token).
- Persisted server-side so they can be looked up, rotated, and revoked.

### Token families and rotation

- Logging in creates a new **refresh token family**.
- Every time a refresh token is exchanged for a new one, the new token joins the same family (rotation).
- A family therefore represents one continuous session, even as individual tokens are rotated.

### Reuse detection

- Only the most recently issued token in a family is valid for exchange.
- If an older (already-rotated) token in the family is presented, treat this as token reuse — a strong signal of theft:
  - Invalidate the entire family immediately.
  - Require the user to log in again to establish a new session.

## Transport

- Send access token using either:
  - `Authorization: Bearer <token>` header, or
  - A cookie with `HttpOnly`, `Secure`, and `SameSite=Strict` set.
- Send refresh token using either:
  - Request body, or
  - A cookie with `HttpOnly`, `Secure`, and `SameSite=Strict` set.
- Protected routes require the **access token**.
- The refresh and logout endpoints require the **refresh token**.
- Do not mix the two: an access token must never be accepted on the refresh endpoint, and vice versa.

## Implementation notes

- Access tokens: issue and verify with `golang-jwt/jwt/v5`.
- Refresh tokens: generate with a `crypto/rand` which is a cryptographically secure random source; store the family ID and rotation chain.

### Tradeoff: stale roles and permissions

Because the access token's `roles` and `permissions` claims are trusted without a database lookup, any change made server-side - a role revoked, a permission granted — does not take effect until the user's current access token expires and is refreshed.

- This is a deliberate consequence of the stateless-JWT design, not a bug.
- The maximum staleness window equals the access token's lifetime.
- If a change needs to take effect immediately (e.g., disabling a compromised account), it must also invalidate the user's refresh token family server-side. This forces re-authentication and prevents the stale access token from being silently renewed with outdated claims once it expires — though the token itself remains valid, with its original claims, until that expiry.

## Logout

Logout must invalidate the session server-side, not just clear cookies client-side.

1. The client calls a logout endpoint, passing the current refresh token.
2. The server invalidates the entire refresh token family associated with that token (same mechanism as reuse detection).
3. The server clears the access and refresh token cookies by setting them with an immediate expiry.

Clearing cookies alone is not sufficient: if a refresh token was previously stolen (e.g., copied from network traffic or extracted from the client before logout), it would remain valid and usable by an attacker even after the legitimate user logs out. Invalidating the family server-side closes that gap.

## CSRF protection (cookie transport only)

If tokens are sent via cookie rather than the `Authorization: Bearer` header, mutating requests (`POST`, `PUT`, `PATCH`, `DELETE`) must also carry a CSRF token. This is a second, independent layer of protection alongside `SameSite=Strict` — useful because `SameSite` enforcement isn't uniform across all browser versions.

**Why it's needed:** browsers automatically attach cookies to requests based on domain, regardless of which site initiated the request. A malicious page on another site can trigger a request to the API, and the browser will still attach the session cookie. A CSRF token proves the request actually originated from the application's own frontend.

**Approach — double-submit cookie:**

1. On login, the server generates a random CSRF token and sends it as a **non-`HttpOnly`** cookie (must be readable by JS, unlike the auth cookies).
2. The frontend reads the cookie value and attaches it as a custom header (e.g., `X-CSRF-Token`) on every mutating request.
3. A middleware compares the header value against the cookie value on each mutating request and rejects the request if they don't match or either is missing.

**Why this works:** a forged cross-site request will still carry the session cookie, but the attacking site cannot read the CSRF cookie value (browsers block cross-origin cookie reads), so it cannot supply a matching header. The request is rejected before it reaches business logic.

**Scope:**

- Required only for cookie-based token transport.
- Not required for the `Authorization: Bearer` header transport — browsers don't auto-attach custom headers, so forged cross-site requests can't replicate them.
- Apply the check only to state-changing methods (`POST`, `PUT`, `PATCH`, `DELETE`); skip `GET`.
- No server-side storage needed — the double-submit pattern is stateless.
