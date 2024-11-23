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

#### Downloader enhanced with chapter metadata preloader

The downloader has been upgraded to include the ability to preload information about next chapters in queue. This enhancement
significantly speeds up the application's performance when loading multiple chapters, especially noticeable when dealing with magazines.

- The variable `download.concurrent_processes` has been renamed to `download.page_batch_size` ([#ccf3f5](https://github.com/sekiju/mdl/commit/ccf3f54c2f22956bd8e281593352528e5a66328f)).
- A new variable `download.preload_next_chapters` has been introduced to control the number of chapters to be preloaded ([#ccf3f5](https://github.com/sekiju/mdl/commit/ccf3f54c2f22956bd8e281593352528e5a66328f)).

The updated `download` section  now looks as follows:

```hcl
download {
  preload_next_chapters = 2
  page_batch_size       = 4
}
```

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
