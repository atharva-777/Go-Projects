URL Shortener (minimal)

Basic RESTful URL shortener service implemented in Go.

Quick start

- Initialize the module (already done here):

  ```bash
  go mod init github.com/<your-username>/url-shortner
  ```

- Fetch dependencies and build:

  ```bash
  cd url-shortner
  go mod tidy
  go build
  ./url-shortner
  ```

API Endpoints

- POST `/api/shorten` with JSON `{ "url": "https://example.com" }` — create short URL
- GET `/{code}` — redirect to original URL (increments visits)
- GET `/api/url/{code}` — get URL info including `visits`
- PUT `/api/url/{code}` with JSON `{ "url": "https://new" }` — update original
- DELETE `/api/url/{code}` — delete short URL

Notes

- This implementation uses an in-memory store; restarting the server loses data.
- To run behind a real domain, update the `short_url` base or build frontend and redirect rules.
