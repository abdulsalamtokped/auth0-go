# How to Run

1. Visit `http://localhost:3010`
2. Hit `http://localhost:3010/api/force-login` to get jwt token
3. Hit `http://localhost:3010/api/private` with header: `authorization: Bearer {{jwtToken}}`