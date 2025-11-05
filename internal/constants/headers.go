package constants

// HTTP Header names used across the application

// HeaderClientEncryption is the header name for enabling client-side encryption
// When set to "true", the server will encrypt vault values with a key derived from
// the API key and vault unique ID, allowing the CLI to decrypt them locally
const HeaderClientEncryption = "X-Enable-Client-Encryption"
