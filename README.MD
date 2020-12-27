### Authentication REST API starter/skeleton/template with Auth0 authentication
#### Go

Featuring Many:Many relationship between organization and users for multi-tenant applications. 
- Echo Framework
- Auth0
- go-pg (migrating to SQLX)
- go redis
- Inspired by https://github.com/ribice/gorsk
- Meant to be paired with https://github.com/jpurdie/frontend_auth

### TODO's:
- Unit testing
- Middleware to use DB connection pool instead of new connection or use Redis to cache profiles and rolse
- Enhance docker integration
- Solidify error handling, logging, and responses

#### DB Pre-requisite:
`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`