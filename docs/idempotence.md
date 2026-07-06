In the context of RESTful APIs and distributed systems, **idempotence** means that making the **same request multiple times** produces the **same final state on the server** as making it just once.

Here is the simplest way to think about it:

- **Idempotent:** If you send the request 1 time or 100 times, the server’s data ends up exactly the same after the first attempt. 
- **Not Idempotent:** If you send the request twice, the server's data changes twice (e.g., two records are created, or a balance is deducted twice).

---

### The Crucial Distinction: State vs. Response

The server's **final state** must be the same. The response (HTTP status code) you get back might differ.

**Example (DELETE):**
- **Request 1:** `DELETE /users/123` → Server soft-deletes the user. Response: `200 OK`.
- **Request 2:** `DELETE /users/123` → User is already deleted. Server makes no further changes. Response: `404 Not Found` (or `204 No Content`).

Even though the responses differ, the operation is **idempotent** because the final state (user is deleted) is identical after both requests.

---

### How It Applies to Your API (from the design doc)

| Method | Idempotent? | Why? |
| :--- | :--- | :--- |
| **`GET`** | ✅ Yes | Reading data doesn't change anything. Running it 10 times reads the same data (unless someone else changed it). |
| **`PUT`** (Full update) | ✅ Yes | Sending the exact same full user object 5 times sets the data to that exact state once. The 5th request is a no-op. |
| **`DELETE`** | ✅ Yes | As explained above. |
| **`PATCH`** (Partial update) | ⚠️ It *can* be, but only if you send **absolute** values. If you send `PATCH /users/1` with `{"age": 30}`, it's idempotent. If you send `{"age": "increment by 1"}`, it is **not** idempotent (running it twice adds 2 years). |
| **`POST`** | ❌ **No (by default)** | Every time you send `POST /orders` with the same cart, it typically creates a *new* order with a new `id`. Request 1 = Order #1; Request 2 = Order #2. The state changes each time. |

---

### Why Do We Care? (Practical Importance)

In a microservices/ERP architecture, networks are unreliable. Your frontend (Next.js) might send a request to create an order, but the connection times out before it receives the response.

- **If the API is idempotent:** The frontend can safely **retry** the request (e.g., after 2 seconds). The server checks the `Idempotency-Key` header, sees it already processed this exact request, and just returns the existing order details. No duplicate charges or duplicate orders.
- **If the API is NOT idempotent:** The frontend cannot safely retry automatically. If it retries, it creates a second duplicate order, which would be a disaster for an ERP.

**In your documentation:** I mentioned the `Idempotency-Key` header. In practice, you give your `POST` endpoints a special header (like `Idempotency-Key: unique-uuid-123`). The server remembers this key and the resulting order. If the client sends the same key again, the server ignores the duplicate and returns the previously created order—turning a naturally non-idempotent `POST` into a safe, idempotent one for critical financial flows.
