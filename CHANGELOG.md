## [1.4.24] - 2026-02-12

### 🚀 Features

- Add API key vault update endpoint (#400)

### 💼 Other

- *(deps)* Upgrade Go version to 1.26 (#404)

### ⚙️ Miscellaneous Tasks

- Update Claude workflows and add ANTHROPIC_BASE_URL (#391)
- Move client publishing logic to standalone publish workflows (#392)
- Fix publish ts client (#393)
## [1.4.23] - 2026-01-18

### 🚀 Features

- Add change password endpoint (#382)

### ⚙️ Miscellaneous Tasks

- Rename build artifact to vault-hub (#386)
## [1.4.22] - 2026-01-18

### 🚀 Features

- Add audit log for reset password (#381)

### 🐛 Bug Fixes

- Show friendly error for duplicate email signup (#383)

### ⚙️ Miscellaneous Tasks

- *(ci)* Update GitLab CI configuration (#380)
- Add pinned frontend build script (#384)
## [1.4.18] - 2025-12-06

### ⚙️ Miscellaneous Tasks

- Update CLAUDE.md documentation for codebase (#363)
## [1.4.16] - 2025-12-03

### 🚀 Features

- Add demo config item to API (#360)

### ⚙️ Miscellaneous Tasks

- Disable credentials persistence in mirror workflow (#358)
## [1.4.14] - 2025-11-25

### ⚙️ Miscellaneous Tasks

- Bump vault-hub-go-client to v1.4.13 (#352)
## [1.4.12] - 2025-11-21

### 🚀 Features

- Add 'Favourite' flag to vault model and API schema (#349)

### ⚙️ Miscellaneous Tasks

- Update Claude workflow to use OAuth token instead of API key (#347)
## [1.4.11] - 2025-11-18

### ⚙️ Miscellaneous Tasks

- Add post-start script for devcontainer setup (#337)
## [1.4.10] - 2025-11-11

### ⚙️ Miscellaneous Tasks

- Add Catnip Feature (#335)
## [1.4.9] - 2025-11-10

### 🚀 Features

- Add audit log source filter and vault options endpoint (#334)
## [1.4.7] - 2025-11-08

### 📚 Documentation

- Merge AGENTS documentation into CLAUDE file (#328)
## [1.4.6] - 2025-11-06

### 🐛 Bug Fixes

- Inline X-Enable-Client-Encryption header handling (#325)

### 🚜 Refactor

- Simplify update-web.sh (#324)

### ⚙️ Miscellaneous Tasks

- Change model from gpt-5 to composer-1 in workflow (#326)
## [1.4.5] - 2025-10-24

### 🚀 Features

- Implement client-side encryption for vault retrieval (#312)
## [1.4.3] - 2025-10-23

### 🚀 Features

- Add pagination parameters to GetVaults API  (#309)
## [1.4.2] - 2025-10-23

### 🚀 Features

- Update build process and add script for frontend submodule (#307)
## [1.4.1] - 2025-10-22

### 🚜 Refactor

- Move frontend to separate repo (#305)
## [1.3.22] - 2025-10-22

### 🚀 Features

- Expose email enabled flag in config api (#303)
## [1.3.20] - 2025-10-22

### 🚀 Features

- Add email-first login flow with magic link and password reset (#292)
## [1.3.19] - 2025-10-17

### 🚀 Features

- Add 404 response for magic link request when email not found (#295)
## [1.3.18] - 2025-10-17

### 🚀 Features

- Add rate limit on email token requests (#296)
## [1.3.17] - 2025-10-15

### 🚀 Features

- Add resend provider with unified email config (#291)

### 🐛 Bug Fixes

- Move magic-link consume route under api namespace (#290)
## [1.3.16] - 2025-10-14

### 🚀 Features

- Add email support (#282)
## [1.3.15] - 2025-10-13

### 🚀 Features

- Add automated code review workflow using Cursor CLI (#281)
## [1.3.14] - 2025-10-12

### 🚀 Features

- Enforce password requirement for email/password signup (#279)
## [1.3.13] - 2025-10-12

### 🚀 Features

- Integrate OIDC configuration fetching in login and signup forms (#277)
## [1.3.12] - 2025-10-12

### ⚙️ Miscellaneous Tasks

- Bump api version to 1.0.1 (#276)
## [1.3.11] - 2025-10-12

### 🚀 Features

- Implement OIDC login functionality (#273)
- Add public configuration endpoint (#274)
## [1.3.10] - 2025-09-28

### 🚀 Features

- Enhance VaultValueEditor with dynamic textarea sizing (#268)

### 📚 Documentation

- Add repository guidelines (#269)

### ⚙️ Miscellaneous Tasks

- Update Claude workflows to use ANTHROPIC_API_KEY and add ANTHROPIC_BASE_URL (#267)
## [1.3.9] - 2025-09-28

### 🚀 Features

- Add vault detail page (#265)

### 🐛 Bug Fixes

- Markdown css style (#264)

### 📚 Documentation

- Update douments and readme (#263)
## [1.3.8] - 2025-09-26

### 🐛 Bug Fixes

- Remove unnecessary condition in GetAPIKeyByHash query (#261)
## [1.3.7] - 2025-09-26

### 🚀 Features

- Add environment variable support for CLI flags (#259)
## [1.3.6] - 2025-09-26

### ⚙️ Miscellaneous Tasks

- Enhance Dockerfiles and workflows to include version and commit for builds (#257)
## [1.3.5] - 2025-09-25

### 🚀 Features

- Integrate embedded static files into the web server (#255)
## [1.3.4] - 2025-09-24

### ⚙️ Miscellaneous Tasks

- Fix build version and commit (#248)
## [1.3.3] - 2025-09-24

### ⚙️ Miscellaneous Tasks

- Enhance Dockerfiles to dynamically retrieve version and commit info (#246)
## [1.3.2] - 2025-09-23

### ⚙️ Miscellaneous Tasks

- Add versioning support in Dockerfiles and workflows (#244)
## [1.3.1] - 2025-09-23

### ⚙️ Miscellaneous Tasks

- Update image tagging strategy in workflows (#242)
## [1.3.0] - 2025-09-21

### 🚀 Features

- Add Documentation page (#238)

### 📚 Documentation

- Enhance README and CLAUDE.md (#239)
## [1.2.10] - 2025-09-20

### 🚀 Features

- Add Features page (#236)
## [1.2.9] - 2025-09-20

### 🚀 Features

- Redirect to dashboard after login and signup (#230)
- Integrate system status display (#233)

### 📚 Documentation

- Update claude document (#234)
## [1.2.8] - 2025-09-20

### 🚀 Features

- Add system and database status (#231)
## [1.2.7] - 2025-09-20

### 🚀 Features

- Using Zustand for state management for audit log, api key, vault (#225)
- *(dashboard)* Update recent audit logs (#226)
- *(dashboard)* Update recent vaults (#228)

### ⚙️ Miscellaneous Tasks

- Hide navigation bar (#227)
## [1.2.6] - 2025-09-15

### 🚀 Features

- *(audit)* Include source in audit log (#219)

### 🚜 Refactor

- *(cli)* Streamline get and list commands (#218)
## [1.2.5] - 2025-09-10

### 🚀 Features

- *(docker)* Add go-cron binary and update Dockerfile for cron support (#214)

### 📚 Documentation

- Update CLAUDE.md (#212)
## [1.2.4] - 2025-09-08

### 🚀 Features

- *(docker)* Add Dockerfile for CLI and update documentation (#203)
## [1.2.3] - 2025-09-07

### 🚀 Features

- *(web)* Show metrics in audit log page (#200)

### ⚙️ Miscellaneous Tasks

- *(web)* Enhance Vite configuration with manual chunking (#201)
## [1.2.2] - 2025-09-06

### 🚀 Features

- *(audit)* Add metrics endpoints for audit logs (#198)
## [1.2.1] - 2025-09-06

### 🚀 Features

- *(web)* Update audit log metrics and refine sidebar in mobile view (#194)

### ⚙️ Miscellaneous Tasks

- Remove duplicated cocker build steps (#193)
- Skip CHANGELOG.md updates in release log (#196)
## [1.2.0] - 2025-09-05

### 🚀 Features

- Add executing follow-up commands (#191)

### ⚙️ Miscellaneous Tasks

- Update CHANGELOG.md for v1.1.6 (#189)
- *(release)* Update commit message format for CHANGELOG (#190)
- Add confirmation prompt in bump script (#192)
## [1.1.6] - 2025-09-05

### ⚙️ Miscellaneous Tasks

- Update CHANGELOG.md for v1.1.5 (#187)
- Add dry run and version increment selection in bump script (#188)
## [1.1.5] - 2025-09-05

### ⚙️ Miscellaneous Tasks

- Update token reference in release workflow (#186) (#186)
## [1.1.4] - 2025-09-05

### ⚙️ Miscellaneous Tasks

- Update CHANGELOG.md for v1.1.3 (#184)
- Update release workflow for automated changelog PRs (#185)
## [1.1.3] - 2025-09-01

### ⚙️ Miscellaneous Tasks

- Revert release workflow (#183)
## [1.1.2] - 2025-09-01

### ⚙️ Miscellaneous Tasks

- Enhance changelog update process in release workflow (#182)
## [1.1.1] - 2025-09-01

### ⚙️ Miscellaneous Tasks

- Update pull request permissions in release workflow (#181)
## [1.1.0] - 2025-08-31

### 🚀 Features

- *(web)* Add version API and display version info in dashboard (#179)

### ⚙️ Miscellaneous Tasks

- Update CHANGELOG.md for v1.0.4 (#177)
- Update release workflow to create and auto-merge changelog PR (#178)
- Format yaml (#180)
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
