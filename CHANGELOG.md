# Changelog

<!--
To add a new release, copy from this template:

## v2.X.Y

Date: YYYY-MM-DD

### What's Changed

#### Big New Feature 1

#### Big New Feature 2

### Minor Changes

### Patch Changes

-->

## v2.0.0

Date: YYYY-MM-DD

### What's Changed

- **Added automatic image formatting:** Support for converting images to the following formats: JPEG, PNG, AVIF, and WebP.
- **Configuration file format changed:** Replaced the YAML config with a new HCL (HashiCorp Configuration Language) format.
- **Rewritten internal API for extensions:** The internal API has been refactored to improve extensibility and maintainability.
- **New command-line arguments:**
  - `--download-chapter`: Specifies the URL of the chapter to download
  - `--session`: Provides the session token for the current service
  - `--config`: Specifies the path to the configuration file (default: `config.hcl`)