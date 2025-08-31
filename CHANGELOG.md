## [1.0.4] - 2025-08-31

### ⚙️ Miscellaneous Tasks

- Fix release workflow (#176)
## [1.0.3] - 2025-08-29

### ⚙️ Miscellaneous Tasks

- *(web)* Update package dependencies (#173)
- Add git-cliff configuration for changelog generation (#174)
- Upgrade git-cliff-action to version 4 (#175)
## [1.0.1] - 2025-08-27

### 🚀 Features

- Add version for server and cli (#172)
## [1.0.0] - 2025-08-27

### ⚙️ Miscellaneous Tasks

- Update release workflow (#170)
- Refine release workflow configuration (#171)
## [0.0.1-beta.1] - 2025-08-24

### 🚀 Features

- Add minimal ping API server with Fiber and OpenAPI codegen scaffolding
- Add CI/CD workflows and frontend setup (#2)
- Add landing page (#7)
- Add theme mode toggle component (#8)
- Support light mode and update icon (#9)
- Add db integration (#11)
- Implement OIDC authentication and user management (#12)
- Implement login, signup, and logout (#13)
- Add openapi dependency (#14)
- Enhance dockerfile (#32)
- Implement sign up and login (#62)
- Publish ts fetch client (#63)
- Flatten request parameters (#65)
- Implement login function (#64)
- Add GetCurrentUser API endpoint and JWT authentication middleware (#74)
- Get user info after login and refactor auth context (#75)
- Add configuration management API endpoints and models (#94)
- *(web)* Add dashboard page (#95)
- *(backend)* List vault by lite item (#100)
- *(web)* Add create vault component (#99)
- *(backend)* Add audit log api (#109)
- *(web)* Refine audit log page (#110)
- *(vault)* Enforce unique vault names per user (#112)
- *(docs)* Add CLAUDE.md for project guidance and update .gitignore (#114)
- *(api_key)* Add API key apis (#115)
- *(web)* Add api-key page (#118)
- *(web)* Refine vault, audit log, api key page (#119)
- Add cli api (#145)
- Add health check route to public routes (#146)
- Add GitHub Actions workflow for building and publishing Go client (#149)
- Add CLI for VaultHub with list and get commands (#150)
- *(cli)* Add CLI application (#164)
- Enhance vault management with edit and view modals (#165)

### 🐛 Bug Fixes

- Improve error message handling (#84)
- Fix vault api response (#98)
- Api key authentication in header and fix yaml linting issues (#160)
- Fix cli api tag (#162)
- Update public route in middleware (#163)
- Fix create vault and enhance audit log display (#166)

### 🚜 Refactor

- Standardize API method names (#29)
- Rename configuration to vault (#96)
- Update API to use unique_id instead of id for vault operations (#97)
- *(audit_log)* Replace vault_id with vault object in audit log (#111)
- *(api)* Use camel case in api schema (#116)
- *(api)* Split api.yaml and maintain go generate (#130)
- Refactor project (#147)

### 🎨 Styling

- Enhance eslint rules (#113)

### ⚙️ Miscellaneous Tasks

- Add gitlab ci/cd configuration (#3)
- Upload artifact upload (#4)
- Add dockerfile (#6)
- Bump version (#10)
- Enable force push for gitlab sync (#30)
- Configure npm authentication in gitlab ci and add base docker image (#31)
- Update GitHub Actions workflow for publishing client (#33)
- Fix npm publishing (#34)
- Concurrently run jobs (#35)
- Use npm registry (#36)
- Only sync main branch
- Add code review actions (#61)
- Add claude GitHub actions (#93)
- *(workflows)* Add allowed_bots parameter to Claude workflows (#136)
- Update claude code token (#144)
- Fix gitlab ci (#148)
- Fix license generation and allow PR runs (#156)
- Ensure git history in publish go client workflow (#158)
- Correct version tagging format in Go client publish workflow (#161)
- Add release workflow (#167)
