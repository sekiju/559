application {
  check_updates = true
  max_parallel_downloads = 4
}

output {
  directory      = "downloads"
  clean_on_start = false
  file_format    = "auto" # auto | jpeg | png | avif | webp
}

site {
  "shonenjumpplus.com" {
    cookies = ""
  }

  "comic-zenon.com" {
    cookies = ""
  }

  "pocket.shonenmagazine.com" {
    cookies = ""
  }

  "comic-gardo.com" {
    cookies = ""
  }

  "magcomi.com" {
    cookies = ""
  }

  "comic-action.com" {
    cookies = ""
  }

  "comic-days.com" {
    cookies = ""
  }

  "kuragebunch.com" {
    cookies = ""
  }

  "viewer.heros-web.com" {
    cookies = ""
  }

  "www.cmoa.jp" {
    cookies = ""
  }

  "www.corocoro.jp" {
    cookies = ""
  }
}
