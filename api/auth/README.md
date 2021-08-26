# Auth service

Auth service manages JWT auth. Additional external dependency is Redis that is necessary to
keep `refresh_token`s.

# Session management

To make jwt tokens work like a session there was implemented simple idea that we can have short-live
auth tokens and long-live refresh tokens. More details: [https://blog.hasura.io/best-practices-of-using-jwt-with-graphql/](The Ultimate Guide to handling JWTs on frontend clients)

# Usage

There are multiple endpoints to create and refresh jwt tokens:
1. `/refresh`: refresh token when `refresh_token` cookie is valid
2. `/*/login`: login with external services like github, facebook and so on

Login workflow:
1. Try to `/refresh`, if there is no errors just use token from response, else:
2. Perform `/*/login`, after get redirected back perform `1`
