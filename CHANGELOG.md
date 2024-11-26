# Changelog

<!--
To add a new release, copy from this template:

## v1.X.Y

Date: YYYY-MM-DD

### What's Changed

#### Big New Feature 1

#### Big New Feature 2

### Minor Changes

### Patch Changes

-->

## v1.1.0

Date: YYYY-MM-DD

### What's Changed

### Configuration Changes

All variables have been renamed. For a detailed overview of the changes, please refer to the [example.config.hcl](https://github.com/sekiju/mdl/blob/c6bcfaf5ce28b6b73abebb6a6db97e25803f9f1e/example.config.hcl) file. Note that only the section names remain unchanged; all other variables have been updated.

- The `download` section has been removed.
  - The variable `concurrent_downloads` has been renamed to `max_parallel_downloads` and moved to the `application` section.

- `check_for_updates` has been renamed to `check_updates`.
- `dir` has been renamed to `directory`.
- `clean_dir` has been renamed to `clean_on_start`.
- `format` has been renamed to `file_format`.
- In the `site` section, `cookie_string` has been renamed to `cookie`.

#### Support for new Websites

- [www.corocoro.jp](https://www.corocoro.jp/) ([#4f4d590](https://github.com/sekiju/mdl/commit/4f4d590d606371455b803af38007edbeec047fad))
- [storia.takeshobo.co.jp](https://storia.takeshobo.co.jp/) ([#5205970](https://github.com/sekiju/mdl/commit/520597093fb45e00602c78c78e34829df4d43284))

## v1.0.1

Date: 2024-11-23

### Patch Changes

- Fix the `unsupported provider for hostname:` error when the user launches the app from
  Explorer ([#5ba31cc](https://github.com/sekiju/mdl/commit/5ba31cc023d1abb9f92adfacb8319d2310ae2760)).

## v1.0.0

Date: 2024-11-22

### What's Changed

- **Added automatic image formatting:** Support for converting images to the following formats: JPEG, PNG, AVIF and WebP
- **Configuration file format updated:** Replaced the YAML config with a new HCL (HashiCorp Configuration Language) format
- **Rewritten internal API for extensions:** The internal API has been refactored to improve extensibility and maintainability
- **New command-line arguments:**
    - `--cookie (string)`: Provides the cookie string for the current session
    - `--config (string)`: Specifies the path to the configuration file (default: `config.hcl`)
- **Support for multi-download:** Now you can pass multiple links, and all of them will be downloaded
- **Better performance:** Optimizations across the system have resulted in a +56% performance boost

#### Supported Websites

- [comic-walker.com](https://comic-walker.com)
- [shonenjumpplus.com](https://shonenjumpplus.com)
- [comic-zenon.com](https://comic-zenon.com)
- [pocket.shonenmagazine.com](https://pocket.shonenmagazine.com)
- [comic-gardo.com](https://comic-gardo.com)
- [magcomi.com](https://magcomi.com)
- [tonarinoyj.jp](https://tonarinoyj.jp)
- [comic-ogyaaa.com](https://comic-ogyaaa.com)
- [comic-action.com](https://comic-action.com)
- [comic-days.com](https://comic-days.com)
- [comic-growl.com](https://comic-growl.com)
- [comic-earthstar.com](https://comic-earthstar.com)
- [comicborder.com](https://comicborder.com)
- [comic-trail.com](https://comic-trail.com)
- [kuragebunch.com](https://kuragebunch.com)
- [viewer.heros-web.com](https://viewer.heros-web.com)
- [www.sunday-webry.com](https://www.sunday-webry.com)
- [www.cmoa.jp](https://www.cmoa.jp)
